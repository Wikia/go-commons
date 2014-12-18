package perfmonitoring

import (
	"fmt"
	"strings"

	"github.com/Wikia/go-commons/settings"
	"github.com/influxdb/influxdb/client"
)

/*Note that this object is designed to be used in a single thread
If you have a multithreaded application, then it is recommended to create a new PerfMonitoring
instance for each thread/request. The InfluxDB client can be reused in all threads.
*/
type PerfMonitoring struct {
	seriesName     string
	metrics        map[string]interface{}
	influxdbClient *client.Client
}

func NewInfluxdbClient() (*client.Client, error) {
	settings := settings.GetSettings()
	influxConfig := new(client.ClientConfig)
	influxConfig.Host = fmt.Sprintf("%s:%d", settings.InfluxDB.Host, settings.InfluxDB.UdpPort)
	influxConfig.IsUDP = true
	influxClient, err := client.NewClient(influxConfig)
	if err != nil {
		return nil, err
	}

	return influxClient, nil
}

func NewPerfMonitoring(influxClient *client.Client, appName string, seriesName string) *PerfMonitoring {
	perfMon := new(PerfMonitoring)
	perfMon.seriesName = fmt.Sprintf("%s_%s", strings.ToLower(appName), strings.ToLower(seriesName))
	perfMon.metrics = make(map[string]interface{})
	perfMon.influxdbClient = influxClient
	return perfMon
}

func (perfMon *PerfMonitoring) Set(name string, value interface{}) {
	perfMon.metrics[name] = value
}

func (perfMon *PerfMonitoring) Inc(name string, inc int) {
	value := perfMon.metrics[name]
	perfMon.metrics[name] = value.(int) + inc
}

func (perfMon *PerfMonitoring) Get(name string) interface{} {
	return perfMon.metrics[name]
}

func (perfMon *PerfMonitoring) Push() error {
	keys := make([]string, 0, len(perfMon.metrics))
	values := [][]interface{}{make([]interface{}, 0, len(perfMon.metrics))}
	for k, v := range perfMon.metrics {
		keys = append(keys, k)
		values[0] = append(values[0], v)
	}

	series := new(client.Series)
	series.Name = perfMon.seriesName
	series.Columns = keys
	series.Points = values

	err := perfMon.influxdbClient.WriteSeriesOverUDP([]*client.Series{series})
	perfMon.metrics = make(map[string]interface{}) //Reset metrics after push

	return err
}
