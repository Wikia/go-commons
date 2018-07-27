package tracing

import (
	"net/http"

	"context"
)

const (
	XTraceId           = "X-Trace-Id"
	XClientBeaconId    = "X-Client-Beacon-Id"
	XClientDeviceId    = "X-Client-Device-Id"
	XClientIp          = "X-Client-Ip"
	XUserId            = "X-User-Id"
	XSpanId            = "X-Span-Id"
	XParentSpanId      = "X-Parent-Span-Id"
	XWikiaUserId       = "X-Wikia-UserId"
	XForwardedFor      = "X-Forwarded-For"
	XSJCShieldsHealthy = "X-SJC-shields-healthy"
)

var ContextHeaderFields = map[string]string{
	USER_ID:             XUserId,
	WIKIA_USER_ID:       XWikiaUserId,
	BEACON_ID:           XClientBeaconId,
	CLIENT_IP:           XClientIp,
	DEVICE_ID:           XClientDeviceId,
	FORWARDED:           XForwardedFor,
	PARENT_SPAN_ID:      XParentSpanId,
	TRACE_ID:            XTraceId,
	X_SJC_SHIELD_STATUS: XSJCShieldsHealthy,
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

func AddHttpHeadersFromContext(headers http.Header, c context.Context) http.Header {
	for header, val := range GetHeadersFromContextAsMap(c) {
		if len(headers[header]) == 0 {
			headers[header] = []string{val}
		}
	}
	return headers
}
