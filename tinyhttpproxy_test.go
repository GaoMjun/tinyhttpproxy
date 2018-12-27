package main

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	s := "GET http://baidu.com/ HTTP/1.1\r\nHost: baidu.com\r\nUser-Agent: curl/7.54.0\r\nAccept: */*\r\nProxy-Connection: Keep-Alive\r\n\r\n"

	r, err := http.ReadRequest(bufio.NewReader(strings.NewReader(s)))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(r.Host)
}
