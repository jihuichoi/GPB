package main

import (
	"time"

	"github.com/gorilla/websocket"
)

type client struct {
	// socket is the web socket for this client
	socket *websocket.Conn

	// send is a channel on which messages are sent
	// 사용자가 받은 메세지. room 의 forward chan 에 메세지가 들어오면
	// 각 사용자의 send chan 에 메세지 전달. 그러면 클라이언트 앱이 write 메서드를 통해 사용자의 화면에 글을 뿌림(write)
	// send chan []byte

	// ch2: 사용자 정보도 함께 전달하기 위해 변경
	send chan *message

	// room is the room this client is chatting in
	room *room

	// userData holds information about the user
	// from room.go userData: objx.MustFromBase64(authCookie.Value),
	userData map[string]interface{}
}

// 유저의 행동으로서 read 가 아니라, 클라이언트 프로그램의 행동으로서 read
// 즉, 사용자가 글을 쓰면, 클라이언트 앱이 그 내용을 읽어서(read) room 의 forward chan 으로 전달
func (c *client) read() {
	defer c.socket.Close()
	for {
		// _, msg, err := c.socket.ReadMessage()
		// if err != nil {
		// 	return
		// }
		// c.room.forward <- msg // room 의 forward 채널로 읽은 msg 를 전달

		// ch2: json 으로 변경
		var msg *message
		// err := c.socket.ReadJSON(&msg)
		// if err != nil {
		// 	return
		// }
		if err := c.socket.ReadJSON(&msg); err != nil {
			return
		}
		msg.When = time.Now()
		msg.Name = c.userData["name"].(string)
		msg.AvatarURL, _ = c.room.avatar.GetAvatarURL(c)
		// c.userData["avatar_url"] 이 nil 일 경우 string type 에 대입하면 panic 이 발생하므로 미리 확인해준다.
		// if avatarURL, ok := c.userData["avatar_url"]; ok {
		// 	msg.AvatarURL = avatarURL.(string)
		// }
		c.room.forward <- msg
	}
}

// 마찬가지로, 사용자가 글을 작성하는 것이 아니라,
// 클라이언트 앱이 forward chan 에서 각 클라이언트의 send 채널로 메세지를 전달하면
// send 채널에 있는 메세지를 화면에 write 한다는 의미
func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		// err := c.socket.WriteMessage(websocket.TextMessage, msg)

		// ch2: json 형식으로 변경
		err := c.socket.WriteJSON(msg)
		if err != nil {
			return
		}
	}
}
