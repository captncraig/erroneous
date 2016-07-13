package erroneous

import (
	"github.com/pkg/errors"
	"github.com/twinj/uuid"

	"bytes"
	"fmt"
	"os"
	"strings"
	"time"
)

var machineName string

func init() {
	if name, err := os.Hostname(); err == nil {
		machineName = name
	}
}

// ErrorLogger is the core interface through wich loggers can be stored.
type ErrorLogger struct {
	Store
	ApplicationName string
	MachineName     string
}

// Default is the Logger that all package-level functions will use
var Default = &ErrorLogger{Store: &MemoryStore{}}

// Use sets the default logger that package-level functions will use to record errors
func Use(logger *ErrorLogger) {
	Default = logger
}

func (l *ErrorLogger) LogError(e error) {
	err := &Error{
		GUID:            uuid.NewV4(),
		ApplicationName: l.ApplicationName,
		Count:           1,
		CreationDate:    time.Now(),
		Message:         e.Error(),
	}
	if l.MachineName != "" {
		err.MachineName = l.MachineName
	} else {
		err.MachineName = machineName
	}
	err.Detail = GetDetail(e)
	l.Store.LogError(err) //TODO: catch failure
}

func LogError(e error) {
	Default.LogError(e)
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func GetDetail(e error) string {
	buf := &bytes.Buffer{}
	buf.WriteString(e.Error())
	buf.WriteRune('\n')
	err, ok := e.(stackTracer)
	if ok {
		fmt.Fprintf(buf, "%+v", err.StackTrace())
	} else {
		//just format stack trace ourselves
		st := callers().stackTrace()
		for len(st) > 0 && strings.HasSuffix(st[0].file(), "erroneous/logger.go") {
			st = st[1:]
		}
		fmt.Fprintf(buf, "%+v", st)
	}
	return buf.String()
}
