// Copyright (C) 2019 tisnyo <tisnyo@gmail.com>.
//
// package convert can helps convert types to another type.
// license that can be found in the LICENSE file.
package convert

import (
	"encoding/base64"
	"encoding/binary"
	"strconv"
	"strings"
)

// IntToStr converts int to string.
func IntToStr(value int) string {
	return strconv.Itoa(value)
}

// Int64ToStr converts int64 to string.
func Int64ToStr(value int64) string {
	return strconv.FormatInt(value, 10)
}

// Int32ToStr converts int32 to string.
func Int32ToStr(value int32) string {
	return strconv.FormatInt(int64(value), 10)
}

// Uint64ToStr converts uint64 to string.
func Uint64ToStr(value uint64) string {
	return strconv.FormatUint(value, 10)
}

// ByteToStr converts byte to string.
func ByteToStr(value byte) string {
	return strconv.Itoa(int(value))
}

// Float32ToStr converts float32 to string.
func Float32ToStr(value float32) string {
	return strconv.FormatFloat(float64(value), 'f', -1, 32)
}

// Float64ToStr converts float64 to string.
func Float64ToStr(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 32)
}

// BoolToStr converts bool to string.
func BoolToStr(value bool) string {
	if value {
		return "true"
	}
	return "false"
}

// StrToInt converts string to int.
func StrToInt(value string) (int, error) {
	return strconv.Atoi(value)
}

// StrToInt32 converts string to int32.
func StrToInt32(value string) (int32, error) {
	v, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return 0, nil
	}
	return int32(v), nil
}

// StrToInt64 converts string to int64.
func StrToInt64(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}

// StrToUint32 converts string to uint32.
func StrToUint32(value string) (uint32, error) {
	v, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return 0, nil
	}
	return uint32(v), nil
}

// StrToUint64 converts string to uint64.
func StrToUint64(value string) (uint64, error) {
	return strconv.ParseUint(value, 10, 64)
}

// StrToByte converts string to byte.
func StrToByte(value string) (byte, error) {
	v, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return 0, nil
	}
	return byte(v), nil
}

// StrToFloat32 converts string to float32.
func StrToFloat32(value string) (float32, error) {
	v, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return 0, nil
	}
	return float32(v), nil
}

// StrToFloat64 converts string to float64.
func StrToFloat64(value string) (float64, error) {
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, nil
	}
	return v, nil
}

// StrToBool converts string to bool.
func StrToBool(value string) (bool, error) {
	return strconv.ParseBool(strings.ToLower(value))
}

// ToBase64 converts an input string to base64 string.
func EncodeBase64(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

// DecodeBase64 decode a base64 string.
func DecodeBase64(input string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(input)
}

// Length2Bytes converts an int64 value to a byte array.
func Length2Bytes(len int64, buffer []byte) []byte {
	binary.BigEndian.PutUint64(buffer, uint64(len))
	return buffer
}

// Bytes2Length converts a byte array to an int64 value.
func Bytes2Length(ret []byte) int64 {
	return int64(binary.BigEndian.Uint64(ret))
}
