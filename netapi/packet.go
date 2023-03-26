package netapi

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/gorilla/websocket"
)

var G_packets = []gopacket.Packet{}

// 使用默认配置
var upgrader = websocket.Upgrader{}

func GetPakcet(c *gin.Context) {
	// 网络过滤器
	filter := c.DefaultQuery("filter", "")

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
	// handle, err := pcap.OpenLive("\\Device\\NPF_{6859611A-A240-4DE0-8FC1-B7645ABB5A84}", 65535, true, -1*time.Second)
	handle, err := pcap.OpenLive(G_device, 65535, true, -1*time.Second)
	if err != nil {
		log.Println(err)
	}
	defer handle.Close()

	// 设置过滤器
	err = handle.SetBPFFilter(filter)
	if err != nil {
		log.Println(err)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	// 处理ws退出
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

		// 先暂存起来，为实现流追踪做准备
		G_packets = append(G_packets, p)

		err = ws.WriteJSON(p.String())
		if err != nil {
			log.Println(err)
			break
		}
	}

	// log.Println(packetSource)

	// device, ok := c.Get("device")
	// fmt.Println(device, ok)

	// fmt.Println(G_device)
	select {}
}

func GetPacketTrace(c *gin.Context) {
	// 通过四元组来进行筛选
	srcip := c.DefaultQuery("srcip", "")
	dstip := c.DefaultQuery("dstip", "")
	// srcport := c.DefaultQuery("srcport", "")
	// dstport := c.DefaultQuery("dstport", "")

	// fmt.Println(G_packets)
	// fmt.Println(srcip, dstip)

	ans := []string{}

	for _, p := range G_packets {
		// 实现ipv4追踪
		if ipv4 := p.Layer(layers.LayerTypeIPv4); ipv4 != nil {
			ipv4Packet, _ := ipv4.(*layers.IPv4)

			if (ip2String(ipv4Packet.SrcIP) == srcip && ip2String(ipv4Packet.DstIP) == dstip) ||
				(ip2String(ipv4Packet.SrcIP) == dstip && ip2String(ipv4Packet.DstIP) == srcip) {
				ans = append(ans, p.String())
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": ans,
	})
}

func HandleWsQuit(ws *websocket.Conn) {
	for {
		_, data, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// fmt.Println(string(data))

		if string(data) == "quit" {
			ws.Close()
			return
		}
	}
}

func ip2String(ip []byte) (ans string) {
	for _, b := range ip {
		ans += strconv.Itoa(int(b))
		ans += "."
	}
	// 去除最后的点
	ans = ans[:len(ans)-1]
	return
}
