package perfmonitoring

import (
	"testing"
)

func TestPerfMonGetSet(t *testing.T) {
	influxClient, err := NewInfluxdbClient()
	if err != nil {
		t.Fatal("Could not create InfluxClient", err)
	}
	perfMon := NewPerfMonitoring(influxClient, "go_commons_unit_tests", "metrics")

	perfMon.Set("testcolumn1", 5)
	if perfMon.Get("testcolumn1").(int) != 5 {
		t.Fatal("Invalid value of testcolumn1")
	}
	perfMon.Push()

	if perfMon.Get("testcolumn1") != nil {
		t.Fatal("Values not cleared after push")
	}
}
