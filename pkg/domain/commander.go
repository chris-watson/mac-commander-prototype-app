package domain

import "time"

type Commander interface {
	Ping(host string) (PingResult, error)
	GetSystemInfo() (SystemInfo, error)
}

type PingResult struct {
	Successful bool          `json:"successful"`
	Time       time.Duration `json:"time"`
}

type SystemInfo struct {
	Hostname  string `json:"hostname"`
	IPAddress string `json:"ip_address"`
}
