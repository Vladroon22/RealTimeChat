package server

import (
	"fmt"
	"net/http"
	"os"
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
	server  *http.Server
	logg    *logrus.Logger
}

func New() *Server {
	return &Server{
		clients: make(map[string]*websocket.Conn),
		logg:    logrus.New(),
		server: &http.Server{
			Addr: ":8080",
		},
	}
}

func (s *Server) Start() {
	s.logg.Fatalln(s.server.ListenAndServe())
}

func (s *Server) GetName(conn *websocket.Conn) string {
	mesType, name, err := conn.ReadMessage()
	if len(name) == 0 || err != nil {
		s.logg.Errorln("ERROR Read: ", err.Error())
		return ""
	}

	s.clients[string(name)] = conn
	s.logg.Infoln("\n" + string(name) + " from " + conn.RemoteAddr().String() + " has been connected")

	if err := conn.WriteMessage(mesType, name); err != nil {
		s.logg.Errorln(err)
	}

	return string(name)
}

func (s *Server) GetMessage(name string, conn *websocket.Conn) {
	defer conn.Close()
	for {
		mesType, message, err := conn.ReadMessage()
		if len(message) == 0 || err != nil {
			s.logg.Errorln("ERROR Read: ", err.Error())
		}

		if err := conn.WriteMessage(mesType, message); err != nil {
			s.logg.Errorln(err)
		}

		if break_conn(message) {
			s.mutex.Lock()
			delete(s.clients, name)
			s.mutex.Unlock()
			conn.Close()
			fmt.Println("\n" + name + " disconnected!")
			if len(s.clients) == 0 {
				fmt.Println("All clients are disconnected!")
				os.Exit(0)
				break
			}
			break
		}

		go func() {
			s.logg.Infoln(name + ": " + string(message))
		}()
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

func homePage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Home"))
}

func (s *Server) wsEndPoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logg.Errorln(err)
	}

	s.logg.Infoln("Client connected!")

	go s.GetMessage(s.GetName(ws), ws)
}

func (s *Server) SetupEndPoints() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", s.wsEndPoint)
}
