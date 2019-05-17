# package conn

package conn manages connection pool for client side application.

Usage:

```golang
import "github.com/hetianyi/gox/conn"

server := &conn.Server{
    Host: "127.0.0.1",
    Port: 8899,
}
conn.InitServerSettings(server, 100)
c, err := conn.GetConnection(server)
// ...
```