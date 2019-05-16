package conn

import (
	"net"
	"time"
)

var (
	// manager manages connection pools for many servers.
	poolManager map[string]*pool
	// config is a settings which stores each server's max connection size.
	config map[string]int
	// default max connection size.
	defaultMaxConn = 50
)

// SetDefaultMaxConnSize sets default max connection size.
func SetDefaultMaxConnSize(maxConn int) {
	defaultMaxConn = maxConn
}

// InitSettings initializes settings of a server.
// It is better way that initialize a server settings before getting connections from it's pool.
func InitServerSettings(server *Server, maxConn int) {
	s := server.GetConnectionString()
	config[s] = maxConn
	if poolManager[s] == nil {
		poolManager[s] = NewPool(maxConn, &ConnectionFactory{
			server:        server,
			DialogTimeout: time.Second * 15,
		})
	}
}

// getServerConnPool gets server's connection pool.
func getServerConnPool(server *Server) *pool {
	p := poolManager[server.GetConnectionString()]
	if p == nil {
		InitServerSettings(server, defaultMaxConn)
		p = poolManager[server.GetConnectionString()]
	}
	return p
}

// GetConnection tries to get a connection from it's connection pool.
func GetConnection(server *Server) (net.Conn, error) {
	return getServerConnPool(server).GetConnection()
}

// ReturnConnection returns connection to it's connection pool.
func ReturnConnection(server *Server, conn net.Conn, broken bool) {
	p := getServerConnPool(server)
	if broken {
		p.ReturnBrokenConnection(conn)
	} else {
		p.ReturnConnection(conn)
	}
}
