# package gpip

### Description

gpip is a simple tool for frame message transfer based on tcp.
It can transfer text message and byte content such as file with very easy way.

### Usage

***Server side:***

```golang
func main() {
	listener, e1 := net.Listen("tcp", ":8080")
	if e1 != nil {
		panic(e1)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go serverHandler(conn)
	}
}


func serverHandler(conn net.Conn) {
	pip := &gpip.Pip{
		Conn: conn,
	}
	defer pip.Close()
	for {
		err := pip.Receive(&common.Header{}, func(_header interface{}, bodyReader io.Reader, bodyLength int64) error {
			header := _header.(*common.Header)
			bs, _ := json.Marshal(header)
			log.Info("server got message:", string(bs))
			return pip.Send(&common.Header{
				Code: 200,
				Attribute: map[string]string{"Result":"success"},
			}, nil, 0)
		})
		if err != nil {
			log.Error("error receive data:", err)
			break
		}
	}
}

```


Client side:

```golang
func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:"+strconv.Itoa(common.ServerPort))
	if err != nil {
		log.Fatal("error start client:", err)
	}
	pip := &gpip.Pip{
		Conn: conn,
	}
	defer pip.Close()
	err := pip.Send(&common.Meta{
		Code: 1,
		Attribute: nil,
	}, nil, 0)
	if err != nil {
		log.Fatal("error send data:", err)
		break
	}
	err = pip.Receive(&common.Header{}, func(_header interface{}, bodyReader io.Reader, bodyLength int64) error {
		/*header := _header.(*common.Header)
		bs, _ := json.Marshal(header)
		log.Info("client got message:", string(bs))*/
		return nil
	})
	if err != nil {
		log.Error("error:", err)
	}
}
```