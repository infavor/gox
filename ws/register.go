package ws

import (
	"errors"
	"github.com/gorilla/websocket"
	"github.com/hetianyi/gox/logger"
	"net/http"
)

var upgrader = websocket.Upgrader{} // use default options

// Handle handles and serve an http connection as websocket.
func Handle(w http.ResponseWriter, r *http.Request, handler func(messageType int, p []byte) (int, []byte, error)) error {
	if handler == nil {
		return errors.New("nil websocket handler")
	}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	defer func() {
		logger.Debug("close websocket connection ", &c)
		c.Close()
	}()
	logger.Debug("start a new websocket connection ", &c)
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			return err
		}
		rmt, data, err := handler(mt, message)
		if err != nil {
			return err
		}
		err = c.WriteMessage(rmt, data)
		if err != nil {
			return err
		}
	}
}
