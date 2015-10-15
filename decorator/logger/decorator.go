package logger

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httputil"

	"bitbucket.org/lazadaweb/go-kapusta"
)

// ILogger logger interface
type ILogger interface {
	Debugf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// Decorator returns DecoratorFunc that logs before and after request
func Decorator(logger ILogger, dumpRequests bool) kapusta.DecoratorFunc {
	return func(c kapusta.IClient) kapusta.IClient {
		return kapusta.ClientFunc(func(r *http.Request) (*http.Response, error) {
			logger.Debugf("start request: %v", r.URL)
			var bodyBytes []byte

			if dumpRequests {
				// We should preserve Body before Do and restore it after Do, because Body is io.ReadCloser and
				// it will have already been read after Do. DumpRequestOut is going to read Body again.

				if r.Body != nil {
					bodyBytes, _ = ioutil.ReadAll(r.Body)
				}
				// Restore the io.ReadCloser to its original state
				r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
			}

			response, err := c.Do(r)

			if dumpRequests {
				// Restore the io.ReadCloser to its original state again, because it has been already read by Do
				r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

				if requestDump, requestReadErr := httputil.DumpRequestOut(r, true); requestReadErr == nil {
					logger.Debugf("dump request to %v:\n%s", r.URL, string(requestDump))
				} else {
					logger.Debugf("can't dump request to %v\nError: %v", r.URL, requestReadErr)
				}
			}

			if err != nil {
				logger.Errorf("done request: %v, error: %v", r.URL, err)
				return nil, err
			}

			logger.Debugf("done request: %v, HTTP status: %s", r.URL, response.Status)

			if dumpRequests {
				if responseDump, responseReadErr := httputil.DumpResponse(response, true); responseReadErr == nil {
					logger.Debugf("dump response from %v:\n%s", r.URL, string(responseDump))
				} else {
					logger.Debugf("can't dump response from %v\nError: %v", r.URL, responseReadErr)
				}
			}

			return response, err
		})
	}
}
