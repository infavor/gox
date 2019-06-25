package conn

import (
	"net"
	"time"
)

var (
	// manager manages connection pools for many servers.
	poolManager map[string]*pool
	// config is a settings which stores each server's max connection size.
	config map[string]uint
	// default max connection size.
	defaultMaxConn         uint = 50
	defaultMaxConnIdleTime      = time.Minute * 5
)

func init() {
	poolManager = make(map[string]*pool)
	config = make(map[string]uint)
}

// SetDefaultMaxConnSize sets default max connection size.
func SetDefaultMaxConnSize(maxConn uint) {
	defaultMaxConn = maxConn
}

// SetDefaultMaxConnIdleTime sets default max connection expire time.
func SetDefaultMaxConnIdleTime(connMaxIdleTime time.Duration) {
	defaultMaxConnIdleTime = connMaxIdleTime
}

// InitServerSettings initializes settings of a server.
// It is better way that initialize a server settings before getting connections from it's pool.
func InitServerSettings(server *Server, maxConn uint, connMaxIdleTime time.Duration) {
	s := server.GetConnectionString()
	config[s] = maxConn
	if poolManager[s] == nil {
		poolManager[s] = NewPool(maxConn, &ConnectionFactory{
			Server:          server,
			ConnMaxIdleTime: connMaxIdleTime,
			DialogTimeout:   time.Second * 15,
		})
	}
}

// getServerConnPool gets server's connection pool.
func getServerConnPool(server *Server) *pool {
	p := poolManager[server.GetConnectionString()]
	if p == nil {
		InitServerSettings(server, defaultMaxConn, defaultMaxConnIdleTime)
		p = poolManager[server.GetConnectionString()]
	}
	return p
}

// GetConnection tries to get a connection from it's connection pool.
func GetConnection(server *Server) (*net.Conn, error) {
	return getServerConnPool(server).GetConnection()
}

// ReturnConnection returns connection to it's connection pool.
func ReturnConnection(server *Server, conn *net.Conn, broken bool) {
	p := getServerConnPool(server)
	if broken {
		p.ReturnBrokenConnection(conn)
	} else {
		p.ReturnConnection(conn)
	}
}
