package app_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/keskad/webhook-conversion-service/pkg/app"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestMainFuncMocked(t *testing.T) {
	buf := bytes.Buffer{}
	logrus.SetOutput(&buf)

	err := app.Main("../../example-config.yaml", ":8080", app.MockServer{})
	assert.Nil(t, err)
	assert.Contains(t, buf.String(), "/api/webhook")
	assert.Contains(t, buf.String(), "/test")
}

func TestMainFunc_CallsWebhook(t *testing.T) {
	// create a mocked upstream
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("How-to-fix", "Implement and grassroot direct democracy")
		res.WriteHeader(200)
		_, writeErr := res.Write([]byte("capitalism means X for 99%"))
		assert.Nil(t, writeErr)
	}))
	defer func() { testServer.Close() }()

	// create config that includes test server
	config := "endpoints:\n    - path: /api/webhook\n      targetUrl: " + testServer.URL + "\n      replacements:\n          - from: X\n            to: poverty"
	_ = ioutil.WriteFile("../../.build/test.yaml", []byte(config), 0755)

	// create application instance
	go app.Main("../../.build/test.yaml", "127.0.0.1:8081", app.HttpServer{})
	time.Sleep(time.Second * 2)

	// perform a test request through the reverse proxy
	req, _ := http.NewRequest(http.MethodGet, "http://127.0.0.1:8081/api/webhook", nil)
	response, reqErr := http.DefaultClient.Do(req)

	assert.Nil(t, reqErr)
	responseB, _ := ioutil.ReadAll(response.Body)

	// by checking if EQUALS we check:
	//   - if it does contain
	//   - if there are no null bytes in the string
	assert.Equal(t, string(responseB), "capitalism means poverty for 99%")
}

func TestMainFunc_ProcessesGiteaWebhook(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		buff, _ := io.ReadAll(req.Body)
		_, _ = res.Write(buff)
	}))
	defer func() { testServer.Close() }()

	// config file
	config := "endpoints:\n    - path: /gitea\n      targetUrl: " + testServer.URL + "\n      replacements:\n          - from: example.org\n            to: domain.org"
	_ = ioutil.WriteFile("../../.build/test.yaml", []byte(config), 0755)

	// create application instance
	go app.Main("../../.build/test.yaml", "127.0.0.1:8081", app.HttpServer{})
	time.Sleep(time.Second * 2)

	// get test body
	w, _ := os.Open("main_test_testdata_gitea.json")

	// perform a test request through the reverse proxy
	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8081/gitea", w)
	assert.Nil(t, err)

	response, reqErr := http.DefaultClient.Do(req)
	assert.Nil(t, reqErr)

	buff, _ := io.ReadAll(response.Body)
	responseAsStr := string(buff)

	// check if it is a valid JSON
	var i interface{}
	assert.Nil(t, json.Unmarshal(buff, &i), fmt.Sprintf("Expected, that response will be a valid JSON. Response: %v", responseAsStr))

	// check if example.org is not present anymore
	assert.NotContainsf(t, responseAsStr, "https://git.example.org", "Expected, that git.example.org will be no longer present - as it should be replaced to 'git.domain.org'")
	assert.Contains(t, responseAsStr, "https://git.domain.org/my-org/", "Expected, that URL will be changed")
}
