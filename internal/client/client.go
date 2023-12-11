package client

import (
	"TestNationalSystems/internal/Message"
	"TestNationalSystems/internal/loging"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net/url"
	"os"
	"time"
)

var (
	counter = 1
)

func NewClient(address string, log *loging.Logger) {
	u := url.URL{Scheme: "ws", Host: address, Path: "/"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal().Err(err)
	}
	defer conn.Close()
	var name string
	fmt.Print("Enter your name: ")
	_, err = fmt.Scanln(&name)
	if err != nil {
		return
	}
	// for 1st we must send info
	clientInfo := Message.ClientInfo{Name: name}
	message := Message.ResponseMessage{
		IDMessage: counter,
		Type:      Message.ResponseInfoStr,
		Result:    fmt.Sprintf(`{"name": "%s"}`, clientInfo.Name),
	}
	err = conn.WriteJSON(message)
	counter++
	if err != nil {
		log.Error().Err(err)
		return
	}

	// Handle incoming messages from the server
	go func() {
		for {
			// bcz base struct of Messages is similar we can just use ReqRead, we need only type
			_, data, err := conn.ReadMessage()
			if err != nil {
				log.Error().Err(err)
			}
			var msg Message.BaseMessage
			err = json.Unmarshal(data, &msg)
			if err != nil {
				return
			}
			if msg.Type == Message.RequestInfoStr {
				clientInfo := Message.ClientInfo{Name: name}
				message := Message.ResponseMessage{
					IDMessage: msg.IDMessage,
					Type:      Message.ResponseInfoStr,
					Result:    clientInfo.Name,
				}
				err := conn.WriteJSON(message)
				counter++
				if err != nil {
					log.Error().Err(err)
					return
				}
			}
			if msg.Type == Message.RequestPingStr {
				message := Message.ResponseMessage{
					IDMessage: msg.IDMessage,
					Type:      Message.ResponsePingStr,
					Result:    fmt.Sprintf(`timestamp:%s`, time.Now()),
				}
				err := conn.WriteJSON(message)
				if err != nil {
					log.Error().Err(err)
					return
				}
				counter++
			}
			if msg.Type == Message.RequestMessageStr {
				var mes Message.RequestMessage
				if err = json.Unmarshal(data, &mes); err != nil {
					log.Error().Err(err)
				}
				log.Info().Msgf("%d - %s", msg.IDMessage, mes.Data)
			}
			if msg.Type == Message.ResponseInfoStr {
				log.Error().Err(errors.New("client cant receive response info, smth went wrong"))
				continue
			}
			if msg.Type == Message.ResponsePingStr {
				log.Error().Err(errors.New("client cant receive response ping, smth went wrong"))
				continue
			}
			if msg.Type == Message.ResponseMessageStr {
				log.Info().Msg("Received response, nothing to do with that")
				counter++
			}
		}
	}()

	for {
		reader := bufio.NewReader(os.Stdin)
		text, _, err := reader.ReadLine()

		if err != nil {
			log.Error().Err(err)
			continue
		}
		msg := Message.RequestMessage{
			IDMessage: counter,
			Type:      Message.RequestMessageStr,
			Data:      fmt.Sprintf(fmt.Sprintf(`{"data": "%s"}`, text)),
		}
		err = conn.WriteJSON(msg)
		if err != nil {
			log.Error().Err(err)
			continue
		}
	}

}
