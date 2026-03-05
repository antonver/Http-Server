package main

import (
	"fmt"
	"net"

	"github.com/antonver/Http-Server/internal/request"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on port :42069")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		fmt.Println("accepted connection")

		line, err := request.RequestFromReader(conn)
		if err != nil{
			fmt.Println("Error reading line:", err)
		}
		
		fmt.Printf(`Request line: \n
	- Method: %s\n
	- Target: %s\n
	- Version: %s\n
	Headers:\n
	`,
		line.RequestLine.Method,
		line.RequestLine.RequestTarget,
		line.RequestLine.HttpVersion,
	
)		
		for key, val := range line.Headers{
			fmt.Printf("%s: %s", key, val)
		}
		
		conn.Close()
		fmt.Println("connection closed")

	}
}
