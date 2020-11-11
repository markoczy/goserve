package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	host   string
	port   int
	folder string
	server http.Handler
)

func initFlags() {
	hostPtr := flag.String("host", "localhost", "the designated host")
	portPtr := flag.Int("port", 7890, "designated port")
	folderPtr := flag.String("folder", ".", "the path to serve")
	flag.Parse()
	host, port, folder = *hostPtr, *portPtr, *folderPtr
}

func main() {
	initFlags()

	// Handle HTTP: 1. Log, 2. Serve file
	server = http.FileServer(http.Dir(folder))
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request detected %s %s from %s\n", r.Method, r.URL, r.RemoteAddr)
		enableCors(&w)
		server.ServeHTTP(w, r)
	}))

	log.Printf("Serving folder \"%s\" on \"%s:%d\"\n", folder, host, port)

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), nil)
	if err != nil {
		panic(err)
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
