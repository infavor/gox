// Copyright (C) 2019 tisnyo <tisnyo@gmail.com>.
//
// package logger is used for initializing logrus.
package logger

import (
	"bytes"
	"fmt"
	"github.com/hetianyi/gox"
	. "github.com/logrusorgru/aurora"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
)

var initialized = false

type Config struct {
	Formatter logrus.Formatter
	Level     log.Level
	Out       io.Writer
}

// Init initialize logrus logger.
func Init(config *Config) {
	if IsInitialized() {
		log.Warn("logger has already initialized")
	}
	if config == nil {
		config = &Config{
			Formatter: new(DefaultTextFormatter),
			Level:     logrus.DebugLevel,
			Out:       os.Stdout,
		}
	}
	initialized = true
	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(&log.JSONFormatter{})
	// log.SetFormatter(&log.TextFormatter{})
	log.SetFormatter(gox.TValue(config.Formatter == nil, new(DefaultTextFormatter), config.Formatter).(log.Formatter))

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(gox.TValue(config.Out == nil, os.Stdout, config.Out).(io.Writer))

	// Only log the warning severity or above.
	log.SetLevel(config.Level)
}

func IsInitialized() bool {
	return initialized
}

// default text formatter.
type DefaultTextFormatter struct {
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
	b.WriteString(Cyan(gox.GetLongLongDateString(entry.Time)).String())
	b.WriteString(getCaller())
	b.WriteString(BrightBlue(entry.Message).String())
	b.WriteString("\n")
	return b.Bytes(), nil
}

func getCaller() string {
	_, file, line, success := runtime.Caller(8)
	if success {
		return Magenta(strings.Join([]string{" [", file[strings.LastIndex(file, "/")+1:], ":", strconv.Itoa(line), "] "}, "")).String()
	}
	return " [known] "
}

func changeLevelColor(l uint8) string {
	if l == 'T' {
		return strings.Join([]string{"[", string(l), "] "}, "")
	}
	if l == 'D' {
		return BrightBlack(strings.Join([]string{"[", string(l), "] "}, "")).String()
	}
	if l == 'I' {
		return Blue(strings.Join([]string{"[", string(l), "] "}, "")).String()
	}
	if l == 'W' {
		return Yellow(strings.Join([]string{"[", string(l), "] "}, "")).String()
	}
	if l == 'E' {
		return Red(strings.Join([]string{"[", string(l), "] "}, "")).String()
	}
	if l == 'F' || l == 'P' {
		return SlowBlink(BgRed(strings.Join([]string{"[", string(l), "] "}, ""))).String()
	}
	return strings.Join([]string{"[", string(l), "] "}, "")
}
