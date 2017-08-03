package main

import (
	"flag"
	"log"
	"net/http"
)

type clientMem struct {
	clients map[int]*Client
}

var addr = flag.String("addr", ":8081", "ip:port")
var mems = &clientMem{clients: make(map[int]*Client)}
var memsNum = 1000

func chatRoot(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	log.Println(r.URL.Path)
	log.Println(r.Method)
	http.ServeFile(w, r, "home.html")
}

func main() {
	flag.Parse()
	log.Printf("listen ip_port: %s\n", *addr)
	http.HandleFunc("/", chatRoot)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		connect(w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Print("start fail")
	}
}
