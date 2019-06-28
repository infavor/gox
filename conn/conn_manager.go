package conn

import (
	"net"
	"time"
)

var (
	poolManager            map[string]*pool      // manager manages connection pools for many servers.
	config                 map[string]uint       // config is a settings which stores each server's max connection size.
	defaultMaxConn         uint             = 50 // default max connection size.
	defaultMaxConnIdleTime                  = time.Minute * 5
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
func InitServerSettings(server Server, maxConn uint, connMaxIdleTime time.Duration) {
	s := server.ConnectionString()
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
func getServerConnPool(server Server) *pool {
	p := poolManager[server.ConnectionString()]
	if p == nil {
		InitServerSettings(server, defaultMaxConn, defaultMaxConnIdleTime)
		p = poolManager[server.ConnectionString()]
	}
	return p
}

// GetConnection tries to get a connection from it's connection pool.
func GetConnection(server Server) (*net.Conn, interface{}, error) {
	return getServerConnPool(server).GetConnection()
}

// ReturnConnection returns connection to it's connection pool.
func ReturnConnection(server Server, conn *net.Conn, attr interface{}, broken bool) {
	p := getServerConnPool(server)
	if broken {
		p.ReturnBrokenConnection(conn)
	} else {
		p.ReturnConnection(conn, attr)
	}
}
