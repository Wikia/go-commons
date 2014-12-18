go-commons/perfmonitoring
==========

Common library for storing data in InfluxDB.

Sample usage:
```go
perfMon, err := NewPerfMonitoring("my_app", "metrics")
if err != nil {
	t.Fatal(err)
}

perfMon.Set("ResponseTime", []interface{}{measuredValue})

perfMon.Push()
```

You can also use the timer object which automates time tracking for a given scope:
```go
func closeTimer(timer *perfmonitoring.Timer) {
	err := timer.Close()
	if err != nil {
		logger.GetLogger().ErrorErr(err)
	}
}

func TimerMeasuredFunc() {

	timer := NewTimer(perfMon, "test_time")//Starts time measurement
	defer closeTimer(timer) //After function finishes measured time will be pushed to InfluxDB

	//Your code goes here
}
```