package convert_test

import (
	"fmt"
	"github.com/hetianyi/gox/convert"
	"os"
	"testing"
)

func TestAll(t *testing.T) {
	fmt.Println(convert.IntToStr(1))
	fmt.Println(convert.UintToStr(1))

	fmt.Println(convert.Int8ToStr(1))
	fmt.Println(convert.Uint8ToStr(1))

	fmt.Println(convert.Int16ToStr(1))
	fmt.Println(convert.Uint16ToStr(65535))

	fmt.Println(convert.Int32ToStr(1))
	fmt.Println(convert.Uint32ToStr(1))

	fmt.Println(convert.Int64ToStr(1))
	fmt.Println(convert.Uint64ToStr(1))

	fmt.Println(convert.ByteToStr(112))
	fmt.Println(convert.Float32ToStr(12.11))
	fmt.Println(convert.Float64ToStr(12.11))
	fmt.Println(convert.BoolToStr(true), convert.BoolToStr(false))

	fmt.Println(convert.StrToInt("123"))
	fmt.Println(convert.StrToUint("123"))

	fmt.Println(convert.StrToInt8("123"))
	fmt.Println(convert.StrToUint8("123"))

	fmt.Println(convert.StrToInt16("123"))
	fmt.Println(convert.StrToUint16("123"))

	fmt.Println(convert.StrToInt32("123"))
	fmt.Println(convert.StrToUint32("123"))

	fmt.Println(convert.StrToInt64("123"))
	fmt.Println(convert.StrToUint64("123"))

	fmt.Println(convert.StrToByte("123"))
	fmt.Println(convert.StrToFloat32("123.123333"))
	fmt.Println(convert.StrToFloat64("123.123333"))
	fmt.Println(convert.StrToBool("true"))
	fmt.Println(convert.StrToBool("false"))

	f, _ := os.Open("D:/zhigan.png")
	info, _ := f.Stat()
	fmt.Println(info.Name())
}
