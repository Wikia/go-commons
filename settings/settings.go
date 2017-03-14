package settings

import (
	"os"
)

const ENV_VAR = "WIKIA_ENVIRONMENT"

type InfluxDBSettings struct {
	Host    string
	UdpPort int
}

type Settings struct {
	InfluxDB *InfluxDBSettings
}

var settings *Settings

func init() {
	settings = new(Settings)

	influxDBSettings := new(InfluxDBSettings)
	settings.InfluxDB = influxDBSettings

	env := os.Getenv(ENV_VAR)
	if env == "prod" {
		influxDBSettings.Host = "prod.app-metrics-etl.service.sjc.consul"
		influxDBSettings.UdpPort = 4444
	} else if env == "staging" {
		influxDBSettings.Host = "staging.app-metrics-etl.service.sjc.consul"
		influxDBSettings.UdpPort = 5552
	} else {
		influxDBSettings.Host = "prod.app-metrics-etl.service.sjc.consul"
		influxDBSettings.UdpPort = 5551
	}
}

func GetSettings() *Settings {
	return settings
}
