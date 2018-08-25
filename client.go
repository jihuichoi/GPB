package main

import "github.com/gorilla/websocket"

// client represents a single chatting user.
type client struct {
	// socket is the web socket for this client.
	socket *websocket.Conn

	// send is a channel on which messages are sent.
	// 사용자가 받은 메세지. room 의 forward chan 에 메세지가 들어오면
	// 각 사용자의 send chan 에 메세지 전달. 그러면 클라이언트 앱이 write 메서드를 통해 사용자의 화면에 글을 뿌림(write)
	send chan []byte

	// room is the room this client is chatting in.
	room *room
}

// 유저의 행동으로서 read 가 아니라, 클라이언트 프로그램의 행동으로서 read 인듯..
// 즉, 사용자가 글을 쓰면, 클라이언트 앱이 그 내용을 읽어서(read) room 의 forward chan 으로 전달
func (c *client) read() {
	defer c.socket.Close()
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		c.room.forward <- msg
	}
}

// 마찬가지로, 사용자가 글을 작성하는 것이 아니라, 클라이언트 앱이 forward chan 에 있는 메세지를
// 화면에 write 한다는 의미인듯.
func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
