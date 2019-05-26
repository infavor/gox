package compressx

import (
	"archive/tar"
	"compress/gzip"
	"github.com/hetianyi/gox"
	"github.com/hetianyi/gox/file"
	"io"
	"log"
	"os"
)

// TODO bug in mime and https://github.com/mholt/archiver
type FORMAT byte

const (
	GZIP FORMAT = 1
	ZLIB FORMAT = 1
)

// Compress compresses src file/dir as target file 'dest'.
func Compress(dest string, format FORMAT, src ...string) error {
	if format == GZIP {
		return startGZIP(dest, src...)
	}
	return nil
}

// startGZIP compress in gzip.
func startGZIP(dest string, src ...string) error {
	// create output: .tar.gz
	destOut, err := file.CreateFile(dest)
	if err != nil {
		return err
	}
	defer destOut.Close()
	tw := tar.NewWriter(gzip.NewWriter(destOut))

	for _, root := range src {
		fi, err := file.GetFile(root)
		if err != nil {
			return err
		}
		if err := gzipCompress(fi, tw, ""); err != nil {
			return err
		}
	}
	if err := tw.Close(); err != nil {
		log.Fatal(err)
	}
	return nil
}

func gzipCompress(parent *os.File, tw *tar.Writer, prefix string) error {
	parentInfo, err := parent.Stat()
	if err != nil {
		return err
	}
	if parentInfo.IsDir() {
		prefix = gox.TValue(prefix == "", "", prefix+"/").(string) + parentInfo.Name()
		infos, err := file.ListFiles(parent.Name())
		if err != nil {
			return err
		}
		if infos == nil || len(infos) == 0 {
			header := &tar.Header{
				Name:     prefix,
				Mode:     0600,
				Size:     0,
				Typeflag: tar.TypeDir,
			}
			tw.WriteHeader(header)
			tw.Write(nil)
			return nil
		}
		for _, info := range infos {
			subFile, err := file.GetFile(parent.Name() + string(os.PathSeparator) + info.Name())
			if err != nil {
				return err
			}
			if err := gzipCompress(subFile, tw, prefix); err != nil {
				return err
			}
		}
	} else {
		header := &tar.Header{
			Name: gox.TValue(prefix == "", "", prefix+"/").(string) + parentInfo.Name(),
			Mode: 0600,
			Size: parentInfo.Size(),
		}
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		if _, err = io.Copy(tw, parent); err != nil {
			return err
		}
	}
	return nil
}
