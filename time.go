// Copyright (C) 2019 tisnyo <tisnyo@gmail.com>.
//
// This file provides some functions about time operation.
package gox

import (
	"bytes"
	"strconv"
	"sync"
	"time"
)

var (
	lock      = *new(sync.Mutex)
	increment = 0
)

// GetDateString gets short date format like '2018-11-11'.
func GetDateString(t time.Time) string {
	var buff bytes.Buffer
	buff.WriteString(strconv.Itoa(GetYear(t)))
	buff.WriteString("-")
	buff.WriteString(format2(GetMonth(t)))
	buff.WriteString("-")
	buff.WriteString(format2(GetDay(t)))
	return buff.String()
}

// GetLongDateString gets date format like '2018-11-11 12:12:12'.
func GetLongDateString(t time.Time) string {
	var buff bytes.Buffer
	buff.WriteString(strconv.Itoa(GetYear(t)))
	buff.WriteString("-")
	buff.WriteString(format2(GetMonth(t)))
	buff.WriteString("-")
	buff.WriteString(format2(GetDay(t)))
	buff.WriteString(" ")
	buff.WriteString(format2(GetHour(t)))
	buff.WriteString(":")
	buff.WriteString(format2(GetMinute(t)))
	buff.WriteString(":")
	buff.WriteString(format2(GetSecond(t)))
	return buff.String()
}

// GetShortDateString gets time format like '12:12:12'.
func GetShortDateString(t time.Time) string {
	var buff bytes.Buffer
	buff.WriteString(format2(GetHour(t)))
	buff.WriteString(":")
	buff.WriteString(format2(GetMinute(t)))
	buff.WriteString(":")
	buff.WriteString(format2(GetSecond(t)))
	return buff.String()
}

// GetLongLongDateString gets short date format like '2018-11-11 12:12:12,233'.
func GetLongLongDateString(t time.Time) string {
	var buff bytes.Buffer
	buff.WriteString(strconv.Itoa(GetYear(t)))
	buff.WriteString("-")
	buff.WriteString(format2(GetMonth(t)))
	buff.WriteString("-")
	buff.WriteString(format2(GetDay(t)))
	buff.WriteString(" ")
	buff.WriteString(format2(GetHour(t)))
	buff.WriteString(":")
	buff.WriteString(format2(GetMinute(t)))
	buff.WriteString(":")
	buff.WriteString(format2(GetSecond(t)))
	buff.WriteString(",")
	buff.WriteString(format3(GetMillionSecond(t)))
	return buff.String()
}

// GetTimestamp gets current timestamp in milliseconds.
func GetTimestamp(t time.Time) int64 {
	return t.UnixNano() / 1e6
}

// CreateTime returns the local Time corresponding to the given Unix time,
// sec seconds and nsec nanoseconds since January 1, 1970 UTC. It is valid to pass nsec outside the range [0, 999999999].
// Not all sec values have a corresponding time value. One such value is 1<<63-1 (the largest int64 value).
func CreateTime(millis int64) time.Time {
	return time.Unix(millis, 0)
}

// GetNanosecond gets current timestamp in Nanosecond.
func GetNanosecond(t time.Time) int64 {
	return t.UnixNano()
}

// GetYear gets year number.
func GetYear(t time.Time) int {
	return t.Year()
}

// GetMonth gets month number.
func GetMonth(t time.Time) int {
	return int(t.Month())
}

// GetDay gets the day of the month.
func GetDay(t time.Time) int {
	return t.Day()
}

// GetHour gets hour number.
func GetHour(t time.Time) int {
	return t.Hour()
}

// GetMinute gets minute number.
func GetMinute(t time.Time) int {
	return t.Minute()
}

// GetSecond gets second number.
func GetSecond(t time.Time) int {
	return t.Second()
}

// GetMillionSecond gets millionSecond number.
func GetMillionSecond(t time.Time) int {
	return t.Nanosecond() / 1e6
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

// GetHumanReadableDuration gets a duration between times
// and returns format like '01:12:31' (?hour:?minute:?second).
func GetHumanReadableDuration(start time.Time, end time.Time) string {
	v := GetTimestamp(end)/1000 - GetTimestamp(start)/1000 // seconds
	h := v / 3600
	m := v % 3600 / 60
	s := v % 60
	return format2(int(h)) + ":" + format2(int(m)) + ":" + format2(int(s))
}

// GetLongHumanReadableDuration gets a duration between times
// and returns format like '1d 3h 12m 11s' (?day ?hour ?minute ?second).
func GetLongHumanReadableDuration(start time.Time, end time.Time) string {
	v := int(GetTimestamp(end)/1000 - GetTimestamp(start)/1000) // seconds
	return strconv.Itoa(v/86400) + "d " + strconv.Itoa(v%86400/3600) + "h " + strconv.Itoa(v%3600/60) + "m " + strconv.Itoa(v%60) + "s"
}
