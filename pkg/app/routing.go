package app

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
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

		logrus.Debug("Creating request with incomingReq.Body as body")
		req, err := http.NewRequest(incomingReq.Method, endpoint.TargetUrl+"?"+incomingReq.URL.RawQuery, createMutator(&endpoint, incomingReq.Body))
		logrus.Infof("Creating request to '%s'", endpoint.TargetUrl)
		req.Form = incomingReq.Form
		req.MultipartForm = incomingReq.MultipartForm
		req.PostForm = incomingReq.PostForm

		for header, headerValues := range incomingReq.Header {
			for _, headerValue := range headerValues {
				// gzip, deflate, br etc. encoding would make content processing not possible - try to disable it
				h := strings.ToLower(header)
				if h == "accept-encoding" || h == "content-encoding" || h == "content-length" {
					logrus.Debugf("Skipping request Header: %s", header)
					continue
				}

				logrus.Debugf("Forward request Header: %s = %s", header, headerValue)
				req.Header.Set(header, headerValue)
			}
		}

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
					// we are modifying content, so the content-length header must be wiped and rewritten again
					// it will be automatically generated
					if strings.ToLower(header) == "content-length" {
						logrus.Debug("Skipping content-length header")
						continue
					}

					logrus.Debugf("Forward response Header: %s = %s", header, headerValue)
					w.Header().Add(header, headerValue)
				}
			}
		}

		// body
		logrus.Debug("Copying response.Body")
		_, wErr := io.Copy(w, createMutator(&endpoint, response.Body))
		if wErr != nil {
			_, _ = fmt.Fprintf(w, "Cannot read response body: %s", wErr.Error())
		}
	}
}
