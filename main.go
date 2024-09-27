package main

import (
	"fmt"
	"net"
	"time"

	"netcat/srm"
)

func main() {
	data, err := srm.Read("./bitri9")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(data))
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		// handle error
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
		}
		go handleConnection(conn)
	}
}

func handleConnection(c net.Conn) {
	f := "Welcome to TCP-Chat!\n         _nnnn_\n        dGGGGMMb\n       @p~qp~~qMb\n       M|@||@) M|\n       @,----.JM|\n      JS^\\__/  qKL\n     dZP        qKRb\n    dZP          qKKb\n   fZP            SMMb\n   HZM            MMMM\n   FqM            MMMM\n __| \".        |\\dS\"qML\n |    `.       | `' \\Zq\n_)      \\.___.,|     .'\n\\____   )MMMMMP|   .'\n     `-'       `--'\n[ENTER YOUR NAME]:"
	c.Write([]byte(f))

	g := string(time.Now().Local().Format("[" + time.DateTime + "] ["))

	b := make([]byte, 1)
	for {
		c.Read(b)
		if b[0] == '\n' {
			break
		}
		g += string(b)
	}
	var conec struct {
		name string
	}
	conec.name = g
	c.Write([]byte(g + "]"))
}
