package trace

import (
	"fmt"

	"bitbucket.org/lazadaweb/go-kapusta"
	"bitbucket.org/lazadaweb/go-trace"
)

// Decorator decorates kapusta client with Span info
func Decorator(id gotrace.Span, name string) kapusta.DecoratorFunc {
	var forwardedAppsHeader string

	if id.ForwardedApps == "" {
		forwardedAppsHeader = name
	} else {
		forwardedAppsHeader = fmt.Sprintf("%s,%s", id.ForwardedApps, name)
	}

	return kapusta.HeadersDecorator(map[string]string{
		gotrace.TraceIDHeader:       fmt.Sprintf("%X", id.TraceID),
		gotrace.SpanIDHeader:        fmt.Sprintf("%X", id.SpanID),
		gotrace.ParentSpanIDHeader:  fmt.Sprintf("%X", id.ParentSpanID),
		gotrace.ForwardedAppsHeader: forwardedAppsHeader,
	})
}
