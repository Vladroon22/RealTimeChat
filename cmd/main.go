package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Vladroon22/RealTimeChat/internal/server"
)

func main() {
	fmt.Println("Server is listening --> 127.0.0.1:8000")
	http.HandleFunc("/ws", server.New().ServeHTTP)
	log.Fatal(http.ListenAndServe("127.0.0.1:8000", nil))
}
