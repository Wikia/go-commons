package perfmonitoring

import (
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockedPerfMonitoring struct {
	mock.Mock
}

func (perfMonitoring *MockedPerfMonitoring) Set(name string, value interface{}) {

	perfMonitoring.Mock.Called(name, value)
}

func (perfMonitoring *MockedPerfMonitoring) Push() error {

	args := perfMonitoring.Mock.Called()
	return args.Error(0)
}

func CloseTimer(t *testing.T, timer *Timer) {
	err := timer.Close()
	if err != nil {
		t.Fatal("Could not close timer", err)
	}
}

func TimerMeasuredFunc(t *testing.T, perfMon *MockedPerfMonitoring) {

	timer := NewTimer(perfMon, "test_time")
	timer.AddValue("additional_value", "Value")
	defer CloseTimer(t, timer) //Recommended practice is to use defer for closing the timer

	time.Sleep(time.Duration(10) * time.Millisecond)
}

func TestTimer(t *testing.T) {

	perfMon := new(MockedPerfMonitoring)
	perfMon.On("Push").Return(nil)
	perfMon.On("Set", "additional_value", "Value").Return()
	perfMon.On("Set", "test_time", mock.AnythingOfType("int64")).Return()

	TimerMeasuredFunc(t, perfMon)

	perfMon.Mock.AssertExpectations(t)
}
