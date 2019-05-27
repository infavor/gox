package ws

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var upgrader = websocket.Upgrader{} // use default options

// Handle handles and serve an http connection as websocket.
func Handle(w http.ResponseWriter, r *http.Request, handler func(messageType int, p []byte) ([]byte, error)) {
	if handler == nil {
		log.Error("websocket handler cannot be nil")
		return
	}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("websocket upgrade error:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Error("websocket upgrade error:", err)
			break
		}
		log.Debug("recv: %s", message)

		data, err := handler(mt, message)
		if err != nil {
			log.Error("error handle websocket message:", err)
			break
		}

		err = c.WriteMessage(mt, data)
		if err != nil {
			log.Error("error write websocket message:", err)
			break
		}
	}
}
