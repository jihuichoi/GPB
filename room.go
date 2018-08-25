package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type room struct {
	// forward is a channel that holds incoming messages
	// that should be forwarded to the other clients.
	forward chan []byte

	// join is a channel for clients wishing to join the room.
	join chan *client

	// leave is a channel for clients wishing to leave the room.
	leave chan *client

	// clients holds all current clients in this room
	clients map[*client]bool
}

func (r *room) run() {
	for { // 무한 루프 돌면서 아래 select 문을 반복
		select {
		case client := <-r.join: // join 채널에 클라이언트가 들어오면
			// joining
			r.clients[client] = true
		case client := <-r.leave: // leave 채널에 클라이언트가 들어오면
			// leaving
			delete(r.clients, client)
		case msg := <-r.forward: // forward 채널에 메세지가 들어오면
			// forward message to all clients
			for client := range r.clients {
				client.send <- msg // 각 클라이언트의 send 채널로 메세지 전달
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

// 웹소켓 연결을 위해 http 연결을 웹소켓용으로 업그레이드?
var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	// 클라이언트와 소켓 연결 생성
	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	r.join <- client                     // room 입장을 위해 join 채널에 클라이언트를 전달
	defer func() { r.leave <- client }() // 웹소켓 종료시 클라이언트가 룸에서 떠남을 기록
	go client.write()                    // 클라이언트 화면에 메세지를 뿌림
	client.read()                        // 클라이언트의 입력을 기다림
}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}
