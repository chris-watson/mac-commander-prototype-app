package service

import (
	"log"

	"github.com/chris-watson/mac-windows-installer-app/pkg/domain"
)

type Commander interface {
	Ping(host string) (domain.PingResult, error)
	GetSystemInfo() (domain.SystemInfo, error)
}

type CommanderService struct {
	cmdr Commander
}

func NewCommanderService(cmdr Commander) *CommanderService {
	return &CommanderService{cmdr: cmdr}
}

func (c *CommanderService) HandlePing(host string) (domain.PingResult, error) {
	result, err := c.cmdr.Ping(host)
	if err != nil {
		log.Printf("failed to ping host %s: %v", host, err)
	}
	return result, err
}

func (c *CommanderService) GetSystemInfo() (domain.SystemInfo, error) {
	info, err := c.cmdr.GetSystemInfo()
	if err != nil {
		log.Printf("failed to get system info: %v", err)
	}
	return info, err
}
