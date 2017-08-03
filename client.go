package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	send chan []byte
	uid  int
}

type User struct {
	Uid int
}

var totalNum sync.Mutex
var writeTimeOut = 10 * time.Second
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func connect(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	memsNum++
	client := &Client{conn: conn, uid: memsNum, send: make(chan []byte, 1024)}

	mems.clients[memsNum] = client
	go sendMessage(memsNum, client)
	go client.readData()
	go client.writeData()
}

func (c *Client) writeData() {

}

func (c *Client) readData() {

}

func (c *Client) afterConnect() {
}

func sendMessage(num int, c *Client) {
	c.conn.SetWriteDeadline(time.Now().Add(writeTimeOut))
	w, err := c.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return
	}
	str := User{Uid: num}
	jsonStr, _ := json.Marshal(str)
	log.Println(jsonStr)
	w.Write(jsonStr)
	w.Close()
}
