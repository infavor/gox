// Copyright (C) 2019 tisnyo <tisnyo@gmail.com>.
//
// package convert can helps convert types to another type.
// license that can be found in the LICENSE file.
package convert

import "strconv"

// Int2Str converts int to string.
func IntToStr(value int) string {
	return strconv.Itoa(value)
}

// Int2Str converts int to string.
func Int64ToStr(value int64) string {
	return strconv.FormatInt(value, 10)
}

// Int2Str converts int to string.
func Uint64ToStr(value uint64) string {
	return strconv.FormatUint(value, 10)
}

// Int2Str converts int to string.
func ByteToStr(value byte) string {
	return strconv.Itoa(int(value))
}

// Int2Str converts int to string.
func Float32ToStr(value float32) string {
	return strconv.FormatFloat(float64(value), 'b', -1, 32)
}

// Str2Int converts string to int.
func StrToInt(value string) (int, error) {
	return strconv.Atoi(value)
}
