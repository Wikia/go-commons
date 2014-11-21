go-commons/logger
==========

Common library for logging which stores the logs in elasticsearch to be accessable through kibana. The logs are stored in JSON format with the
@timestamp field automatically set to the invocation time.

Sample usage:
```go
//Init the logger and set application name and minimum log level with which messages will be stored
err := InitLogger("My application name", LOG_LEVEL_DEBUG)
if err != nil {
    panic(err)
}

//Global function to get the logger after it has been initialized. The logger is thread safe.
logger := GetLogger()

//Simple message logging
logger.Debug("debug message") //This will create a JSON with the @message field set to the passed argument and @severity to debug
logger.Info("info message") //This will create a JSON with the @message field set to the passed argument and @severity to info

//Custom map logging - allows more control over fields being stored in the JSON object
errorMap := make(map[string]interface{})
errorMap["testField"] = "testValue"
logger.ErrorMap(errorMap) //This will create a JSON with the @testField set to "testValue"
```
