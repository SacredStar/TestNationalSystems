package Message

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
