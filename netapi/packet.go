package netapi

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	// "github.com/google/gopacket"
	"github.com/gorilla/websocket"

	// "github.com/google/gopacket/layers"

	"github.com/google/gopacket/pcap"
)

// 使用默认配置
var upgrader = websocket.Upgrader{}

func GetPakcet(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
	}
	defer ws.Close()

	// packet_chan := make(chan, 1024)

	device, ok := c.Get("device")
	if !ok {
		log.Println("device not exists in gin context")
	}

	// 打开设备
	handle, err := pcap.OpenLive(device.(string), 65535, true, -1*time.Second)
	if err != nil {
		log.Println(err)
	}
	defer handle.Close()

}
