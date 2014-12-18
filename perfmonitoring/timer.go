package perfmonitoring

import (
	"time"
)

/*
This class enables easy tracking of time needed to perform a task and then storage
of the result in influx db
*/

type Timer struct {
	columnName string
	perfMon    IPerfMonitoring
	startTime  time.Time
}

type IPerfMonitoring interface {
	Set(columnName string, value interface{})
	Push() error
}

func NewTimer(perfMon IPerfMonitoring, columnName string) *Timer {
	timer := new(Timer)
	timer.columnName = columnName
	timer.perfMon = perfMon
	timer.startTime = time.Now()
	return timer
}

func (timer *Timer) AddValue(columnName string, value interface{}) {
	timer.perfMon.Set(columnName, value)
}

func (timer *Timer) Close() error {

	measuredTime := time.Now().Sub(timer.startTime).Nanoseconds() / 1000
	timer.perfMon.Set(timer.columnName, measuredTime)
	return timer.perfMon.Push()
}
