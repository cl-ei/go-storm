package main

import (
	"../api"
	_ "../config"
	"github.com/wonderivan/logger"
	"time"
)

var ClientsMap = make(map[int]*api.WsClient)

func main() {
	logger.Info("Starting ws source proc.")
	time.Sleep(100 * time.Millisecond)

	roomId := 11382758

	clientPointer := api.CreateWsConnection(roomId)
	ClientsMap[roomId] = clientPointer

	for {
		time.Sleep(3 * time.Second)
		// println("wait.")
	}
	println("Shutdown !")
}
