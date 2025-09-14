package main

import (
	"log"
	"net"

	"github.com/mrvin/tasks-go/tcp-echo/internal/protocol"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal("Error: ", err)
	}
	defer listener.Close()
	log.Print("Server started on port 8000")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print("Error accepting connection: ", err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	request, err := protocol.ReceiveRequest(conn)
	if err != nil {
		log.Print("Error receive request: ", err)
		response := protocol.Response{
			Type:    protocol.TypeResponse,
			ErrorNo: protocol.Failure,
			Buffer:  []byte(err.Error()),
		}
		if err := protocol.SendResponse(conn, response); err != nil {
			log.Print("Error send response: ", err)
		}

		return
	}

	log.Print("Received request: ", request.Str)

	response := protocol.Response{
		Type:    protocol.TypeResponse,
		ErrorNo: protocol.Success,
		Buffer:  []byte("Echo: " + request.Str),
	}
	if err := protocol.SendResponse(conn, response); err != nil {
		log.Print("Error send response: ", err)
		return
	}

	log.Print("Response sent successfully")
}
