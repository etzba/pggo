package main

import (
	"github.com/etzba/pggo/pkg/logger"
	"github.com/etzba/pggo/server"
)

func main() {
	logger := logger.New()
	server := server.New(logger, ":8080")
	if err := server.Run(); err != nil {
		panic(err)
	}
}
