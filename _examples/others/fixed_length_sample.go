package main

import (
	"fmt"
	"github.com/smallnest/goframe"
	"io"
	"net"
	"time"
)

func main() {
	//server endpoint
	go func() {
		lis, err := net.Listen("tcp", ":1234")
		defer lis.Close()

		if err != nil {
			panic(err)
		}

		for {
			conn, err := lis.Accept()
			if err != nil {
				panic(err)
			}

			frameConn := goframe.NewFixedLengthFrameConn(12, conn)
			go func(conn goframe.FrameConn) {
				for  {
					rbuf, err := frameConn.ReadFrame()
					if err != nil {
						if err == io.EOF {
							return
						}
						panic(err)
					}

					fmt.Printf("%s\r\n", string(rbuf))
				}
			}(frameConn)
		}
	}()

	//wait for server
	time.Sleep(3 * time.Second)

	//client endpoint
	conn, err := net.Dial("tcp", "127.0.0.1:1234")
	defer conn.Close()
	if err != nil {
		panic(err)
	}
	frameConn := goframe.NewFixedLengthFrameConn(12, conn)
	err = frameConn.WriteFrame([]byte("Hello,World!"))
	if err != nil {
		panic(err)
	}

	time.Sleep(3 * time.Second)
}
