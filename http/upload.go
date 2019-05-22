package http

import (
	"bytes"
	"container/list"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"hash"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const (
	ContentDispositionPattern = "^Content-Disposition: form-data; name=\"([^\"]*)\"(; filename=\"([^\"]*)\".*)?$"
	ContentTypePattern        = "^multipart/form-data; boundary=(.*)$"
)

type FileFormReader struct {
	request          *http.Request
	unReadableBuffer           *bytes.Buffer
	atomByte         []byte
	newLineBytesPair []byte
	buffer		 []byte
	newLineBuffer    *bytes.Buffer
}

type FileInfo struct {
	Index    int    `json:"index"`
	FileName string `json:"fileName"`
	Path     string `json:"path"`
}

type HttpUploadResponse struct {
	Status   string              `json:"status"`   // handler result status
	FormData map[string][]string `json:"formData"` // form data
	FileInfo []FileInfo          `json:"fileInfo"` // files info for all uploaded file.
}

func (reader *FileFormReader) Unread(read []byte) {
	reader.unReadableBuffer.Write(read)
}

func (reader *FileFormReader) Read(buff []byte) (int, error) {
	// if buffer of FileFormReader has bytes cached before(from Unread()), then read it first.
	if reader.unReadableBuffer.Len() > 0 {
		if len(buff) <= reader.unReadableBuffer.Len() {
			return reader.unReadableBuffer.Read(buff)
		} else {
			offsetPos, err := reader.unReadableBuffer.Read(buff)
			if err != nil {
				return 0, err
			}
			// read directly from reader
			len, err := reader.request.Body.Read(buff[offsetPos:])
			if err != nil && err != io.EOF {
				return 0, err
			}
			return offsetPos + len, err
		}
	}
	// read directly from reader
	return reader.request.Body.Read(buff)
}

var newLineMarker = []byte{13, 10}

type FileUploadHandler struct {
	writer      http.ResponseWriter
	request     *http.Request
	onTextField func((name string, value string)
}

type StageUploadStatus struct {
	readBodySize  uint64
	sliceReadSize int64
	md            hash.Hash
	sliceMd5      hash.Hash
	fileParts     *list.List
	fileName      string
	index         int
	out           *os.File
	path          string
}


// beginUpload begin to read request entity and parse form field
func (handler *FileUploadHandler) beginUpload() (*HttpUploadResponse, error) {
	var formReader = &FileFormReader{
		request: handler.request,
		unReadableBuffer:  new(bytes.Buffer),
		atomByte: make([]byte, 1),
		newLineBytesPair: make([]byte, 2),
		newLineBuffer: new(bytes.Buffer),
		buffer:make([]byte, 1024*30),
	}

	var ret = &HttpUploadResponse{
		FormData: make(map[string][]string),
	}

	var fileStages list.List
	var fileIndex = 0

	headerContentType := handler.request.Header["Content-Type"]
	contentType := ""
	if headerContentType != nil || len(headerContentType) == 0 {
		contentType = headerContentType[0]
	}
	if mat, _ := regexp.Match(ContentTypePattern, []byte(contentType)); mat {
		boundary := regexp.MustCompile(ContentTypePattern).ReplaceAllString(contentType, "${1}")
		paraSeparator := "--" + boundary
		endSeparator := "--" + boundary + "--"
		for {
			line, e := readNextLine(formReader)
			if e != nil {
				return nil, e
			}
			// if it is paraSeparator, then start read new form text field or file field
			if paraSeparator == line {
				contentDisposition, e1 := readNextLine(formReader)
				if e1 != nil {
					return nil, e1
				} else {
					mat1, e := regexp.Match(ContentDispositionPattern, []byte(contentDisposition))
					if e != nil {
						return nil, e
					}
					paramName := ""
					paramValue := ""
					if mat1 {
						paramName = regexp.MustCompile(ContentDispositionPattern).ReplaceAllString(contentDisposition, "${1}")
					}

					paramContentType, e2 := readNextLine(formReader)
					if e2 != nil {
						return nil, e2
					} else {
						if paramContentType == "" { // read text parameter field
							param, e3 := readNextLine(formReader)
							if e3 != nil {
								return nil, e3
							} else {
								paramValue = param
								handler.onTextField(paramName, paramValue)
							}
						} else { // parse content type
							mat2, _ := regexp.Match(ContentDispositionPattern, []byte(contentDisposition))
							if e != nil {
								return nil, e
							}
							fileName := ""
							if mat2 {
								fileName = regexp.MustCompile(ContentDispositionPattern).ReplaceAllString(contentDisposition, "${3}")
							}

							_, e3 := readNextLine(formReader)
							if e3 != nil {
								return nil, e3
							} else { // read file body
								stageUploadStatus, e4 := readFileBody(formReader, paraSeparator)
								if e4 != nil {
									return nil, e4
								}
								fileStages.PushBack(stageUploadStatus)
								stageUploadStatus.index = fileIndex
								handler.onTextField(paramName, stageUploadStatus.path)
								stageUploadStatus.fileName = fileName
								fileIndex++
							}
						}
					}

				}
			} else if endSeparator == line {
				// form stream hit end
				break
			} else {
				logger.Error("unknown line")
			}
		}
	}
	// copy string list to string[]
	for k, v := range handler.params {
		if v != nil {
			tmp := make([]string, v.Len())
			index := 0
			for ele := v.Front(); ele != nil; ele = ele.Next() {
				tmp[index] = ele.Value.(string)
				index++
			}
			ret.FormData[k] = tmp
		}
	}

	fInfo := make([]FileInfo, fileStages.Len())
	k := 0
	for ele := fileStages.Front(); ele != nil; ele = ele.Next() {
		stage := ele.Value.(*StageUploadStatus)
		info := &FileInfo{
			Index:    stage.index,
			FileName: stage.fileName,
			Path:     stage.path,
		}
		fInfo[k] = *info
		k++
	}
	ret.FileInfo = fInfo
	ret.Status = "success"
	return ret, nil
}

// readNextLine reads next form field meta string.
func readNextLine(reader *FileFormReader) (string, error) {
	for {
		len, err := reader.Read(reader.atomByte)
		if err != nil && err != io.EOF {
			return "", err
		}
		if len != 1 {
			return "", errors.New("error read from stream[0]")
		}
		reader.newLineBytesPair[0] = reader.newLineBytesPair[1]
		reader.newLineBytesPair[1] = reader.atomByte[0]
		reader.newLineBuffer.Write(reader.atomByte)
		if bytes.Equal(newLineMarker, reader.newLineBytesPair) {
			return string(reader.newLineBuffer.Bytes()[0 : reader.newLineBuffer.Len() - 2]), nil
		}
	}
}

// readFileBody reads a file body part.
func readFileBody(reader *FileFormReader, separator string) (*StageUploadStatus, error) {
	if oe != nil {
		return nil, oe
	}
	stateUploadStatus := &StageUploadStatus{
		readBodySize:  0,
		sliceReadSize: 0,
		sliceMd5:      md5.New(),
		md:            md,
		fileParts:     list.New(),
		out:           out,
	}
	separator = "\r\n" + separator
	buff2, _ := bridgev2.MakeBytes(int64(len(separator)), true, 1024, true)
	tail, _ := bridgev2.MakeBytes(int64(len(separator)*2), true, 1024, true)
	for {
		len1, e1 := reader.Read(reader.buffer)
		if e1 != nil {
			if e1 != io.EOF {
				return nil, e1
			}
		}
		if len1 == 0 {
			return nil, errors.New("read file body failed1")
		}
		// whether buff1 contains separator
		i1 := bytes.Index(reader.buffer, []byte(separator))
		if i1 != -1 {
			out.Write(reader.buffer[0:i1])
			e8 := handleStagePartFile(reader.buffer[0:i1], stateUploadStatus)
			if e8 != nil {
				return nil, e8
			}
			reader.Unread(reader.buffer[i1+2 : len1]) // skip "\r\n"
			break
		} else {
			len2, e2 := reader.Read(buff2)
			if e2 != nil {
				if e2 != io.EOF {
					return nil, e2
				}
			}
			if len2 == 0 {
				return nil, errors.New("read file body failed2")
			}
			// []byte tail is last bytes of buff1 and first bytes of buff2 in case broken separator.
			if len1 >= len(separator) {
				ByteCopy(tail, 0, len(separator), reader.buffer[len1-len(separator):len1])
			}
			if len2 >= len(separator) {
				ByteCopy(tail, len(separator), len(tail), buff2[0:len(separator)])
			}

			i2 := bytes.Index(tail, []byte(separator))
			if i2 != -1 {
				if i2 < len(separator) {
					e8 := handleStagePartFile(reader.buffer[0:len1-i2], stateUploadStatus)
					if e8 != nil {
						return nil, e8
					}
					reader.Unread(reader.buffer[len1-i2+2 : len1])
					reader.Unread(buff2[0:len2])
				} else {
					e8 := handleStagePartFile(reader.buffer[0:len1], stateUploadStatus)
					if e8 != nil {
						return nil, e8
					}
					reader.Unread(buff2[i2-len(separator)+2 : len2])
				}
				break
			} else {
				e8 := handleStagePartFile(reader.buffer[0:len1], stateUploadStatus)
				if e8 != nil {
					return nil, e8
				}
				reader.Unread(buff2[0:len2])
			}
		}
	}
	stateUploadStatus.out.Close()
	if stateUploadStatus.sliceReadSize > 0 {
		sliceCipherStr := stateUploadStatus.sliceMd5.Sum(nil)
		sMd5 := hex.EncodeToString(sliceCipherStr)
		stateUploadStatus.sliceMd5.Reset()
		e10 := libcommon.MoveTmpFileTo(sMd5, stateUploadStatus.out)
		if e10 != nil {
			return nil, e10
		}

		tmpPart := &app.PartDO{Md5: sMd5, Size: stateUploadStatus.sliceReadSize}
		stateUploadStatus.fileParts.PushBack(tmpPart)
		app.UpdateDiskUsage(stateUploadStatus.sliceReadSize)
	}
	sliceCipherStr := md.Sum(nil)
	sMd5 := hex.EncodeToString(sliceCipherStr)

	finalFile := &app.FileVO{
		Md5:        sMd5,
		PartNumber: stateUploadStatus.fileParts.Len(),
		Group:      app.Group,
		Instance:   app.InstanceId,
		Finish:     1,
		FileSize:   0,
		Flag:       flag,
	}
	parts := make([]app.PartDO, stateUploadStatus.fileParts.Len())
	index := 0
	var totalSize int64 = 0
	for ele := stateUploadStatus.fileParts.Front(); ele != nil; ele = ele.Next() {
		parts[index] = *ele.Value.(*app.PartDO)
		totalSize += parts[index].Size
		index++
	}
	finalFile.Parts = parts
	finalFile.FileSize = totalSize
	// stoe := libservice.StorageAddFile(sMd5, app.Group, stateUploadStatus.fileParts)
	stoe := libservicev2.InsertFile(finalFile, nil)
	if stoe != nil {
		return nil, stoe
	}
	// mark the file is multi part or single part
	if stateUploadStatus.fileParts.Len() > 1 {
		stateUploadStatus.path = app.Group + "/" + app.InstanceId + "/M/" + sMd5
	} else {
		stateUploadStatus.path = app.Group + "/" + app.InstanceId + "/S/" + sMd5
	}
	app.UpdateUploads()
	return stateUploadStatus, nil
}



// ByteCopy copies bytes
func ByteCopy(src []byte, start int, end int, cp []byte) {
	for i := range src {
		if i >= start && i < end {
			src[i] = cp[i]
		} else {
			break
		}
	}
}
