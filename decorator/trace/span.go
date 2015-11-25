package trace

import (
	"fmt"

	"bitbucket.org/lazadaweb/go-trace"
	"godep.lzd.co/go-kapusta"
	"godep.lzd.co/go-kapusta/decorator"
)

// SpanDecorator decorates kapusta client with Span info
func SpanDecorator(id gotrace.Span, name string) kapusta.DecoratorFunc {
	var forwardedAppsHeader string

	if id.ForwardedApps == "" {
		forwardedAppsHeader = name
	} else {
		forwardedAppsHeader = fmt.Sprintf("%s,%s", id.ForwardedApps, name)
	}

	return decorator.HeadersDecorator(map[string]string{
		gotrace.TraceIDHeader:       fmt.Sprintf("%X", id.TraceID),
		gotrace.SpanIDHeader:        fmt.Sprintf("%X", id.SpanID),
		gotrace.ParentSpanIDHeader:  fmt.Sprintf("%X", id.ParentSpanID),
		gotrace.ForwardedAppsHeader: forwardedAppsHeader,
	})
}
