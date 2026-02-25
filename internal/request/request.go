package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type State int;

type Request struct {
	RequestLine RequestLine
	state State
}

const (
	StateRequestLine = iota
	StateDone
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}



func RequestFromReader(reader io.Reader) (*Request, error){
	buff := make([]byte, 0, 8)
	tmp := make([]byte, 256)
	req := &Request{state: StateRequestLine}
	for{
		n, err := reader.Read(tmp)
		if err != nil{
			return &Request{}, err
		}
		buff = append(buff, tmp[:n]...)
		consumed, err := req.parse(buff)
		if err != nil{
			return  &Request{}, err
		}
		if consumed > 0{
			buff = buff[consumed:]
		}
		if req.state == StateDone{
			return req, nil
		}

	}
	
}

func parseRequestLine(data string) (RequestLine, int, error){
	n := strings.Index(data, "\r\n")
	if n == -1 {
		return RequestLine{}, 0, nil
	}
	line := strings.Fields(data[:n])
	fmt.Println(line)
	method := line[0]
	if len(line) != 3{
		return RequestLine{}, 0, errors.New("Request line cant have more or less than 3 elements")
	}
	if strings.ToUpper(method) != method{
		return RequestLine{}, 0, errors.New("Method is not in the upper case")
	}
	endpoint := line[1]
	if endpoint[0] != '/' {
		return RequestLine{}, 0, errors.New("Path must contain one / at the start")
	}
	
	protocol := line[2]
	protocolParts := strings.Split(protocol, "/")
	if len(protocolParts) != 2 || protocolParts[0] != "HTTP" || protocolParts[1] != "1.1" {
		return RequestLine{},0, errors.New("Protocol is not equal to HTTP/1.1")
	}
	version := protocolParts[1]

	result := RequestLine{
		HttpVersion: version,
		RequestTarget: endpoint,
		Method: method,
	}

	return result, n+2, nil
}


func (r *Request) parse(data []byte) (int, error){
	switch r.state{
	case StateRequestLine:
		res, n, err := parseRequestLine(string(data))
		if err!= nil{
			return 0, err}
		if n == 0 {
			return 0, nil
		}
		r.RequestLine = res
		r.state = StateDone
		return n, nil
		
	case StateDone:
		return 0,nil
	
}
return 0, nil
	}

