package utils

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type User struct {
	sync.Mutex
	max   int
	USERS map[string]net.Conn
}

// NewUser creates and initializes a new User instance
func NewUser(max int) *User {
	return &User{
		max:   max,
		USERS: make(map[string]net.Conn),
	}
}

var (
	port   = "8989"
	netcat = NewUser(10)
)

func HandleConnection(conn *net.Conn, PORT string) {
	port = PORT
	if netcat.max == 0 {
		(*conn).Write([]byte("Room is full only 10 people allowed"))
		return
	}
	greeting := "Welcome to TCP-Chat!\n         _nnnn_\n        dGGGGMMb\n       @p~qp~~qMb\n       M|@||@) M|\n       @,----.JM|\n      JS^\\__/  qKL\n     dZP        qKRb\n    dZP          qKKb\n   fZP            SMMb\n   HZM            MMMM\n   FqM            MMMM\n __| \".        |\\dS\"qML\n |    `.       | `' \\Zq\n_)      \\.___.,|     .'\n\\____   )MMMMMP|   .'\n     `-'       `--'"
	defer (*conn).Close()
	(*conn).Write([]byte(greeting))
	name := login((conn))
	if name == "" {
	}
	chat(conn, &name)
	disconect(conn, name)
	netcat.Lock()
	delete(netcat.USERS, name)
	netcat.max++
	netcat.Unlock()
}

func login(conn *net.Conn) string {
	connFile, err := os.OpenFile("netcat-connection_"+port+".log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatalln(err)
	}
	defer connFile.Close()
	chatFile, err := os.OpenFile("netcat-chat_"+port+".log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatalln(err)
	}
	defer chatFile.Close()
	date := time.Now().Format(time.DateTime)
	username := ""
	buffer := make([]byte, 1024)
	(*conn).Write([]byte("\n[ENTER YOUR NAME]:"))
	for {
		n, err := (*conn).Read(buffer)
		if err != nil {
			return ""
		}
		username += string(buffer[:n-1])
		if strings.Contains(string(buffer), "\n") {
			break
		}
	}
	status := checkUsername(username, conn)
	if status != "" {
		(*conn).Write([]byte(status))
		return login(conn)
	} else {
		oldchat, _ := os.ReadFile("netcat-chat_" + port + ".log")
		(*conn).Write(oldchat)
		(*conn).Write([]byte("[" + date + "][" + username + "]:"))
	}
	for user, Conn := range netcat.USERS {
		if user != username {
			Conn.Write([]byte("\n" + username + " has joined our chat...\n[" + date + "][" + user + "]:"))
		}
	}

	connFile.Write([]byte(username + " has joined our chat...\n"))
	return username
}

func checkUsername(username string, conn *net.Conn) string {
	netcat.Lock()
	defer netcat.Unlock()

	if len(username) < 3 {
		return "username too small"
	}
	if netcat.USERS[username] != nil {
		return "username already used"
	}
	if len(username) > 25 {
		return "username too long"
	}
	if !validchars(username) {
		return "only use latin letters and \"-\""
	}
	if netcat.max == 0 {
		return "room is full"
	}
	netcat.USERS[username] = (*conn)
	netcat.max--
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

func chat(Conn *net.Conn, name *string) {
	chatFile, err := os.OpenFile("netcat-chat_"+port+".log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatalln(err)
	}
	msgPrefix := "[" + time.Now().Format(time.DateTime) + "][" + (*name) + "]:"
	msg := ""
	buffer := make([]byte, 1024)
	for {
		n, err := (*Conn).Read(buffer)
		if err != nil {
			return
		}
		msg += string(buffer[:n-1])
		if strings.Contains(string(buffer), "\n") {
			break
		}
	}
	msg = strings.TrimSpace(msg)
	if !Validmsg(msg) {
		(*Conn).Write([]byte("\033[F\033[2K[" + time.Now().Format(time.DateTime) + "][" + (*name) + "]:"))
	} else {
		for name, conn := range netcat.USERS {
			if conn != (*Conn) {
				conn.Write([]byte("\a\n" + msgPrefix + msg + "\n"))
			}
			conn.Write([]byte("[" + time.Now().Format(time.DateTime) + "][" + name + "]:"))

		}
		fmt.Fprint(chatFile, msgPrefix+msg+"\n")
	}
	chatFile.Close()
	chat(Conn, name)
}

func disconect(conn *net.Conn, name string) {
	connFile, err := os.OpenFile("netcat-connection_"+port+".log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatalln(err)
	}
	defer connFile.Close()
	for user, c := range netcat.USERS {
		if c != (*conn) {
			c.Write([]byte("\n" + name + " has left our chat...\n[" + time.Now().Format(time.DateTime) + "][" + user + "]"))
		}
	}
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
