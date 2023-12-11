package server

import (
	"TestNationalSystems/Message"
	"TestNationalSystems/internal/loging"
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

type ClientInfo struct {
	IdConnection     int
	IdClient         int
	CurrentMessageId int
	conn             *websocket.Conn
}

var (
	upgrader = websocket.Upgrader{}
	clients  = make(map[int]*ClientInfo)
	mutex    sync.Mutex
	counter  = 1
)

// ServeWs  handles websocket requests from the peer.
func ServeWs(w http.ResponseWriter, r *http.Request, logger *loging.Logger) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer conn.Close()

	client, err := ProcessClientInfo(conn)
	if err != nil {
		logger.Error().Err(err)
		return
	} else {
		//increment counter bcz client already register
		counter++
	}

	// Send ping every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	//reading
	for {
		select {
		case <-ticker.C:
			// Do ping request for every tick
			client.CurrentMessageId++
			msg := Message.RequestMessage{
				IDMessage: client.CurrentMessageId,
				Type:      Message.RequestPingStr,
			}
			err := conn.WriteJSON(msg)
			logger.Info().Msgf("S:%d:%d:%v", client.IdConnection, client.IdClient, msg)
			if err != nil {
				logger.Error().Err(err)
				mutex.Lock()
				delete(clients, client.IdConnection)
				mutex.Unlock()
				conn.Close()
				return
			}
		default:
			_, data, err := conn.ReadMessage()
			if err != nil {
				logger.Error().Err(err)
				mutex.Lock()
				delete(clients, client.IdConnection)
				mutex.Unlock()
				conn.Close()
				return
			}
			//we can use any type of msg,we need only type for processing
			var msg Message.BaseMessage
			err = json.Unmarshal(data, &msg)
			if err != nil {
				logger.Error().Err(err)
				continue
			}
			logger.Info().Msgf("R:%d:%d:%v", client.IdConnection, client.IdClient, msg)
			switch msg.Type {
			case Message.RequestMessageStr:
				var mes Message.RequestMessage
				err := json.Unmarshal(data, &mes)
				if err != nil {
					return
				}
				// change id message for current user identification
				msg.IDMessage = client.IdClient
				mutex.Lock()
				for id, cl := range clients {
					if id != client.IdClient {
						err = cl.conn.WriteJSON(msg)
						logger.Info().Msgf("S:%d:%d:%v", cl.IdConnection, cl.IdClient, msg)
						if err != nil {
							logger.Error().Err(err)
							// can delete client if we cant send message to he
							delete(clients, id)
						}
					}
				}
				mutex.Unlock()
				//Send response for first client
				msg := Message.ResponseMessage{
					IDMessage: client.CurrentMessageId,
					Type:      Message.ResponseMessageStr,
					Result:    "true",
				}
				err = client.conn.WriteJSON(msg)
				if err != nil {
					return
				}
				logger.Info().Msgf("S:%d:%d:%v", client.IdConnection, client.IdClient, msg)
				client.CurrentMessageId++
			case Message.RequestPingStr:
				//Server cant recieve that message,just log error
				logger.Error().Err(errors.New("server cant receive req ping"))
			case Message.RequestInfoStr:
				logger.Error().Err(errors.New("server cant receive req info"))
			case Message.ResponsePingStr:
				//Nothing to do?
				client.CurrentMessageId++
			case Message.ResponseInfoStr:
				logger.Error().Err(errors.New("duplicate,we can receive only one resp info from client"))
			case Message.ResponseMessageStr:
				logger.Error().Err(errors.New("server cant receive response mess"))
			default:
				log.Println("unknown message type:", msg.Type)
			}
		}
	}
}

func ProcessClientInfo(conn *websocket.Conn) (*ClientInfo, error) {
	// Send client info request
	msg := Message.RequestMessage{
		IDMessage: 1,
		Type:      Message.RequestInfoStr,
	}
	err := conn.WriteJSON(msg)
	if err != nil {
		conn.Close()
		return nil, err
	}

	// Wait for client info response
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			return &ClientInfo{}, err
		}
		var msg Message.BaseMessage
		if err = json.Unmarshal(data, &msg); err != nil {
			conn.Close()
			return &ClientInfo{}, err
		}
		if msg.Type == Message.ResponseInfoStr {
			var mes Message.ResponseMessage
			err := json.Unmarshal(data, &mes)
			if err != nil {
				return nil, err
			}
			var info ClientInfo
			err = json.Unmarshal([]byte(mes.Result), &info)
			if err != nil {
				//log.Println("decode info_data:", err)
				conn.Close()
				return nil, err
			}

			mutex.Lock()
			cl := &ClientInfo{
				IdConnection:     counter,
				IdClient:         counter,
				CurrentMessageId: 1,
				conn:             conn,
			}
			clients[counter] = cl
			counter++
			mutex.Unlock()
			return cl, nil
		}
	}

}
