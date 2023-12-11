package main

import (
	"TestNationalSystems/internal/application"
	"TestNationalSystems/internal/cli"
	"TestNationalSystems/internal/loging"
	"log"
)

//var addr = flag.String("addr", ":8080", "http service address")
//usage : -c=true -address=127.0.0.1 - client
//base run without params start a server 127.0.0.1 7623

func main() {
	log.Println("Initialising  application...")
	log.Println("Getting logger of  application...")
	logger := loging.StartLog(-1, "./logs/log.txt")
	logger.Info().Msg("Collecting cli info...")
	role, address, err := cli.GetCliInfo()
	if err != nil {
		logger.Error().Err(err)
		return
	}
	logger.Info().Msgf("Create application instance:role - %s, address - %s", role, address)
	app, err := application.NewApplication(logger, role, address)
	if err != nil {
		logger.Fatal().Err(err).Msg("Application can't start")
	}
	app.Run()
}
