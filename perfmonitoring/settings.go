package perfmonitoring

import (
    "os"
)

const ENV_VAR = "WIKIA_ENVIRONMENT"

type Settings struct {
    Host string
    UdpPort int
}

var isDev bool

func init() {
    env:= os.Getenv(ENV_VAR)
    isDev = env == "dev" || env == ""
}

func getSettings() *Settings {
    settings := new(Settings)

    if isDev {
        settings.Host = "graph-s3"
        settings.UdpPort = 5551
    } else {
        settings.Host = "graph-s3"
        settings.UdpPort = 4444
    }

    return settings
}
