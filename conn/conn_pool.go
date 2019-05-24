// Copyright (C) 2019 tisnyo <tisnyo@gmail.com>.
//
// conn is connection pool designed for limiting the size of connections to server.
// license that can be found in the LICENSE file.
package conn

import (
	"container/list"
	"errors"
	"github.com/hetianyi/gox/logger"
	"github.com/sirupsen/logrus"
	"net"
	"strconv"
	"sync"
	"time"
)

func init() {
	logger.Init(nil)
}

// pool is a connection pool.
type pool struct {
	maxSize     int
	currentSize int
	connList    *list.List
	listLock    *sync.Mutex
	// registeredConnMap stores the connection's max idle deadline
	registeredConnMap map[*net.Conn]time.Time
	connFactory       *ConnectionFactory
}

// ConnectionFactory is a factory which creates connection for specific server。
type ConnectionFactory struct {
	Server          *Server
	ConnMaxIdleTime time.Duration
	DialogTimeout   time.Duration
}

// Server defines a server connection info.
type Server struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

// NewPool creates a connection pool.
func NewPool(size int, connFactory *ConnectionFactory) *pool {
	if size <= 0 {
		panic(errors.New("size must be a positive number"))
	}
	if connFactory == nil {
		panic(errors.New("connFactory can not be nil"))
	}

	p := &pool{
		maxSize:           size,
		currentSize:       0,
		connList:          list.New(),
		listLock:          new(sync.Mutex),
		connFactory:       connFactory,
		registeredConnMap: make(map[*net.Conn]time.Time),
	}
	go p.expireConnections()
	return p
}

// createConn creates a connection.
func (fac *ConnectionFactory) createConn() (*net.Conn, error) {
	d := net.Dialer{Timeout: fac.DialogTimeout}
	conn, err := d.Dial("tcp", fac.Server.Host+":"+strconv.Itoa(fac.Server.Port))
	if err != nil {
		return nil, err
	}
	logrus.Debug("create new connection:", &conn)
	return &conn, nil
}

// GetConnection gets a connection from pool,
func (p *pool) GetConnection() (*net.Conn, error) {
	p.listLock.Lock()
	defer p.listLock.Unlock()
	if p.connList.Len() > 0 {
		return p.connList.Remove(p.connList.Front()).(*net.Conn), nil
	}
	if p.currentSize >= p.maxSize {
		return nil, errors.New("connection pool is full")
	}
	c, err := p.connFactory.createConn()
	if err != nil {
		return nil, err
	}
	p.currentSize++
	p.registeredConnMap[c] = time.Now().Add(p.connFactory.ConnMaxIdleTime)
	return c, nil
}

// ReturnConnection returns a healthy connection
func (p *pool) ReturnConnection(c *net.Conn) {
	p.listLock.Lock()
	defer func() {
		p.currentSize--
		p.listLock.Unlock()
	}()
	if c != nil {
		p.registeredConnMap[c] = time.Now().Add(p.connFactory.ConnMaxIdleTime)
		p.connList.PushBack(c)
	}
}

// ReturnBrokenConnection returns a broken connection.
func (p *pool) ReturnBrokenConnection(conn *net.Conn) {
	p.listLock.Lock()
	defer func() {
		p.currentSize--
		p.listLock.Unlock()
	}()
	if conn != nil {
		(*conn).Close()
		conn = nil
	}
}

// GetConnectionString gets the server's connection string.
func (s *Server) GetConnectionString() string {
	return s.Host + ":" + strconv.Itoa(s.Port)
}

// ReturnBrokenConnection returns a broken connection.
func (p *pool) expireConnections() {
	t := time.NewTicker(p.connFactory.ConnMaxIdleTime)
	for {
		now := time.Now()
		var next *list.Element
		for e := p.connList.Front(); e != nil; e = next {
			c := e.Value.(*net.Conn)
			next = e.Next()
			if p.registeredConnMap[c].Unix() <= now.Unix() {
				p.connList.Remove(e)
				logrus.Debug("expire connection:", &c)
				p.ReturnBrokenConnection(c)
			}
		}
		<-t.C
	}
}
