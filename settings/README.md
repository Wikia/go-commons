go-commons/settings
==========

Common settings for modules in go-commons. Detects whether the lib is operating in development
or production environment by checking the WIKIA_ENVIRONMENT environment variable.

Sample usage:
```go
s := settings.GetSettings()
host := s.InfluxDB.Host
```
