// Copyright (C) 2019 tisnyo <tisnyo@gmail.com>.
//
// A Pip can used for frame data transfer.
// It can transfer text header content and byte content such as file.
package gpip

import (
	"encoding/json"
	"github.com/hetianyi/gox/convert"
	"github.com/json-iterator/go"
	"io"
	"net"
	"reflect"
)

const (
	// frameHeadSize frame head byte array size
	frameHeadSize = 8
)

// pipFrame defines a pipe request instance.
type pipFrame struct {
	// head bytes of frame
	frameHead []byte
	// frame header object
	Header interface{}
	// bytes array of frame header
	headerBody []byte
	// frame body length
	bodyLength int64
	// frame body reader
	bodyReader io.Reader
}

// SetBodyReader set PipRequest's body input stream.
func (frame *pipFrame) SetBodyReader(reader io.Reader, length int64) {
	if length < 0 {
		panic("body length must be positive number")
	}
	frame.bodyReader = reader
	frame.bodyLength = length
}

// GetMeta set PipRequest's body input stream.
func (frame *pipFrame) GetMeta(initialObj interface{}) error {
	return DeserializeFromObject(frame.headerBody, initialObj)
}

// Pip is a tcp pipe manager which controls the reception and sending of frames.
type Pip struct {
	Conn net.Conn
}

// Close close the Pip.
func (pip *Pip) Close() {
	if pip.Conn != nil {
		pip.Conn.Close()
	}
}

// Receive receives a frame by it's pip.
// headerObject is an interface which will be used to load header data,
// handler is a callback function with which you can handle a data frame,
// parameter filledMetaObject of function handler and headerObject are the same object.
// filledHeaderObject is data-filled header Object.
func (pip *Pip) Receive(
	headerObject interface{},
	handler func(filledHeaderObject interface{}, bodyReader io.Reader, bodyLength int64) error,
) error {
	headerLenBytes := make([]byte, frameHeadSize)
	if _, err := io.ReadFull(pip.Conn, headerLenBytes); err != nil {
		return err
	}
	bodyLenBytes := make([]byte, frameHeadSize)
	if _, err := io.ReadFull(pip.Conn, bodyLenBytes); err != nil {
		return err
	}
	headerLen := convert.Bytes2Length(headerLenBytes)
	bodyLen := convert.Bytes2Length(bodyLenBytes)
	headerBs := make([]byte, headerLen)
	if _, err := io.ReadFull(pip.Conn, headerBs); err != nil {
		return err
	}
	frame := &pipFrame{
		frameHead:  append(headerLenBytes, bodyLenBytes...),
		headerBody: headerBs,
		bodyLength: bodyLen,
	}
	err := frame.GetMeta(headerObject)
	if err != nil {
		return err
	}
	if bodyLen > 0 {
		handler(headerObject, io.LimitReader(pip.Conn, bodyLen), bodyLen)
	} else {
		handler(headerObject, nil, 0)
	}
	return nil
}

// Send sends a frame by it's pip.
func (pip *Pip) Send(header interface{}, bodyReader io.Reader, bodyLength int64) error {
	frame := &pipFrame{
		Header: header,
	}
	frame.SetBodyReader(bodyReader, bodyLength)
	autoFillFrame(frame)
	if _, err := pip.Conn.Write(frame.frameHead); err != nil {
		return err
	}
	if _, err := pip.Conn.Write(frame.headerBody); err != nil {
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

// DeserializeFromType deserialize an byte array to an interface by type.
func DeserializeFromType(data []byte, p reflect.Type) (interface{}, error) {
	o := reflect.New(p).Interface()
	err := json.Unmarshal(data, &o)
	return o, err
}

// DeserializeFromObject deserialize an byte array to an interface by type.
func DeserializeFromObject(data []byte, obj interface{}) error {
	return json.Unmarshal(data, &obj)
}

// autoFillFrame fills frame head bytes and header body bytes.
func autoFillFrame(frame *pipFrame) error {
	headerBs, err := Serialize(frame.Header)
	if err != nil {
		return err
	}
	frame.headerBody = headerBs
	frameHead := make([]byte, frameHeadSize*2)
	convert.Length2Bytes(int64(len(headerBs)), frameHead[0:frameHeadSize])
	convert.Length2Bytes(frame.bodyLength, frameHead[frameHeadSize:])
	frame.frameHead = frameHead
	return nil
}
