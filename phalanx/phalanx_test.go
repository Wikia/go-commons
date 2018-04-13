package phalanx

import (
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/Wikia/go-commons/tracing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"
)

type PhalanxTestSuite struct {
	suite.Suite
	ctx           context.Context
	requestData   url.Values
	httpResponse  *http.Response
	apiClientMock *ApiClientMock
	phalanxClient PhalanxClient
}

func (suite *PhalanxTestSuite) SetupTest() {
	suite.ctx = tracing.NewTestContext()

	suite.requestData = url.Values{}
	suite.requestData.Add(typeKey, CheckTypeName)
	suite.requestData.Add(contentKey, "SomeUserName")

	suite.httpResponse = &http.Response{Status: "200"}

	suite.apiClientMock = new(ApiClientMock)

	suite.phalanxClient = &Client{apiClient: suite.apiClientMock}
}

func (suite *PhalanxTestSuite) TestRetriesSuccess() {
	suite.apiClientMock.On("Call", "POST", checkEndpoint, suite.requestData,
		tracing.AddHttpHeadersFromContext(http.Header{}, suite.ctx)).Return(suite.httpResponse, nil).Once()
	suite.apiClientMock.On("GetBody", suite.httpResponse).Return([]byte(checkOk), nil).Once()

	result, err := suite.phalanxClient.CheckName(suite.ctx, suite.requestData.Get(contentKey))

	assert.True(suite.T(), result)
	assert.Nil(suite.T(), err)
	assert.True(suite.T(), suite.apiClientMock.AssertNumberOfCalls(suite.T(), "Call", 1))
}

func (suite *PhalanxTestSuite) TestRetriesFailedTwice() {
	suite.apiClientMock.On("Call", "POST", checkEndpoint, suite.requestData,
		tracing.AddHttpHeadersFromContext(http.Header{}, suite.ctx)).Return(suite.httpResponse, errors.New("error")).Twice()
	suite.apiClientMock.On("Call", "POST", checkEndpoint, suite.requestData,
		tracing.AddHttpHeadersFromContext(http.Header{}, suite.ctx)).Return(suite.httpResponse, nil).Once()
	suite.apiClientMock.On("GetBody", suite.httpResponse).Return([]byte(checkOk), nil).Once()

	result, err := suite.phalanxClient.CheckName(suite.ctx, suite.requestData.Get(contentKey))

	assert.True(suite.T(), result)
	assert.Nil(suite.T(), err)
	assert.True(suite.T(), suite.apiClientMock.AssertNumberOfCalls(suite.T(), "Call", 3))
}

func (suite *PhalanxTestSuite) TestRetriesFailed() {
	suite.apiClientMock.On("Call", "POST", checkEndpoint, suite.requestData,
		tracing.AddHttpHeadersFromContext(http.Header{}, suite.ctx)).Return(suite.httpResponse, errors.New("error"))

	result, err := suite.phalanxClient.CheckName(suite.ctx, suite.requestData.Get(contentKey))

	assert.False(suite.T(), result)
	assert.NotNil(suite.T(), err)
	assert.True(suite.T(), suite.apiClientMock.AssertNumberOfCalls(suite.T(), "Call", 3))
}

func TestPhalanxTestSuite(t *testing.T) {
	suite.Run(t, new(PhalanxTestSuite))
}
