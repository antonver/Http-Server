package request

import (
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/antonver/Http-Server/internal/headers"
)

type State int;

type Request struct {
	RequestLine RequestLine
	Headers headers.Headers
	state State
	Body []byte
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}


const (
	requestStateParsingRequestLine = iota
	requestStateParsingHeaders
	requestStateParsingBody
	requestStateParsingDone
)




func RequestFromReader(reader io.Reader) (*Request, error){
	const bufferSize = 1024
	buff := make([]byte, bufferSize)
	req := &Request{state: requestStateParsingRequestLine,
	Headers: headers.NewHeaders(),
	}
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
		return req, parseErr
	}
	copy(buff, buff[consumed:readToIndex])
	readToIndex -= consumed
	if req.state == requestStateParsingDone{
		if len(req.Headers) == 0{
			return req, errors.New("We must have at least one header")
		}
		return req, nil
	}
	if readErr != nil{
		if readErr == io.EOF{
			return nil, io.ErrUnexpectedEOF 
		}
		return req, readErr
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
	totalBytesParsed := 0
	for r.state != requestStateParsingDone{
		n, err := r.parseSingle(data[totalBytesParsed:])
		totalBytesParsed += n
		if err != nil{
			return 0, err
		}
		if n == 0{
			return totalBytesParsed, nil
		}
	}
	return totalBytesParsed, nil

}




func (r *Request) parseSingle(data []byte)(int, error){
	switch r.state{
	case requestStateParsingRequestLine:
		res, n, err := parseRequestLine(string(data))
		if err!= nil{
			return 0, err}
		if n == 0 {
			return 0, nil
		}
		r.RequestLine = res
		r.state = requestStateParsingHeaders
		return n, nil
	case requestStateParsingHeaders:
    		n, done, err := r.Headers.Parse(data)
			if done{
				r.state = requestStateParsingBody
				return n, nil
			}
			if err != nil{
				return 0, err
			}
			if n == 0 {
				return n, nil
        }
			return n, nil
			
	case requestStateParsingBody:
		content_length := r.Headers.Get("Content-Length") 
		n, err := r.parseBody(data, content_length)
		if err!= nil{
			return 0, err
		}
		if n == 0{
			return 0, nil
		}
		return n, nil
	case requestStateParsingDone:
		return 0, errors.New("Parsing is already done")

}
			return 0, nil
	}



func (r *Request) parseBody(data []byte, content_length string) (int, error){
	sData := string(data)
	content_length_int, ok := strconv.Atoi(content_length)
	if content_length == ""{
			r.state = requestStateParsingDone
			return 0, nil
		}
	if ok != nil{
		return 0, errors.New("Content_length header is not a number")
	}
	if content_length_int > len(data){
		return 0, nil
	}
	if content_length_int < len(data){
		return 0, errors.New("Content length cann't be bigger than length precised in content-length header")
	} 
	if content_length_int != len(sData){
		return 0, errors.New("Body length differ from length precised in content-length header")
	}
	r.state = requestStateParsingDone
	r.Body = data
	return len(data), nil
}
