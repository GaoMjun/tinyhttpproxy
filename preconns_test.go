package main

import (
	"log"
	"testing"
)

func TestCreateConns(t *testing.T) {
	conns, err := PreCreateConns("baidu.com", 3)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(conns)
	for _, conn := range conns {
		conn.Close()
	}
}
