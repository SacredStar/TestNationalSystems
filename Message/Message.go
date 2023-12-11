package Message

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

type BaseMessage struct {
	IDMessage int    `json:"id"`
	Type      string `json:"type"`
	Data      string `json:"data"`
}

type RequestMessage struct {
	IDMessage int    `json:"id"`
	Type      string `json:"type"`
	Data      string `json:"data"`
}

type ResponseMessage struct {
	IDMessage int    `json:"id"`
	Type      string `json:"type"`
	Result    string `json:"result"`
}

type ClientInfo struct {
	Name string `json:"name"`
}

const (
	RequestInfoStr     = "  function: info"
	RequestMessageStr  = "function: message"
	RequestPingStr     = "function: ping"
	ResponseInfoStr    = "response: info"
	ResponseMessageStr = " response: message"
	ResponsePingStr    = "response: ping"
)

func ReadResponseMessageServer(conn *websocket.Conn, counter int, print bool) (ResponseMessage, error) {
	_, message, err := conn.ReadMessage()
	if err != nil {
		return ResponseMessage{}, err
	}
	var msg ResponseMessage
	if err = json.Unmarshal(message, &msg); err != nil {
		conn.Close()
		return ResponseMessage{}, err
	}
	if print {
		log.Info().Msgf("R: ID_Conn = %d : ID_CL = %d  MSG =%s", msg.IDMessage, msg.IDMessage, msg)
	}
	return msg, nil
}
