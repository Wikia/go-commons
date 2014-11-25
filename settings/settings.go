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
	IsDev    bool
}

var settings *Settings

func init() {
	env := os.Getenv(ENV_VAR)
	settings = new(Settings)
	settings.IsDev = env == "dev" || env == ""

	influxDBSettings := new(InfluxDBSettings)
	settings.InfluxDB = influxDBSettings

	if settings.IsDev {
		influxDBSettings.Host = "graph-s3"
		influxDBSettings.UdpPort = 5551
	} else {
		influxDBSettings.Host = "graph-s3"
		influxDBSettings.UdpPort = 4444
	}
}

func GetSettings() *Settings {
	return settings
}
