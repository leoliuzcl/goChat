package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
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
	Uid     int
	ConType string
}

const (
	pingWait       = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	writeTimeOut   = 10 * time.Second
	maxMessage     = 10240000
	maxMessageSize = 10240000
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)
var totalNum sync.Mutex
var upgrader = websocket.Upgrader{
	ReadBufferSize:  10240000,
	WriteBufferSize: 10240000,
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
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeTimeOut))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			log.Printf("writeData.....%s, uid is:%d", message, c.uid)
			w.Write(message)
			if err := w.Close(); err != nil {
				return
			}
			log.Printf("writeData.end....%s", message)
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeTimeOut))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (c *Client) readData() {
	defer func() {
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("error: %v", err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		log.Printf("receive msg...%s", message)
		broadcast <- message
	}
}

func sendMessage(num int, c *Client) {
	c.conn.SetWriteDeadline(time.Now().Add(writeTimeOut))
	w, err := c.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return
	}
	str := User{Uid: num, ConType: "connect"}
	jsonStr, _ := json.Marshal(str)
	os.Stdout.Write(jsonStr)
	w.Write(jsonStr)
	w.Close()
}
