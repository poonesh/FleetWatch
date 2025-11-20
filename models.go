package main

type HeartbeatRequest struct {
	SentAt string `json:"sent_at"`
}

type UploadStatsRequest struct {
	SentAt     string `json:"sent_at"`
	UploadTime int    `json:"upload_time"`
}

type GetDeviceStatsResponse struct {
	AverageUploadTime string  `json:"avg_upload_time"`
	Uptime            float64 `json:"uptime"`
}

type ErrorResponse struct {
	Msg string `json:"msg"`
}
