package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sync"
	"time"
)

type DeviceData struct {
	Heartbeats      []time.Time
	UploadTimes     []int
	FirstHeartbeat  *time.Time
	LatestHeartbeat *time.Time
	TotalUploadTime int64
	UploadCount     int
}

type DeviceManager struct {
	devices map[string]*DeviceData
	mu      sync.RWMutex
}

func NewDeviceManager() *DeviceManager {
	return &DeviceManager{
		devices: make(map[string]*DeviceData),
	}
}

func (s *DeviceManager) LoadFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open device file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV: %w", err)
	}

	if len(records) < 2 {
		return fmt.Errorf("CSV file must have at least a header and one device")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for i, record := range records {
		if i == 0 {
			continue
		}
		if len(record) > 0 && record[0] != "" {
			s.devices[record[0]] = &DeviceData{
				Heartbeats:  make([]time.Time, 0),
				UploadTimes: make([]int, 0),
			}
		}
	}

	return nil
}

func (s *DeviceManager) IsValid(deviceID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.devices[deviceID]
	return exists
}

func (s *DeviceManager) RecordHeartbeat(deviceID string, sentAt time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if data, exists := s.devices[deviceID]; exists {
		data.Heartbeats = append(data.Heartbeats, sentAt)

		if data.FirstHeartbeat == nil || sentAt.Before(*data.FirstHeartbeat) {
			data.FirstHeartbeat = &sentAt
		}

		if data.LatestHeartbeat == nil || sentAt.After(*data.LatestHeartbeat) {
			data.LatestHeartbeat = &sentAt
		}
	}
}

func (s *DeviceManager) RecordUploadTime(deviceID string, sentAt time.Time, uploadTime int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if data, exists := s.devices[deviceID]; exists {
		data.UploadTimes = append(data.UploadTimes, uploadTime)
		data.TotalUploadTime += int64(uploadTime)
		data.UploadCount++
	}
}

func (s *DeviceManager) CalculateUptime(deviceID string) (float64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, exists := s.devices[deviceID]
	if !exists {
		return 0, fmt.Errorf("device not found")
	}

	if data.FirstHeartbeat == nil || data.LatestHeartbeat == nil || len(data.Heartbeats) < 2 {
		return 0, fmt.Errorf("insufficient data: need at least 2 heartbeats")
	}

	totalMinutes := data.LatestHeartbeat.Sub(*data.FirstHeartbeat).Minutes()
	if totalMinutes == 0 {
		return 100.0, nil
	}

	heartbeatCount := float64(len(data.Heartbeats))
	uptime := (heartbeatCount / totalMinutes) * 100

	return uptime, nil
}

func (s *DeviceManager) CalculateAverageUploadTime(deviceID string) (time.Duration, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, exists := s.devices[deviceID]
	if !exists {
		return 0, fmt.Errorf("device not found")
	}

	if data.UploadCount == 0 {
		return 0, fmt.Errorf("no upload data available")
	}

	avgNanoseconds := data.TotalUploadTime / int64(data.UploadCount)
	return time.Duration(avgNanoseconds), nil
}
