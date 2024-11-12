package adapter

import (
	"bytes"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/chris-watson/mac-windows-installer-app/pkg/domain"
)

type WindowsCommander struct{}

func NewWindowsCommander() domain.Commander {
	return &WindowsCommander{}
}

func (c *WindowsCommander) GetSystemInfo() (domain.SystemInfo, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return domain.SystemInfo{}, err
	}

	// We collect all IP addresses to allow for multiple network interfaces
	// but for this prototype application we will only return the first one
	//
	// Note: This will only return local network IPs (e.g., 192.168.1.x, 10.0.x.x).
	// We cannot detect the public IP address if the host is:
	// - Behind a NAT/home router
	// - Behind a corporate proxy
	// For public IP detection, an external service would be needed.
	ips, err := c.getIPAddresses()
	if err != nil {
		return domain.SystemInfo{}, err
	}

	if len(ips) == 0 {
		return domain.SystemInfo{
			Hostname:  hostname,
			IPAddress: "127.0.0.1",
		}, nil
	}

	return domain.SystemInfo{
		Hostname:  hostname,
		IPAddress: ips[0].String(),
	}, nil
}

func (c *WindowsCommander) getIPAddresses() ([]net.IP, error) {

	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var ips []net.IP
	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip != nil {
				ips = append(ips, ip)
			}
		}
	}

	return ips, nil

}

func (c *WindowsCommander) Ping(host string) (domain.PingResult, error) {
	cmd := exec.Command("ping", "-n", "1", host) // this is platform specific
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return domain.PingResult{
			Successful: false,
			Time:       0,
		}, err
	}

	output := out.String()
	if strings.Contains(output, "TTL=") {
		// Extract the time from the output
		timeIndex := strings.Index(output, "time=")
		if timeIndex != -1 {
			timeStr := output[timeIndex+5:]
			timeStr = strings.Split(timeStr, "ms")[0]
			timeStr = strings.TrimSpace(timeStr)
			timeValue, err := time.ParseDuration(timeStr + "ms")
			if err == nil {
				return domain.PingResult{
					Successful: true,
					Time:       timeValue,
				}, nil
			}
		}
	}

	return domain.PingResult{
		Successful: false,
		Time:       0,
	}, nil
}
