package structs

import (
	"bufio"

	"fmt"
	"net"
)

func AcceptSocket(head *Head) {
	fmt.Println("Start server...")

	// listen on port 8000
	ln, _ := net.Listen("tcp", ":4567")

	// accept connection

	for {
		conn, _ := ln.Accept()

		fmt.Println("Connection Aquired")

		// run loop forever (or until ctrl-c)

		r := bufio.NewReader(conn)
		for {
			// get message, output

			value, err := r.ReadString('\n')
			// message, e := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				// fmt.Println(err)
				conn.Close()
				break

			}
			// fmt.Println("Message Received:", string(value))
			if value != "e" {
				head.ApplyRotation(value)
			}
		}
	}

}
