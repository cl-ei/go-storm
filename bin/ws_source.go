package main

import (
	"../config"
	"fmt"
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
	time.Sleep(5 * time.Second)

	fmt.Print("Starting ws source proc. \n\tConfig: ", config.CONFIG)
	h := Handler{}
	h.init()
	println(len(h.clients))

	println("Shutdown !")
}
