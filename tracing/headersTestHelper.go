package tracing

import (
	"context"
	"net/http"
)

func GetTestHeadersAsMap() map[string]string {
	headers := map[string]string{}

	headers[XTraceId] = "12345"
	headers[XClientBeaconId] = "9876"
	headers[XClientDeviceId] = "546845"
	headers[XClientIp] = "1.1.1.1"
	headers[XUserId] = "2.2.2.2"
	headers[XSpanId] = "54gyr54g45"
	headers[XParentSpanId] = "54gyr54g45"
	headers[XSJCShieldsHealthy] = "0"

	return headers
}

func NewTestContext() context.Context {
	c := context.TODO()
	for header, value := range GetTestHeadersAsMap() {
		c = context.WithValue(c, header, value)
	}
	return c
}

func NewTestRequest() *http.Request {
	request, _ := http.NewRequest("POST", "/theMiddleOfNowhere", nil)

	for header, value := range GetTestHeadersAsMap() {
		request.Header.Add(header, value)
	}

	return request
}
