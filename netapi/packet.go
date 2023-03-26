package netapi

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/gorilla/websocket"
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

	// device, ok := G_device
	// if !ok {
	// 	log.Println("device not exists in gin context")
	// }

	log.Println(G_device)
	if G_device == "" {
		log.Println("device not set")
	}

	// 打开设备
	handle, err := pcap.OpenLive("\\Device\\NPF_{6859611A-A240-4DE0-8FC1-B7645ABB5A84}", 65535, true, -1*time.Second)
	if err != nil {
		log.Println(err)
	}
	defer handle.Close()

	// 设置过滤器
	filter := ""
	err = handle.SetBPFFilter(filter)
	if err != nil {
		log.Println(err)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	go HandleWsQuit(ws)

	for p := range packetSource.Packets() {
		// fmt.Println(p.LinkLayer())
		// fmt.Println("----- ----- ----- -----")
		// fmt.Println(p)
		// // fmt.Println(p.NetworkLayer().LayerContents())
		// // fmt.Println(p.NetworkLayer().LayerPayload())
		// fmt.Println(p.TransportLayer().LayerContents())
		// fmt.Println(p.TransportLayer().LayerPayload())

		// GetPacketInfo(p)
		// fmt.Println("----- ----- ----- -----")

		// data, err := json.Marshal(p)
		// if err != nil {
		// 	log.Println(err)
		// }

		// fmt.Println(p.String())

		ws.WriteJSON(p.String())

	}

	// log.Println(packetSource)

	// device, ok := c.Get("device")
	// fmt.Println(device, ok)

	// fmt.Println(G_device)
	select {}
}

func HandleWsQuit(ws *websocket.Conn) {
	for {
		_, data, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
		}
		// fmt.Println(string(data))

		if string(data) == "quit" {
			ws.Close()
			return
		}
	}
}
