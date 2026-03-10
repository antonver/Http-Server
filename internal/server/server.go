package server

import (
	"errors"
	"fmt"
	"net"
	"sync/atomic"

	"github.com/antonver/Http-Server/internal/request"
	"github.com/antonver/Http-Server/internal/response"
)

type Server struct {
	isRunning atomic.Bool
	port      int
	listener  net.Listener
	handler   Handler
}

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w *response.Writer, req *request.Request)


func Serve(port int, h Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	server := Server{
		listener: listener,
		port:     port,
		handler:  h,
	}
	server.isRunning.Store(true)
	go server.listen()
	return &server, nil
}

func (s *Server) Close() error {
	defer s.listener.Close()
	is_changed := s.isRunning.CompareAndSwap(true, false)
	if !is_changed {
		return errors.New("Server was already closed")
	}
	return nil
}

func (s *Server) listen() {
	for s.isRunning.Load() {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Println("Error to establish connection")
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	r, err := request.RequestFromReader(conn)
	writer := response.Writer{
		W: conn,
	}
	s.handler(&writer, r)
	if err != nil {
		errorHandler(conn, err)
		return
	}
	
	fmt.Printf("Client: %s was served!!!", conn.RemoteAddr().String())
}

func errorHandler(conn net.Conn, parsingError error) {
		errorText := parsingError.Error()
		response.WriteStatusLine(conn, 400)
		defaultHeaders := response.GetDefaultHeaders(len(errorText))
		response.WriteHeaders(conn, defaultHeaders)
		conn.Write([]byte(errorText))
		fmt.Printf("Error happened: %s", errorText)
}
