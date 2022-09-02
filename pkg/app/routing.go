package app

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

// createAllRoutes is creating a set of routes basing on the whole application configuration
func createAllRoutes(config Config) map[string]func(w http.ResponseWriter, req *http.Request) {
	funcs := make(map[string]func(w http.ResponseWriter, req *http.Request))
	for _, endpoint := range config.Endpoints {
		logrus.Infof("Registering endpoint: '%s' -> '%s'", endpoint.Path, endpoint.TargetUrl)
		funcs[endpoint.Path] = createRoute(endpoint)
	}
	return funcs
}

func createMutator(endpoint *Endpoint, body io.ReadCloser) io.ReadCloser {
	return StreamMutator{body, endpoint.Replacements}
}

// createRoute is creating a single route for single endpoint
func createRoute(endpoint Endpoint) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, incomingReq *http.Request) {
		c := http.Client{}
		c.Timeout = endpoint.getTimeout()

		req, err := http.NewRequest(incomingReq.Method, endpoint.TargetUrl+"?"+incomingReq.URL.RawQuery, createMutator(&endpoint, incomingReq.Body))
		logrus.Infof("Creating request to '%s'", endpoint.TargetUrl)
		req.Form = incomingReq.Form
		req.MultipartForm = incomingReq.MultipartForm
		req.PostForm = incomingReq.PostForm

		if err != nil {
			_, _ = fmt.Fprintf(w, "Cannot proxy request: %s", err.Error())
			return
		}

		response, pErr := c.Do(req)
		if pErr != nil {
			_, _ = fmt.Fprintf(w, "Cannot proxy request: %s", pErr.Error())
			return
		}

		if len(response.Header) > 0 {
			// headers
			for header, headerValues := range response.Header {
				for _, headerValue := range headerValues {
					w.Header().Add(header, headerValue)
				}
			}
		}

		// body
		_, wErr := io.Copy(w, createMutator(&endpoint, response.Body))
		if wErr != nil {
			_, _ = fmt.Fprintf(w, "Cannot read response body: %s", wErr.Error())
		}
	}
}
