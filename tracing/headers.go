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

	if c.Value(TRACE_ID) != nil {
		headers[XTraceId] = c.Value(TRACE_ID).(string)
	}
	if c.Value(BEACON_ID) != nil {
		headers[XClientBeaconId] = c.Value(BEACON_ID).(string)
	}
	if c.Value(DEVICE_ID) != nil {
		headers[XClientDeviceId] = c.Value(DEVICE_ID).(string)
	}
	if c.Value(CLIENT_IP) != nil {
		headers[XClientIp] = c.Value(CLIENT_IP).(string)
	}
	if c.Value(USER_ID) != nil {
		headers[XUserId] = c.Value(USER_ID).(string)
	}
	if c.Value(WIKIA_USER_ID) != nil {
		headers[XWikiaUserId] = c.Value(WIKIA_USER_ID).(string)
	}
	if c.Value(FORWARDED) != nil {
		headers[XForwardedFor] = c.Value(FORWARDED).(string)
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
