package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
)

var (
	host       string
	port       int
	folder     string
	logHeaders bool
	server     http.Handler
)

func initFlags() {
	hostPtr := flag.String("host", "localhost", "the designated host")
	portPtr := flag.Int("port", 7890, "designated port")
	folderPtr := flag.String("folder", ".", "the path to serve")
	logHeadersPtr := flag.Bool("log-headers", false, "wether client headers should be logged")

	flag.Parse()
	host, port, folder, logHeaders = *hostPtr, *portPtr, *folderPtr, *logHeadersPtr
}

func main() {
	initFlags()

	// Handle HTTP: 1. Log, 2. Serve file
	server = http.FileServer(http.Dir(folder))
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request detected %s %s from %s\n", r.Method, r.URL, r.RemoteAddr)
		if logHeaders {
			for k, v := range r.Header {
				log.Printf("Client Header: %s : %s", k, strings.Join(v, " "))
			}
		}
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
