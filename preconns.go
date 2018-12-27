package main

import (
	"net"
	"time"
)

func PreCreateConns(addr string, n int) (conns []net.Conn, err error) {
	dialer := &net.Dialer{Timeout: time.Second * 10}

	connCh := make(chan net.Conn)
	countCh := make(chan bool)
	count := 0

	for i := 0; i < n; i++ {
		go func() {
			var (
				err1 error
				conn net.Conn
			)

			conn, err1 = dialer.Dial("tcp", addr)
			if err1 == nil {
				connCh <- conn
			} else {
				err = err1
			}

			countCh <- true
		}()
	}

	for {
		select {
		case conn := <-connCh:
			conns = append(conns, conn)
		case <-countCh:
			count++
			if count >= n {
				goto END
			}
		}
	}

END:
	if len(conns) > 0 {
		err = nil
	}

	return
}
