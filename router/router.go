package router

import (
	"github.com/gin-gonic/gin"

	"example/netapi"
)

func GetRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.GET("/device", netapi.GetAllDevices)
		api.POST("/device/choose", netapi.ChooseDevice)
	}
	return r
}
