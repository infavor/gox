package compressx_test

import (
	"github.com/hetianyi/gox/logger"
	"github.com/mholt/archiver"
	"github.com/sirupsen/logrus"
	"testing"
)

func init() {
	logger.Init(nil)
}

func TestCompress(t *testing.T) {
	//logrus.Info(compressx.Compress("D:/tmp.tar.gz", compressx.GZIP,
	//							"D:\\tmp"))

	logrus.Info(archiver.Archive([]string{"D:/tmp"}, "D:/tmp.tar.gz"))

}
