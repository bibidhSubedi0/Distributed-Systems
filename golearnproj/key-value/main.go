package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

type keyvalue struct {
	data  map[string]string
	mutex sync.Mutex
}

func listen() {
	kvl := keyvalue{
		data: make(map[string]string),
	}
	var port = ":5000"
	var protocal = "tcp"
	listener, err := net.Listen(protocal, port)
	if err != nil {
		panic(err)
	}

	defer listener.Close()
	fmt.Println("Listening on port :", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error ", err)
		}
		go handelconnection(conn, &kvl)
	}
}

func handelconnection(conn net.Conn, kvl *keyvalue) {
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	fmt.Println("Client connected from : ", clientAddr)
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Printf("Received from %s: %s\n", clientAddr, text)
		// set key=hello value=hi

		parts := strings.Split(text, " ")
		switch parts[0] {
		case "set":
			k := strings.Split(parts[1], "=")[1]
			v := strings.Split(parts[2], "=")[1]

			{
				kvl.mutex.Lock()
				kvl.data[k] = v
				kvl.mutex.Unlock()
			}

			fmt.Println("Set ", v, " at ", k, " by ", clientAddr)

		case "get":
			k := strings.Split(parts[1], "=")[1]

			kvl.mutex.Lock()
			v, ok := kvl.data[k]
			kvl.mutex.Unlock()

			if ok {
				conn.Write([]byte(v + "\n"))
			} else {
				fmt.Println("Key", k, "not found")
			}
		case "replicate":
			fmt.Println("replicate")
		default:
			fmt.Println("Invalid command")
		}
	}

}

func main() {
	listen()
}
