// see https://steveperkins.com/go-java-programmers-simple-types/
package hashcode

import (
	"errors"
	"github.com/infavor/gox"
	"math"
	"reflect"
)

// HashCode returns the interface's hashcode.
func HashCode(o interface{}) int32 {
	if o == nil {
		return 0
	}
	t := reflect.TypeOf(o)
	/*for {
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
			continue
		}
		break
	}*/
	switch t.Kind() {
	case reflect.String:
		return stringHashCode(o.(string))
	case reflect.Int8:
		return int8HashCode(o.(int8))
	case reflect.Uint8:
		return uint8HashCode(o.(uint8))
	case reflect.Int16:
		return int16HashCode(o.(int16))
	case reflect.Uint16:
		return uint16HashCode(o.(uint16))
	case reflect.Int:
		return intHashCode(o.(int))
	case reflect.Uint:
		return uintHashCode(o.(uint))
	case reflect.Int32:
		return int32HashCode(o.(int32))
	case reflect.Uint32:
		return uint32HashCode(o.(uint32))
	case reflect.Int64:
		return int64HashCode(o.(int64))
	case reflect.Uint64:
		return uint64HashCode(o.(uint64))
	case reflect.Float32:
		return float32HashCode(o.(float32))
	case reflect.Float64:
		return float64HashCode(o.(float64))
	case reflect.Bool:
		return boolHashCode(o.(bool))
	}
	panic(errors.New("not support type " + t.String()))
}

// stringHashCode returns a string's hashcode.
// This function is copied form Java String.hashCode().
func stringHashCode(o string) int32 {
	if o == "" {
		return 0
	}
	var h int32 = 0
	for _, r := range o {
		h = 31*h + r
	}
	return h
}

func byteHashCode(o byte) int32 {
	return int32(o)
}

func runeHashCode(o rune) int32 {
	return int32(o)
}

func int8HashCode(o int8) int32 {
	return int32(o)
}

func uint8HashCode(o uint8) int32 {
	return int32(o)
}

func int16HashCode(o int16) int32 {
	return int32(o)
}

func uint16HashCode(o uint16) int32 {
	return int32(o)
}

func intHashCode(o int) int32 {
	return int32(o)
}

func uintHashCode(o uint) int32 {
	return int32(o)
}

func int32HashCode(o int32) int32 {
	return o
}

func uint32HashCode(o uint32) int32 {
	return int32(o)
}

func int64HashCode(o int64) int32 {
	return (int32)(o ^ (o >> 32))
}

func uint64HashCode(o uint64) int32 {
	return (int32)(o ^ (o >> 32))
}

func float32HashCode(o float32) int32 {
	return uint32HashCode(math.Float32bits(o))
}

func float64HashCode(o float64) int32 {
	return uint64HashCode(math.Float64bits(o))
}

// see https://stackoverflow.com/questions/3912303/boolean-hashcode
func boolHashCode(o bool) int32 {
	return gox.TValue(o, int32(1231), int32(1237)).(int32)
}
