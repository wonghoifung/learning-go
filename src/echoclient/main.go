package main

import(
	"fmt"
	"time"
	"net"
)

const RECV_BUF_LEN = 10240

func main() {
	conn,err := net.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	buf := make([]byte, RECV_BUF_LEN)

	for i := 0; i < 5; i++ {
		msg := fmt.Sprintf("helloworld, %03d", i)
		n, err := conn.Write([]byte(msg))
		if err != nil {
			println("write error:", err.Error())
			break
		}
		fmt.Println(msg)

		n, err = conn.Read(buf)
		if err !=nil {
			println("read error:", err.Error())
			break
		}
		fmt.Println(string(buf[0:n]))

		time.Sleep(time.Second)
	}
}

