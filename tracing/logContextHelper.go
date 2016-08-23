package tracing

import (
	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
)

const (
	USER_ID        = "user_id"
	WIKIA_USER_ID  = "wikia_user_id"
	HTTP_METHOD    = "http_method"
	HTTP_DOMAIN    = "http_url_domain"
	HTTP_PATH      = "http_url_path"
	HTTP_PARAM     = "http_url_query_param"
	HTTP_URL       = "http_url"
	BEACON_ID      = "beacon_id"
	CLIENT_IP      = "client_ip"
	DEVICE_ID      = "device_id"
	FORWARDED      = "forwarder_for"
	SPAN_ID        = "span_id"
	PARENT_SPAN_ID = "parent_span_id"
	TRACE_ID       = "trace_id"

	WIKI_ID     = "wiki_id"
	ENVIRONMENT = "environment"
	DATA_CENTER = "datacenter"
)

var ContextFields = map[string]string{
	USER_ID: XUserId,
	WIKIA_USER_ID: XWikiaUserId,
	HTTP_METHOD: "",
	HTTP_DOMAIN: "",
	HTTP_PATH: "",
	HTTP_PARAM: "",
	HTTP_URL: "",
	BEACON_ID: XClientBeaconId,
	CLIENT_IP: XClientIp,
	DEVICE_ID: XClientDeviceId,
	FORWARDED: XForwardedFor,
	SPAN_ID: "",
	PARENT_SPAN_ID: XParentSpanId,
	TRACE_ID: XTraceId,

	WIKI_ID: "",
	ENVIRONMENT: "",
	DATA_CENTER: "",
}

func WithContext(c context.Context) *log.Entry {
	fields := log.Fields{}

	for val, _ := range ContextFields {
		if c.Value(val) != nil {
			fields[val] = c.Value(val).(string)
		}
	}

	return log.WithFields(log.Fields{
		"@fields": fields,
	})
}
