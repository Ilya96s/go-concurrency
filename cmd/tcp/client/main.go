package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	address := flag.String("address", "localhost:3223", "server address")
	flag.Parse()

	conn, err := net.Dial("tcp", *address)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to", *address)
	fmt.Println("Type commands (SET key value / GET key / DEL key). Type EXIT to quit.")

	userInput := bufio.NewScanner(os.Stdin)
	serverReader := bufio.NewReader(conn)

	for {
		fmt.Print("> ")
		if !userInput.Scan() {
			break
		}
		line := userInput.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Отправка команды на сервер
		_, err := conn.Write([]byte(line + "\n"))
		if err != nil {
			fmt.Println("Error writing to server:", err)
			return
		}

		// Чтение ответа от  сервера
		resp, err := serverReader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from server:", err)
			return
		}

		fmt.Print(resp)
		if strings.HasPrefix(strings.ToUpper(line), "EXIT") {
			break
		}
	}
}
