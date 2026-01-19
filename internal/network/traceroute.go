package network

import (
	"fmt"
	"net"
	"time"
)

// Hop traceroute hop info
type Hop struct {
	Number  int           `json:"hop"`
	Address string        `json:"address"`
	RTT     time.Duration `json:"rtt"`
	Host    string        `json:"host,omitempty"`
}

// Traceroute performs simple TCP traceroute
func Traceroute(target string, maxHops int) ([]Hop, error) {
	var hops []Hop

	// TCP traceroute on port 80 or 443
	ports := []int{80, 443, 22}

	for ttl := 1; ttl <= maxHops; ttl++ {
		hop := Hop{Number: ttl}

		for _, port := range ports {
			address := net.JoinHostPort(target, fmt.Sprintf("%d", port))
			start := time.Now()

			conn, err := net.DialTimeout("tcp", address, 2*time.Second)
			if err != nil {
				// Timeout or error - continue
				continue
			}

			hop.RTT = time.Since(start)
			hop.Address = conn.RemoteAddr().String()
			conn.Close()

			// Resolve hostname
			names, err := net.LookupAddr(target)
			if err == nil && len(names) > 0 {
				hop.Host = names[0]
			}

			hops = append(hops, hop)

			// Reached target
			return hops, nil
		}

		hops = append(hops, hop)
	}

	return hops, nil
}

// SimplePing performs simple ping check
func SimplePing(target string) (time.Duration, error) {
	// TCP ping (ICMP requires admin privileges)
	ports := []int{80, 443, 22, 21}

	for _, port := range ports {
		start := time.Now()
		address := net.JoinHostPort(target, fmt.Sprintf("%d", port))

		conn, err := net.DialTimeout("tcp", address, 3*time.Second)
		if err != nil {
			continue
		}
		conn.Close()
		return time.Since(start), nil
	}

	return 0, fmt.Errorf("target not responding")
}

// PrintTraceroute prints traceroute results
func PrintTraceroute(hops []Hop) {
	fmt.Println("\n\033[36m╔══════════════════════════════════════════════════════╗\033[0m")
	fmt.Println("\033[36m║                    TRACEROUTE                          ║\033[0m")
	fmt.Println("\033[36m╚══════════════════════════════════════════════════════╝\033[0m")

	for _, h := range hops {
		if h.Address != "" {
			host := h.Host
			if host == "" {
				host = h.Address
			}
			fmt.Printf("  \033[33m%2d\033[0m  %-40s  \033[32m%v\033[0m\n", h.Number, host, h.RTT.Round(time.Millisecond))
		} else {
			fmt.Printf("  \033[33m%2d\033[0m  \033[31m* * * (timeout)\033[0m\n", h.Number)
		}
	}
	fmt.Println()
}
