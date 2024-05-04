package main

import (
	"fmt"
	"log"

	"github.com/Vladroon22/RealTimeChat/internal/server"
)

func main() {
	fmt.Println("Server is listening --> localhost:8080")
	srv := server.New()
	srv.SetupEndPoints()
	log.Fatal(srv.Server.ListenAndServe())
}
