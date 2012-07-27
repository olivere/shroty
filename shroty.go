package main

import (
	"code.google.com/p/gorilla/mux"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"net/http"
)

const (
	defaultHttpAddr = ":8910"
)

var (
	staticDir *string

	httpAddr = flag.String("http", "", "http server address, e.g. 192.168.2.1:8910")
)

func shorten(w http.ResponseWriter, r *http.Request) {
	url := mux.Vars(r)["url"]

	fmt.Fprintf(w, "Shorten %v", url)
}

func open(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func info(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func static(w http.ResponseWriter, r *http.Request) {
	file := mux.Vars(r)["file"]
	
	if len(file) == 0 {
		file = "index.html"
	}

	staticFile := path.Join(*staticDir, file)

	fi, e := os.Stat(staticFile)
	if e != nil {
		log.Println("File error: ", e)
		http.NotFound(w, r)
	} else if fi.IsDir() {
		log.Println("File not found: ", staticFile)
		http.NotFound(w, r)
	} else {
		http.ServeFile(w, r, staticFile)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: shroty -http="+defaultHttpAddr+"\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	cwd, e := os.Getwd()
	if e != nil {
		log.Fatalln("Cannot get current directory")
	}

	cwd, _ = filepath.Abs(cwd)

	s := path.Join(cwd, "public")
	if _, e := os.Stat(s); e != nil {
		log.Fatalln("Cannot find /public directory")
	}
	staticDir = &s

	router := mux.NewRouter()
	
	router.HandleFunc("/s/{url:(.*$)}", shorten)
	router.HandleFunc("/go/{id:([a-zA-Z0-9]+$)}", open)
	router.HandleFunc("/go/{id:([a-zA-Z0-9]+$)+}", info)
	router.HandleFunc("/{file:(.*$)}", static)

	server := *httpAddr
	if (server == "") {
		server = defaultHttpAddr
	}

	httpServ := &http.Server{
		Addr: server,
		Handler: router,
	}

	log.Printf("Listening on %v\n", server)

	httpServ.ListenAndServe()
}
