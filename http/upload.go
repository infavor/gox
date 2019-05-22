package http

import (
	"bytes"
	"container/list"
	"errors"
	"fmt"
	"hash"
	"io"
	"net/http"
	"os"
	"regexp"
)

const (
	ContentTypePattern        = "^multipart/form-data; boundary=(.*)$"
	ContentDispositionPattern = "^Content-Disposition: form-data; name=\"([^\"]*)\"(; filename=\"([^\"]*)\".*)?$"
)

var (
	RegexContentTypePattern        = regexp.MustCompile(ContentTypePattern)
	RegexContentDispositionPattern = regexp.MustCompile(ContentDispositionPattern)
)

type FileFormReader struct {
	request          *http.Request
	unReadableBuffer *bytes.Buffer
	atomByte         []byte
	newLineBytesPair []byte
	buffer           []byte
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
	writer               http.ResponseWriter
	request              *http.Request
	params               map[string]*list.List
	boundary             string
	paraBoundary         string
	endParaBoundary      string
	separator            []byte
	separatorTestBuffer  []byte
	separatorMergeBuffer []byte
	onTextField          func(name string, value string)
	onFileField          func(name string, value string)
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
		request:          handler.request,
		unReadableBuffer: new(bytes.Buffer),
		atomByte:         make([]byte, 1),
		newLineBytesPair: make([]byte, 2),
		newLineBuffer:    new(bytes.Buffer),
		buffer:           make([]byte, 1024*30),
	}
	var ret = &HttpUploadResponse{
		FormData: make(map[string][]string),
	}

	var fileStages list.List
	var fileIndex = 0

	headerContentType := handler.request.Header["Content-Type"]
	contentType := ""
	if headerContentType != nil && len(headerContentType) > 0 {
		contentType = headerContentType[0]
	}
	if RegexContentTypePattern.Match([]byte(contentType)) {
		handler.boundary = RegexContentDispositionPattern.ReplaceAllString(contentType, "${1}")
		handler.paraBoundary = "--" + handler.boundary
		handler.endParaBoundary = "--" + handler.boundary + "--"
		handler.separator = []byte("\r\n" + handler.boundary)
		handler.separatorTestBuffer = make([]byte, len(handler.separator))
		handler.separatorMergeBuffer = make([]byte, len(handler.separator)*2)
		for {
			line, err := readNextLine(formReader)
			if err != nil {
				return nil, err
			}
			// if it is paraSeparator, then start read new form text field or file field
			if handler.paraBoundary == line {
				contentDisposition, err := readNextLine(formReader)
				if err != nil {
					return nil, err
				} else {
					mat1, err := regexp.Match(ContentDispositionPattern, []byte(contentDisposition))
					if err != nil {
						return nil, err
					}
					paramName := ""
					paramValue := ""
					if mat1 {
						paramName = regexp.MustCompile(ContentDispositionPattern).ReplaceAllString(contentDisposition, "${1}")
					}

					paramContentType, err := readNextLine(formReader)
					if err != nil {
						return nil, err
					} else {
						if paramContentType == "" { // read text parameter field
							param, err := readNextLine(formReader)
							if err != nil {
								return nil, err
							} else {
								paramValue = param
								handler.onTextField(paramName, paramValue)
							}
						} else { // parse content type
							mat2, err := regexp.Match(ContentDispositionPattern, []byte(contentDisposition))
							if err != nil {
								return nil, err
							}
							fileName := ""
							if mat2 {
								fileName = regexp.MustCompile(ContentDispositionPattern).ReplaceAllString(contentDisposition, "${3}")
							}
							fmt.Println(fileName)

							_, err = readNextLine(formReader) // read blank line
							if err != nil {
								return nil, err
							} else { // read file body
								err := handler.readFileBody(formReader)
								if err != nil {
									return nil, err
								}
								handler.onTextField(paramName, "")
								fileIndex++
							}
						}
					}

				}
			} else if handler.endParaBoundary == line {
				// form stream hit end
				break
			} else {
				fmt.Println("unknown line")
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
			return string(reader.newLineBuffer.Bytes()[0 : reader.newLineBuffer.Len()-2]), nil
		}
	}
}

// readFileBody reads a file body part.
func (handler *FileUploadHandler) readFileBody(reader *FileFormReader) error {
	separatorLength := len(handler.separator)
	for {
		len1, err := reader.Read(reader.buffer)
		if err != nil && err != io.EOF {
			return err
		}
		if len1 == 0 {
			return errors.New("read file body failed[0]")
		}
		// whether buff1 contains separator
		pos := bytes.Index(reader.buffer, handler.separator)
		if pos != -1 {
			// out.Write(reader.buffer[0:pos])
			// write bytes......
			///-----------------
			reader.Unread(reader.buffer[pos+2 : len1]) // skip "\r\n"
			break
		} else {
			len2, err := reader.Read(handler.separatorTestBuffer)
			if err != nil {
				if err != io.EOF {
					return err
				}
			}
			if len2 == 0 {
				return errors.New("read file body failed[1]")
			}
			// []byte tail is last bytes of buff1 and first bytes of buff2 in case broken separator.
			//
			if len1 >= separatorLength {
				ByteCopy(handler.separatorMergeBuffer, 0, separatorLength, reader.buffer[len1-separatorLength:len1])
			}
			if len2 >= separatorLength {
				ByteCopy(handler.separatorMergeBuffer, separatorLength, len(handler.separatorMergeBuffer), handler.separatorTestBuffer[0:separatorLength])
			}

			i2 := bytes.Index(handler.separatorMergeBuffer, handler.separator)
			if i2 != -1 {
				if i2 < separatorLength {
					reader.Unread(reader.buffer[len1-i2+2 : len1])
					reader.Unread(handler.separatorTestBuffer[0:len2])
				} else {
					reader.Unread(handler.separatorTestBuffer[i2-separatorLength+2 : len2])
				}
				break
			} else {
				reader.Unread(handler.separatorTestBuffer[0:len2])
			}
		}
	}
	return nil
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
