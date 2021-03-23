package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var (
	host   string
	port   int
	folder string
	server http.Handler
	mode   serverMode
)

type serverMode int

const (
	modeFile = serverMode(iota)
	modeApiTrace
)

func initFlags() {
	hostPtr := flag.String("host", "localhost", "the designated host")
	portPtr := flag.Int("port", 7890, "designated port")
	folderPtr := flag.String("folder", ".", "the path to serve")
	modePtr := flag.String("mode", "fileserver", "The server mode ('fileserver', 'apitrace'")

	flag.Parse()

	switch strings.ToLower(*modePtr) {
	case "fileserver":
		mode = modeFile
	case "apitrace":
		mode = modeApiTrace
	default:
		flag.Usage()
		panic("Invalid Mode selected: " + *modePtr)
	}
	host, port, folder = *hostPtr, *portPtr, *folderPtr
}

func main() {
	initFlags()

	// Handle HTTP: 1. Log, 2. Serve file
	if mode == modeFile {
		server = http.FileServer(http.Dir(folder))
	} else {
		server = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			d, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println("ERROR while reading Body:", err.Error())
			}
			log.Println("*** Body:")
			body := "<empty>"
			if len(d) > 0 {
				body = string(d)
			}
			log.Printf("***   %s\n", body)
		})
	}
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("********************* BEGIN REQUEST *************************")
		log.Printf("*** Request: %s %s from %s\n", r.Method, r.URL, r.RemoteAddr)
		log.Println("*** Headers:")
		for k, v := range r.Header {
			log.Printf("***   %s : %s", k, strings.Join(v, " "))
		}
		enableCors(&w)
		server.ServeHTTP(w, r)
		log.Println("********************* END REQUEST ***************************")
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
