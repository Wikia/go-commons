package perfmonitoring

import (
	"testing"
)

func TestPerfMonGetSet(t *testing.T) {
	perfMon, err := NewPerfMonitoring("go_commons_unit_tests", "metrics")
	if err != nil {
		t.Fatal(err)
	}

	perfMon.Set("testcolumn1", []interface{}{5})

	if perfMon.Get("testcolumn1")[0].(int) != 5 {
		t.FailNow()
	}
	perfMon.Push()
}
