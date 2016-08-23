package tracing

import (
	"golang.org/x/net/context"
	"net/http"
)

const (
	XTraceId        string = "X-Trace-Id"
	XClientBeaconId string = "X-Client-Beacon-Id"
	XClientDeviceId string = "X-Client-Device-Id"
	XClientIp       string = "X-Client-Ip"
	XUserId         string = "X-User-Id"
	XSpanId         string = "X-Span-Id"
	XParentSpanId   string = "X-Parent-Span-Id"
	XWikiaUserId	string = "X-Wikia-UserId"
	XForwardedFor	string = "X-Forwarded-For"
)

func GetHeadersFromContextAsMap(c context.Context) map[string]string {
	headers := map[string]string{}

	for key, val := range ContextFields {
		if val != "" && c.Value(key) != nil {
			headers[val] = c.Value(key).(string)
		}
	}

	if c.Value(SPAN_ID) != nil {
		headers[XParentSpanId] = c.Value(SPAN_ID).(string)
	}

	return headers
}

func ReturnRequestWithHeadersGivenAsMap(r *http.Request, dataMap map[string]string) {
	for header, value := range dataMap {
		r.Header.Add(header, value)
	}
}
