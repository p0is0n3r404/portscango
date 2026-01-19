package detect

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

// HTTPInfo HTTP service information
type HTTPInfo struct {
	StatusCode   int               `json:"status_code"`
	Server       string            `json:"server"`
	PoweredBy    string            `json:"powered_by,omitempty"`
	ContentType  string            `json:"content_type"`
	Headers      map[string]string `json:"headers"`
	Title        string            `json:"title,omitempty"`
	Technologies []string          `json:"technologies,omitempty"`
	WAF          string            `json:"waf,omitempty"`
}

// GetHTTPInfo retrieves HTTP service information
func GetHTTPInfo(host string, port int, useSSL bool) (*HTTPInfo, error) {
	scheme := "http"
	if useSSL || port == 443 || port == 8443 {
		scheme = "https"
	}

	url := fmt.Sprintf("%s://%s:%d/", scheme, host, port)

	// Create client with custom transport
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).DialContext,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // Don't follow redirects
		},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "PortScanGO/4.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	info := &HTTPInfo{
		StatusCode: resp.StatusCode,
		Headers:    make(map[string]string),
	}

	// Extract headers
	for key, values := range resp.Header {
		info.Headers[key] = strings.Join(values, ", ")
	}

	// Server header
	info.Server = resp.Header.Get("Server")
	info.PoweredBy = resp.Header.Get("X-Powered-By")
	info.ContentType = resp.Header.Get("Content-Type")

	// Read body for title and tech detection
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1024*100)) // 100KB max
	if err == nil {
		bodyStr := string(body)
		info.Title = extractTitle(bodyStr)
		info.Technologies = detectTechnologies(bodyStr, info.Headers)
	}

	// WAF detection
	info.WAF = detectWAF(resp.Header, info.Server)

	return info, nil
}

// extractTitle extracts page title from HTML
func extractTitle(html string) string {
	html = strings.ToLower(html)
	start := strings.Index(html, "<title>")
	if start == -1 {
		return ""
	}
	start += 7
	end := strings.Index(html[start:], "</title>")
	if end == -1 {
		return ""
	}
	return strings.TrimSpace(html[start : start+end])
}

// detectTechnologies detects web technologies
func detectTechnologies(body string, headers map[string]string) []string {
	var techs []string
	bodyLower := strings.ToLower(body)

	// CMS Detection
	if strings.Contains(bodyLower, "wp-content") || strings.Contains(bodyLower, "wordpress") {
		techs = append(techs, "WordPress")
	}
	if strings.Contains(bodyLower, "drupal") {
		techs = append(techs, "Drupal")
	}
	if strings.Contains(bodyLower, "joomla") {
		techs = append(techs, "Joomla")
	}

	// JS Frameworks
	if strings.Contains(bodyLower, "react") || strings.Contains(bodyLower, "reactdom") {
		techs = append(techs, "React")
	}
	if strings.Contains(bodyLower, "vue.js") || strings.Contains(bodyLower, "vuejs") {
		techs = append(techs, "Vue.js")
	}
	if strings.Contains(bodyLower, "angular") {
		techs = append(techs, "Angular")
	}
	if strings.Contains(bodyLower, "jquery") {
		techs = append(techs, "jQuery")
	}

	// Server technologies from headers
	if server, ok := headers["Server"]; ok {
		serverLower := strings.ToLower(server)
		if strings.Contains(serverLower, "nginx") {
			techs = append(techs, "Nginx")
		}
		if strings.Contains(serverLower, "apache") {
			techs = append(techs, "Apache")
		}
		if strings.Contains(serverLower, "iis") {
			techs = append(techs, "IIS")
		}
	}

	if poweredBy, ok := headers["X-Powered-By"]; ok {
		if strings.Contains(strings.ToLower(poweredBy), "php") {
			techs = append(techs, "PHP")
		}
		if strings.Contains(strings.ToLower(poweredBy), "asp") {
			techs = append(techs, "ASP.NET")
		}
	}

	return techs
}

// detectWAF detects Web Application Firewall
func detectWAF(headers http.Header, server string) string {
	serverLower := strings.ToLower(server)

	// Cloudflare
	if headers.Get("CF-Ray") != "" || strings.Contains(serverLower, "cloudflare") {
		return "Cloudflare"
	}

	// AWS WAF
	if headers.Get("X-Amz-Cf-Id") != "" {
		return "AWS CloudFront"
	}

	// Akamai
	if headers.Get("X-Akamai-Transformed") != "" {
		return "Akamai"
	}

	// Sucuri
	if strings.Contains(serverLower, "sucuri") {
		return "Sucuri"
	}

	// Incapsula
	if headers.Get("X-Iinfo") != "" {
		return "Incapsula"
	}

	// ModSecurity
	if strings.Contains(serverLower, "mod_security") {
		return "ModSecurity"
	}

	return ""
}

// PrintHTTPInfo prints HTTP info
func PrintHTTPInfo(info *HTTPInfo) {
	fmt.Println("\n\033[35m╔══════════════════════════════════════════════════════╗\033[0m")
	fmt.Println("\033[35m║           HTTP SERVICE INFORMATION                     ║\033[0m")
	fmt.Println("\033[35m╚══════════════════════════════════════════════════════╝\033[0m")

	fmt.Printf("  \033[36mStatus:\033[0m       %d\n", info.StatusCode)
	if info.Server != "" {
		fmt.Printf("  \033[36mServer:\033[0m       %s\n", info.Server)
	}
	if info.PoweredBy != "" {
		fmt.Printf("  \033[36mPowered By:\033[0m   %s\n", info.PoweredBy)
	}
	if info.Title != "" {
		fmt.Printf("  \033[36mTitle:\033[0m        %s\n", info.Title)
	}
	if len(info.Technologies) > 0 {
		fmt.Printf("  \033[36mTechnologies:\033[0m %s\n", strings.Join(info.Technologies, ", "))
	}
	if info.WAF != "" {
		fmt.Printf("  \033[33m⚠ WAF:\033[0m        %s\n", info.WAF)
	}
	fmt.Println()
}
