package tracing

import (
	"context"

	log "github.com/sirupsen/logrus"
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

	X_SJC_SHIELD_STATUS = "sjc_shields_status"

	ContextFieldsKey = "@fields"

	shield_unhealthy_value = "0"
	shield_unhealthy_log   = "unhealthy"
)

var contextFields = [17]string{
	USER_ID,
	WIKIA_USER_ID,
	HTTP_METHOD,
	HTTP_DOMAIN,
	HTTP_PATH,
	HTTP_PARAM,
	HTTP_URL,
	BEACON_ID,
	CLIENT_IP,
	DEVICE_ID,
	FORWARDED,
	SPAN_ID,
	PARENT_SPAN_ID,
	TRACE_ID,

	WIKI_ID,
	ENVIRONMENT,
	DATA_CENTER,
}

func WithContext(c context.Context) *log.Entry {
	fields := log.Fields{}

	for _, val := range contextFields {
		if c.Value(val) != nil {
			fields[val] = c.Value(val).(string)
		}
	}

	if c.Value(X_SJC_SHIELD_STATUS) != nil && c.Value(X_SJC_SHIELD_STATUS) == shield_unhealthy_value {
		fields[X_SJC_SHIELD_STATUS] = shield_unhealthy_log
	}

	return log.WithFields(log.Fields{
		ContextFieldsKey: fields,
	})
}
