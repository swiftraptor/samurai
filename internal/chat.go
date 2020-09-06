package internal

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Client represents a user session
type Client struct {
	Id         string
	Username   string
	Connection *websocket.Conn
	Sender     chan Message
}

type MessageRoom struct {
	FriendlyName string
	RoomId       string
	HostId       string
	Users        []string
	InviteCode   string
}

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()),
)

func generateInviteCode() string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	buf := make([]byte, 6)
	for i := range buf {
		buf[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(buf)
}

func (room MessageRoom) inRoom(userId string) bool {
	for _, user := range room.Users {
		if user == userId {
			return true
		}
	}
	return false
}

func MakeRoom(hostId string, friendlyName string) MessageRoom {
	return MessageRoom{
		RoomId:       uuid.New().String(),
		HostId:       hostId,
		Users:        make([]string, 16),
		InviteCode:   generateInviteCode(),
		FriendlyName: friendlyName,
	}
}

func (client *Client) handleMessage(server *Server, message Message) {
	var response Message
	response.Body = make(map[string]interface{})
	switch message.Header {

	case RegisterRequest:
		if username, ok := message.Body["username"]; ok {
			client.Username = fmt.Sprintf("%v", username)
			response = Registered()
		} else {
			response = RegisterFailed()
		}
		break

	case InfoRequest:
		response = Info(client.Username)
		break

	case JoinRoomRequest:
		if roomId, ok := message.Body["roomId"]; ok {
			if room, roomExists := server.Rooms[fmt.Sprintf("%v", roomId)]; roomExists {
				room.Users = append(room.Users, client.Username)
				// todo announce
				response = JoinedRoom()
			}
		}

	case ListRoomUsersRequest:
		if roomId, ok := message.Body["roomId"]; ok {
			if room, roomExists := server.Rooms[fmt.Sprintf("%v", roomId)]; roomExists {
				if room.inRoom(client.Id) {
					response = RoomUserList(room.Users)
				}
			}
		}

	case NewRoomRequest:
		if friendlyName, ok := message.Body["friendlyName"]; ok {
			room := MakeRoom(client.Id, fmt.Sprintf("%v", friendlyName))
			server.Rooms[room.RoomId] = &room
			response = NewRoomOk(room.InviteCode, room.RoomId)
		}
	default:
		response = InvalidMessage()
		break
	}

	go client.SendMessage(response)
}

func (client *Client) ReadStream(server Server) {
	client.Connection.SetReadLimit(512)
	for {
		message := Message{}
		err := client.Connection.ReadJSON(&message)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("msg in %s\n", message.Header)
		client.handleMessage(&server, message) // hmm?

	}
}

func (client *Client) WriteStream(server Server) {

	for {

		select {

		case msg := <-client.Sender:
			err := client.Connection.WriteJSON(msg)
			if err != nil {
				log.Println(err)
			}
			log.Printf("Sent message (header: %s)\n", msg.Header)
			break
		}
	}
}

func (client *Client) SendMessage(message Message) {
	client.Sender <- message
}

// {"Header": "JOIN", "Body": {"username":"Adil"}}
// {"Header":"Info", "Body":{}}
