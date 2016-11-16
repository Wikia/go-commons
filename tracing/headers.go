package tracing

import (
	"golang.org/x/net/context"
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
	XSJCShieldsHealthy string = "X-SJC-shields-healthy"
)

var ContextHeaderFields = map[string]string{
	USER_ID: XUserId,
	WIKIA_USER_ID: XWikiaUserId,
	BEACON_ID: XClientBeaconId,
	CLIENT_IP: XClientIp,
	DEVICE_ID: XClientDeviceId,
	FORWARDED: XForwardedFor,
	PARENT_SPAN_ID: XParentSpanId,
	TRACE_ID: XTraceId,
	X_SJC_SHIELD: XSJCShieldsHealthy,
}

func GetHeadersFromContextAsMap(c context.Context) map[string]string {
	headers := map[string]string{}

	for key, val := range ContextHeaderFields {
		if c.Value(key) != nil {
			headers[val] = c.Value(key).(string)
		}
	}

	if c.Value(SPAN_ID) != nil {
		headers[XParentSpanId] = c.Value(SPAN_ID).(string)
	}

	return headers
}
