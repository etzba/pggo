package main

import (
	"github.com/etzba/pggo/server"
)

func main() {
	server := server.New(":8080")
	if err := server.Run(); err != nil {
		panic(err)
	}
}
