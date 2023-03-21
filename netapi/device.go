package netapi

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/gopacket/pcap"
)

// 查询所有网络设备
func GetAllDevices(c *gin.Context) {
	// 找到所有网络设备
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"data": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": devices,
	})
}

// 选择网络设备
func ChooseDevice(c *gin.Context) {
	device := c.PostForm("device")

	// 直接设置到context中
	if device != "" {
		c.Set("device", device)
	}

	if v, ok := c.Get("device"); ok {
		log.Println(v)
	} else {
		log.Println("device not set!")
	}
}
