package utils

import (
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	mu    sync.Mutex
	max   = 10
	USERS = make(map[string]net.Conn)
	port  = "8989"
)

func HandleConnection(conn net.Conn, PORT string) {
	port = PORT
	if max == 0 {
		conn.Write([]byte("Room is full only 10 people allowed"))
		return
	}
	greeting := "Welcome to TCP-Chat!\n         _nnnn_\n        dGGGGMMb\n       @p~qp~~qMb\n       M|@||@) M|\n       @,----.JM|\n      JS^\\__/  qKL\n     dZP        qKRb\n    dZP          qKKb\n   fZP            SMMb\n   HZM            MMMM\n   FqM            MMMM\n __| \".        |\\dS\"qML\n |    `.       | `' \\Zq\n_)      \\.___.,|     .'\n\\____   )MMMMMP|   .'\n     `-'       `--'"
	defer conn.Close()
	conn.Write([]byte(greeting))
	name := login(conn)
	chat(conn, name)
	disconect(conn, name)
}

func login(conn net.Conn) string {
	connFile, err := os.OpenFile("netcat-connection_"+port+".log", os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatalln(err)
	}
	defer connFile.Close()
	chatFile, err := os.Open("netcat-chat_" + port + ".log")
	if err != nil {
		log.Fatalln(err)
	}
	defer chatFile.Close()
	date := time.Now().Format(time.DateTime)
	username := ""
	buffer := make([]byte, 1)
	conn.Write([]byte("\n[ENTER YOUR NAME]:"))
	for {
		conn.Read(buffer)
		if buffer[0] == '\n' {
			break
		}
		username += string(buffer)
	}
	status := checkUsername(username, conn)
	if status != "" {
		conn.Write([]byte(status))
		return login(conn)
	} else {
		oldchat, _ := os.ReadFile("netcat-chat_" + port + ".log")
		conn.Write(oldchat)
		conn.Write([]byte("[" + date + "][" + username + "]:"))
	}
	for user, Conn := range USERS {
		if user != username {
			Conn.Write([]byte("\n" + username + " has joined our chat...\n[" + date + "][" + user + "]:"))
		}
	}
	connFile.Write([]byte(username + " has joined our chat...\n"))
	return username
}

func checkUsername(username string, conn net.Conn) string {
	mu.Lock()
	defer mu.Unlock()

	if len(username) < 3 {
		return "username too small"
	}
	if USERS[username] != nil {
		return "username already used"
	}
	if len(username) > 25 {
		return "username too long"
	}
	if !validchars(username) {
		return "only use latin letters and \"-\""
	}
	if max == 0 {
		return "room is full"
	}
	USERS[username] = conn
	max--
	return ""
}

func validchars(s string) bool {
	for _, v := range s {
		if !((v >= 'a' && v <= 'z') || (v >= 'A' && v <= 'Z') || v == '-') {
			return false
		}
	}
	return true
}

func chat(Conn net.Conn, name string) {
	defer delete(USERS, name)
	chatFile, err := os.OpenFile("netcat-chat_"+port+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatalln(err)
	}
	msgPrefix := "[" + time.Now().Format(time.DateTime) + "][" + name + "]:"
	msg := ""
	buffer := make([]byte, 1)
	for {
		_, err := Conn.Read(buffer)
		if err != nil {
			return
		}
		if buffer[0] == '\n' {
			break
		}
		msg += string(buffer)
	}
	msg = strings.TrimSpace(msg)
	if !Validmsg(msg) {
		Conn.Write([]byte("message too long\n"))
		Conn.Write([]byte("[" + time.Now().Format(time.DateTime) + "][" + name + "]:"))
	} else {
		for name, conn := range USERS {
			if conn != Conn {
				conn.Write([]byte("\n" + msgPrefix + msg + "\n"))
			}
			conn.Write([]byte("[" + time.Now().Format(time.DateTime) + "][" + name + "]:"))

		}
		chatFile.Write([]byte(msgPrefix + msg + "\n"))
	}
	chatFile.Close()
	chat(Conn, name)
}

func disconect(conn net.Conn, name string) {
	mu.Lock()
	connFile, err := os.OpenFile("netcat-connection_"+port+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatalln(err)
	}
	defer connFile.Close()
	for user, c := range USERS {
		if c != conn {
			c.Write([]byte("\n" + name + " has left our chat...\n[" + time.Now().Format(time.DateTime) + "][" + user + "]"))
		}
	}
	mu.Unlock()

	connFile.Write([]byte(name + " has left our chat...\n"))
}

func Validmsg(msg string) bool {
	if len(msg) > 255 {
		return false
	}
	if msg == "" {
		return false
	}
	for _, v := range msg {
		if (v < 32 || v > 126) && (v < 128 || v > 255) {
			return false
		}
	}
	return true
}
