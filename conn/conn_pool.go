// Copyright (C) 2019 tisnyo <tisnyo@gmail.com>.
//
// conn is connection pool designed for limiting the size of connections to server.
// license that can be found in the LICENSE file.
package conn

import (
	"container/list"
	"errors"
	"net"
	"strconv"
	"sync"
	"time"
)

// pool is a connection pool.
type pool struct {
	maxSize     int
	currentSize int
	connList    *list.List
	listLock    *sync.Mutex
	connFactory *ConnectionFactory
}

// ConnectionFactory is a factory which creates connection for specific server。
type ConnectionFactory struct {
	server        *Server
	DialogTimeout time.Duration
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
	return &pool{
		maxSize:     size,
		currentSize: 0,
		connList:    list.New(),
		listLock:    new(sync.Mutex),
		connFactory: connFactory,
	}
}

// createConn creates a connection.
func (fac *ConnectionFactory) createConn() (net.Conn, error) {
	d := net.Dialer{Timeout: fac.DialogTimeout}
	conn, err := d.Dial("tcp", fac.server.Host+":"+strconv.Itoa(fac.server.Port))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// GetConnection gets a connection from pool,
func (p *pool) GetConnection() (net.Conn, error) {
	p.listLock.Lock()
	defer p.listLock.Unlock()
	if p.connList.Len() > 0 {
		return p.connList.Remove(p.connList.Front()).(net.Conn), nil
	}
	if p.currentSize >= p.maxSize {
		return nil, errors.New("connection pool is full")
	}
	conn, err := p.connFactory.createConn()
	if err != nil {
		return nil, err
	}
	p.currentSize++
	return conn, nil
}

// ReturnConnection returns a healthy connection
func (p *pool) ReturnConnection(conn net.Conn) {
	p.listLock.Lock()
	defer p.listLock.Unlock()
	if conn != nil {
		p.connList.PushBack(conn)
	}
}

// ReturnBrokenConnection returns a broken connection.
func (p *pool) ReturnBrokenConnection(conn net.Conn) {
	if conn != nil {
		conn.Close()
		conn = nil
	}
}

// GetConnectionString gets the server's connection string.
func (s *Server) GetConnectionString() string {
	return s.Host + ":" + strconv.Itoa(s.Port)
}
