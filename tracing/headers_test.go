package tracing

import (
	"github.com/go-playground/lars"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestShouldCreateRequestUsingContextDataWhereHeadersAreNotEmpty(t *testing.T) {
	req := NewTestRequest()
	res := httptest.NewRecorder()

	c := lars.NewContext(lars.New())
	c.RequestStart(res, req)
	defer c.RequestEnd()
	c = ContextSetHandlerTest(c)

	newReq, _ := http.NewRequest("POST", "/theMiddleOfNowhere", nil)
	ReturnRequestWithHeadersGivenAsMap(newReq, GetHeadersFromContextAsMap(c))

	expected := GetTestHeadersAsMap()

	assert.Equal(t, expected[XTraceId], newReq.Header.Get(XTraceId), "X Trace Id Header Should Have Different Value")
	assert.Equal(t, expected[XClientBeaconId], newReq.Header.Get(XClientBeaconId), "X Client Beacon Id Header Should Have Different Value")
	assert.Equal(t, expected[XClientDeviceId], newReq.Header.Get(XClientDeviceId), "X Client Device Id Header Should Have Different Value")
	assert.Equal(t, expected[XClientIp], newReq.Header.Get(XClientIp), "X Client Ip Header Should Have Different Value")
	assert.Equal(t, expected[XUserId], newReq.Header.Get(XUserId), "X User Id Header Should Have Different Value")
	assert.Equal(t, expected[XParentSpanId], newReq.Header.Get(XParentSpanId), "X Parent Span Id Header Should Have Different Value")
}

func TestShouldCreateRequestUsingContextDataWhereHeadersAreEmpty(t *testing.T) {
	req, _ := http.NewRequest("POST", "/theMiddleOfNowhere", nil)
	res := httptest.NewRecorder()

	c := lars.NewContext(lars.New())
	c.RequestStart(res, req)
	defer c.RequestEnd()
	c = ContextSetHandlerTest(c)

	newReq, _ := http.NewRequest("POST", "/theMiddleOfNowhere", nil)
	ReturnRequestWithHeadersGivenAsMap(newReq, GetHeadersFromContextAsMap(c))

	assert.Empty(t, newReq.Header.Get(XTraceId), "X Trace Id Header Should Be Empty")
	assert.Empty(t, newReq.Header.Get(XClientBeaconId), "X Client Beacon Id Header Should Be Empty")
	assert.Empty(t, newReq.Header.Get(XClientDeviceId), "X Client Device Id Header Should Be Empty")
	assert.Empty(t, newReq.Header.Get(XClientIp), "X Client Ip Header Should Be Empty")
	assert.Empty(t, newReq.Header.Get(XUserId), "X User Id Header Should Be Empty")
	assert.Empty(t, newReq.Header.Get(XParentSpanId), "X Parent Span Id Header Should Be Empty")
}

func ContextSetHandlerTest(c *lars.Ctx) *lars.Ctx {
	r := c.Request()

	c.WithValue(TRACE_ID, r.Header.Get(XTraceId))
	c.WithValue(BEACON_ID, r.Header.Get(XClientBeaconId))
	c.WithValue(DEVICE_ID, r.Header.Get(XClientDeviceId))
	c.WithValue(CLIENT_IP, r.Header.Get(XClientIp))
	c.WithValue(USER_ID, r.Header.Get(XUserId))
	c.WithValue(SPAN_ID, r.Header.Get(XSpanId))
	c.WithValue(PARENT_SPAN_ID, r.Header.Get(XParentSpanId))

	return c
}
