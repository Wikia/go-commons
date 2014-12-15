package perfmonitoring

import (
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	influxClient, err := NewInfluxdbClient()
	if err != nil {
		t.Fatal("Could not create InfluxClient", err)
	}
	perfMon := NewPerfMonitoring(influxClient, "go_commons_unit_tests", "metrics")

	timer := NewTimer(perfMon, "test_time")
	timer.AddValue("additional_value", "Value")

	time.Sleep(time.Duration(10) * time.Millisecond)

	err = timer.Close()
	if err != nil {
		t.Fatal("Could not close timer", err)
	}

	if perfMon.Get("test_time").(int64) < 10 {
		t.Fatal("Invalid value of column test_time")
	}

	if perfMon.Get("additional_value") != "Value" {
		t.Fatal("Invalid value of column additional_value")
	}
}
