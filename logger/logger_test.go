package logger

import (
	"encoding/json"
	"testing"
)

type TestLogConsumer struct {
	InvocationCount int
	Testing         *testing.T
}

func (logConsumer *TestLogConsumer) Log(data string, logLevel int) {
	logMap := make(map[string]string)
	err := json.Unmarshal([]byte(data), &logMap)
	if err != nil {
		logConsumer.Testing.Fatal(err)
	}

	if logConsumer.InvocationCount == 0 {
		if logMap["@message"] != "debug message" || logLevel != LOG_LEVEL_DEBUG {
			logConsumer.Testing.Fatal("incorrect log message or log level: ", logMap["@message"])
		}
	} else if logConsumer.InvocationCount == 1 {
		if logMap["@message"] != "info message" || logLevel != LOG_LEVEL_INFO {
			logConsumer.Testing.Fatal("incorrect log message or log level: ", logMap["@message"])
		}
	} else if logConsumer.InvocationCount == 2 {
		if logMap["testField"] != "testValue" || logLevel != LOG_LEVEL_ERROR {
			logConsumer.Testing.Fatal("incorrect log message or log level: ", logMap["testField"])
		}
	} else {
		logConsumer.Testing.Fatal("too many log messages logged")
	}

	logConsumer.InvocationCount = logConsumer.InvocationCount + 1
}

func TestLogger(t *testing.T) {
	err := InitLogger("UnitTests", LOG_LEVEL_DEBUG)
	if err != nil {
		t.Fatal(err)
	}

	logger := GetLogger()
	if logger == nil {
		t.FailNow()
	}
	logConsumer := new(TestLogConsumer)
	logConsumer.Testing = t
	logger.AddLogConsumer(logConsumer)

	logger.Debug("debug message")
	logger.Info("info message")
	errorMap := make(map[string]interface{})
	errorMap["testField"] = "testValue"
	logger.ErrorMap(errorMap)
	if logConsumer.InvocationCount != 3 {
		t.Fatal("incorrent number of log invocations: ", logConsumer.InvocationCount)
	}
}
