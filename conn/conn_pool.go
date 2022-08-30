// Copyright (C) 2019 tisnyo <tisnyo@gmail.com>.
//
// conn is connection pool designed for limiting the size of connections to server.
// license that can be found in the LICENSE file.
package conn

import (
	"container/list"
	"errors"
	"github.com/infavor/gox"
	"github.com/infavor/gox/convert"
	"github.com/infavor/gox/logger"
	"github.com/infavor/gox/timer"
	"net"
	"sync"
	"time"
)

// pool is a connection pool.
type pool struct {
	maxSize               uint
	currentSize           uint
	connList              *list.List
	listLock              *sync.Mutex
	registeredConnMap     map[*net.Conn]time.Time   // registeredConnMap stores the connection's max idle deadline
	registeredConnAttrMap map[*net.Conn]interface{} // registeredConnAttrMap stores attributes with this connection
	connFactory           *ConnectionFactory
}

// ConnectionFactory is a factory which creates connection for specific server。
type ConnectionFactory struct {
	Server          Server
	ConnMaxIdleTime time.Duration
	DialogTimeout   time.Duration
}

// Server defines a server connection info.
type Server interface {
	ConnectionString() string
	GetHost() string
	GetPort() uint16
}

type PlainServer struct {
	Host string `json:"host"`
	Port uint16 `json:"port"`
}

// GetConnectionString gets the server's connection string.
func (ps *PlainServer) ConnectionString() string {
	return ps.Host + ":" + convert.Uint16ToStr(ps.Port)
}

// GetHost returns server's host.
func (ps *PlainServer) GetHost() string {
	return ps.Host
}

// GetPort returns server's port.
func (ps *PlainServer) GetPort() uint16 {
	return ps.Port
}

// NewPool creates a connection pool.
func NewPool(size uint, connFactory *ConnectionFactory) *pool {
	if size <= 0 {
		panic(errors.New("size must be a positive number"))
	}
	if connFactory == nil {
		panic(errors.New("connFactory can not be nil"))
	}

	p := &pool{
		maxSize:               size,
		currentSize:           0,
		connList:              list.New(),
		listLock:              new(sync.Mutex),
		connFactory:           connFactory,
		registeredConnMap:     make(map[*net.Conn]time.Time),
		registeredConnAttrMap: make(map[*net.Conn]interface{}),
	}
	go p.expireConnections()
	return p
}

// createConn creates a connection.
func (fac *ConnectionFactory) createConn() (*net.Conn, error) {
	d := net.Dialer{Timeout: fac.DialogTimeout}
	conn, err := d.Dial("tcp", fac.Server.ConnectionString())
	if err != nil {
		return nil, err
	}
	logger.Debug("create new connection to ", fac.Server.ConnectionString())
	return &conn, nil
}

// GetConnection gets a connection from pool,
func (p *pool) GetConnection() (*net.Conn, interface{}, error) {
	p.listLock.Lock()
	defer p.listLock.Unlock()
	if p.connList.Len() > 0 {
		co := p.connList.Remove(p.connList.Front()).(*net.Conn)
		return co, p.registeredConnAttrMap[co], nil
	}
	if p.currentSize >= p.maxSize {
		return nil, nil, errors.New("connection pool is full")
	}
	c, err := p.connFactory.createConn()
	if err != nil {
		return nil, nil, err
	}
	p.currentSize++
	p.registeredConnMap[c] = time.Now().Add(p.connFactory.ConnMaxIdleTime)
	return c, nil, nil
}

// ReturnConnection returns a healthy connection
func (p *pool) ReturnConnection(c *net.Conn, attr interface{}) {
	p.listLock.Lock()
	defer p.listLock.Unlock()
	if c != nil {
		p.registeredConnMap[c] = time.Now().Add(p.connFactory.ConnMaxIdleTime)
		p.registeredConnAttrMap[c] = attr
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
	}
}

// ReturnBrokenConnection returns a broken connection.
func (p *pool) expireConnections() {
	timer.Start(0, p.connFactory.ConnMaxIdleTime, 0, func(t *timer.Timer) {
		now := time.Now()
		var next *list.Element
		p.listLock.Lock()
		gox.Try(func() {
			for e := p.connList.Front(); e != nil; e = next {
				c := e.Value.(*net.Conn)
				next = e.Next()
				if p.registeredConnMap[c].Unix() <= now.Unix() {
					p.connList.Remove(e)
					delete(p.registeredConnMap, c)
					delete(p.registeredConnAttrMap, c)
					logger.Debug("expire connection:", &c)
					p.currentSize--
					if c != nil {
						(*c).Close()
					}
				}
			}
		}, func(e interface{}) {})
		p.listLock.Unlock()
	})
}
