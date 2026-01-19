package target

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

// ParseTargets parses multiple targets
func ParseTargets(input string) ([]string, error) {
	// Read from file
	if strings.HasPrefix(input, "@") {
		return parseFromFile(input[1:])
	}

	// CIDR check
	if strings.Contains(input, "/") {
		return parseCIDR(input)
	}

	// IP range check (192.168.1.1-50)
	if strings.Contains(input, "-") && !strings.Contains(input, ".") == false {
		return parseIPRange(input)
	}

	// Comma-separated list
	if strings.Contains(input, ",") {
		targets := strings.Split(input, ",")
		for i := range targets {
			targets[i] = strings.TrimSpace(targets[i])
		}
		return targets, nil
	}

	// Single target
	return []string{input}, nil
}

// parseFromFile reads targets from file
func parseFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	var targets []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			targets = append(targets, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return targets, nil
}

// parseCIDR converts CIDR notation to IP list
func parseCIDR(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, fmt.Errorf("invalid CIDR: %w", err)
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); incrementIP(ip) {
		ips = append(ips, ip.String())
	}

	// Remove first (network) and last (broadcast) addresses
	if len(ips) > 2 {
		ips = ips[1 : len(ips)-1]
	}

	return ips, nil
}

// parseIPRange parses IP range (192.168.1.1-50)
func parseIPRange(input string) ([]string, error) {
	// Find last octet range
	lastDot := strings.LastIndex(input, ".")
	if lastDot == -1 {
		return nil, fmt.Errorf("invalid IP range")
	}

	baseIP := input[:lastDot+1]
	rangeStr := input[lastDot+1:]

	parts := strings.Split(rangeStr, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid range format")
	}

	start, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid start: %w", err)
	}

	end, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid end: %w", err)
	}

	var ips []string
	for i := start; i <= end; i++ {
		ips = append(ips, fmt.Sprintf("%s%d", baseIP, i))
	}

	return ips, nil
}

// incrementIP increments IP address by one
func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
