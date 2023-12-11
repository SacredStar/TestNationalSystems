package cli

import (
	"errors"
	"flag"
	"fmt"
)

func GetCliInfo() (role string, addressPort string, err error) {
	isClient := flag.Bool("c", false, "run as client? Default run as a server")
	address := flag.String("address", "127.0.0.1", "network interface or server address")
	port := flag.Int("port", 7623, "server port")
	flag.Parse()
	if *address == "" {
		return "", "", errors.New("requires an address/interface to start")
	}
	addressPort = fmt.Sprintf("%s:%d", *address, *port)
	if *isClient {
		return "client", addressPort, nil
	} else {
		return "server", addressPort, nil
	}
}
