package perfmonitoring

import (
    "testing"
)

func TestPerfMonGetSet(t *testing.T) {
    perfMon, err := NewPerfMonitoring("testapp", "testseries")
    if err != nil {
        t.Fatal(err)
    }

    perfMon.Set("testname", 5)

    if perfMon.Get("testname").(int) != 5 {
        t.FailNow()
    }
}
