package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDeviceManager_IsValid(t *testing.T) {
	service := NewDeviceManager()
	service.devices["valid-device"] = &DeviceData{
		Heartbeats:  make([]time.Time, 0),
		UploadTimes: make([]int, 0),
	}

	assert.True(t, service.IsValid("valid-device"))
	assert.False(t, service.IsValid("invalid-device"))
}

func TestDeviceManager_RecordHeartbeat(t *testing.T) {
	service := NewDeviceManager()
	service.devices["device-1"] = &DeviceData{
		Heartbeats:  make([]time.Time, 0),
		UploadTimes: make([]int, 0),
	}

	t1, _ := time.Parse(time.RFC3339, "2025-11-19T10:00:00Z")
	t2, _ := time.Parse(time.RFC3339, "2025-11-19T10:01:00Z")

	service.RecordHeartbeat("device-1", t1)
	service.RecordHeartbeat("device-1", t2)

	data := service.devices["device-1"]
	assert.Equal(t, 2, len(data.Heartbeats))
	assert.Equal(t, t1, *data.FirstHeartbeat)
	assert.Equal(t, t2, *data.LatestHeartbeat)
}

func TestDeviceManager_RecordUploadTime(t *testing.T) {
	service := NewDeviceManager()
	service.devices["device-1"] = &DeviceData{
		Heartbeats:  make([]time.Time, 0),
		UploadTimes: make([]int, 0),
	}

	t1, _ := time.Parse(time.RFC3339, "2025-11-19T10:00:00Z")

	service.RecordUploadTime("device-1", t1, 1000000)
	service.RecordUploadTime("device-1", t1, 1000000)

	data := service.devices["device-1"]
	assert.Equal(t, 2, data.UploadCount)
}

func TestDeviceManager_CalculateUptime_Success(t *testing.T) {
	service := NewDeviceManager()
	service.devices["device-1"] = &DeviceData{
		Heartbeats:  make([]time.Time, 0),
		UploadTimes: make([]int, 0),
	}

	t1, _ := time.Parse(time.RFC3339, "2025-11-19T10:00:00Z")
	t2, _ := time.Parse(time.RFC3339, "2025-11-19T10:05:00Z")

	service.RecordHeartbeat("device-1", t1)
	service.RecordHeartbeat("device-1", t2)

	uptime, err := service.CalculateUptime("device-1")

	assert.NoError(t, err)
	assert.Equal(t, 40.0, uptime)
}

func TestDeviceManager_CalculateUptime_InsufficientData(t *testing.T) {
	service := NewDeviceManager()
	service.devices["device-1"] = &DeviceData{
		Heartbeats:  make([]time.Time, 0),
		UploadTimes: make([]int, 0),
	}

	t1, _ := time.Parse(time.RFC3339, "2025-11-19T10:00:00Z")
	service.RecordHeartbeat("device-1", t1)

	_, err := service.CalculateUptime("device-1")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient data")
}

func TestDeviceManager_CalculateUptime_DeviceNotFound(t *testing.T) {
	service := NewDeviceManager()

	_, err := service.CalculateUptime("nonexistent")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "device not found")
}

func TestDeviceManager_CalculateAverageUploadTime_Success(t *testing.T) {
	service := NewDeviceManager()
	service.devices["device-1"] = &DeviceData{
		Heartbeats:  make([]time.Time, 0),
		UploadTimes: make([]int, 0),
	}

	t1, _ := time.Parse(time.RFC3339, "2025-11-19T10:00:00Z")

	service.RecordUploadTime("device-1", t1, 310000000000)
	service.RecordUploadTime("device-1", t1, 190000000000)

	avgTime, err := service.CalculateAverageUploadTime("device-1")

	assert.NoError(t, err)
	assert.Equal(t, 250*time.Second, avgTime)
}

func TestDeviceManager_CalculateAverageUploadTime_NoData(t *testing.T) {
	service := NewDeviceManager()
	service.devices["device-1"] = &DeviceData{
		Heartbeats:  make([]time.Time, 0),
		UploadTimes: make([]int, 0),
	}

	_, err := service.CalculateAverageUploadTime("device-1")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no upload data available")
}

func TestDeviceManager_CalculateAverageUploadTime_DeviceNotFound(t *testing.T) {
	service := NewDeviceManager()

	_, err := service.CalculateAverageUploadTime("nonexistent")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "device not found")
}
