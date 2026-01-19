package detection

import (
	"fmt"
	"strings"

	"portscango/internal/scanner"
)

// Vulnerability vulnerability info
type Vulnerability struct {
	Port        int    `json:"port"`
	Service     string `json:"service"`
	Severity    string `json:"severity"` // LOW, MEDIUM, HIGH, CRITICAL
	Description string `json:"description"`
	CVE         string `json:"cve,omitempty"`
}

// VulnChecker vulnerability checker
type VulnChecker struct {
	Results []scanner.Result
}

// NewVulnChecker creates a new checker
func NewVulnChecker(results []scanner.Result) *VulnChecker {
	return &VulnChecker{Results: results}
}

// Check checks for vulnerabilities
func (vc *VulnChecker) Check() []Vulnerability {
	var vulns []Vulnerability

	for _, r := range vc.Results {
		// Port-based checks
		portVulns := vc.checkPortVulns(r)
		vulns = append(vulns, portVulns...)

		// Banner-based checks
		if r.Banner != "" {
			bannerVulns := vc.checkBannerVulns(r)
			vulns = append(vulns, bannerVulns...)
		}
	}

	return vulns
}

// checkPortVulns checks port-based vulnerabilities
func (vc *VulnChecker) checkPortVulns(r scanner.Result) []Vulnerability {
	var vulns []Vulnerability

	// Dangerous ports
	switch r.Port {
	case 21:
		vulns = append(vulns, Vulnerability{
			Port:        21,
			Service:     "FTP",
			Severity:    "MEDIUM",
			Description: "FTP is a clear-text protocol. Data can be transferred without encryption.",
		})

	case 23:
		vulns = append(vulns, Vulnerability{
			Port:        23,
			Service:     "Telnet",
			Severity:    "HIGH",
			Description: "Telnet is an insecure protocol. SSH should be used instead.",
		})

	case 445:
		vulns = append(vulns, Vulnerability{
			Port:        445,
			Service:     "SMB",
			Severity:    "HIGH",
			Description: "SMB port is open. Check for EternalBlue (MS17-010) vulnerability.",
			CVE:         "CVE-2017-0144",
		})

	case 3389:
		vulns = append(vulns, Vulnerability{
			Port:        3389,
			Service:     "RDP",
			Severity:    "MEDIUM",
			Description: "RDP is open. Check for BlueKeep vulnerability.",
			CVE:         "CVE-2019-0708",
		})

	case 6379:
		vulns = append(vulns, Vulnerability{
			Port:        6379,
			Service:     "Redis",
			Severity:    "CRITICAL",
			Description: "Redis often runs without authentication. Unauthorized access risk!",
		})

	case 27017:
		vulns = append(vulns, Vulnerability{
			Port:        27017,
			Service:     "MongoDB",
			Severity:    "CRITICAL",
			Description: "MongoDB does not require authentication by default. Data leak risk!",
		})

	case 11211:
		vulns = append(vulns, Vulnerability{
			Port:        11211,
			Service:     "Memcached",
			Severity:    "HIGH",
			Description: "Memcached can be used in DDoS amplification attacks.",
		})
	}

	return vulns
}

// checkBannerVulns checks banner-based vulnerabilities
func (vc *VulnChecker) checkBannerVulns(r scanner.Result) []Vulnerability {
	var vulns []Vulnerability
	banner := strings.ToLower(r.Banner)

	// Old SSH versions
	if strings.Contains(banner, "ssh") {
		if strings.Contains(banner, "openssh_4") || strings.Contains(banner, "openssh_5") {
			vulns = append(vulns, Vulnerability{
				Port:        r.Port,
				Service:     "SSH",
				Severity:    "HIGH",
				Description: "Old OpenSSH version detected. Update recommended.",
			})
		}
	}

	// Apache old version
	if strings.Contains(banner, "apache/2.2") {
		vulns = append(vulns, Vulnerability{
			Port:        r.Port,
			Service:     "HTTP",
			Severity:    "MEDIUM",
			Description: "Old Apache version (2.2.x) detected.",
		})
	}

	// PHP old version
	if strings.Contains(banner, "php/5") {
		vulns = append(vulns, Vulnerability{
			Port:        r.Port,
			Service:     "HTTP",
			Severity:    "HIGH",
			Description: "PHP 5.x detected. End-of-life version!",
		})
	}

	// ProFTPD
	if strings.Contains(banner, "proftpd") {
		vulns = append(vulns, Vulnerability{
			Port:        r.Port,
			Service:     "FTP",
			Severity:    "MEDIUM",
			Description: "ProFTPD detected. Check for CVE-2015-3306.",
			CVE:         "CVE-2015-3306",
		})
	}

	return vulns
}

// PrintVulns prints vulnerabilities with colors
func PrintVulns(vulns []Vulnerability) {
	if len(vulns) == 0 {
		fmt.Println("\033[32m[✓] No known vulnerabilities detected.\033[0m")
		return
	}

	fmt.Printf("\n\033[31m[!] %d vulnerabilities detected:\033[0m\n\n", len(vulns))

	for _, v := range vulns {
		var color string
		switch v.Severity {
		case "CRITICAL":
			color = "\033[31m" // Red
		case "HIGH":
			color = "\033[91m" // Light red
		case "MEDIUM":
			color = "\033[33m" // Yellow
		case "LOW":
			color = "\033[36m" // Cyan
		}

		fmt.Printf("%s[%s]\033[0m Port %d (%s)\n", color, v.Severity, v.Port, v.Service)
		fmt.Printf("    └─ %s\n", v.Description)
		if v.CVE != "" {
			fmt.Printf("    └─ \033[35mCVE: %s\033[0m\n", v.CVE)
		}
		fmt.Println()
	}
}
