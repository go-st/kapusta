package decorator

import (
	"fmt"
	"net/http"
	"net/url"

	"godep.lzd.co/go-kapusta"
)

// HeaderDecorator returns a DecoratorFunc that adds the given HTTP header to every request done by a Client.
func HeaderDecorator(name, value string) kapusta.DecoratorFunc {
	return HeadersDecorator(map[string]string{name: value})
}

// HeadersDecorator returns a DecoratorFunc that adds the given HTTP headers to every request done by a Client.
func HeadersDecorator(values map[string]string) kapusta.DecoratorFunc {
	return func(c kapusta.IClient) kapusta.IClient {
		return kapusta.ClientFunc(func(r *http.Request) (*http.Response, error) {
			for key, value := range values {
				r.Header.Add(key, value)
			}
			return c.Do(r)
		})
	}
}

// RecoverDecorator returns a DecoratorFunc that recovers panic and convert it to error
func RecoverDecorator() kapusta.DecoratorFunc {
	return func(c kapusta.IClient) kapusta.IClient {
		return kapusta.ClientFunc(func(r *http.Request) (res *http.Response, err error) {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("recovered panic: %v", r)
				}
			}()
			return c.Do(r)
		})
	}
}

// BaseURLDecorator returns a DecoratorFunc that replaces scheme and host in request from specified baseURL
func BaseURLDecorator(baseURL string) kapusta.DecoratorFunc {
	parsed, err := url.Parse(baseURL)

	if err != nil {
		panic(err)
	}

	return func(c kapusta.IClient) kapusta.IClient {
		return kapusta.ClientFunc(func(r *http.Request) (*http.Response, error) {
			r.URL.Scheme = parsed.Scheme
			r.URL.Host = parsed.Host

			return c.Do(r)
		})
	}
}
