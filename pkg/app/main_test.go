package app_test

import (
	"bytes"
	"github.com/keskad/webhook-conversion-service/pkg/app"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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
		_, _ = res.Write([]byte("capitalism means X for 99%\n\n\n\n\n\n\n\n\n\n\n"))
	}))
	defer func() { testServer.Close() }()

	// create config that includes test server
	config := "endpoints:\n    - path: /api/webhook\n      targetUrl: " + testServer.URL + "\n      replacements:\n          - from: X\n            to: poverty"
	_ = ioutil.WriteFile("../../.build/test.yaml", []byte(config), 0755)

	// create application instance
	go app.Main("../../.build/test.yaml", "127.0.0.1:8081", app.HttpServer{})
	time.Sleep(time.Second * 5)

	// perform a test request through the reverse proxy
	req, _ := http.NewRequest(http.MethodGet, "http://127.0.0.1:8081/api/webhook", nil)
	response, reqErr := http.DefaultClient.Do(req)

	assert.Nil(t, reqErr)
	responseB, _ := ioutil.ReadAll(response.Body)
	assert.Contains(t, string(responseB), "capitalism means poverty for 99%")
}
