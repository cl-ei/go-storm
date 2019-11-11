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
var HeartBeatCoroutineStarted = false

func StartHeartBeatGoroutine() {
	if HeartBeatCoroutineStarted {
		return
	}
	HeartBeatCoroutineStarted = true

	go func() {
		var (
			startTime  int64
			costTime   int64
			sleepTime  int64
			flushCount int32
		)
		for {
			startTime = time.Now().UnixNano()
			flushCount = 0

			for _, client := range ClientsMap {
				if client.IsClosed {
					continue
				}

				if err := client.SendHeartBeatPackage(); err != nil {
					client.Close()
					continue
				}
				flushCount++
			}
			costTime = time.Now().UnixNano() - startTime
			if costTime < 20*1e9 {
				sleepTime = 20*1e9 - costTime
			} else {
				sleepTime = 0
			}
			logger.Info(
				"HeartBeat cost: %.3f, sleep time: %.2f, flushCount: %d",
				float64(costTime)/1e9, float64(sleepTime)/1e9, flushCount,
			)
			time.Sleep(time.Duration(sleepTime) * time.Nanosecond)
		}
	}()
}

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

func CreateConnections() {
	liveRooms := GetMonitorLiveRooms()
	if liveRooms == nil {
		logger.Error("Cannot get monitor live rooms. ")
		return
	}
	logger.Info("Result: ", liveRooms[:100])
	for index, roomId := range liveRooms {
		go func() {
			clientPointer := api.CreateWsConnection(roomId)
			ClientsMap[roomId] = clientPointer
		}()

		if index%500 == 0 {
			time.Sleep(1 * time.Second)
		}
	}
}

func main() {
	logger.Info("Starting ws source proc.")
	time.Sleep(100 * time.Millisecond)
	StartHeartBeatGoroutine()
	CreateConnections()

	for {
		time.Sleep(10 * time.Second)

		activeClientsCount := 0
		for _, v := range ClientsMap {
			if !v.IsClosed {
				activeClientsCount++
			}
		}
		logger.Warn("Ws client count: %d/%d", activeClientsCount, len(ClientsMap))
	}
}
