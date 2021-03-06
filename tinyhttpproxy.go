package main

import (
	"fmt"
	"log"
	"net"
	_ "net/http/pprof"
	"strings"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	var (
		err      error
		l        net.Listener
		connPool = NewConnPool()
	)
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	// go func() {
	// 	err = http.ListenAndServe(":6060", nil)
	// 	log.Panicln(err)
	// }()

	l, err = net.Listen("tcp", ":8123")
	if err != nil {
		return
	}

	for {
		c, err := l.Accept()
		if err != nil {
			return
		}

		conn := NewConnWithTimeout(c)
		go handleConn(conn, connPool)
	}
}

func handleConn(c *ConnWithTimeout, connPool *ConnPool) {
	var (
		err     error
		request *Request
		address string
		server  net.Conn
	)
	defer func() {
		c.Close()
		if server != nil {
			server.Close()
		}
		if err != nil {
			log.Println(err)
		}
	}()

	request, err = NewRequest(c)
	if err != nil {
		return
	}

	address = request.HttpRequest.Host
	if strings.Index(address, ":") == -1 {
		address = address + ":80"
	}
	log.Println(address)

	server = connPool.GetConn(address)
	if server == nil {
		conns, err := PreCreateConns(address, 1)
		if err != nil {
			return
		}

		if len(conns) > 1 {
			connPool.AddConns(address, conns[1:])
		}

		server = conns[0]
	}

	if request.HttpRequest.Method == "CONNECT" {
		fmt.Fprint(c, "HTTP/1.1 200 Connection established\r\n\r\n")
	} else {
		server.Write(request.Bytes())
	}

	Pipe(c, server)
}
