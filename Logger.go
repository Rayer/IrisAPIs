package ai_rec_dna

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
	"time"
)

var defaultLogger *defLogger

func init() {
	defaultLogger = ExistingLoggerWithMeta(logrus.New(), DnaContextMeta{})
}

type defLogger struct {
	logrus.FieldLogger
	meta DnaContextMeta
}

func ExistingLoggerWithMeta(logger logrus.FieldLogger, meta DnaContextMeta) *defLogger {
	return &defLogger{
		FieldLogger: logger,
		meta:        meta,
	}
}

func (d *defLogger) Log() logrus.FieldLogger {
	return d.WithField(DnaMeta, d.meta)
}

type LoggerFormat struct {
	EventTime string `json:"event_time"`
	Level     string `json:"level"`
	Filename  string `json:"filename"`
	Func      string `json:"func"`
	Message   string `json:"message"`
	RequestId string `json:"request_id"`
	RemoteIp  string `json:"remote_ip"`
	ExecInfo  string `json:"exec_info"`
}

func (l *LoggerFormat) Format(entry *logrus.Entry) ([]byte, error) {
	var contextMeta DnaContextMeta
	meta, convertible := entry.Data[DnaMeta].(DnaContextMeta)
	if convertible {
		contextMeta = meta
	}
	execInfo := ""
	contextExecInfo, convertible := entry.Data[ExecInfo].(string)
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
		RequestId: contextMeta.CorrelationId,
		RemoteIp:  contextMeta.IpAddress,
		ExecInfo:  execInfo,
	}
	bytes, err := json.Marshal(loggerFormat)
	return append(bytes, '\n'), err
}

type DnaContextMeta struct {
	CorrelationId string
	IpAddress     string
	SiteId        string
	BidObjId      string
}

func GetLogger(ctx context.Context) logrus.FieldLogger {
	ret := ctx.Value(DnaMetaLoggerKey)
	if ret == nil {
		//This behavior should be change?
		defaultLogger.Error("No logger found in context, use default.")
		return defaultLogger
	}
	return ret.(*defLogger).Log()
}

func GetMeta(ctx context.Context) DnaContextMeta {
	ret := ctx.Value(DnaMetaLoggerKey)
	if ret != nil {
		return ret.(*defLogger).meta
	}
	return DnaContextMeta{}
}