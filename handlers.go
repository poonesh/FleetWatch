package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func HeartbeatHandler(deviceManager *DeviceManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		deviceID := c.Param("device_id")

		if !deviceManager.IsValid(deviceID) {
			return c.JSON(http.StatusNotFound, ErrorResponse{
				Msg: "Device not found",
			})
		}

		var req HeartbeatRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, ErrorResponse{
				Msg: "invalid request body",
			})
		}

		if req.SentAt == "" {
			return c.JSON(http.StatusBadRequest, ErrorResponse{
				Msg: "sent_at is required",
			})
		}

		sentAt, err := time.Parse(time.RFC3339, req.SentAt)
		if err != nil {
			return c.JSON(http.StatusBadRequest, ErrorResponse{
				Msg: "sent_at must be a valid RFC3339 date-time",
			})
		}

		deviceManager.RecordHeartbeat(deviceID, sentAt)

		c.Logger().Infof("Heartbeat received: device_id=%s sent_at=%s", deviceID, req.SentAt)

		return c.NoContent(http.StatusNoContent)
	}
}

func UploadStatsHandler(deviceManager *DeviceManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		deviceID := c.Param("device_id")

		if !deviceManager.IsValid(deviceID) {
			return c.JSON(http.StatusNotFound, ErrorResponse{
				Msg: "Device not found",
			})
		}

		var req UploadStatsRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, ErrorResponse{
				Msg: "invalid request body",
			})
		}

		if req.SentAt == "" {
			return c.JSON(http.StatusBadRequest, ErrorResponse{
				Msg: "sent_at is required",
			})
		}

		sentAt, err := time.Parse(time.RFC3339, req.SentAt)
		if err != nil {
			return c.JSON(http.StatusBadRequest, ErrorResponse{
				Msg: "sent_at must be a valid RFC3339 date-time",
			})
		}

		if req.UploadTime == 0 {
			return c.JSON(http.StatusBadRequest, ErrorResponse{
				Msg: "upload_time is required",
			})
		}

		deviceManager.RecordUploadTime(deviceID, sentAt, req.UploadTime)

		c.Logger().Infof("Stats received: device_id=%s sent_at=%s upload_time=%d", deviceID, req.SentAt, req.UploadTime)

		return c.NoContent(http.StatusNoContent)
	}
}

func GetStatsHandler(deviceManager *DeviceManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		deviceID := c.Param("device_id")

		if !deviceManager.IsValid(deviceID) {
			return c.JSON(http.StatusNotFound, ErrorResponse{
				Msg: "Device not found",
			})
		}

		uptime, err := deviceManager.CalculateUptime(deviceID)
		if err != nil {
			return c.JSON(http.StatusNoContent, nil)
		}

		avgUploadTime, err := deviceManager.CalculateAverageUploadTime(deviceID)
		if err != nil {
			return c.JSON(http.StatusNoContent, nil)
		}

		response := GetDeviceStatsResponse{
			AverageUploadTime: avgUploadTime.String(),
			Uptime:            uptime,
		}

		c.Logger().Infof("Stats sent: device_id=%s avg_upload_time=%s uptime=%.2f", deviceID, response.AverageUploadTime, response.Uptime)

		return c.JSON(http.StatusOK, response)
	}
}
