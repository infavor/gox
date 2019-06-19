package ws_test

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/hetianyi/gox/logger"
	"github.com/hetianyi/gox/ws"
	"net/http"
	"testing"
	"time"
)

func init() {
	logger.Init(nil)
}

func TestHandle(t *testing.T) {
	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/", htmlHandler)

	s := &http.Server{
		Addr: ":8080",
		// ReadTimeout:    10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      0,
		MaxHeaderBytes:    1 << 20,
	}
	fmt.Println("server listening on port 8080")
	s.ListenAndServe()
}

func wsHandler(writer http.ResponseWriter, request *http.Request) {
	err := ws.Handle(writer, request, func(messageType int, p []byte) (int, []byte, error) {
		logger.Info("recv: ", string(p))
		if messageType == websocket.PingMessage {
			return websocket.PongMessage, p, nil
		}
		return messageType, []byte("Got it!"), nil
	})
	if err != nil {
		logger.Error(err)
	}
}
func htmlHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.innerHTML = message;
        output.appendChild(d);
    };
    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("ws://localhost:8080/ws");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output"></div>
</td></tr></table>
</body>
</html>
`))
}
