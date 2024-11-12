package model

import "github.com/chris-watson/mac-windows-installer-app/pkg/domain"

func TransformPingResponse(res domain.PingResult) PingResponse {
	return PingResponse{
		Success: res.Successful,
		Time:    res.Time.String(),
	}
}

func TransformSystemInfoResponse(res domain.SystemInfo) SystemInfo {
	return SystemInfo{
		Hostname:  res.Hostname,
		IPAddress: res.IPAddress,
	}
}
