package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func main() {
	e := echo.New()
	e.Logger.SetLevel(log.INFO)

	deviceService := NewDeviceService()
	if err := deviceService.LoadFromFile("deviceid.csv"); err != nil {
		e.Logger.Fatal(err)
	}

	api := e.Group("/api/v1")

	api.POST("/devices/:device_id/heartbeat", HeartbeatHandler(deviceService))
	api.POST("/devices/:device_id/stats", UploadStatsHandler(deviceService))
	api.GET("/devices/:device_id/stats", GetStatsHandler(deviceService))

	e.Logger.Fatal(e.Start(":6733"))
}
