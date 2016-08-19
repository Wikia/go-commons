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

var ContextFields = [17]string{
	"user_id",
	"wikia_user_id",
	"http_method",
	"http_url_domain",
	"http_url_path",
	"http_url_query_param",
	"http_url",
	"beacon_id",
	"client_ip",
	"device_id",
	"forwarder_for",
	"span_id",
	"parent_span_id",
	"trace_id",

	"wiki_id",
	"environment",
	"datacenter",
}

func WithContext(c context.Context) *log.Entry {
	fields := log.Fields{}

	for _, val := range ContextFields {
		if c.Value(val) != nil {
			fields[val] = c.Value(val).(string)
		}
	}

	return log.WithFields(fields)
}
