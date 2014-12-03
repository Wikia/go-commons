package perfmonitoring

import (
	"testing"
)

func TestPerfMonGetSet(t *testing.T) {
	influxClient, err := NewInfluxdbClient()
	if err != nil {
		t.Fatal(err)
	}
	perfMon := NewPerfMonitoring(influxClient, "go_commons_unit_tests", "metrics")

	perfMon.Set("testcolumn1", []interface{}{5})

	if perfMon.Get("testcolumn1")[0].(int) != 5 {
		t.FailNow()
	}
	perfMon.Push()
}
