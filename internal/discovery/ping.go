package discovery

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// Host discovered host info
type Host struct {
	IP       string        `json:"ip"`
	Hostname string        `json:"hostname,omitempty"`
	IsAlive  bool          `json:"is_alive"`
	RTT      time.Duration `json:"rtt"`
}

// PingSweep performs TCP ping sweep on a network range
func PingSweep(hosts []string, threads int, timeout time.Duration) []Host {
	var results []Host
	var mutex sync.Mutex
	var wg sync.WaitGroup

	semaphore := make(chan struct{}, threads)

	for _, host := range hosts {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(h string) {
			defer wg.Done()
			defer func() { <-semaphore }()

			result := tcpPing(h, timeout)

			mutex.Lock()
			results = append(results, result)
			mutex.Unlock()
		}(host)
	}

	wg.Wait()
	return results
}

// tcpPing performs TCP ping on common ports
func tcpPing(host string, timeout time.Duration) Host {
	result := Host{
		IP:      host,
		IsAlive: false,
	}

	// Try common ports
	ports := []int{80, 443, 22, 21, 25, 3389}

	for _, port := range ports {
		start := time.Now()
		address := net.JoinHostPort(host, fmt.Sprintf("%d", port))

		conn, err := net.DialTimeout("tcp", address, timeout)
		if err == nil {
			conn.Close()
			result.IsAlive = true
			result.RTT = time.Since(start)

			// Try to resolve hostname
			names, err := net.LookupAddr(host)
			if err == nil && len(names) > 0 {
				result.Hostname = names[0]
			}
			break
		}
	}

	return result
}

// PrintDiscoveryResults prints discovered hosts
func PrintDiscoveryResults(hosts []Host) {
	aliveCount := 0

	fmt.Println("\n\033[36m╔══════════════════════════════════════════════════════╗\033[0m")
	fmt.Println("\033[36m║             HOST DISCOVERY RESULTS                     ║\033[0m")
	fmt.Println("\033[36m╚══════════════════════════════════════════════════════╝\033[0m")
	fmt.Println()

	fmt.Printf("%-18s %-8s %-12s %s\n", "IP ADDRESS", "STATUS", "RTT", "HOSTNAME")
	fmt.Println("─────────────────────────────────────────────────────────")

	for _, h := range hosts {
		if h.IsAlive {
			aliveCount++
			hostname := h.Hostname
			if hostname == "" {
				hostname = "-"
			}
			fmt.Printf("\033[32m%-18s %-8s %-12s %s\033[0m\n",
				h.IP, "UP", h.RTT.Round(time.Millisecond), hostname)
		}
	}

	fmt.Println()
	fmt.Printf("\033[34m[*] Hosts scanned: %d\033[0m\n", len(hosts))
	fmt.Printf("\033[32m[*] Hosts alive: %d\033[0m\n", aliveCount)
	fmt.Printf("\033[31m[*] Hosts down: %d\033[0m\n", len(hosts)-aliveCount)
}

// GetAliveHosts returns only alive hosts
func GetAliveHosts(hosts []Host) []string {
	var alive []string
	for _, h := range hosts {
		if h.IsAlive {
			alive = append(alive, h.IP)
		}
	}
	return alive
}
