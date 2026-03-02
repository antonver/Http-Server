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

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}


const (
	Initialized = iota
	Done
)




func RequestFromReader(reader io.Reader) (*Request, error){
	const bufferSize = 1024
	buff := make([]byte, bufferSize)
	req := &Request{state: Initialized}
	readToIndex := 0
for {
	if readToIndex == len(buff){
		new_buff := make([]byte, len(buff) * 2)
		copy(new_buff, buff) 
		buff = new_buff
	}
	n, readErr := reader.Read(buff[readToIndex:])
    readToIndex += n
	consumed, parseErr := req.parse(buff[:readToIndex])
	if parseErr != nil{
		return nil, parseErr
	}
	copy(buff, buff[consumed:readToIndex])
	readToIndex -= consumed
	if req.state == Done{
		return req, nil
	}
	if readErr != nil{
		if readErr == io.EOF{
			return nil, io.ErrUnexpectedEOF 
		}
		return nil, readErr
	}
}
	
}

func parseRequestLine(data string) (RequestLine, int, error){
	n := strings.Index(data, "\r\n")
	if n == -1{
		return RequestLine{}, 0, nil
	}
	line := strings.Fields(data[:n])
	if len(line) != 3{
		return RequestLine{}, 0, errors.New("Request line must consist of 3 elements")
	}
	fmt.Println(line)
	method := line[0]
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
	case Initialized:
		res, n, err := parseRequestLine(string(data))
		if err!= nil{
			return 0, err}
		if n == 0 {
			return 0, nil
		}
		r.RequestLine = res
		r.state = Done
		return n, nil
		
	case Done:
		return 0,errors.New("Parsing is already done")
	
}
return 0, nil
	}
