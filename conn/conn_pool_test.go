package conn_test

import (
	"container/list"
	"encoding/json"
	"fmt"
	"github.com/hetianyi/gox"
	"github.com/hetianyi/gox/conn"
	"github.com/hetianyi/gox/gpip"
	"github.com/hetianyi/gox/logger"
	"io"
	"log"
	"net"
	"testing"
	"time"
)

type Header struct {
	Code      byte              `json:"code"`
	Attribute map[string]string `json:"attrs"`
}

func init() {
	logger.Init(&logger.Config{
		Level: logger.DebugLevel,
	})
}

func TestServer(t *testing.T) {
	listener, err := net.Listen("tcp", ":8899")
	if err != nil {
		log.Fatal("error start server:", err)
	}
	logger.Info("server is listening on port: 8899")
	for {
		conn, err := listener.Accept()
		logger.Info("accepted new connection")
		if err != nil {
			logger.Error("error accept new connection:", err)
			continue
		}
		gox.Try(func() {
			go serverHandler(conn)
		}, func(i interface{}) {
			logger.Error("connection error:", err)
		})
	}
}

func TestClient(t *testing.T) {
	server := &conn.Server{
		Host: "127.0.0.1",
		Port: 8899,
	}
	conn.InitServerSettings(server, 1000, time.Second*5)

	cache := list.New()

	// TODO bug fix error: connection pool is full
	for i := 0; i < 10000000; i++ {
		if i > 0 && i%500 == 0 {
			gox.WalkList(cache, func(item interface{}) bool {
				conn.ReturnConnection(server, item.(*net.Conn),
					"01092391231231231023sdkasdasdaksdkasjdajsdjasdjalsjdlasjdljalsd01092391231231231023sdkasdasdaksdkasjdajsdjasdjalsjdlasjdljalsd01092391231231231023sdkasdasdaksdkasjdajsdjasdjalsjdlasjdljalsd01092391231231231023sdkasdasdaksdkasjdajsdjasdjalsjdlasjdljalsd01092391231231231023sdkasdasdaksdkasjdajsdjasdjalsjdlasjdljalsd01092391231231231023sdkasdasdaksdkasjdajsdjasdjalsjdlasjdljalsd01092391231231231023sdkasdasdaksdkasjdajsdjasdjalsjdlasjdljalsd01092391231231231023sdkasdasdaksdkasjdajsdjasdjalsjdlasjdljalsd01092391231231231023sdkasdasdaksdkasjdajsdjasdjalsjdlasjdljalsd01092391231231231023sdkasdasdaksdkasjdajsdjasdjalsjdlasjdljalsd",
					false)
				return false
			})
			cache = list.New()
			time.Sleep(time.Second * 7)
		}
		c, _, err := conn.GetConnection(server)
		if err != nil {
			logger.Error("error: ", err)
		}
		cache.PushBack(c)
	}
}

func TestAttr(t *testing.T) {
	server := &conn.Server{
		Host: "127.0.0.1",
		Port: 8899,
	}
	conn.InitServerSettings(server, 1000, time.Second*5)

	cache := list.New()

	// TODO bug fix error: connection pool is full
	for i := 1; i <= 100; i++ {
		if i > 0 && i%100 == 0 {
			gox.WalkList(cache, func(item interface{}) bool {
				conn.ReturnConnection(server, item.(*net.Conn),
					"123",
					false)
				return false
			})
			cache = list.New()
		}
		c, _, err := conn.GetConnection(server)
		if err != nil {
			logger.Error("error: ", err)
		}
		cache.PushBack(c)
	}
	for i := 0; i < 100; i++ {
		_, a, _ := conn.GetConnection(server)
		fmt.Println(a)
	}
}

func serverHandler(conn net.Conn) {
	pip := &gpip.Pip{
		Conn: conn,
	}
	defer pip.Close()
	for {
		err := pip.Receive(&Header{}, func(_header interface{}, bodyReader io.Reader, bodyLength int64) error {
			header := _header.(*Header)
			bs, _ := json.Marshal(header)
			logger.Info("server got message:", string(bs))
			return pip.Send(&Header{
				Code:      200,
				Attribute: map[string]string{"Result": "success"},
			}, nil, 0)
		})
		if err != nil {
			logger.Error("error receive data:", err)
			break
		}
	}
}
