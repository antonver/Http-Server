package main

import (
	"strings"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"strconv"
	"github.com/antonver/Http-Server/internal/headers"
	"github.com/antonver/Http-Server/internal/request"
	"github.com/antonver/Http-Server/internal/response"
	"github.com/antonver/Http-Server/internal/server"
	"net/http"
	"io"
	"crypto/sha256"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w *response.Writer, req *request.Request){
	fmt.Println(req.RequestLine.RequestTarget)
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/video") {
		err := w.WriteStatusLine(200, "OK")
		if err != nil{
				fmt.Println(err.Error())
		}
		headers := headers.NewHeaders()
		headers["content-type"] = "video/mp4"
		file, err := os.ReadFile("../../assets/vim.mp4")
		if err != nil{
			log.Printf("Error: %v", err)
			return
		}
		err = w.WriteHeaders(headers)
		if err != nil{
			log.Printf("Error: %v", err)
			return
		}
		_, err = w.W.Write(file)
		if err != nil{
			log.Printf("Error: %v", err)
		}
		return
	}
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin"){
		target := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
		handlerChunked(w, target)
		return
	}
	if req.RequestLine.RequestTarget == "/yourproblem"{
		err := w.WriteStatusLine(400, "Bad Request")
		if err != nil{
		fmt.Println(err.Error())}
		headers := headers.NewHeaders()
		body := `<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`
		headers["content-length"] = strconv.Itoa(len(body))
		headers["content-type"] = "text/html"
		err = w.WriteHeaders(headers)
		if err != nil{
		fmt.Println(err.Error())
	}
		n, err := w.WriteBody([]byte(body))
		if err != nil{
			fmt.Printf("During body(with length: %d) sending happend erros: %s", n, err.Error())
		}
		return
	}

	if req.RequestLine.RequestTarget == "/myproblem"{
		err := w.WriteStatusLine(500, "Internal Server Error")
		if err != nil{
		fmt.Println(err.Error())
	}
		headers := headers.NewHeaders()
		body := `<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`
		headers["content-length"] = strconv.Itoa(len(body))
		headers["content-type"] = "text/html"
		err = w.WriteHeaders(headers)
		if err != nil{
			fmt.Print(err.Error())
		}
		n, err := w.WriteBody([]byte(body))
		if err != nil{
			fmt.Printf("During body(with length: %d) sending happend erros: %s", n, err.Error())
		}
		return
	}
		err := w.WriteStatusLine(200, "OK")
		if err != nil{
				fmt.Println(err.Error())
		}
		headers := headers.NewHeaders()
		body := `<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`
		headers["content-length"] = strconv.Itoa(len(body))
		headers["content-type"] = "text/html"
		err = w.WriteHeaders(headers)
		if err != nil{
			fmt.Print(err.Error())
		}
		n, err := w.WriteBody([]byte(body))
		if err != nil{
			fmt.Printf("During body(with length: %d) sending happend erros: %s", n, err.Error())
		}
}


func handlerChunked(w *response.Writer, target string){
	h := headers.NewHeaders()
	if target == "html"{
		h["Trailer"] = "x-content-sha256, x-content-length"
	}

	w.WriteStatusLine(200, "I am proxy machine of chunked data")
	h["transfer-encoding"] = "chunked"
	w.WriteHeaders(h)
	rsp, err := http.Get(fmt.Sprintf("https://httpbin.org/%s", target))
	defer rsp.Body.Close()
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	buff := make([]byte, 1024)
	buffAll := make([]byte, 0, 5024)
	for {
	readToIndex, err := rsp.Body.Read(buff)
	buffAll = append(buffAll, buff[:readToIndex]...)
	info := buff[:readToIndex]
	if len(info) > 0{
		_, err = w.WriteChunkedBody(info)
			if err != nil{
				log.Printf("Error: %v", err)
				return
			}
		}
	if err == io.EOF{
		_, err := w.WriteChunkedBodyDone()
		log.Printf("Error: %v", err)
		if target != "html"{
			_, err = w.W.Write([]byte("\r\n"))
			if err != nil{
				log.Printf("Error: %v", err)
			}
			return 
		}
		break
	}
	if err != nil{
		log.Printf("Error: %v", err)
		return
	}
}
trailers := headers.NewHeaders()
hash := sha256.Sum256(buffAll)
trailers["x-content-sha256"] = fmt.Sprintf("%x", hash)
trailers["x-content-length"] = fmt.Sprintf("%d", len(buffAll))
err = w.WriteTrailers(trailers)
if err != nil{
	log.Printf("Error: %v", err)
}
}

