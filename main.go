package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func main() {
	e := echo.New()
	e.Logger.SetLevel(log.INFO)

	deviceManager := NewDeviceManager()
	if err := deviceManager.LoadFromFile("deviceid.csv"); err != nil {
		e.Logger.Fatal(err)
	}

	api := e.Group("/api/v1")

	api.POST("/devices/:device_id/heartbeat", HeartbeatHandler(deviceManager))
	api.POST("/devices/:device_id/stats", UploadStatsHandler(deviceManager))
	api.GET("/devices/:device_id/stats", GetStatsHandler(deviceManager))

	e.Logger.Fatal(e.Start(":6733"))
}
