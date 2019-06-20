// Copyright (C) 2019 tisnyo <tisnyo@gmail.com>.
//
// package logger is used for initializing logrus.
package logger

import (
	"bytes"
	"fmt"
	"github.com/hetianyi/gox"
	"github.com/hetianyi/gox/convert"
	"github.com/hetianyi/gox/file"
	. "github.com/logrusorgru/aurora"
	"github.com/mattn/go-colorable"
	"github.com/mholt/archiver"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Level type
type Level uint32

// These are the different logging levels. You can set the logging level to log
// on your instance of logger, obtained with `logrus.New()`.
const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel Level = iota
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel
)

const (
	HOUR int = iota
	DAY
	MONTH
	YEAR
	MB64
	MB128
	MB256
	MB512
	MB1024

	colorFlag  = "\033\\[([0-9]+;)?[0-9]+m"
	archiveExt = ".tar.gz"
)

var (
	initialized        = false
	write2File         bool
	alwaysWriteConsole bool
	rollingPolicy      []int
	curWriteLen        int64
	lastWriteTime      time.Time
	logDirectory       string
	LogFileName        string
	curOut             *os.File
	lock               *sync.Mutex

	timePolicy = HOUR
	sizePolicy = 0

	colorPattern = regexp.MustCompile(colorFlag)
)

func init() {
	lock = new(sync.Mutex)
}

type Config struct {
	Formatter          logrus.Formatter
	Level              Level
	Write2File         bool
	AlwaysWriteConsole bool // 是否总是将日志写入控制台
	RollingFileDir     string
	RollingFileName    string
	RollingPolicy      []int
}

type logWriter struct {
	colorableStdout io.Writer
}

func (w *logWriter) Write(p []byte) (int, error) {
	defer func() {
		lastWriteTime = time.Now()

	}()
	writeP := colorPattern.ReplaceAll(p, []byte(""))
	defer func() {
		curWriteLen += int64(len(writeP))
	}()
	if !write2File {
		return w.colorableStdout.Write(gox.TValue(runtime.GOOS == "linux", p, p).([]byte))
	}
	now := time.Now()
	triggerExchange(now)
	if curOut != nil {
		if alwaysWriteConsole {
			return w.colorableStdout.Write(gox.TValue(runtime.GOOS == "linux", p, p).([]byte))
		}
		return curOut.Write(writeP)
	}
	return w.colorableStdout.Write(gox.TValue(runtime.GOOS == "linux", p, p).([]byte))
}

// Init initialize logrus logger.
func Init(config *Config) {
	lastWriteTime = time.Now()
	if IsInitialized() {
		fmt.Println("logger has already initialized")
		return
	}
	if config == nil {
		config = &Config{
			Formatter:          new(DefaultTextFormatter),
			Level:              InfoLevel,
			Write2File:         false,
			RollingFileDir:     "./",
			RollingFileName:    "app",
			RollingPolicy:      []int{YEAR, MB1024},
			AlwaysWriteConsole: true,
		}
	}
	write2File = config.Write2File
	if write2File && (config.RollingFileName == "") {
		config.RollingFileName = "app"
	}
	logDirectory, _ = file.AbsPath(config.RollingFileDir)
	logDirectory = file.FixPath(logDirectory)
	LogFileName = config.RollingFileName
	LogFileName = file.FixPath(LogFileName)
	rollingPolicy = config.RollingPolicy
	alwaysWriteConsole = config.AlwaysWriteConsole
	parsePolicy()
	initialized = true
	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(&log.JSONFormatter{})
	// log.SetFormatter(&log.TextFormatter{})
	logrus.SetFormatter(gox.TValue(config.Formatter == nil, new(DefaultTextFormatter), config.Formatter).(logrus.Formatter))

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(&logWriter{colorable.NewColorableStdout()})

	// Only log the warning severity or above.
	var l = uint32(config.Level)
	logrus.SetLevel(logrus.Level(l))
}

func IsInitialized() bool {
	return initialized
}

func getCaller() string {
	_, file, line, success := runtime.Caller(9)
	if success {
		return Yellow(strings.Join([]string{" [", file[strings.LastIndex(file, "/")+1:], ":", strconv.Itoa(line), "] "}, "")).String()
	}
	return " [unknown] "
}

func changeLevelColor(l uint8) string {
	if l == 'T' {
		return strings.Join([]string{"[", string(l), "] "}, "")
	}
	if l == 'D' {
		return BrightBlack(strings.Join([]string{"[", string(l), "] "}, "")).String()
	}
	if l == 'I' {
		return BrightGreen(strings.Join([]string{"[", string(l), "] "}, "")).String()
	}
	if l == 'W' {
		return BrightYellow(strings.Join([]string{"[", string(l), "] "}, "")).String()
	}
	if l == 'E' {
		return Red(strings.Join([]string{"[", string(l), "] "}, "")).String()
	}
	if l == 'F' || l == 'P' {
		return SlowBlink(BgRed(strings.Join([]string{"[", string(l), "] "}, ""))).String()
	}
	return strings.Join([]string{"[", string(l), "] "}, "")
}

func Trace(args ...interface{}) {
	logrus.Trace(args...)
}

func Debug(args ...interface{}) {
	logrus.Debug(args...)
}

func Info(args ...interface{}) {
	logrus.Info(args...)
}

func Warn(args ...interface{}) {
	logrus.Warn(args...)
}

func Error(args ...interface{}) {
	logrus.Error(args...)
}

func Fatal(args ...interface{}) {
	logrus.Fatal(args...)
}

func Panic(args ...interface{}) {
	logrus.Panic(args...)
}

func triggerExchange(t time.Time) {

	if curOut != nil {
		if !isChanged(t) {
			return
		}
		lock.Lock()
		defer lock.Unlock()
		if !isChanged(t) {
			return
		}
	}

	var buffer bytes.Buffer
	buffer.WriteString(logDirectory)
	buffer.WriteString(string(os.PathSeparator))
	buffer.WriteString(LogFileName)
	buffer.WriteString("-")
	buffer.WriteString(strings.ReplaceAll(gox.GetDateString(t), "-", ""))
	if timePolicy == HOUR {
		buffer.WriteString(gox.TValue(t.Hour() < 10, "0", "").(string))
		buffer.WriteString(convert.IntToStr(t.Hour()))
	}
	if sizePolicy != 0 {
		buffer.WriteString("-part1")
	}
	buffer.WriteString(".log")
	newfile := buffer.String()
	// 不限制文件大小，则一定是日期变化
	if sizePolicy == 0 {
		newOut, err := file.AppendFile(newfile)
		if err != nil {
			return
		}
		if curOut != nil {
			curOut.Close()
			// 压缩历史日志
			go compressOldFile(curOut.Name())
		}
		fmt.Println("create new log file:", newfile)
		curWriteLen = 0
		curOut = newOut
		return
	}
	// 限制文件大小
	index := 1
	for file.Exists(newfile) || file.Exists(newfile+archiveExt) {
		buffer.Reset()
		buffer.WriteString(newfile[0:strings.LastIndex(newfile, "-")])
		buffer.WriteString("-part")
		buffer.WriteString(convert.IntToStr(index))
		buffer.WriteString(".log")
		newfile = buffer.String()
		index++
	}
	// fmt.Println("create new log file:", newfile)
	newOut, err := file.AppendFile(newfile)
	if err != nil {
		return
	}
	if curOut != nil {
		curOut.Close()
		// 压缩历史日志
		go compressOldFile(curOut.Name())
	}
	curWriteLen = 0
	curOut = newOut
}

func parsePolicy() {
	if rollingPolicy == nil || len(rollingPolicy) == 0 {
		return
	}
	for _, p := range rollingPolicy {
		if p >= HOUR && p <= YEAR {
			if timePolicy < p {
				timePolicy = p
			}
		}
		if p >= MB64 && p <= MB1024 {
			if sizePolicy < p {
				sizePolicy = p
			}
		}
	}
}

func isChanged(t time.Time) bool {
	changed := false
	// hour changed
	if (timePolicy == HOUR && (gox.GetYear(lastWriteTime) != gox.GetYear(t) ||
		gox.GetMonth(lastWriteTime) != gox.GetMonth(t) ||
		gox.GetDay(lastWriteTime) != gox.GetDay(t) ||
		gox.GetHour(lastWriteTime) != gox.GetHour(t))) ||
		(timePolicy == DAY && (gox.GetYear(lastWriteTime) != gox.GetYear(t) ||
			gox.GetMonth(lastWriteTime) != gox.GetMonth(t) ||
			gox.GetDay(lastWriteTime) != gox.GetDay(t))) ||
		(timePolicy == MONTH && (gox.GetYear(lastWriteTime) != gox.GetYear(t) ||
			gox.GetMonth(lastWriteTime) != gox.GetMonth(t))) ||
		(timePolicy == YEAR && (gox.GetYear(lastWriteTime) != gox.GetYear(t))) {
		changed = true
	}
	if !changed && sizePolicy != 0 {
		if (sizePolicy == MB64 && curWriteLen >= 2<<25) ||
			(sizePolicy == MB128 && curWriteLen >= 2<<26) ||
			(sizePolicy == MB256 && curWriteLen >= 2<<27) ||
			(sizePolicy == MB512 && curWriteLen >= 2<<28) ||
			(sizePolicy == MB1024 && curWriteLen >= 2<<29) {
			changed = true
		}
	}
	return changed
}

func compressOldFile(path string) {
	// fmt.Println("压缩日志：", path)
	fileName := filepath.Base(path) + archiveExt
	dir := filepath.Dir(path)
	gox.Try(func() {
		err := archiver.Archive([]string{path}, dir+string(os.PathSeparator)+fileName)
		if err != nil {
			fmt.Println("err while compressing log file:", err)
		} else {
			file.Delete(path)
		}
	}, func(e interface{}) {
		fmt.Println("err while compressing log file:", e)
	})
}

func FakeWriteLen(len1 int64) {
	curWriteLen = len1
}
