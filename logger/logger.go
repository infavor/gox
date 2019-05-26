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
	"time"
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
	b.WriteString(GetLongLongDateString(entry.Time))
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

// GetLongLongDateString gets short date format like '2018-11-11 12:12:12,233'.
func GetLongLongDateString(t time.Time) string {
	var buff bytes.Buffer
	buff.WriteString(strconv.Itoa(t.Year()))
	buff.WriteString("-")
	buff.WriteString(format2(int(t.Month())))
	buff.WriteString("-")
	buff.WriteString(format2(t.Day()))
	buff.WriteString(" ")
	buff.WriteString(format2(t.Hour()))
	buff.WriteString(":")
	buff.WriteString(format2(t.Minute()))
	buff.WriteString(":")
	buff.WriteString(format2(t.Second()))
	buff.WriteString(",")
	buff.WriteString(format3(t.Nanosecond() / 1e6))
	return buff.String()
}
