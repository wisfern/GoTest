package main

import (
	"encoding/json"
	"log"
	"github.com/gorilla/websocket"
	"time"
)

type ClientInfo struct {
	name string
	conn *websocket.Conn
	send chan ChatMsg
}

func (c *ClientInfo) writer() {
	for message := range c.send {
		msg_for_send, err := json.Marshal(message)
		if err != nil {
			log.Fatalln(c.name + " send msg error: " + err.Error())
			break
		}
		err = c.conn.WriteMessage(websocket.TextMessage, msg_for_send)
		if err != nil {
			log.Fatalln(c.name + " send msg error: " + err.Error())
			break
		}
	}
	c.conn.Close()
}

func (c *ClientInfo) reader(server *ChatRoomServerInfo) bool {
	mt, message, err := c.conn.ReadMessage()
	if err != nil {
		log.Println(c.name, "log out,read err:", err, ",mt=", mt)
		return false
	}

	var cmsg ChatMsg
	cmsg.Name = c.name
	cmsg.Msg = string(message)
	cmsg.MsgType = "txt"
	cmsg.Time = time.Now().Format("2006-01-02 15:04:05")

	server.broadcast <- cmsg
	return true
}