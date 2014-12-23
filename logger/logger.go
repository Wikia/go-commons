/*
Logger writing logs to kibana (elasticsearch). Uses the Syslog internally.
Formats messages as JSON.
*/

package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"log/syslog"
	"runtime"
	"strings"
	"time"
)

//Main logger structure
type Logger struct {
	logLevel     int
	appName      string
	errorLogger  *log.Logger
	warnLogger   *log.Logger
	infoLogger   *log.Logger
	debugLogger  *log.Logger
	logConsumers []LogConsumer
}

//Enables defining additional log handlers, mainly used for unit tests
type LogConsumer interface {
	Log(data string, logLevel int)
}

const (
	LogLevelError = iota
	LogLevelWarn  = iota
	LogLevelInfo  = iota
	LogLevelDebug = iota
)

var logger *Logger

//Should be called before using the logger - sets the app name and minimum log level
//of messages that should be handled
func InitLogger(appName string, logLevel int) error {
	logger = new(Logger)
	logger.logLevel = logLevel
	logger.appName = appName

	var err error
	logger.errorLogger, err = syslog.NewLogger(syslog.LOG_ERR|syslog.LOG_USER, 0)
	if err != nil {
		return err
	}

	logger.warnLogger, err = syslog.NewLogger(syslog.LOG_WARNING|syslog.LOG_USER, 0)
	if err != nil {
		return err
	}

	logger.infoLogger, err = syslog.NewLogger(syslog.LOG_INFO|syslog.LOG_USER, 0)
	if err != nil {
		return err
	}

	logger.debugLogger, err = syslog.NewLogger(syslog.LOG_DEBUG|syslog.LOG_USER, 0)
	if err != nil {
		return err
	}

	return nil
}

func GetLogger() *Logger {
	if logger == nil {
		panic("You cannot call GetLogger before initializing it by calling InitLogger")
	}
	return logger
}

func (logger *Logger) AddLogConsumer(logConsumer LogConsumer) {
	logger.logConsumers = append(logger.logConsumers, logConsumer)
}

func (logger *Logger) Error(message string) {
	logger.logMessage(message, LogLevelError, logger.errorLogger)
}

func (logger *Logger) ErrorErr(err error) {
	if err == nil {
		return
	}
	_, file, line, _ := runtime.Caller(1)
	fileSlices := strings.Split(file, "/")
	file = fileSlices[len(fileSlices)-1]

	message := fmt.Sprintf("%s:%d %s", file, line, err.Error())
	logger.logMessage(message, LogLevelError, logger.errorLogger)
}

func (logger *Logger) ErrorMap(entry map[string]interface{}) {
	logger.logMap(entry, LogLevelError, logger.errorLogger)
}

func (logger *Logger) Warn(message string) {
	logger.logMessage(message, LogLevelWarn, logger.warnLogger)
}

func (logger *Logger) WarnMap(entry map[string]interface{}) {
	logger.logMap(entry, LogLevelWarn, logger.warnLogger)
}

func (logger *Logger) Info(message string) {
	logger.logMessage(message, LogLevelInfo, logger.infoLogger)
}

func (logger *Logger) InfoMap(entry map[string]interface{}) {
	logger.logMap(entry, LogLevelInfo, logger.infoLogger)
}

func (logger *Logger) Debug(message string) {
	logger.logMessage(message, LogLevelDebug, logger.debugLogger)
}

func (logger *Logger) DebugMap(entry map[string]interface{}) {
	logger.logMap(entry, LogLevelDebug, logger.debugLogger)
}

func (logger *Logger) logMap(entry map[string]interface{}, level int, logLogger *log.Logger) {
	if logger.logLevel >= level {
		j := logger.prepareMapJson(entry, level)
		logLogger.Print(j)
		logger.notifyLogConsumers(j, level)
	}
}

func (logger *Logger) logMessage(message string, level int, logLogger *log.Logger) {
	if logger.logLevel >= level {
		j := logger.prepareMapJsonFromMessage(message, level)
		logLogger.Print(j)
		logger.notifyLogConsumers(j, level)
	}
}

func (logger *Logger) prepareMapJsonFromMessage(message string, logLevel int) string {

	entry := make(map[string]interface{})
	entry["@message"] = message

	return logger.prepareMapJson(entry, logLevel)
}

func (logger *Logger) prepareMapJson(entry map[string]interface{}, logLevel int) string {

	entry["@timestamp"] = time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	result, err := json.Marshal(entry)
	if err != nil {
		panic(err)
	}

	return string(result)
}

func (logger *Logger) notifyLogConsumers(message string, logLevel int) {
	if logger.logConsumers != nil {
		for _, consumer := range logger.logConsumers {
			consumer.Log(message, logLevel)
		}
	}
}
