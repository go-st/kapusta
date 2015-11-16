package trace

import (
	"os"

	"bitbucket.org/lazadaweb/go-kapusta"
	"bitbucket.org/lazadaweb/go-kapusta/decorator"
	"bitbucket.org/lazadaweb/go-trace"
)

// ApplicationDecorator returns a decorator that add information about current application
func ApplicationDecorator(name, version string) kapusta.DecoratorFunc {
	hostname, _ := os.Hostname()

	return decorator.HeadersDecorator(map[string]string{
		gotrace.AppNameHeader:    name,
		gotrace.AppVersionHeader: version,
		gotrace.NodeHeader:       hostname,
	})
}
