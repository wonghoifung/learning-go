package frog_parser

import(
	"frog_codec"
)

const(
	c_header = iota
	c_body
	c_done
	c_error
)

type callback func(decoder *frog_codec.Codec, quit *bool)int

type PackParser struct{
	state int16
	packpos int16
	bodylen int16
	decoder frog_codec.Codec
}

func (pp *PackParser) Parse(data []byte, cb callback, quit *bool)int{
	ndx := 0
	for ndx<len(data) && pp.state!=c_error{
		switch pp.state{
		case c_header:
			if !pp.read_header(data,&ndx){
				break
			}
			if pp.parse_header()!=0{
				pp.state = c_error
				break
			}
			pp.state = c_body

		case c_body:
			if pp.parse_body(data,&ndx){
				pp.state = c_done
			}

		default:
			return -1
		}
		if pp.state==c_error{
			pp.Init()
		}
		if pp.state==c_done{
			cb(&pp.decoder,quit)
			pp.Init()
		}
	}
	return 0
}

// set complete callback first
// then call Init()
func (pp *PackParser) Init(){
	pp.state = c_header;
	pp.packpos = 0;
	pp.bodylen = 0;
	pp.decoder.Data = make([]byte,frog_codec.C_BUFFER_SIZE)
	pp.decoder.Readpos = frog_codec.C_HEADER_SIZE
}

func (pp *PackParser) read_header(data []byte, ndx *int)bool{
	//if *ndx==0{
	//	pp.packpos = 0
	//}
	for pp.packpos<frog_codec.C_HEADER_SIZE && *ndx<len(data){
		copy(pp.decoder.Data[pp.packpos:pp.packpos+1],data[*ndx:*ndx+1])
		pp.packpos++
		*ndx++
	}
	if pp.packpos<frog_codec.C_HEADER_SIZE{
		return false
	}
	return true
}

func (pp *PackParser) parse_header()int{
	if pp.decoder.Data[0]!='G' || pp.decoder.Data[1]!='P'{
		return -1
	}
	cmd := pp.decoder.Cmd()
	if cmd<0 || cmd>=32000{
		return -2
	}
	pp.bodylen = pp.decoder.BodyLen()
	if pp.bodylen<(int16(frog_codec.C_BUFFER_SIZE)-frog_codec.C_HEADER_SIZE){
		return 0
	}
	return -4
}

func (pp *PackParser) parse_body(data []byte, ndx *int)bool{
	needlen := (pp.bodylen + frog_codec.C_HEADER_SIZE) - pp.packpos
	if needlen<=0{
		return true
	}

	buflen := len(data) - *ndx
	if buflen<=0{
		return false
	}

	var cplen int16
	if int16(buflen)<needlen{
		cplen=int16(buflen)
	}else{
		cplen=needlen
	}

	dstart := pp.decoder.Datalen+frog_codec.C_HEADER_SIZE
	copy(pp.decoder.Data[dstart:dstart+cplen],data[int16(*ndx):int16(*ndx)+cplen])
	pp.decoder.Datalen += cplen

	pp.packpos += cplen
	*ndx += int(cplen)

	if pp.packpos<(pp.bodylen+frog_codec.C_HEADER_SIZE){
		return false
	}
	return true
}

