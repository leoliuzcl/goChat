package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
)

type clientMem struct {
	clients map[int]*Client
}

type messageConten struct {
	MsgType int
	Message []byte
	Uid     int
}

var addr = flag.String("addr", ":8081", "ip:port")
var mems = &clientMem{clients: make(map[int]*Client)}
var unRegUID = make(chan int)
var broadcast = make(chan []byte)
var memsNum = 1000

func chatRoot(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "home.html")
}
func handMessage() {
	for {
		select {
		case uid := <-unRegUID:
			if _, ok := mems.clients[uid]; ok {
				delete(mems.clients, uid)
			}
		case message := <-broadcast:
			msgT, uid, messageConten := parseMessage(message)
			log.Printf("msgorgi...%s", message)
			log.Printf("msg...%s", messageConten)
			log.Printf("msgtype...%d", msgT)
			log.Printf("msguid...%d", uid)
			for _, client := range mems.clients {
				if msgT == 1 {
					if client.uid == uid {
						select {
						case client.send <- message:
						}
					}
				} else {
					select {
					case client.send <- message:
					}
				}

			}
		}
	}
}

func parseMessage(msg []byte) (int, int, []byte) {
	fmt.Println(msg)
	msgC := &messageConten{}
	err := json.Unmarshal(msg, &msgC)
	if err != nil {
		fmt.Println("Unmarshal faild")
	}
	msgType := msgC.MsgType
	uid := msgC.Uid
	msgConten := msgC.Message
	return msgType, uid, msgConten
}

func main() {
	flag.Parse()
	log.Printf("listen ip_port: %s\n", *addr)
	http.HandleFunc("/", chatRoot)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		connect(w, r)
	})
	go handMessage()
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Print("start fail")
	}
}
