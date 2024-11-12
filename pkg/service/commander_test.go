package service_test

import (
	"testing"

	"github.com/chris-watson/mac-windows-installer-app/pkg/domain"
	"github.com/chris-watson/mac-windows-installer-app/pkg/service"
)

type MockCommander struct{}

func (m *MockCommander) GetSystemInfo() (domain.SystemInfo, error) {
	return domain.SystemInfo{
		Hostname:  "mock-host",
		IPAddress: "192.168.1.1",
	}, nil
}

func (m *MockCommander) Ping(_ string) (domain.PingResult, error) {
	return domain.PingResult{
		Successful: true,
	}, nil
}

func TestGetSystemInfo(t *testing.T) {
	cmdr := service.NewCommanderService(&MockCommander{})
	info, err := cmdr.GetSystemInfo()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if info.Hostname == "" {
		t.Error("Expected hostname to be non-empty")
	}

	if info.IPAddress == "" {
		t.Error("Expected IP address to be non-empty")
	}
}

func TestPing(t *testing.T) {
	cmdr := service.NewCommanderService(&MockCommander{})
	res, err := cmdr.HandlePing("google.com")
	if err != nil {
		t.Errorf("Expected ping to succeed, got error: %v", err)
	}
	if !res.Successful {
		t.Errorf("Expected successful ping, got %v", res)
	}
}
