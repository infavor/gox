package binlog_test

import (
	"github.com/hetianyi/gox/file"
	"testing"
)

func TestInit(t *testing.T) {
	out, _ := file.CreateFile("D:\\tmp\\godfs\\block")
	out.WriteAt([]byte{0}, 1024)
	out.WriteAt([]byte{1}, 0)
	out.Close()
}
