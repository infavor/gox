package file_test

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/hetianyi/gox/convert"
	"github.com/hetianyi/gox/file"
	"github.com/hetianyi/gox/logger"
	"math/rand"
	"testing"
	"time"
)

func init() {
	logger.Init(nil)
}

func TestCreateEmptyFile(t *testing.T) {
	fi, err := file.CreateFile("D:/tmp/placeholder.txt")
	if err != nil {
		logger.Fatal(err)
	}
	defer fi.Close()
	_, err = fi.Seek(1023, 0)
	if err != nil {
		logger.Fatal(err)
	}
	fi.Write([]byte("\x00"))
	//206,905,344
}

func TestCreateEmptyFile1(t *testing.T) {
	fi, err := file.CreateFile("D:/tmp/placeholder.txt")
	if err != nil {
		logger.Fatal(err)
	}
	defer fi.Close()
	fmt.Println(fi.WriteAt([]byte{222}, 1023))
	//206,905,344
}

func TestCreateEmptyFile2(t *testing.T) {
	fi, err := file.CreateFile("D:/tmp/placeholder.txt")
	if err != nil {
		logger.Fatal(err)
	}
	defer fi.Close()
	fmt.Println(fi.WriteAt([]byte{222}, 1024*1024*1024))
	//206,905,344
}

func Test1(t *testing.T) {
	fmt.Println("HelloFrom")
	fmt.Println("Hello\000From")
}

func TestCrc32(t *testing.T) {
	fmt.Println(file.Crc32("C:\\Users\\hehety\\AppData\\Local\\godfs\\Data\\instance.dat"))
}

func Test2(t *testing.T) {
	var fileLen int64 = 1024
	crc32 := "933736b0"
	instanceId := "5e4d6b56"
	randInt := ""
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 3; i++ {
		randInt += convert.IntToStr(rnd.Intn(10))
	}
	fmt.Println("rand=", randInt)
	var buffer bytes.Buffer
	buffer.WriteString(instanceId)
	buffer.Write(convert.Length2Bytes(fileLen, make([]byte, 8)))
	buffer.WriteString(crc32)
	buffer.WriteString(randInt)

	enc := base64.StdEncoding.EncodeToString(buffer.Bytes())
	fmt.Println(enc)
	// rBNM4lrf6BCAGSzZAAANNrYCTnc900
	bs, _ := base64.StdEncoding.DecodeString(enc)
	fmt.Println(string(bs))
}
