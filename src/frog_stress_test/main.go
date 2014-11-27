package main

import(
	"fmt"
	"net"
	"sync"
	"frog_codec"
	"frog_parser"
)

const RECV_BUF_LEN = 10240

func main() {
	var wg sync.WaitGroup
	for	i:=0;i<10;i++{
		wg.Add(1)
		go connSend(&wg)
	}
	wg.Wait()
	fmt.Println("all done, exit")
}

func connSend(wg *sync.WaitGroup){
	conn,err := net.Dial("tcp", "127.0.0.1:9878")
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	buf := make([]byte, RECV_BUF_LEN)

	var op frog_codec.Codec
	op.Begin(123)
	op.WriteInt32(987)
	op.WriteInt16(654)
	op.WriteCString("helloworld!!!")
	op.WriteInt16(321)
	op.WriteInt32(123456)
	op.End()

	c:=0
	for{
		n,err := conn.Write(op.Data)
		if err != nil{
			fmt.Println("write error:",n,err.Error())
			return
		}
		c+=n
		if c==len(op.Data){
			break
		}
	}

	quit := false
	var parser frog_parser.PackParser
	parser.Init()
	for !quit{
		n,err := conn.Read(buf)
		if err!=nil{
			fmt.Println("read error:",err.Error())
			return
		}

		bb := make([]byte, n)
		copy(bb[0:n],buf[0:n])

		parser.Parse(bb,onMessage,&quit)
	}
	wg.Done()
}

func onMessage(inpack *frog_codec.Codec, quit *bool)int{
	var cmd int16 = inpack.Cmd()
	var bodylen int16 = inpack.BodyLen()
	var i int32 = inpack.ReadInt32()
	var s string = inpack.ReadCString()
	fmt.Println("receive:",cmd,bodylen,i,s)
	*quit=true
	return 0
}



