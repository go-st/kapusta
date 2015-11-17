package monitoring

import (
	"fmt"
	"net/http"
	"time"

	"bitbucket.org/lazadaweb/go-kapusta"
	"bitbucket.org/lazadaweb/go-logger"
	"bitbucket.org/lazadaweb/go-metrics"
)

const longRequestThreshold = 500 * time.Millisecond

// DecoratorConfig configuration for decorator contains of few metrics:
//
// ResponseTimer - measures overall request time
// LongRequestCounter - increments if request took more than LongRequestThreshold time
// NotFoundCounter - increments if response returned 404 http status
// ErrorRequestCounter - increments for all over non 2xx statuses
//
// All of these fields are optional. If they aren`t set Nil(Timer|Counter) will be used.
type DecoratorConfig struct {
	ResponseTimer       metrics.Timer
	LongRequestCounter  metrics.Counter
	ErrorRequestCounter metrics.Counter
	NotFoundCounter     metrics.Counter

	LongRequestThreshold time.Duration
}

// GetResponseTimer returns ResponseTimer if it was set, NilTimer otherwise
func (c DecoratorConfig) GetResponseTimer() metrics.Timer {
	if c.ResponseTimer == nil {
		return metrics.NilTimer{}
	}

	return c.ResponseTimer
}

// GetLongRequestCounter returns LongRequestCounter if it was set, NilCounter otherwise
func (c DecoratorConfig) GetLongRequestCounter() metrics.Counter {
	if c.LongRequestCounter == nil {
		return metrics.NilCounter{}
	}

	return c.LongRequestCounter
}

// GetErrorRequestCounter returns ErrorRequestCounter if it was set, NilCounter otherwise
func (c DecoratorConfig) GetErrorRequestCounter() metrics.Counter {
	if c.ErrorRequestCounter == nil {
		return metrics.NilCounter{}
	}

	return c.ErrorRequestCounter
}

// GetNotFoundCounter returns NotFoundCounter if it was set, NilCounter otherwise
func (c DecoratorConfig) GetNotFoundCounter() metrics.Counter {
	if c.NotFoundCounter == nil {
		return metrics.NilCounter{}
	}

	return c.NotFoundCounter
}

// GetRequestThreshold returns RequestThreshold if it was set, default longRequestThreshold otherwise
func (c DecoratorConfig) GetRequestThreshold() time.Duration {
	if c.LongRequestThreshold > 0 {
		return c.LongRequestThreshold
	}

	return longRequestThreshold
}

// MonitoringDecorator returns decorator which monitors calls to external services
func MonitoringDecorator(serviceName string, configuration DecoratorConfig, logger logger.ILogger) kapusta.DecoratorFunc {
	return func(c kapusta.IClient) kapusta.IClient {
		return kapusta.ClientFunc(func(r *http.Request) (*http.Response, error) {
			t0 := time.Now()
			res, err := c.Do(r)
			apiCallDuration := time.Since(t0)
			configuration.GetResponseTimer().Update(apiCallDuration)
			if apiCallDuration > configuration.GetRequestThreshold() {
				responseDetails := ""
				if res != nil {
					responseDetails += fmt.Sprintf("Response: %d (%d bytes)", res.StatusCode, res.ContentLength)
				} else {
					responseDetails += "Response: nil"
				}
				if err != nil {
					// don't write Error: <nil> in logs - it breaks our filters
					responseDetails += fmt.Sprintf(" Error: %v", err)
				}
				logger.Warningf("Call to %s (%s) took %s: %s", serviceName, r.URL, apiCallDuration, responseDetails)
				configuration.GetLongRequestCounter().Inc(1)
			}
			if err != nil || (res.StatusCode >= http.StatusBadRequest && res.StatusCode != http.StatusNotFound) {
				configuration.GetErrorRequestCounter().Inc(1)
			}
			if res != nil && res.StatusCode == http.StatusNotFound {
				configuration.GetNotFoundCounter().Inc(1)
			}
			return res, err
		})
	}
}
