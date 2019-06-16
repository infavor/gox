package hashcode

import (
	"errors"
	"reflect"
)

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
	case reflect.Int:
		return intHashCode(o.(int))
	case reflect.Int64:
		return int64HashCode(o.(int64))
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

func intHashCode(o int) int32 {
	return int32(o)
}

func int64HashCode(o int64) int32 {
	return (int32)(o ^ (o >> 32))
}
