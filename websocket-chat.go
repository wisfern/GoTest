package main

import (
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

type ChatMsg struct {
	Name    string `json:"name"`
	Msg     string `json:"msg"`
	MsgType string `json:"msg_type"`
	Time    string `json:"time"`
}

type ChatRoomServerInfo struct {
	member     map[string]ClientInfo
	register   chan ClientInfo
	unregister chan ClientInfo
	broadcast  chan ChatMsg
}

// define server global config
var (
	maxClientConnection int    = 1000
	listenAddrPort      string = *flag.String("wsaddr", ":11888", "websocket service address")
)

// init global server info
var server = ChatRoomServerInfo{
	member:     make(map[string]ClientInfo),
	register:   make(chan ClientInfo, maxClientConnection),
	unregister: make(chan ClientInfo, maxClientConnection),
	broadcast:  make(chan ChatMsg, maxClientConnection),
}

func (sr *ChatRoomServerInfo) serverRun() {
	for {
		select {
		case client := <-sr.register: {
				sr.member[client.name] = client
				// notify all chat room member, new user login
				var cm ChatMsg
				cm.Name = "system"
				cm.Msg = client.name+" login!"
				cm.MsgType = "register"
				cm.Time = time.Now().Format("2006-01-02 15:04:05")
				server.broadcast <- cm
			}
		case client := <-server.unregister: {
				if _, ok := server.member[client.name]; ok {
					var cm ChatMsg
					cm.Name = "system"
					cm.Msg = client.name + " logout"
					cm.MsgType = "unregister"
					cm.Time = time.Now().Format("2006-01-02 15:04:05")
					server.broadcast <- cm

					delete(server.member, client.name)
					client.conn.Close()
				}
			}
		case msg := <-server.broadcast: {
				for k, c := range server.member {
					select {
					case c.send <- msg:
						log.Print(c.name + " send msg "+ msg.Msg)
						continue
					default:
						log.Println(c.name + " logout")
						delete(server.member, k)
					}
				}
			}
		} // end of select
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if len(r.Form["name"]) <= 0 {
		log.Println("dont have login name")
		return
	}
	register(w, r)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func register(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	var ci ClientInfo
	ci.name = r.Form["name"][0]
	ci.conn = conn
	ci.send = make(chan ChatMsg, 256)
	log.Println("login name="+r.Form["name"][0], ",ip="+r.RemoteAddr)

	server.register <- ci

	defer func() {
		conn.Close()
		server.unregister <- ci
	}()

	go ci.writer()
	for {
		if ci.reader(&server) == false {
			break
		}
	}
}

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.Println("start server listen in http://localhost" + listenAddrPort + "!")
	go server.serverRun()

	http.HandleFunc("/login", login)
	http.Handle("/", http.FileServer(http.Dir("./html")))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir("fonts"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.ListenAndServe(listenAddrPort, nil)
	log.Println("stop server!")
}