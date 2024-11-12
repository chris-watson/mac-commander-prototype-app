package model

type CommandRequest struct {
	Type    string `json:"type"`
	Payload string `json:"payload,omitempty"`
}

type CommandResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type PingResponse struct {
	Success bool   `json:"success"`
	Time    string `json:"time"`
	Error   string `json:"error,omitempty"`
}

type SystemInfo struct {
	Hostname  string `json:"hostname"`
	IPAddress string `json:"ip_address"`
}
