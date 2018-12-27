package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"strings"
	"sync"
)

func init() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
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

	go func() {
		err = http.ListenAndServe(":6060", nil)
		log.Panicln(err)
	}()

	l, err = net.Listen("tcp", ":8888")
	if err != nil {
		return
	}

	for {
		c, err := l.Accept()
		if err != nil {
			return
		}

		go handleConn(c, connPool)
	}
}

func handleConn(c net.Conn, connPool *ConnPool) {
	var (
		err     error
		request *Request
		address string
		server  net.Conn
		// closed  = make(chan bool)
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
	// log.Println(address)

	// server, err = net.Dial("tcp", address)
	// if err != nil {
	// 	return
	// }

	server = connPool.GetConn(address)
	if server == nil {
		server, err = net.Dial("tcp", address)
		if err != nil {
			return
		}
	}
	// 	server = NewConnWithTimeout(netConn)
	// } else {
	// 	log.Println("get conn from pool ", address)
	// }

	if request.HttpRequest.Method == "CONNECT" {
		fmt.Fprint(c, "HTTP/1.1 200 Connection established\r\n\r\n")
	} else {
		server.Write(request.Bytes())
	}

	// go io.Copy(server, c)
	// go func() {
	// 	_, err1 := io.Copy(c, server)
	// 	if err1, ok := err1.(net.Error); ok && err1.Timeout() {
	// 		close(closed)
	// 		return
	// 	}
	// 	err = err1
	// 	server.Conn.Close()
	// 	server = nil
	// 	close(closed)
	// }()

	// <-closed
	// if server != nil {
	// 	connPool.AddConn(address, server)
	// }

	Pipe(c, server)
}

func Pipe(src io.ReadWriteCloser, dst io.ReadWriteCloser) (int64, int64) {

	var sent, received int64
	var c = make(chan bool)
	var o sync.Once

	close := func() {
		src.Close()
		dst.Close()
		close(c)
	}

	go func() {
		received, _ = io.Copy(src, dst)
		o.Do(close)
	}()

	go func() {
		sent, _ = io.Copy(dst, src)
		o.Do(close)
	}()

	<-c
	return sent, received
}
