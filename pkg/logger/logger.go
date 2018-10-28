package logger

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// Data is a type alias so that it is easy to add additional data to log lines.
type Data map[string]interface{}

// Logger is a logger instance that contains necessary info needed when logging.
type Logger struct {
	zl   zerolog.Logger
	id   string
	data []Data
	err  error
	root []Data
}

const stackSize = 4 << 10 // 4KB

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func init() {
	zerolog.TimestampFieldName = "timestamp"
}

// New returns a new configured Logger instance.
func New() Logger {
	host, _ := os.Hostname() // nolint: gosec

	zl := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Str("hostname", host)

	return Logger{
		zl:   zl.Logger(),
		data: []Data{},
		root: []Data{},
	}
}

// ID adds an identifier to every subsequent log line.
func (log Logger) ID(id string) Logger {
	log.id = id
	return log
}

// Data adds additional data to every subsequent log line. This data is nested
// within the "data" field.
func (log Logger) Data(data Data) Logger {
	log.data = append(log.data, data)
	return log
}

// Err adds additional an error object to every subsequent log line. This is
// meant to be immediately chained to emit the log line.
func (log Logger) Err(err error) Logger {
	log.err = err
	return log
}

// Root adds additional data to every subsequent log line. This data is added to
// the root of the JSON line.
func (log Logger) Root(root Data) Logger {
	log.root = append(log.root, root)
	return log
}

// Info emits a log line with an "info" log level.
func (log Logger) Info(message string, fields ...Data) {
	log.log(log.zl.Info(), message, fields...)
}

// Error emits a log line with an "error" log level.
func (log Logger) Error(message string, fields ...Data) {
	log.log(log.zl.Error(), message, fields...)
}

// Warn emits a log line with an "warn" log level.
func (log Logger) Warn(message string, fields ...Data) {
	log.log(log.zl.Warn(), message, fields...)
}

// Debug emits a log line with an "debug" log level.
func (log Logger) Debug(message string, fields ...Data) {
	log.log(log.zl.Debug(), message, fields...)
}

// Fatal emits a log line with an "fatal" log level. It also calls os.Exit(1)
// afterwards.
func (log Logger) Fatal(message string, fields ...Data) {
	log.log(log.zl.Fatal(), message, fields...)
}

func (log Logger) log(evt *zerolog.Event, message string, fields ...Data) {
	// Merge data fields
	hasData := false
	data := zerolog.Dict()
	for _, field := range append(log.data, fields...) {
		if len(field) != 0 {
			hasData = true
			data = data.Fields(field)
		}
	}

	// Add root fields
	for _, field := range log.root {
		if len(field) != 0 {
			evt = evt.Fields(field)
		}
	}
	// Add id field
	if log.id != "" {
		evt = evt.Str("id", log.id)
	}
	// Add data field
	if hasData {
		evt = evt.Dict("data", data)
	}

	if log.err != nil {
		var stack []byte
		// support pkg/errors stackTracer interface
		if err, ok := log.err.(stackTracer); ok {
			st := err.StackTrace()
			stack = []byte(fmt.Sprintf("%+v", st))
		} else {
			stack = make([]byte, stackSize)
			n := runtime.Stack(stack, true)
			stack = stack[:n]
		}
		f := Data{"message": log.err, "stack": stack}
		evt = evt.Dict("error", zerolog.Dict().Fields(f))
	}

	evt.Int64("nanoseconds", time.Now().UnixNano()).Msg(message)
}
