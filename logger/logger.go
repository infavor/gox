// Copyright (C) 2019 tisnyo <tisnyo@gmail.com>.
//
// package logger is used for initializing logrus.
package logger

import (
	"bytes"
	"fmt"
	"github.com/hetianyi/gox"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
)

// Init initialize logrus logger.
func Init(formatter *MyTextFormatter) {
	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(&log.JSONFormatter{})
	// log.SetFormatter(&log.TextFormatter{})
	log.SetFormatter(gox.TValue(formatter == nil, new(MyTextFormatter), formatter).(*MyTextFormatter))

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

// default text formatter.
type MyTextFormatter struct {
}

// Format formats logs.
func (f *MyTextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	b.WriteString(gox.GetLongLongDateString(entry.Time))
	b.WriteString(" | ")
	b.WriteString(fmt.Sprintf("%-5s", strings.ToUpper(entry.Level.String())))
	b.WriteString(" | ")
	b.WriteString(entry.Message)
	b.WriteString("\n")
	return b.Bytes(), nil
}

func format2(input int) string {
	if input < 10 {
		return "0" + strconv.Itoa(input)
	}
	return strconv.Itoa(input)
}

func format3(input int) string {
	if input < 10 {
		return "00" + strconv.Itoa(input)
	}
	if input < 100 {
		return "0" + strconv.Itoa(input)
	}
	return strconv.Itoa(input)
}
