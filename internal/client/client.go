package client

import (
	"fmt"
	"os"

	"golang.org/x/net/websocket"
)

func GetMessage(conn *websocket.Conn) error {
	buf := make([]byte, 2048)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if break_conn(buf) {
		fmt.Println("\nServer closed connection!")
		conn.Close()
		os.Exit(0)
	}

	fmt.Println("Server:", string(buf[:n]))
	return nil
}

func SendMessage(conn *websocket.Conn) error {
	var source string
	fmt.Print("You: ")
	_, err := fmt.Scanln(&source)
	if err != nil {
		fmt.Println("ERROR Scan: ", err.Error())
		return err
	}

	if n, err := conn.Write([]byte(source)); n == 0 || err != nil {
		fmt.Println(err.Error())
		return err
	}
	if break_conn([]byte(source)) {
		fmt.Println("\nYou closed connection!")
		conn.Close()
		os.Exit(0)
	}

	return nil
}

func SendName(conn *websocket.Conn) {
	var name string
	fmt.Print("What's your name, warrior? ")
	_, err := fmt.Scanln(&name)
	if err != nil {
		fmt.Println("ERROR Scan:", err.Error())
	}

	if n, err := conn.Write([]byte(name)); n == 0 || err != nil {
		fmt.Println(err.Error())
	}
}

func break_conn(text []byte) bool {
	for _, ch := range text {
		if string(ch) == "#" {
			return true
		}
	}
	return false
}

func main() {
	fmt.Println("Connecting to server...")
	conn, err := websocket.Dial("ws://127.0.0.1:8000/ws", "", "127.0.0.1:8080")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()

	SendName(conn)

	fmt.Println("\nYou have been connected to server --> 127.0.0.1:8080")
	fmt.Println("tap '#' to exit")
	fmt.Println()

	for {
		if err := SendMessage(conn); err != nil {
			fmt.Println(err)
			continue
		}
		if err := GetMessage(conn); err != nil {
			fmt.Println(err)
			continue
		}
	}

}
