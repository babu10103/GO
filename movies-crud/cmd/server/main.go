package main

import (
	"flag"
	"movies-crud/internal/router"
	"movies-crud/pkg/log"
	"net"
	"net/http"
)

func main() {
	r := router.NewRouter()

	port := flag.String("port", "8000", "port to run the server on")
	flag.Parse()

	ln, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.ErrorLogger.Fatalf("Port %v i already in use. Please choose another port", *port)
		return
	}
	ln.Close()

	log.InfoLogger.Printf("Starting server at %s port...\n", *port)

	if err := http.ListenAndServe(":"+*port, r); err != nil {
		log.ErrorLogger.Fatalf("Failed to start server: %v", err)
	}
}
