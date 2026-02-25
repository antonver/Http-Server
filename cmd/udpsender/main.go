package main

import (
	"net"
	"fmt"
	"bufio"
	"os"
	"log"
)


func main(){
	server_socket, err := net.ResolveUDPAddr("udp", ":42069")
	if err != nil{
		return
	}
	conn, err := net.DialUDP("udp", nil, server_socket)
	if err != nil{
		return 
	}
	defer conn.Close()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println(">")
		text, err := reader.ReadString('\n')
		if err != nil{
		log.Printf(err.Error())}

		_, err2 := conn.Write([]byte(text))
		if err != nil{
		log.Println(err2.Error())
		}
	}
}