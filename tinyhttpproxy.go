package main

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
)

func init() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
}

func main() {
	var (
		err error
		l   net.Listener
	)
	defer func() {
		if err != nil {
			log.Println(err)
		}
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

		go handleConn(c)
	}
}

func handleConn(c net.Conn) {
	var (
		err    error
		reader = bufio.NewReader(c)

		request = &Request{}

		address string

		server net.Conn
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

	for {
		line, isPrefix, err := reader.ReadLine()
		if err != nil {
			return
		}
		if isPrefix {
			err = errors.New("line is too long")
			return
		}

		request.RawLines = append(request.RawLines, string(line))

		if len(line) == 0 {
			break
		}
	}

	err = request.Parse()
	if err != nil {
		return
	}

	address = request.HttpRequest.Host
	if strings.Index(address, ":") == -1 {
		address = address + ":80"
	}
	log.Println(address)

	server, err = net.Dial("tcp", address)
	if err != nil {
		return
	}

	if request.HttpRequest.Method == "CONNECT" {
		fmt.Fprint(c, "HTTP/1.1 200 Connection established\r\n\r\n")
	} else {
		server.Write(request.Bytes())
	}

	go io.Copy(server, c)
	io.Copy(c, server)
}

type Request struct {
	RawLines    []string
	HttpRequest *http.Request

	rawString string
	rawBytes  []byte
}

func (self *Request) Dump() (s string) {
	if len(self.rawString) > 0 {
		s = self.rawString
		return
	}

	s = strings.Join(self.RawLines, "\r\n")
	s = strings.Join([]string{s, "\r\n"}, "")

	self.rawString = s
	return
}

func (self *Request) DumpHex() (s string) {
	s = hex.Dump([]byte(self.Dump()))
	return
}

func (self *Request) Bytes() (bs []byte) {
	if len(self.rawBytes) > 0 {
		bs = self.rawBytes
		return
	}

	bs = []byte(self.Dump())

	self.rawBytes = bs
	return
}

func (self *Request) Parse() (err error) {
	r, err := http.ReadRequest(bufio.NewReader(strings.NewReader(self.Dump())))
	if err != nil {
		return
	}

	self.HttpRequest = r
	return
}
