package structs

import (
	"bufio"

	"fmt"
	"net"
)

func AcceptSocket(head *Head) {
	fmt.Println("Start server...")

	ln, _ := net.Listen("tcp", ":4567")

	fmt.Println(ln.Addr())

	for {
		conn, _ := ln.Accept()

		fmt.Println("Connection Aquired")

		r := bufio.NewReader(conn)
		for {

			value, err := r.ReadString('\n')

			if err != nil {

				conn.Close()
				break

			}

			if value != "e" {
				head.ApplyRotation(value)
			}
		}
	}

}
