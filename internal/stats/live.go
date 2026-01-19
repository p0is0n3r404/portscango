package stats

import (
	"fmt"
	"sync"
	"time"
)

// LiveStats real-time scan statistics
type LiveStats struct {
	StartTime    time.Time
	TotalPorts   int
	ScannedPorts int
	OpenPorts    int
	ClosedPorts  int
	mutex        sync.Mutex
}

// NewLiveStats creates new stats tracker
func NewLiveStats(totalPorts int) *LiveStats {
	return &LiveStats{
		StartTime:  time.Now(),
		TotalPorts: totalPorts,
	}
}

// IncrementScanned increments scanned port count
func (s *LiveStats) IncrementScanned() {
	s.mutex.Lock()
	s.ScannedPorts++
	s.mutex.Unlock()
}

// IncrementOpen increments open port count
func (s *LiveStats) IncrementOpen() {
	s.mutex.Lock()
	s.OpenPorts++
	s.mutex.Unlock()
}

// GetPortsPerSecond returns current scan rate
func (s *LiveStats) GetPortsPerSecond() float64 {
	elapsed := time.Since(s.StartTime).Seconds()
	if elapsed == 0 {
		return 0
	}
	return float64(s.ScannedPorts) / elapsed
}

// GetProgress returns progress percentage
func (s *LiveStats) GetProgress() float64 {
	if s.TotalPorts == 0 {
		return 0
	}
	return float64(s.ScannedPorts) / float64(s.TotalPorts) * 100
}

// GetETA returns estimated time remaining
func (s *LiveStats) GetETA() time.Duration {
	pps := s.GetPortsPerSecond()
	if pps == 0 {
		return 0
	}
	remaining := s.TotalPorts - s.ScannedPorts
	seconds := float64(remaining) / pps
	return time.Duration(seconds) * time.Second
}

// PrintStats prints current statistics
func (s *LiveStats) PrintStats() {
	progress := s.GetProgress()
	pps := s.GetPortsPerSecond()
	eta := s.GetETA()

	// Clear line and print stats
	fmt.Printf("\r\033[K\033[34m[*] Progress: %.1f%% | Scanned: %d/%d | Open: %d | Rate: %.0f ports/sec | ETA: %s\033[0m",
		progress, s.ScannedPorts, s.TotalPorts, s.OpenPorts, pps, eta.Round(time.Second))
}

// GetSummary returns statistics summary
func (s *LiveStats) GetSummary() string {
	elapsed := time.Since(s.StartTime).Round(time.Millisecond)
	pps := s.GetPortsPerSecond()

	return fmt.Sprintf("Scanned %d ports in %s (%.0f ports/sec), %d open",
		s.ScannedPorts, elapsed, pps, s.OpenPorts)
}
