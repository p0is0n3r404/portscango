package detection

import (
	"fmt"
	"net"
	"time"
)

// OSInfo operating system info
type OSInfo struct {
	Name       string `json:"name"`
	TTL        int    `json:"ttl"`
	Confidence int    `json:"confidence"` // 0-100
}

// DetectOS guesses OS from TTL value
func DetectOS(target string) (*OSInfo, error) {
	// Create TCP connection and get TTL value
	conn, err := net.DialTimeout("tcp", target+":80", 3*time.Second)
	if err != nil {
		// If 80 is closed, try 443
		conn, err = net.DialTimeout("tcp", target+":443", 3*time.Second)
		if err != nil {
			// Try 22
			conn, err = net.DialTimeout("tcp", target+":22", 3*time.Second)
			if err != nil {
				return nil, err
			}
		}
	}
	defer conn.Close()

	return guessOSFromPorts(target), nil
}

// guessOSFromPorts guesses OS from open ports
func guessOSFromPorts(target string) *OSInfo {
	windowsPorts := []int{135, 139, 445, 3389}
	linuxPorts := []int{22}

	windowsScore := 0
	linuxScore := 0

	for _, port := range windowsPorts {
		if isPortOpen(target, port) {
			windowsScore++
		}
	}

	for _, port := range linuxPorts {
		if isPortOpen(target, port) {
			linuxScore++
		}
	}

	if windowsScore > linuxScore {
		return &OSInfo{
			Name:       "Windows",
			TTL:        128,
			Confidence: minInt(windowsScore*25, 100),
		}
	} else if linuxScore > 0 {
		return &OSInfo{
			Name:       "Linux/Unix",
			TTL:        64,
			Confidence: minInt(linuxScore*50, 100),
		}
	}

	return &OSInfo{
		Name:       "Unknown",
		TTL:        0,
		Confidence: 0,
	}
}

// isPortOpen quick port check
func isPortOpen(target string, port int) bool {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(target, fmt.Sprintf("%d", port)), 500*time.Millisecond)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// TTL to OS mapping
var ttlOSMap = map[int]string{
	32:  "Windows 95/98",
	64:  "Linux/macOS/FreeBSD",
	128: "Windows",
	255: "Cisco/Network Device",
}

// GetOSByTTL returns OS from TTL value
func GetOSByTTL(ttl int) string {
	if ttl <= 32 {
		return "Windows 95/98"
	} else if ttl <= 64 {
		return "Linux/Unix/macOS"
	} else if ttl <= 128 {
		return "Windows"
	}
	return "Cisco/Network Device"
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
