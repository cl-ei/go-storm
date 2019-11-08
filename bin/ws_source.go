package main

import (
	_ "../config"
	"github.com/wonderivan/logger"
	"time"
)

type WsClient struct {
}

type Handler struct {
	clients map[int]WsClient
}

func (h *Handler) init() {
	println("Start ws.")
}

func main() {
	logger.Info("Starting ws source proc.")
	time.Sleep(100 * time.Millisecond)

	h := Handler{}
	h.init()
	println(len(h.clients))

	println("Shutdown !")
}
