package main

import (
	"../api"
	. "../data_access"
	"github.com/wonderivan/logger"
	"strconv"
	"strings"
	"time"
)

var ClientsMap = make(map[int32]*api.WsClient)

func GetMonitorLiveRooms() []int32 {
	val, err := RedisClient.Get("LT_MONITOR_LIVE_ROOMS").Result()
	if err != nil {
		return nil
	}
	liveRoomStrs := strings.Split(val, "$")
	returnDataIndex := 0
	returnData := make([]int32, len(liveRoomStrs))
	for _, roomIdStr := range liveRoomStrs {
		intValue, err := strconv.ParseInt(roomIdStr, 10, 32)
		if err == nil {
			returnData[returnDataIndex] = int32(intValue)
			returnDataIndex++
		}
	}
	return returnData
}

func main() {
	logger.Info("Starting ws source proc.")
	time.Sleep(100 * time.Millisecond)

	liveRooms := GetMonitorLiveRooms()
	if liveRooms == nil {
		logger.Error("Cannot get monitor live rooms. ")
		return
	}
	logger.Info("Result: ", liveRooms[:100])

	for _, roomId := range liveRooms[:3000] {
		clientPointer := api.CreateWsConnection(roomId)
		ClientsMap[roomId] = clientPointer
	}

	for {
		time.Sleep(10 * time.Second)

		activeClientsCount := 0
		for _, v := range ClientsMap {
			if !v.IsClosed {
				activeClientsCount++
			}
		}
		logger.Warn("Ws client count: %d", activeClientsCount)
	}
}
