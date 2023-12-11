package application

import (
	"TestNationalSystems/internal/client"
	"TestNationalSystems/internal/loging"
	"TestNationalSystems/internal/server"
	"net/http"
)

type Application struct {
	logger  *loging.Logger
	Server  *http.Server
	role    string
	address string
}

func NewApplication(logger *loging.Logger, role string, address string) (Application, error) {
	logger.Info().Msg("Creating application...")
	return Application{
		logger:  logger,
		role:    role,
		address: address,
		Server:  &http.Server{Addr: address},
	}, nil
}

func (app *Application) Run() {
	if app.role == "client" {
		app.logger.Info().Msgf("Start client connected to %s", app.address)
		client.NewClient(app.address, app.logger)
	} else {
		app.logger.Info().Msgf("Start server on %s", app.address)
		app.StartRouting()
	}
}

func (app *Application) StartRouting() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		server.ServeWs(w, r, app.logger)
	})
	err := app.Server.ListenAndServe()
	if err != nil {
		app.logger.Fatal().Msg("can't start server")
	}
}
