package scanner

import (
	"fmt"
	"net"
	"sort"
	"sync"
	"time"

	"portscango/pkg/ports"
)

// Result represents a single port scan result
type Result struct {
	Port    int    `json:"port"`
	State   string `json:"state"`
	Service string `json:"service"`
	Banner  string `json:"banner,omitempty"`
}

// Scanner port scanner structure
type Scanner struct {
	Target        string
	Ports         []int
	Timeout       time.Duration
	Threads       int
	ServiceDetect bool
	BannerGrab    bool
	Results       []Result
	mutex         sync.Mutex
}

// NewScanner creates a new scanner
func NewScanner(target string, portList []int, timeout time.Duration, threads int) *Scanner {
	return &Scanner{
		Target:  target,
		Ports:   portList,
		Timeout: timeout,
		Threads: threads,
		Results: make([]Result, 0),
	}
}

// WithServiceDetection enables service detection
func (s *Scanner) WithServiceDetection() *Scanner {
	s.ServiceDetect = true
	return s
}

// WithBannerGrab enables banner grabbing
func (s *Scanner) WithBannerGrab() *Scanner {
	s.BannerGrab = true
	return s
}

// Scan scans ports
func (s *Scanner) Scan(progressChan chan<- int) {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, s.Threads)

	for _, port := range s.Ports {
		wg.Add(1)
		semaphore <- struct{}{} // Rate limiting

		go func(p int) {
			defer wg.Done()
			defer func() { <-semaphore }()

			result := s.scanPort(p)
			if result != nil {
				s.mutex.Lock()
				s.Results = append(s.Results, *result)
				s.mutex.Unlock()
			}

			if progressChan != nil {
				progressChan <- 1
			}
		}(port)
	}

	wg.Wait()

	if progressChan != nil {
		close(progressChan)
	}

	// Sort results by port number
	sort.Slice(s.Results, func(i, j int) bool {
		return s.Results[i].Port < s.Results[j].Port
	})
}

// scanPort scans a single port
func (s *Scanner) scanPort(port int) *Result {
	address := net.JoinHostPort(s.Target, fmt.Sprintf("%d", port))
	conn, err := net.DialTimeout("tcp", address, s.Timeout)
	if err != nil {
		return nil
	}
	defer conn.Close()

	result := &Result{
		Port:  port,
		State: "open",
	}

	// Service Detection
	if s.ServiceDetect {
		result.Service = ports.GetServiceName(port)
	}

	// Banner Grabbing
	if s.BannerGrab {
		result.Banner = s.grabBanner(conn, port)
	}

	return result
}

// grabBanner grabs banner from open port
func (s *Scanner) grabBanner(conn net.Conn, port int) string {
	// Special request for HTTP
	if port == 80 || port == 8080 || port == 8000 || port == 8888 {
		conn.Write([]byte("HEAD / HTTP/1.0\r\n\r\n"))
	} else if port == 443 || port == 8443 {
		return "" // Cannot grab simple banner from SSL ports
	}

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return ""
	}

	banner := string(buffer[:n])
	// Get first line and clean
	for i, c := range banner {
		if c == '\n' || c == '\r' {
			banner = banner[:i]
			break
		}
	}

	// Maximum 50 characters
	if len(banner) > 50 {
		banner = banner[:50] + "..."
	}

	return banner
}

// GetResults returns scan results
func (s *Scanner) GetResults() []Result {
	return s.Results
}

// OpenPortCount returns open port count
func (s *Scanner) OpenPortCount() int {
	return len(s.Results)
}
