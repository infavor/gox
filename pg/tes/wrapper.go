package main

import (
	"github.com/hetianyi/gox/file"
	"github.com/hetianyi/gox/pg"
	"io"
)

func main() {
	fi, _ := file.GetFile("D:\\Captures\\FarCry_ 2018-03-22 16_16_08.mp4")
	defer fi.Close()
	inf, _ := fi.Stat()
	out, _ := file.CreateFile("E:\\WorkSpace2018\\godfs\\bin\\123.mp4")
	defer out.Close()
	/*ww := &pg.WrappedReader{Reader: fi}
	pg.NewWrappedReaderProgress(inf.Size(), 50, "FarCry_ 2018-03-22 16_16_08.mp4", pg.Top, ww)
	io.Copy(out, ww)*/
	ww := &pg.WrappedWriter{Writer: out}
	pg.NewWrappedWriterProgress(inf.Size(), 50, "FarCry_ 2018-03-22 16_16_08.mp4", pg.Top, ww)
	io.Copy(ww, fi)
}
