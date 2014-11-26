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
