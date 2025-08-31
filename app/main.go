package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Println("accepted new connection")

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)

		if err != nil {
			log.Println("error happend")
			if errors.Is(err, io.EOF) {
				log.Println("EOF of the conn")

				break
			}

			log.Fatal("error while reading connection happened: ", err)
		}

		log.Printf("read %d bytes from conn\n", n)

		if n == 0 {
			continue
		}

		conn.Write([]byte("+PONG\r\n"))
	}
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	log.Printf("Started server on address: %s\n", l.Addr())

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		handleConnection(conn)
	}
}
