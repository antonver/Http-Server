package response

import (
	"fmt"
	"io"
	"slices"
	"github.com/antonver/Http-Server/internal/headers"
	"errors"
	"strings"
)

type StatusCode int

type writerState int

const (
	StateWriteStatusLine writerState = iota
	StateWriteHeaders
	StateWriteBody
)

type Writer struct {
	W io.Writer
	state writerState
}

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	switch statusCode {
	case StatusOK:
		_, err := w.Write([]byte("HTTP/1.1 200 OK\r\n"))
		if err != nil {
			return err
		}
		return nil

	case StatusBadRequest:
		_, err := w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
		if err != nil {
			return err
		}
		return nil

	case StatusInternalServerError:
		_, err := w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
		if err != nil {
			return err
		}
		return nil
	}
	_, err := w.Write([]byte(fmt.Sprintf("HTTP/1.1 %d \r\n", statusCode)))
	if err != nil {
		return err
	}
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()
	headers["content-length"] = fmt.Sprintf("%d\r\n", contentLen)
	headers["connection"] = "close\r\n"
	headers["content-type"] = "text/plain\r\n"
	return headers
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, val := range headers {
		_, err := w.Write([]byte(fmt.Sprintf("%s: %s", key, val)))
		if err != nil {
			return err
		}
	}
	w.Write([]byte("\r\n"))
	return nil
}

func (w *Writer) WriteStatusLine(statusCode int, reasonPhrase string) error {
	 if w.state != 0{
	return errors.New("Response like must be first!!!")
 }
 w.state = StateWriteStatusLine
	validStatusCodes := []int{1,2,3,4,5}
 if !slices.Contains(validStatusCodes,(statusCode / 100)) {
	return errors.New("Status code is not valid")
 }
 w.W.Write([]byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, reasonPhrase)))
 return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.state != StateWriteStatusLine{
		return errors.New("Headers must be second")
	}
	w.state = StateWriteHeaders
for key, val := range headers {
		_, err := w.W.Write([]byte(fmt.Sprintf("%s: %s\r\n", strings.TrimSpace(key), val)))
		if err != nil {
			return err
		}
	}
	w.W.Write([]byte("\r\n"))
	return nil
	
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.state != StateWriteHeaders{
		return 0, errors.New("Body must be third(last one)")
	}
	w.state = StateWriteBody
	n, err := w.W.Write(p)
	return n, err
}


func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
    chunkHeader := []byte(fmt.Sprintf("%x\r\n", len(p)))
    
    n1, err := w.W.Write(chunkHeader)
    if err != nil { return n1, err }
    
    n2, err := w.W.Write(p)
    if err != nil { return n1 + n2, err }

    n3, err := w.W.Write([]byte("\r\n"))
    return n1 + n2 + n3, err
}

func (w *Writer) WriteChunkedBodyDone() (int, error){
	info := []byte(fmt.Sprintf("%x\r\n", 0))
	n, err := w.W.Write(info)
	return n, err
}


func (w *Writer) WriteTrailers(h headers.Headers) error{

for key, val := range h {
		_, err := w.W.Write([]byte(fmt.Sprintf("%s: %s\r\n", strings.TrimSpace(key), val)))
		if err != nil {
			return err
		}
	}
	w.W.Write([]byte("\r\n"))
	return nil

}