package monitoring

import (
	"bitbucket.org/lazadaweb/go-metrics"
	. "gopkg.in/check.v1"
)

type TestMonitoringDecoratorSuite struct{}

var _ = Suite(&TestMonitoringDecoratorSuite{})

func (s *TestMonitoringDecoratorSuite) TestEmptyConfiguration(c *C) {
	configuration := DecoratorConfig{}
	c.Assert(configuration.GetResponseTimer(), FitsTypeOf, metrics.NilTimer{})
	c.Assert(configuration.GetLongRequestCounter(), FitsTypeOf, metrics.NilCounter{})
	c.Assert(configuration.GetErrorRequestCounter(), FitsTypeOf, metrics.NilCounter{})
	c.Assert(configuration.GetNotFoundCounter(), FitsTypeOf, metrics.NilCounter{})
	c.Assert(configuration.GetRequestThreshold(), Equals, longRequestThreshold)
}

func (s *TestMonitoringDecoratorSuite) TestGetResponseTimer(c *C) {
	timer := metrics.NewTimer()
	configuration := DecoratorConfig{
		ResponseTimer: timer,
	}
	c.Assert(configuration.GetResponseTimer(), Equals, timer)
}

func (s *TestMonitoringDecoratorSuite) TestGetLongRequestCounter(c *C) {
	counter := metrics.NewCounter()
	configuration := DecoratorConfig{
		LongRequestCounter: counter,
	}
	c.Assert(configuration.GetLongRequestCounter(), Equals, counter)
}

func (s *TestMonitoringDecoratorSuite) TestGetErrorRequestCounter(c *C) {
	counter := metrics.NewCounter()
	configuration := DecoratorConfig{
		ErrorRequestCounter: counter,
	}
	c.Assert(configuration.GetErrorRequestCounter(), Equals, counter)
}

func (s *TestMonitoringDecoratorSuite) TestNotFoundCounter(c *C) {
	counter := metrics.NewCounter()
	configuration := DecoratorConfig{
		NotFoundCounter: counter,
	}
	c.Assert(configuration.GetNotFoundCounter(), Equals, counter)
}
