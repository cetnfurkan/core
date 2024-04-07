package server

import (
	"fmt"
	"log"

	"github.com/cetnfurkan/core/config"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type (
	echoServer struct {
		app     *echo.Echo
		cfg     *config.Server
		options []echoServerOption
	}

	echoServerOption func(*echoServer) error
)

func WithControllers(callback func(app *echo.Echo)) echoServerOption {
	return func(server *echoServer) error {
		callback(server.app)
		return nil
	}
}

func WithMiddlewares(middlewares ...echo.MiddlewareFunc) echoServerOption {
	return func(server *echoServer) error {
		server.app.Use(middlewares...)
		return nil
	}
}

func WithLogger(logger echo.Logger) echoServerOption {
	return func(server *echoServer) error {
		server.app.Logger = logger
		return nil
	}
}

func WithSwaggerController() echoServerOption {
	return func(server *echoServer) error {
		server.app.GET("/docs/*", echoSwagger.WrapHandler)
		return nil
	}
}

// NewEchoServer creates a new echo server instance.
//
// It takes a config instance, a database instance and a cache instance
// and returns a new server interface instance.
//
// It will panic if it fails to create a new echo server instance.
func NewEchoServer(cfg *config.Server, opts ...echoServerOption) Server {
	return &echoServer{
		app:     echo.New(),
		cfg:     cfg,
		options: opts,
	}
}

func (server *echoServer) Start() {
	go server.start()
}

func (server *echoServer) start() {

	// Apply options.
	for _, opt := range server.options {
		if err := opt(server); err != nil {
			log.Fatal("Failed to apply option to echo server: ", err)
		}
	}

	serverUrl := fmt.Sprintf(":%d", server.cfg.Port)

	// Start server with port from config.
	server.app.Logger.Fatal(server.app.Start(serverUrl))
}

func (server *echoServer) Stop() error {
	return server.app.Close()
}
