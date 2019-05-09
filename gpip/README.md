# package gpip

### Description

gpip is a simple tool for frame message transfer based on tcp.
It can transfer text message and byte content such as file with very easy way.

### Usage

***Server side:***

```golang
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


func serverHandler(conn net.Conn) {
	pip := &gpip.Pip{
		Conn: conn,
	}
	fmt.Println("server is listening")
	for {
		frame, err := pip.Receive()
		if err != nil {
			panic(err)
		}
		u, err := frame.GetMeta(reflect.TypeOf(&User1{}))
		if err != nil {
			panic(err)
		}
		d, _ := gpip.Serialize(u)
		fmt.Println("server收到消息:", string(d))
		resp := &gpip.PipFrame{
			Meta: &User1{Name: "zhangsan"},
		}
		pip.Send(resp)
	}
}

```


Client side:

```golang
conn, err := net.Dial("tcp", "127.0.0.1:8080")
if err != nil {
    panic(err)
}
pip := &gpip.Pip{
    Conn: conn,
}
for {
    u := &User{Name: "lisi"}
    frame := &gpip.PipFrame{
        Meta: u,
    }
    if err := pip.Send(frame); err != nil {
        panic(err)
    }
    resp, err := pip.Receive()
    if err != nil {
        panic(err)
    }
    u1, err := resp.GetMeta(reflect.TypeOf(&User{}))
    if err != nil {
        panic(err)
    }
    gpip.Serialize(u1)
}
```