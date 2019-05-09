package gpip

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"github.com/json-iterator/go"
	"io"
	"net"
	"reflect"
)

const (
	// frameHeadSize frame head byte array size
	frameHeadSize = 8
)

// PipFrame defines a Pip request instance.
type PipFrame struct {
	frameHead  []byte
	Meta       interface{}
	metaBody   []byte
	bodyLength int64
	bodyReader io.Reader
}

// SetBodyReader set PipRequest's body input stream.
func (frame *PipFrame) SetBodyReader(reader io.Reader, length int64) {
	if length < 0 {
		panic("body length must be positive number")
	}
	frame.bodyReader = reader
	frame.bodyLength = length
}

// GetMeta set PipRequest's body input stream.
func (frame *PipFrame) GetMeta(p reflect.Type) (interface{}, error) {
	return Deserialize(frame.metaBody, p)
}

// Pip is a tcp pipe manager which controls the reception and sending of frames.
type Pip struct {
	Conn              net.Conn
	BodyReaderHandler func(frame *PipFrame, bodyLength int64, bodyReader io.Reader) error
}

// Receive receives a frame by it's pip.
func (pip *Pip) Receive() (*PipFrame, error) {
	metaLenBytes := make([]byte, frameHeadSize)
	if _, err := io.ReadFull(pip.Conn, metaLenBytes); err != nil {
		return nil, err
	}
	bodyLenBytes := make([]byte, frameHeadSize)
	if _, err := io.ReadFull(pip.Conn, bodyLenBytes); err != nil {
		return nil, err
	}
	metaLen := ConvertBytes2Len(&metaLenBytes)
	bodyLen := ConvertBytes2Len(&bodyLenBytes)
	metaBs := make([]byte, metaLen)
	if _, err := io.ReadFull(pip.Conn, metaBs); err != nil {
		return nil, err
	}
	frame := &PipFrame{
		frameHead:  append(metaLenBytes, bodyLenBytes...),
		metaBody:   metaBs,
		bodyLength: bodyLen,
	}
	if bodyLen > 0 {
		if pip.BodyReaderHandler == nil {
			return nil, errors.New("no body reader handler provided")
		}
		err := pip.BodyReaderHandler(frame, bodyLen, io.LimitReader(pip.Conn, bodyLen))
		if err != nil {
			return nil, err
		}
	}
	return frame, nil
}

// Send sends a frame by it's pip.
func (pip *Pip) Send(frame *PipFrame) error {
	autoFillFrame(frame)
	if _, err := pip.Conn.Write(frame.frameHead); err != nil {
		return err
	}
	if _, err := pip.Conn.Write(frame.metaBody); err != nil {
		return err
	}
	if frame.bodyLength > 0 && frame.bodyReader != nil {
		if _, err := io.CopyN(pip.Conn, frame.bodyReader, frame.bodyLength); err != nil {
			return err
		}
	}
	return nil
}

// Serialize serialize an interface to json.
func Serialize(obj interface{}) ([]byte, error) {
	return jsoniter.Marshal(obj)
}

// Deserialize deserialize an byte array to an interface by type.
func Deserialize(data []byte, p reflect.Type) (interface{}, error) {
	o := reflect.New(p).Interface()
	err := json.Unmarshal(data, &o)
	return o, err
}

// autoFillFrame fills frame head bytes and meta body bytes.
func autoFillFrame(frame *PipFrame) error {
	metaBs, err := Serialize(frame.Meta)
	if err != nil {
		return err
	}
	frame.metaBody = metaBs
	frameHead := make([]byte, frameHeadSize*2)
	ConvertLen2Bytes(int64(len(metaBs)), frameHead[0:frameHeadSize])
	ConvertLen2Bytes(frame.bodyLength, frameHead[frameHeadSize:])
	frame.frameHead = frameHead
	return nil
}

// ConvertLen2Bytes converts an int64 value to an 8 bytes array.
func ConvertLen2Bytes(len int64, buffer []byte) {
	binary.BigEndian.PutUint64(buffer, uint64(len))
}

// ConvertBytes2Len converts an 8 bytes array to an int64 value.
func ConvertBytes2Len(ret *[]byte) int64 {
	return int64(binary.BigEndian.Uint64(*ret))
}
