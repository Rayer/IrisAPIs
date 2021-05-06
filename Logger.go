package IrisAPIs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
	"time"
)

var defaultLogger *defLogger

const (
	LoggerMetaString = "LoggerMetaString"
	ExecInfoString   = "ExecInfoString"
	LoggerKey        = "LoggerKey"
)

func init() {
	defaultLogger = ExistingLoggerWithMeta(logrus.New(), LoggerMeta{})
}

type defLogger struct {
	logrus.FieldLogger
	meta LoggerMeta
}

func ExistingLoggerWithMeta(logger logrus.FieldLogger, meta LoggerMeta) *defLogger {
	return &defLogger{
		FieldLogger: logger,
		meta:        meta,
	}
}

func (d *defLogger) Log() logrus.FieldLogger {
	return d.WithField(LoggerMetaString, d.meta)
}

type LoggerFormat struct {
	EventTime string `json:"event_time"`
	Level     string `json:"level"`
	Filename  string `json:"filename"`
	Func      string `json:"func"`
	Message   string `json:"message"`
	ApiKeyRef int    `json:"apiKeyRef"`
	RequestId string `json:"request_id"`
	RemoteIp  string `json:"remote_ip"`
	ExecInfo  string `json:"exec_info"`
}

func (l *LoggerFormat) Format(entry *logrus.Entry) ([]byte, error) {
	var contextMeta LoggerMeta
	meta, convertible := entry.Data[LoggerMetaString].(LoggerMeta)
	if convertible {
		contextMeta = meta
	}
	execInfo := ""
	contextExecInfo, convertible := entry.Data[ExecInfoString].(string)
	if convertible {
		execInfo = contextExecInfo
	}

	loggerFormat := LoggerFormat{
		EventTime: time.Now().Format("2006-01-02T15:04:05.999999"),
		Level:     entry.Level.String(),
		Filename: func(caller *runtime.Frame) string {
			if caller == nil {
				return ""
			}
			return fmt.Sprintf("%s:%d", caller.File, caller.Line)
		}(entry.Caller),
		Func: func(caller *runtime.Frame) string {
			if caller == nil {
				return ""
			}
			return caller.Function
		}(entry.Caller),
		Message:   entry.Message,
		ApiKeyRef: meta.ApiKeyRef,
		RequestId: contextMeta.CorrelationId,
		RemoteIp:  contextMeta.IpAddress,
		ExecInfo:  execInfo,
	}
	bytes, err := json.Marshal(loggerFormat)
	return append(bytes, '\n'), err
}

type LoggerMeta struct {
	CorrelationId string
	IpAddress     string
	ApiKeyRef     int
}

func GetLogger(ctx context.Context) logrus.FieldLogger {
	if ctx != nil && ctx.Value(LoggerKey) != nil {
		return ctx.Value(LoggerKey).(*defLogger).Log()
	}
	defaultLogger.Error("No logger found in context, use default.")
	return defaultLogger
}

func GetMeta(ctx context.Context) LoggerMeta {
	if ctx != nil && ctx.Value(LoggerKey) != nil {
		return ctx.Value(LoggerKey).(*defLogger).meta
	}
	return LoggerMeta{}
}
