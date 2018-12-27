package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"
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

func TestGet(t *testing.T) {
	start := time.Now()
	_, err := http.Get("http://baidu.com")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(time.Since(start))
}

func TestBreakLoop(t *testing.T) {
	closed := make(chan bool)

	go func() {
		defer func() {
			log.Println("return")
		}()

		for {
			select {
			case <-closed:
				return
			}
		}

		time.Sleep(time.Second * 100)
	}()

	time.Sleep(time.Second * 3)
	select {}
}
