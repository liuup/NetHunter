package router

import (
	"github.com/gin-gonic/gin"

	"example/netapi"
)

func GetRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		// 获取所有设备
		api.GET("/device", netapi.GetAllDevices)
		// 选择设备
		api.POST("/device/choose", netapi.ChooseDevice)

		api.GET("/packet", netapi.GetPakcet)

	}
	return r
}
