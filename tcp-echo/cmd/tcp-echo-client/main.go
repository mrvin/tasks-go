package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/mrvin/tasks-go/tcp-echo/internal/protocol"
)

func main() {
	var port, typeRequest int
	var host, str string

	flag.IntVar(&port, "port", 8000, "port") //nolint:mnd
	flag.StringVar(&host, "host", "localhost", "host name")
	flag.IntVar(&typeRequest, "type", int(protocol.TypeRequest), "request type")
	flag.StringVar(&str, "msg", "Hello!", "message")
	flag.Parse()

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatal("Error:", err)
	}
	defer conn.Close()

	request := protocol.Request{
		Type: int32(typeRequest),
		Str:  str,
	}
	if err := protocol.SendRequest(conn, request); err != nil {
		log.Print("Error send request:", err)
		return
	}

	response, err := protocol.ReceiveResponse(conn)
	if err != nil {
		log.Print("Error receive response:", err)
		return
	}

	fmt.Printf("ErrorNo: %d\n", response.ErrorNo)       //nolint:forbidigo
	fmt.Printf("Buffer: %s\n", string(response.Buffer)) //nolint:forbidigo
}
