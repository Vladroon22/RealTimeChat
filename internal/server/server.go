package server

import (
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Server struct {
	mutex   sync.Mutex
	clients map[string]*websocket.Conn
	Server  *http.Server
	logg    *logrus.Logger
}

func homePage(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Home")
}

func (s *Server) wsEndPoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client connected!")
}

func (s *Server) SetupEndPoints() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", s.wsEndPoint)
}

func New() *Server {
	return &Server{
		clients: make(map[string]*websocket.Conn),
		logg:    logrus.New(),
		Server: &http.Server{
			Addr: ":8080",
		},
	}
}
