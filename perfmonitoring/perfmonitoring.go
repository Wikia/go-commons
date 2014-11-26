package perfmonitoring

import (
	"fmt"
	"strings"

	"github.com/Wikia/go-commons/settings"
	"github.com/influxdb/influxdb/client"
)

type PerfMonitoring struct {
	seriesName     string
	metrics        map[string][]interface{}
	influxdbClient *client.Client
}

func NewPerfMonitoring(appName string, seriesName string) (*PerfMonitoring, error) {
	perfMon := new(PerfMonitoring)
	perfMon.seriesName = fmt.Sprintf("%s_%s", strings.ToLower(appName), strings.ToLower(seriesName))
	perfMon.metrics = make(map[string][]interface{})
	settings := settings.GetSettings()
	influxConfig := new(client.ClientConfig)
	influxConfig.Host = fmt.Sprintf("%s:%d", settings.InfluxDB.Host, settings.InfluxDB.UdpPort)
	influxConfig.IsUDP = true
	influxClient, err := client.NewClient(influxConfig)
	if err != nil {
		return nil, err
	}

	perfMon.influxdbClient = influxClient
	return perfMon, nil
}

func (perfMon *PerfMonitoring) Set(name string, value []interface{}) {
	perfMon.metrics[name] = value
}

func (perfMon *PerfMonitoring) Inc(name string, inc int) {
	value := perfMon.metrics[name][0]
	perfMon.metrics[name][0] = value.(int) + inc
}

func (perfMon *PerfMonitoring) Get(name string) []interface{} {
	return perfMon.metrics[name]
}

func (perfMon *PerfMonitoring) Push() error {
	keys := make([]string, 0, len(perfMon.metrics))
	values := make([][]interface{}, 0, len(perfMon.metrics))
	for k, v := range perfMon.metrics {
		keys = append(keys, k)
		values = append(values, v)
	}

	series := new(client.Series)
	series.Name = perfMon.seriesName
	series.Columns = keys
	series.Points = values

	err := perfMon.influxdbClient.WriteSeriesOverUDP([]*client.Series{series})

	return err
}
