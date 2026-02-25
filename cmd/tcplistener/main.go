package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer listener.Close()

	// Используй Println, чтобы был перенос строки
	fmt.Println("Server is listening on port :42069")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Задание просит "Prints a message... that a connection has been accepted"
		fmt.Println("accepted connection")

		lines := getLinesChannel(conn)
		
		// Читаем строки пока канал не закроется
		for line := range lines {
			fmt.Println(line)
		}
		
		// 1. ВЫНЕСЛИ ЭТО ИЗ ЦИКЛА
		fmt.Println("connection closed")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	linesChan := make(chan string)
	
	go func() {
		defer close(linesChan)
		// 2. ДОБАВИЛИ ЗАКРЫТИЕ СОЕДИНЕНИЯ
		defer f.Close() 

		buff := make([]byte, 8)
		remainder := "" // "Хвост" с прошлого чтения

		for {
			n, err := f.Read(buff)
			
			// 3. СНАЧАЛА ОБРАБАТЫВАЕМ ДАННЫЕ
			if n > 0 {
				chunk := remainder + string(buff[:n])
				parts := strings.Split(chunk, "\n")
				
				// Сохраняем новый хвост
				remainder = parts[len(parts)-1]
				
				// Отправляем все полные строки
				for _, s := range parts[:len(parts)-1] {
					linesChan <- s
				}
			}

			// ПОТОМ ПРОВЕРЯЕМ ОШИБКИ
			if err != nil {
				// Если EOF и остался кусочек текста без \n в конце
				if err == io.EOF && remainder != "" {
					linesChan <- remainder
				}
				return // Выходим из горутины
			}
		}
	}()
	
	return linesChan
}