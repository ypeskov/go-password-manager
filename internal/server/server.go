package server

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"ypeskov/go-orgfin/internal/logger"
	"ypeskov/go-orgfin/internal/routes"

	_ "github.com/joho/godotenv/autoload"

	"ypeskov/go-orgfin/internal/config"
	"ypeskov/go-orgfin/internal/database"
)

type Server struct {
	port int
	Db   database.Service
}

var NewServer *Server

func New(cfg *config.Config, logger *logger.Logger) *http.Server {
	port, _ := strconv.Atoi(cfg.Port)
	NewServer = &Server{
		port: port,
		Db:   database.New(cfg),
	}

	routesInstance := routes.RegisterRoutes(logger)

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      routesInstance.Echo,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
