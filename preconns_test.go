package main

import (
	"log"
	"testing"
	"time"
)

func TestCreateConns(t *testing.T) {
	start := time.Now()
	conns, err := PreCreateConns("baidu.com:80", 4)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(time.Since(start))
	for _, conn := range conns {
		conn.Close()
	}
}
