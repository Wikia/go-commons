package tracing

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"net/http"
	"testing"
)

func TestShouldCreateRequestUsingContextDataWhereHeadersAreNotEmpty(t *testing.T) {
	req := NewTestRequest()

	c := context.TODO()
	c = ContextSetHandlerTest(c, req)

	newReq, _ := http.NewRequest("POST", "/theMiddleOfNowhere", nil)
	SetContextHeaders(&newReq.Header, c)

	expected := GetTestHeadersAsMap()

	assert.Equal(t, expected[XTraceId], newReq.Header.Get(XTraceId), "X Trace Id Header Should Have Different Value")
	assert.Equal(t, expected[XClientBeaconId], newReq.Header.Get(XClientBeaconId), "X Client Beacon Id Header Should Have Different Value")
	assert.Equal(t, expected[XClientDeviceId], newReq.Header.Get(XClientDeviceId), "X Client Device Id Header Should Have Different Value")
	assert.Equal(t, expected[XClientIp], newReq.Header.Get(XClientIp), "X Client Ip Header Should Have Different Value")
	assert.Equal(t, expected[XUserId], newReq.Header.Get(XUserId), "X User Id Header Should Have Different Value")
	assert.Equal(t, expected[XParentSpanId], newReq.Header.Get(XParentSpanId), "X Parent Span Id Header Should Have Different Value")
	assert.Equal(t, expected[XSJCShieldsHealthy], newReq.Header.Get(XSJCShieldsHealthy), "X SJC Shields Healthy header should have different value")
}

func TestShouldCreateRequestUsingContextDataWhereHeadersAreEmpty(t *testing.T) {
	req, _ := http.NewRequest("POST", "/theMiddleOfNowhere", nil)

	c := context.TODO()
	c = ContextSetHandlerTest(c, req)

	newReq, _ := http.NewRequest("POST", "/theMiddleOfNowhere", nil)
	SetContextHeaders(&newReq.Header, c)

	assert.Empty(t, newReq.Header.Get(XTraceId), "X Trace Id Header Should Be Empty")
	assert.Empty(t, newReq.Header.Get(XClientBeaconId), "X Client Beacon Id Header Should Be Empty")
	assert.Empty(t, newReq.Header.Get(XClientDeviceId), "X Client Device Id Header Should Be Empty")
	assert.Empty(t, newReq.Header.Get(XClientIp), "X Client Ip Header Should Be Empty")
	assert.Empty(t, newReq.Header.Get(XUserId), "X User Id Header Should Be Empty")
	assert.Empty(t, newReq.Header.Get(XParentSpanId), "X Parent Span Id Header Should Be Empty")
}

func ContextSetHandlerTest(c context.Context, r *http.Request) context.Context {

	c = context.WithValue(c, TRACE_ID, r.Header.Get(XTraceId))
	c = context.WithValue(c, TRACE_ID, r.Header.Get(XTraceId))
	c = context.WithValue(c, BEACON_ID, r.Header.Get(XClientBeaconId))
	c = context.WithValue(c, DEVICE_ID, r.Header.Get(XClientDeviceId))
	c = context.WithValue(c, CLIENT_IP, r.Header.Get(XClientIp))
	c = context.WithValue(c, USER_ID, r.Header.Get(XUserId))
	c = context.WithValue(c, SPAN_ID, r.Header.Get(XSpanId))
	c = context.WithValue(c, PARENT_SPAN_ID, r.Header.Get(XParentSpanId))
	c = context.WithValue(c, X_SJC_SHIELD_STATUS, r.Header.Get(XSJCShieldsHealthy))

	return c
}
