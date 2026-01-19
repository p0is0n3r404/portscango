package scanner

// Profile scan profile
type Profile struct {
	Name      string
	Ports     []int
	Threads   int
	TimeoutMs int
}

// QuickProfile quick scan profile
var QuickProfile = Profile{
	Name:      "quick",
	Ports:     []int{21, 22, 23, 25, 80, 110, 143, 443, 445, 993, 995, 1433, 3306, 3389, 5432, 5900, 8080, 8443, 27017},
	Threads:   200,
	TimeoutMs: 500,
}

// NormalProfile normal scan profile
var NormalProfile = Profile{
	Name:      "normal",
	Ports:     nil, // Top 100
	Threads:   100,
	TimeoutMs: 1000,
}

// AggressiveProfile aggressive scan profile
var AggressiveProfile = Profile{
	Name:      "aggressive",
	Ports:     nil, // 1-65535
	Threads:   500,
	TimeoutMs: 200,
}

// GetProfileByName returns profile by name
func GetProfileByName(name string) *Profile {
	switch name {
	case "quick":
		return &QuickProfile
	case "normal":
		return &NormalProfile
	case "aggressive":
		return &AggressiveProfile
	default:
		return nil
	}
}

// GenerateFullPortRange generates all ports
func GenerateFullPortRange() []int {
	ports := make([]int, 65535)
	for i := 0; i < 65535; i++ {
		ports[i] = i + 1
	}
	return ports
}
