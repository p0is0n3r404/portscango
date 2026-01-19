package network

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"
)

// SSLInfo SSL certificate info
type SSLInfo struct {
	CommonName   string   `json:"common_name"`
	Issuer       string   `json:"issuer"`
	ValidFrom    string   `json:"valid_from"`
	ValidTo      string   `json:"valid_to"`
	SANs         []string `json:"sans,omitempty"`
	Protocol     string   `json:"protocol"`
	IsExpired    bool     `json:"is_expired"`
	DaysToExpiry int      `json:"days_to_expiry"`
}

// GetSSLInfo retrieves SSL/TLS certificate info
func GetSSLInfo(host string, port int) (*SSLInfo, error) {
	address := fmt.Sprintf("%s:%d", host, port)

	// Create TLS connection
	conn, err := tls.DialWithDialer(
		&net.Dialer{Timeout: 5 * time.Second},
		"tcp",
		address,
		&tls.Config{
			InsecureSkipVerify: true, // Accept self-signed certificates
		},
	)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Get certificate info
	state := conn.ConnectionState()
	if len(state.PeerCertificates) == 0 {
		return nil, fmt.Errorf("no certificate found")
	}

	cert := state.PeerCertificates[0]
	now := time.Now()

	info := &SSLInfo{
		CommonName:   cert.Subject.CommonName,
		Issuer:       cert.Issuer.CommonName,
		ValidFrom:    cert.NotBefore.Format("2006-01-02"),
		ValidTo:      cert.NotAfter.Format("2006-01-02"),
		SANs:         cert.DNSNames,
		Protocol:     tlsVersionString(state.Version),
		IsExpired:    now.After(cert.NotAfter),
		DaysToExpiry: int(cert.NotAfter.Sub(now).Hours() / 24),
	}

	return info, nil
}

// tlsVersionString converts TLS version to string
func tlsVersionString(version uint16) string {
	switch version {
	case tls.VersionTLS10:
		return "TLS 1.0"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS13:
		return "TLS 1.3"
	default:
		return "Unknown"
	}
}

// PrintSSLInfo prints SSL info
func PrintSSLInfo(info *SSLInfo) {
	fmt.Println("\n\033[35m╔══════════════════════════════════════════════════════╗\033[0m")
	fmt.Println("\033[35m║           SSL/TLS CERTIFICATE INFO                     ║\033[0m")
	fmt.Println("\033[35m╚══════════════════════════════════════════════════════╝\033[0m")

	fmt.Printf("  \033[36mCommon Name:\033[0m  %s\n", info.CommonName)
	fmt.Printf("  \033[36mIssuer:\033[0m       %s\n", info.Issuer)
	fmt.Printf("  \033[36mValid From:\033[0m   %s\n", info.ValidFrom)
	fmt.Printf("  \033[36mValid To:\033[0m     %s\n", info.ValidTo)
	fmt.Printf("  \033[36mProtocol:\033[0m     %s\n", info.Protocol)

	if info.IsExpired {
		fmt.Printf("  \033[31m⚠ CERTIFICATE EXPIRED!\033[0m\n")
	} else if info.DaysToExpiry < 30 {
		fmt.Printf("  \033[33m⚠ Certificate expires in %d days\033[0m\n", info.DaysToExpiry)
	} else {
		fmt.Printf("  \033[32m✓ Certificate valid (%d days remaining)\033[0m\n", info.DaysToExpiry)
	}

	if len(info.SANs) > 0 {
		fmt.Printf("  \033[36mSAN:\033[0m          %v\n", info.SANs)
	}
	fmt.Println()
}
