package main

import(
	"net"
	"fmt"
	"io"
)

const RECV_BUF_LEN = 10240

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:9999")
	if err != nil {
		panic("listen error:" + err.Error())
	}
	fmt.Println("starting echo server...")

	for{
		conn, err := listener.Accept() 
		if err != nil {
			panic("accept error:" + err.Error())
		}
		fmt.Println("connected:", conn.RemoteAddr())
		go sessionHandler(conn)
	}
}

func sessionHandler(conn net.Conn) {
	buf := make([]byte, RECV_BUF_LEN)
	defer conn.Close()

	for{
		n, err := conn.Read(buf);
		switch err {
		case nil:
			conn.Write( buf[0:n] )
		case io.EOF:
			fmt.Printf("eof: %s \n", err);
			return
		default:
			fmt.Printf("error: %s \n", err);
			return
		}
	}
}
