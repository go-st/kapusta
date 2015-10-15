package balancer

import (
	"fmt"
	"net/http"
	"net/url"

	"bitbucket.org/lazadaweb/discovery/balancer"
	"bitbucket.org/lazadaweb/go-kapusta"
)

// Decorator returns a DecoratorFunc that gets URL from Balancer and replace host and scheme in request URL
func Decorator(balancer balancer.LoadBalancer) kapusta.DecoratorFunc {
	return func(c kapusta.IClient) kapusta.IClient {
		return kapusta.ClientFunc(func(r *http.Request) (*http.Response, error) {
			balancedURL, err := balancer.Next()

			if err != nil {
				return nil, fmt.Errorf("error while getting URL from locator: %s", err)
			}

			parsed, err := url.Parse(balancedURL)

			if err != nil {
				return nil, fmt.Errorf("error while parsing URL %s: %s", balancedURL, err)
			}

			r.URL.Scheme = parsed.Scheme
			r.URL.Host = parsed.Host

			return c.Do(r)
		})
	}
}
