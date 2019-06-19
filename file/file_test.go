package file_test

import (
	"fmt"
	"github.com/hetianyi/gox/file"
	"github.com/hetianyi/gox/logger"
	"testing"
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
