package infra

import (
	"runtime"
)

type Platform int

const (
	Unknown Platform = iota
	Windows
	MacOS
	Linux
)

func (p Platform) String() string {
	switch p {
	case Windows:
		return "windows"
	case MacOS:
		return "darwin"
	case Linux:
		return "linux"
	default:
		return "unknown"
	}
}

func getPlatform() Platform {
	switch runtime.GOOS {
	case "windows":
		return Windows
	case "darwin":
		return MacOS
	case "linux":
		return Linux
	default:
		return Unknown
	}
}
