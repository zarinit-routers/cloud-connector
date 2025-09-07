package server

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zarinit-routers/cloud-connector/server/handlers"
	"github.com/zarinit-routers/middleware/auth"
	"github.com/zarinit-routers/middleware/cors"
)

const (
	ENV_PORT = "PORT"
)

func getPort() (int, error) {
	str := os.Getenv(ENV_PORT)
	if str == "" {
		return 0, fmt.Errorf("env variable %q is not set", ENV_PORT)
	}

	port, err := strconv.Atoi(str)
	if err != nil {
		return port, fmt.Errorf("failed parse port: %s", err)
	}
	if port <= 0 {
		return port, fmt.Errorf("port value is invalid (%d)", port)
	}
	return port, nil
}

func getAddr() (string, error) {
	port, err := getPort()
	if err != nil {
		return "", fmt.Errorf("bad port specified: %s", err)
	}
	return fmt.Sprintf(":%d", port), nil
}
func Serve() error {
	addr, err := getAddr()
	if err != nil {
		return err
	}

	srv := gin.Default()
	srv.Use(cors.Middleware([]string{
		"http://localhost:3000", // For development purposes
	}))
	api := srv.Group("/api/clients")
	api.GET("/", auth.Middleware(), handlers.GetClientsHandler())
	api.GET("/:id", auth.Middleware(), handlers.GetSingleClientHandler())
	api.POST("/tags/add", auth.Middleware())
	api.POST("/tags/remove", auth.Middleware())
	return srv.Run(addr)
}
