// Copyright (C) 2019 tisnyo <tisnyo@gmail.com>.
//
// package file contains some file operations.
// license that can be found in the LICENSE file.
package file

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"syscall"
)

// GetFile opens file in readonly mode.
func GetFile(path string) (*os.File, error) {
	return os.Open(path)
}

// GetFileMd5 gets file's md5.
func GetFileMd5(fi string) (string, error) {
	md := md5.New()
	f, e := GetFile(fi)
	if e != nil {
		return "", e
	} else {
		defer f.Close()
		_, e1 := io.Copy(md, f)
		if e1 != nil {
			return "", e1
		}
		md5 := hex.EncodeToString(md.Sum(nil))
		return md5, nil
	}
}

// CopyFile copies src file to dest file.
func CopyFile(src string, dest string) (bool, error) {
	srcFile, err := GetFile(src)
	defer srcFile.Close()
	if err != nil {
		return false, err
	} else {
		destFile, err := os.OpenFile(dest, syscall.O_CREAT|os.O_WRONLY|syscall.O_TRUNC, 0660)
		if err != nil {
			return false, err
		} else {
			defer destFile.Close()
			_, err = io.Copy(destFile, srcFile)
			if err != nil {
				return false, nil
			}
			return true, nil
		}
	}
}

// CopyFileTo copies a file to target dir.
func CopyFileTo(src string, dir string) (bool, error) {
	srcFile, err := GetFile(src)
	defer srcFile.Close()
	if err != nil {
		return false, err
	} else {
		fileInfo, _ := srcFile.Stat()
		destFile, err := os.OpenFile(FixPath(dir)+string(os.PathSeparator)+fileInfo.Name(),
			syscall.O_CREAT|os.O_WRONLY|syscall.O_TRUNC, 0660)
		if err != nil {
			return false, err
		} else {
			defer destFile.Close()
			_, err = io.Copy(destFile, srcFile)
			if err != nil {
				return false, nil
			}
			return true, nil
		}
	}
}

// Exists check whether the file exists.
func Exists(path string) bool {
	fi, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	if fi == nil {
		return false
	}
	return true
}

// Delete deletes a file or directory.
// special: it returns true if the file not exists.
func Delete(path string) bool {
	if !Exists(path) {
		return true
	}
	err := os.Remove(path)
	return nil == err
}

// DeleteAll deletes file or directory.
// if it is a directory, it will try to delete all files below.
// special: it returns true if the file not exists.
func DeleteAll(path string) bool {
	return os.RemoveAll(path) == nil
}

// CreateFile creates a new file, truncating it if it already exists
func CreateFile(path string) (*os.File, error) {
	return os.Create(path)
}

// AppendFile opens a file to append or creates it if it is not exist.
func AppendFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
}

// CreateDir creates a new directory.
func CreateDir(path string) error {
	return os.Mkdir(path, 0755)
}

// CreateDirs creates a directory named path,
// along with any necessary parents, and returns nil, or else returns an error.
// The permission bits perm (before umask) are used for all directories that MkdirAll creates.
// If path is already a directory, MkdirAll does nothing and returns nil.
func CreateDirs(path string) error {
	return os.MkdirAll(path, 0755)
}

// IsFile1 checks whether the path represents a file.
func IsFile1(path string) bool {
	fi, err := os.Stat(path)
	if nil == err {
		return !fi.IsDir()
	} else {
		return false
	}
}

// IsFile2 checks whether the file represents a file.
func IsFile2(file *os.File) bool {
	fi, err := file.Stat()
	if nil == err {
		return !fi.IsDir()
	} else {
		return false
	}
}

// IsDir1 checks whether the path represents a file.
func IsDir1(path string) bool {
	fi, err := os.Stat(path)
	if nil == err {
		return fi.IsDir()
	} else {
		return false
	}
}

// IsDir2 checks whether the path represents a file.
func IsDir2(file *os.File) bool {
	fi, err := file.Stat()
	if nil == err {
		return fi.IsDir()
	} else {
		return false
	}
}

// MoveFile renames (moves) oldpath to newpath.
func MoveFile(src string, dest string) error {
	return os.Rename(src, dest)
}

// ChangeWorkDir changes the current working directory to the named directory.
func ChangeWorkDir(path string) error {
	return os.Chdir(path)
}

// GetWorkDir returns a rooted path name corresponding to the current directory.
// If the current directory can be reached via multiple paths (due to symbolic links),
// Getwd may return any one of them.
func GetWorkDir() (string, error) {
	return os.Getwd()
}

// GetTempDir returns the default directory to use for temporary files.
// On Unix systems, it returns $TMPDIR if non-empty, else /tmp. On Windows,
// it uses GetTempPath, returning the first non-empty value from %TMP%, %TEMP%, %USERPROFILE%,
// or the Windows directory. On Plan 9, it returns /tmp.
// The directory is neither guaranteed to exist nor have accessible permissions.
func GetTempDir() string {
	return os.TempDir()
}

// IsAbsPath reports whether the path is absolute.
func IsAbsPath(path string) bool {
	return filepath.IsAbs(path)
}

// GetFileExt returns the file name extension used by path.
// The extension is the suffix beginning at the final dot in the final slash-separated element of path;
// it is empty if there is no dot.
func GetFileExt(filePath string) string {
	return path.Ext(filePath)
}

// FixPath fixed file path in a simple way, examples:
//    "/aaa/aa\\bb\\cc/d/////"     ->  "/aaa/aa/bb/cc/d"
//    "E:/aaa/aa\\bb\\cc/d////e/"  ->  "E:/aaa/aa/bb/cc/d/e"
//    ""                           ->  "."
//    "/"                          ->  "/"
func FixPath(input string) string {
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return "."
	}
	// replace windows path separator '\' to '/'
	replacement := strings.Replace(input, "\\", "/", -1)

	for {
		if strings.Contains(replacement, "//") {
			replacement = strings.Replace(replacement, "//", "/", -1)
			continue
		}
		if replacement == "/" {
			return replacement
		}
		len := len(replacement)
		if len <= 0 {
			break
		}
		if replacement[len-1:] == "/" {
			replacement = replacement[0 : len-1]
		} else {
			break
		}
	}
	return replacement
}
