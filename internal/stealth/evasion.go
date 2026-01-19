package stealth

import (
	"math/rand"
	"time"
)

// ShuffleOrder randomizes port order for IDS evasion
func ShuffleOrder(ports []int) []int {
	shuffled := make([]int, len(ports))
	copy(shuffled, ports)

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	return shuffled
}

// RandomDelay returns a random delay between min and max
func RandomDelay(minMs, maxMs int) time.Duration {
	if minMs >= maxMs {
		return time.Duration(minMs) * time.Millisecond
	}

	rand.Seed(time.Now().UnixNano())
	delay := rand.Intn(maxMs-minMs) + minMs
	return time.Duration(delay) * time.Millisecond
}

// StealthConfig stealth scan configuration
type StealthConfig struct {
	RandomOrder bool
	MinDelayMs  int
	MaxDelayMs  int
	Enabled     bool
}

// DefaultStealthConfig returns default stealth config
func DefaultStealthConfig() *StealthConfig {
	return &StealthConfig{
		RandomOrder: true,
		MinDelayMs:  100,
		MaxDelayMs:  500,
		Enabled:     false,
	}
}

// ApplyDelay applies random delay if stealth mode is enabled
func (sc *StealthConfig) ApplyDelay() {
	if sc.Enabled {
		delay := RandomDelay(sc.MinDelayMs, sc.MaxDelayMs)
		time.Sleep(delay)
	}
}
