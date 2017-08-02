package main

import (
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8081", "ip:port")

func chatRoot(h http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	log.Println(r.URL.Path)
}

func main() {
	flag.Parse()
	log.Printf("listen ip_port: %s\n", *addr)
	http.HandleFunc("/", chatRoot)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Print("start fail")
	}
}
