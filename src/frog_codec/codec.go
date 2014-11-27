package frog_codec

import(
	"bytes"
	"encoding/binary"
)

const(
	C_DEFAULT_VERSION = 1
	C_DEFAULT_SUBVERSION = 1
	C_HEADER_SIZE int16 = 9
	C_BUFFER_SIZE int32 = 1024 * 12
)

type Codec struct{
	Datalen int16
	Readpos int16
	Data []byte
}

func (op *Codec) WriteInt16(val int16){
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, int16(val))
	copy(op.Data[op.Datalen+C_HEADER_SIZE:op.Datalen+C_HEADER_SIZE+2],buf.Bytes()[0:2])
	op.Datalen += 2
}

func (op *Codec) WriteInt32(val int32){
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, int32(val))
	copy(op.Data[op.Datalen+C_HEADER_SIZE:op.Datalen+C_HEADER_SIZE+4],buf.Bytes()[0:4])
	op.Datalen += 4
}

func (op *Codec) WriteInt64(val int64){
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, int64(val))
	copy(op.Data[op.Datalen+C_HEADER_SIZE:op.Datalen+C_HEADER_SIZE+8],buf.Bytes()[0:8])
	op.Datalen += 8
}

func (op *Codec) WriteCString(val string){
	var mybytes []byte = make([]byte,len(val)+1)
	copy(mybytes[0:len(val)],[]byte(val))
	var zero [1]byte = [1]byte{0}
	copy(mybytes[len(val):len(val)+1],zero[0:1])

	var strlen int16 = int16(len(mybytes))
	op.WriteInt32(int32(strlen))

	copy(op.Data[op.Datalen+C_HEADER_SIZE:op.Datalen+C_HEADER_SIZE+strlen],mybytes[0:strlen])
	op.Datalen += strlen
}

func (op *Codec) ReadInt16() int16{
	var ret int16
	var mybytes []byte = make([]byte,2)
	copy(mybytes[0:2],op.Data[op.Readpos:op.Readpos+2])
	buf := bytes.NewBuffer(mybytes)
	binary.Read(buf,binary.BigEndian,&ret)
	op.Readpos += 2
	return ret
}

func (op *Codec) ReadInt32() int32{
	var ret int32
	var mybytes []byte = make([]byte,4)
	copy(mybytes[0:4],op.Data[op.Readpos:op.Readpos+4])
	buf := bytes.NewBuffer(mybytes)
	binary.Read(buf,binary.BigEndian,&ret)
	op.Readpos += 4
	return ret
}

func (op *Codec) ReadInt64() int64{
	var ret int64
	var mybytes []byte = make([]byte,8)
	copy(mybytes[0:8],op.Data[op.Readpos:op.Readpos+8])
	buf := bytes.NewBuffer(mybytes)
	binary.Read(buf,binary.BigEndian,&ret)
	op.Readpos += 8
	return ret
}

func (op *Codec) ReadCString() string{
	var strlen int32 = op.ReadInt32()
	var mybytes []byte = make([]byte,strlen)
	copy(mybytes[0:strlen],op.Data[op.Readpos:op.Readpos+int16(strlen)])
	var str []byte = mybytes[0:strlen-1]
	op.Readpos += int16(strlen)
	return string(str)
}

func (op *Codec) Cmd() int16{
	var ret int16
	var mybytes []byte = make([]byte,2)
	copy(mybytes[0:2],op.Data[2:4])
	buf := bytes.NewBuffer(mybytes)
	binary.Read(buf,binary.BigEndian,&ret)
	return ret
}

func (op *Codec) BodyLen() int16{
	var ret int16
	var mybytes []byte = make([]byte,2)
	copy(mybytes[0:2],op.Data[6:8])
	buf := bytes.NewBuffer(mybytes)
	binary.Read(buf,binary.BigEndian,&ret)
	return ret
}

func (op *Codec) Begin(cmd int16){
	op.Datalen = 0
	op.Readpos = int16(C_HEADER_SIZE)
	op.Data = make([]byte,C_BUFFER_SIZE)

	var flag1 byte = 'G'
	var flag2 byte = 'P'
	var version byte = C_DEFAULT_VERSION
	var subversion byte = C_DEFAULT_SUBVERSION
	var checkcode byte = 0 // by default

	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, byte(flag1))
	binary.Write(buf, binary.BigEndian, byte(flag2))
	binary.Write(buf, binary.BigEndian, int16(cmd))
	binary.Write(buf, binary.BigEndian, byte(version))
	binary.Write(buf, binary.BigEndian, byte(subversion))
	binary.Write(buf, binary.BigEndian, int16(op.Datalen))
	binary.Write(buf, binary.BigEndian, byte(checkcode))

	copy(op.Data[0:C_HEADER_SIZE],buf.Bytes()[0:C_HEADER_SIZE])
}

func (op *Codec) End(){
	var checkcode byte = 0 // by default

	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, int16(op.Datalen))
	binary.Write(buf, binary.BigEndian, byte(checkcode))

	copy(op.Data[6:9],buf.Bytes()[0:3])

	packsize := int16(C_HEADER_SIZE) + op.Datalen
	pack := make([]byte, packsize)
	copy(pack[0:packsize],op.Data[0:packsize])
	op.Data = pack
}

