package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Server struct {
	mutex   sync.Mutex
	clients map[string]*websocket.Conn
}

func New() *Server {
	return &Server{
		clients: make(map[string]*websocket.Conn),
	}
}

func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	defer conn.Close()

	for {
		go srv.GetMessage(srv.GetName(conn), conn)

		srv.SendMessage(conn)
	}
}

func (srv *Server) GetMessage(name string, conn *websocket.Conn) {
	for {
		tp, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("ERROR Read: ", err)
		}
		fmt.Printf("\n%s: %s\n", name, string(message))

		if break_conn(message) {
			srv.mutex.Lock()
			delete(srv.clients, name)
			srv.mutex.Unlock()
			conn.Close()
			fmt.Println("\n" + name + " disconnected!")
			if len(srv.clients) == 0 {
				fmt.Println("All clients are disconnected!")
				os.Exit(0)
				break
			}
			break
		}

		srv.BroadCast(name, tp, message)
	}
}

func (srv *Server) SendMessage(conn *websocket.Conn) {
	var mess string
	fmt.Print("You: ")
	_, err := fmt.Scanln(&mess)
	if err != nil {
		fmt.Println("ERROR Scan: ", err)
		return
	}

	if err := conn.WriteMessage(websocket.TextMessage, []byte(mess)); err != nil {
		fmt.Println("ERROR Write: ", err)
		return
	}

	if break_conn([]byte(mess)) {
		conn.Close()
		for names := range srv.clients {
			srv.mutex.Lock()
			delete(srv.clients, names)
			srv.mutex.Unlock()
		}
		fmt.Println("\nServer(You) closed coonection!")
		os.Exit(0)
	}

}

func (srv *Server) GetName(conn *websocket.Conn) string {
	var name string
	if err := conn.ReadJSON(&name); err != nil {
		fmt.Println("ERROR Read: ", err)
		return ""
	}

	srv.mutex.Lock()
	srv.clients[string(name)] = conn
	srv.mutex.Unlock()

	fmt.Println("\n" + name + " from " + conn.RemoteAddr().String() + " has been connected")
	return name
}

func (srv *Server) BroadCast(sender string, tp int, mess []byte) {
	for names, conns := range srv.clients {
		if names != sender {
			if err := conns.WriteMessage(websocket.TextMessage, []byte(mess)); err != nil {
				fmt.Println("ERROR Write: ", err)
				return
			}
		}
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
