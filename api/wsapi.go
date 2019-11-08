package api

import (
	"bytes"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/wonderivan/logger"
	"log"
	"math"
	"math/rand"
	"net/url"
	"time"
)

const (
	PackageHeaderLength = 16
	Message             = 7
	// HeartBeat = 2
)

var HeartBeatPackage = []byte{
	3:  16,
	5:  16,
	7:  1,
	11: 2,
	15: 1,
}

type WsClient struct {
	roomId      int
	onMessage   func([]byte)
	onError     func(string)
	onConnected func(*websocket.Conn)
	onClose     func()
	conn        *websocket.Conn
	extra       string
	isClosed    bool
}

func GenPacket(act int, payload string) []byte {
	var (
		payloadBytes = []byte(payload) // payload.encode("utf-8")
		packetLen    = int32(PackageHeaderLength + len(payloadBytes))
		headers      = []byte{15: 0}
		join         = make([][]byte, 2)
	)

	headers[0] = byte((packetLen >> 24) & 0xFF)
	headers[1] = byte((packetLen >> 16) & 0xFF)
	headers[2] = byte((packetLen >> 8) & 0xFF)
	headers[3] = byte(packetLen & 0xFF)
	headers[4] = 0
	headers[5] = 16
	headers[6] = 0
	headers[7] = 1
	headers[8] = 0
	headers[9] = 0
	headers[10] = 0
	headers[11] = byte(act)
	headers[15] = 1

	join[0] = headers
	join[1] = payloadBytes

	return bytes.Join(join, []byte(""))
}

func JoinRoom(roomId int) (p []byte) {
	uid := int(1E15 + math.Floor(2E15*rand.Float64()))
	packJsonStr := fmt.Sprintf("{\"uid\":%d,\"roomid\":%d}", uid, roomId)
	return GenPacket(Message, packJsonStr)
}

func ParseMessage(m []byte, handler func([]byte)) {
	var (
		MLen           uint32
		CurrentMessage []byte
	)

	for {
		if len(m) < 16 {
			return
		}
		MLen = uint32(m[0])<<24 + uint32(m[1])<<16 + uint32(m[2])<<8 + uint32(m[3])
		CurrentMessage = m[16:MLen]
		if m[7] != 1 {
			handler(CurrentMessage)
		}
		m = m[MLen:]
	}
}

func (c *WsClient) Connect() {
	u := url.URL{
		Scheme: "ws",
		Host:   "broadcastlv.chat.bilibili.com:2244",
		Path:   "/sub",
	}
	logger.Info("connecting to %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Fatal: ", err)
		c.Close()
		return
	}

	c.conn = conn

	if c.onConnected != nil {
		c.onConnected(conn)
	} else {
		pkg := JoinRoom(c.roomId)
		if err := c.conn.WriteMessage(websocket.BinaryMessage, pkg); err != nil {
			logger.Error("Error happened when join room!")
			c.Close()
			return
		} else {
			// logger.Info("Join room success!")
		}
	}

	go c.ReadMessage()
	go c.HeartBeat()
	// logger.Info("Ws connected!")
}

func (c *WsClient) Close() {
	if c.isClosed {
		// logger.Info("Close a closed Ws Client.")
		return
	}

	c.isClosed = true
	defer func() {
		if c.onClose != nil {
			c.onClose()
		}
	}()

	if err := c.conn.Close(); err != nil {
		logger.Error("Error happened when close ws.")
	}
	logger.Info("Connection closed.")
}

func (c *WsClient) ReadMessage() {
	var (
		errMessage string
	)

	defer func() {
		c.Close()
		// logger.Info("Exit goroutine ReadMessage.")
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if c.isClosed {
				return
			}

			errMessage = fmt.Sprintf("Error happened when read msg: %s", err)
			if c.onError != nil {
				c.onError(errMessage)
			} else {
				logger.Error(errMessage)
			}
			return
		}
		ParseMessage([]byte(message), func(m []byte) {
			if c.onMessage != nil {
				c.onMessage(m)
			} else {
				logger.Info("R -> \n\t%s", m)
			}
		})
	}
}

func (c *WsClient) HeartBeat() {
	defer func() {
		c.Close()
		logger.Info("Exit goroutine HeartBeat.")
	}()
	for {
		time.Sleep(20 * time.Second)
		if err := c.conn.WriteMessage(websocket.BinaryMessage, HeartBeatPackage); err != nil {
			return
		} else {
			// logger.Info("Send heartbeat.")
		}
	}
}

func CreateWsConnection(roomId int) *WsClient {
	client := WsClient{roomId: roomId}
	client.Connect()
	return &client
}
