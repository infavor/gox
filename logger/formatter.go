package logger

import (
	"bytes"
	"fmt"
	"github.com/hetianyi/gox"
	. "github.com/logrusorgru/aurora"
	"github.com/sirupsen/logrus"
	"strings"
)

// SimpleTextFormatter is default text formatter.
// [I] 2019-12-12 12:12:12,221 [xxx.go] xxx
type DefaultTextFormatter struct {
}

// SimpleTextFormatter is the simple version of log format.
// [I] 12:12:12 xxx
type SimpleTextFormatter struct {
}

// Format formats logs.
func (f *DefaultTextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	b.WriteString(fmt.Sprintf("%-1s", changeLevelColor(strings.ToUpper(entry.Level.String())[0])))
	b.WriteString(BrightCyan(gox.GetLongLongDateString(entry.Time)).String())
	b.WriteString(getCaller())
	b.WriteString(White(entry.Message).String())
	b.WriteString("\n")
	return b.Bytes(), nil
}

// Format formats logs.
func (f *SimpleTextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	b.WriteString(fmt.Sprintf("%-1s", changeLevelColor(strings.ToUpper(entry.Level.String())[0])))
	b.WriteString(BrightCyan(gox.GetShortDateString(entry.Time)).String())
	b.WriteString(" ")
	b.WriteString(White(entry.Message).String())
	b.WriteString("\n")
	return b.Bytes(), nil
}
