package application

import (
	"os"

	"bitbucket.org/lazadaweb/go-kapusta"
	"bitbucket.org/lazadaweb/go-trace"
)

// Decorator returns a decorator that add information about current application
func Decorator(name, version string) kapusta.DecoratorFunc {
	hostname, _ := os.Hostname()

	return kapusta.HeadersDecorator(map[string]string{
		gotrace.AppNameHeader:    name,
		gotrace.AppVersionHeader: version,
		gotrace.AppNodeHeader:    hostname,
	})
}
