package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var (
	host    string
	port    int
	folder  string
	server  http.Handler
	mode    serverMode
	https   bool
	certKey string
	cert    string
)

type serverMode int

const (
	modeFile = serverMode(iota)
	modeApiTrace
)

func initFlags() {
	hostPtr := flag.String("host", "localhost", "the designated host")
	portPtr := flag.String("port", "<default>", "designated port (must be int or '<default>')")
	folderPtr := flag.String("folder", ".", "the path to serve")
	modePtr := flag.String("mode", "fileserver", "The server mode ('fileserver', 'apitrace'")
	httpsPtr := flag.Bool("tls", false, "Serve as HTTPS (i.e. TLS)")
	certPtr := flag.String("cert", "server.crt", "Path to TLS Certificate")
	certKeyPtr := flag.String("cert-key", "server.key", "Path to TLS Certificate key")

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
	if *portPtr != "<default>" {
		var err error
		port, err = strconv.Atoi(*portPtr)
		if err != nil {
			flag.Usage()
			panic("Port could not be parsed, please provide int value or '<default>''")
		}
	} else {
		port = 80
		if *httpsPtr {
			port = 443
		}
	}
	host, folder, https, cert, certKey = *hostPtr, *folderPtr, *httpsPtr, *certPtr, *certKeyPtr
}

func main() {
	initFlags()

	// Handle HTTP: 1. Log, 2. Serve file
	if mode == modeFile {
		log.Printf("Serving folder \"%s\" on \"%s:%d\"\n", folder, host, port)
		server = http.FileServer(http.Dir(folder))
	} else {
		server = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			log.Printf("Running API Debugger on \"%s:%d\"\n", host, port)
		})
	}
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("********************* BEGIN REQUEST *************************")
		log.Printf("*** Request: %s %s from %s\n", r.Method, r.URL, r.RemoteAddr)
		log.Println("*** Headers:")
		for k, v := range r.Header {
			log.Printf("***   %s : %s", k, strings.Join(v, " "))
		}
		log.Printf("***   Referer : %s", r.Referer())

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
		enableCors(&w)
		server.ServeHTTP(w, r)
		log.Println("********************* END REQUEST ***************************")
	}))

	var err error
	if https {
		err = http.ListenAndServeTLS(fmt.Sprintf("%s:%d", host, port), cert, certKey, nil)
	} else {
		err = http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), nil)

	}
	if err != nil {
		panic(err)
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
