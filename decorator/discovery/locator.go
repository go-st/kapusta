package discovery

import (
	"fmt"
	"net/http"
	"net/url"

	"bitbucket.org/lazadaweb/discovery/discovery"
	"godep.lzd.co/go-kapusta"
)

// LocatorDecorator returns a DecoratorFunc that gets URL from Locator and replace host and scheme in request URL
func LocatorDecorator(locator discovery.IServiceLocator) kapusta.DecoratorFunc {
	return func(c kapusta.IClient) kapusta.IClient {
		return kapusta.ClientFunc(func(r *http.Request) (*http.Response, error) {
			locatedURL, err := locator.Locate()

			if err != nil {
				return nil, fmt.Errorf("error while getting URL from locator: %s", err)
			}

			parsed, err := url.Parse(locatedURL)

			if err != nil {
				return nil, fmt.Errorf("error while parsing URL %s: %s", locatedURL, err)
			}

			r.URL.Scheme = parsed.Scheme
			r.URL.Host = parsed.Host

			return c.Do(r)
		})
	}
}
