package main

import (
	"IrisAPIs"
	"bytes"
	"fmt"
	"github.com/docker/distribution/uuid"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime"
	"strings"
	"time"
)

var gLogger *logrus.Logger

func init() {
	gLogger = logrus.New()
	gLogger.SetLevel(logrus.DebugLevel)
	gLogger.SetFormatter(&IrisAPIs.LinearLoggerFormat{})
	gLogger.Debug("Logger initialized")
}

func InjectLoggerMiddleware(logger logrus.FieldLogger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ipAddr := ctx.ClientIP()

		correlationId := ""

		correlationIdHeader := ctx.Request.Header["X-Correlation-ID"]
		requestIdHeader := ctx.Request.Header["X-Request-ID"]

		if len(correlationIdHeader) == 0 && len(requestIdHeader) == 0 {
			correlationId = uuid.Generate().String()
		} else {
			for _, v := range [][]string{correlationIdHeader, requestIdHeader} {
				if len(v) != 0 {
					correlationId = v[0]
				}
			}
		}

		apiKeyRef, exist := ctx.Get(ApiKeyRef)
		if !exist {
			apiKeyRef = -1
		}

		meta := IrisAPIs.LoggerMeta{
			CorrelationId: correlationId,
			IpAddress:     ipAddr,
			//TODO: 這個值會空，因為還沒跑到下面，想想怎麼弄
			ApiKeyRef: apiKeyRef.(int),
		}
		serviceLogger := IrisAPIs.ExistingLoggerWithMeta(logger, meta)
		ctx.Set(IrisAPIs.LoggerKey, serviceLogger)
	}
}

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

// RecoveryWithLogger This section is taken from gin source code
// Modify this to be compatible with DnaLogger
// RecoveryWithWriter returns a middleware for a given writer that recovers from any panics and writes a 500 if there was one.
func RecoveryWithLogger(logger logrus.FieldLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				if logger != nil {
					stack := stack(3)
					stackLogger := logger.WithField(IrisAPIs.ExecInfoString, string(stack))
					httpRequest, _ := httputil.DumpRequest(c.Request, false)
					headers := strings.Split(string(httpRequest), "\r\n")
					for idx, header := range headers {
						current := strings.Split(header, ":")
						if current[0] == "Authorization" {
							headers[idx] = current[0] + ": *"
						}
					}

					if brokenPipe {
						stackLogger.Errorf("%s\n%s%s", err, string(httpRequest), "\033[0m")
					} else {
						stackLogger.Errorf("[Recovery] %s panic recovered:\n%s\n%s",
							timeFormat(time.Now()), err, "\033[0m")
					}
				}

				// If the connection is dead, we can't write a status to it.
				if brokenPipe {
					_ = c.Error(err.(error)) // nolint: errcheck
					c.Abort()
				} else {
					c.AbortWithStatus(http.StatusInternalServerError)
				}
			}
		}()
		c.Next()
	}
}

// stack returns a nicely formatted stack frame, skipping skip frames.
func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		_, _ = fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		_, _ = fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}

func timeFormat(t time.Time) string {
	var timeString = t.Format("2006/01/02 - 15:04:05")
	return timeString
}
