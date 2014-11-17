package perfmonitoring

import (
    "fmt"
    "strings"
    "github.com/influxdb/influxdb/client"
)

type PerfMonitoring struct {
    seriesName string
    metrics map[string]interface{}
    influxdbClient *client.Client
}

func NewPerfMonitoring(appName string, seriesName string) (*PerfMonitoring, error) {
    perfMon := new(PerfMonitoring)
    perfMon.seriesName = fmt.Sprintf("%s_%s", strings.ToLower(appName), strings.ToLower(seriesName))
    perfMon.metrics = make(map[string]interface{})
    settings := getSettings()
    influxConfig := new(client.ClientConfig)
    influxConfig.Host = fmt.Sprintf("%s:%d", settings.Host, settings.UdpPort)
    influxConfig.IsUDP = true
    influxClient, err := client.NewClient(influxConfig)
    if err != nil {
        return nil, err
    } else {
        perfMon.influxdbClient = influxClient
        return perfMon, nil
    }
}

func (perfMon *PerfMonitoring) Set(name string, value interface{}) {
	perfMon.metrics[name] = value
}

func (perfMon *PerfMonitoring) Inc(name string, inc int) {
	value := perfMon.metrics[name]
    perfMon.metrics[name] = value.(int) + inc;
}

func (perfMon *PerfMonitoring) Get(name string) interface{} {
    return perfMon.metrics[name]
}
