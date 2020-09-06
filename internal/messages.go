package internal

import "time"

type MessageHeader string

const (
	RegisterRequest       MessageHeader = "REGISTER"
	RegisterResponse      MessageHeader = "REGISTER_OK"
	RegisterFailure       MessageHeader = "REGISTER_FAILURE"
	Join                  MessageHeader = "JOIN"
	InfoRequest           MessageHeader = "INFO"
	InfoResponse          MessageHeader = "INFO_RESPONSE"
	Invalid               MessageHeader = "INVALID"
	JoinRoomRequest       MessageHeader = "JOIN_ROOM"
	JoinRoomResponse      MessageHeader = "JOINED_ROOM"
	NewRoomRequest        MessageHeader = "NEW_ROOM"
	NewRoomResponse       MessageHeader = "NEW_ROOM_OK"
	ListRoomUsersRequest  MessageHeader = "LIST_ROOM_USERS"
	ListRoomUsersResponse MessageHeader = "ROOM_USER_LIST"
)

type Message struct {
	Header MessageHeader
	Body   map[string]interface{} //  YUCK
}

func Registered() Message {
	return Message{
		Header: RegisterResponse,
	}
}

func Info(username string) Message {
	serverTime := time.Now().String()
	return Message{
		Header: InfoResponse,
		Body: map[string]interface{}{
			"serverTime": serverTime,
			"username":   username,
		},
	}
}

func RegisterFailed() Message {
	return Message{
		Header: RegisterFailure,
		Body: map[string]interface{}{
			"message": "Failed to register",
		},
	}
}

func InvalidMessage() Message {
	return Message{
		Header: Invalid,
	}
}

func JoinedRoom() Message {
	return Message{
		Header: JoinRoomResponse,
	}
}

func RoomUserList(users []string) Message {
	return Message{
		Header: ListRoomUsersResponse,
		Body: map[string]interface{}{
			"users": users,
		},
	}
}

func NewRoomOk(inviteCode string, roomId string) Message {
	return Message{
		Header: NewRoomResponse,
		Body: map[string]interface{}{
			"inviteCode": inviteCode,
			"roomId":     roomId,
		},
	}
}
