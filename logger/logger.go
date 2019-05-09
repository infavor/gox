package logger

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

func InitLogger() {
	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(&log.JSONFormatter{})
	// log.SetFormatter(&log.TextFormatter{})
	log.SetFormatter(new(MyTextFormatter))

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

type MyTextFormatter struct {
}

func (f *MyTextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	b.WriteString(common.GetLongLongDateString(entry.Time))
	b.WriteString(" | ")
	b.WriteString(fmt.Sprintf("%-5s", strings.ToUpper(entry.Level.String())))
	b.WriteString(" | ")
	b.WriteString(entry.Message)
	b.WriteString("\n")
	return b.Bytes(), nil
}
