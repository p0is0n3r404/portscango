package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"portscango/internal/detect"
	"portscango/internal/detection"
	"portscango/internal/discovery"
	"portscango/internal/export"
	"portscango/internal/network"
	"portscango/internal/notify"
	"portscango/internal/output"
	"portscango/internal/scanner"
	"portscango/internal/stealth"
	"portscango/internal/target"
	"portscango/pkg/ports"
)

const version = "4.0.0"

const banner = `
██████╗  ██████╗ ██████╗ ████████╗███████╗ ██████╗ █████╗ ███╗   ██╗ ██████╗  ██████╗ 
██╔══██╗██╔═══██╗██╔══██╗╚══██╔══╝██╔════╝██╔════╝██╔══██╗████╗  ██║██╔════╝ ██╔═══██╗
██████╔╝██║   ██║██████╔╝   ██║   ███████╗██║     ███████║██╔██╗ ██║██║  ███╗██║   ██║
██╔═══╝ ██║   ██║██╔══██╗   ██║   ╚════██║██║     ██╔══██║██║╚██╗██║██║   ██║██║   ██║
██║     ╚██████╔╝██║  ██║   ██║   ███████║╚██████╗██║  ██║██║ ╚████║╚██████╔╝╚██████╔╝
╚═╝      ╚═════╝ ╚═╝  ╚═╝   ╚═╝   ╚══════╝ ╚═════╝╚═╝  ╚═╝╚═╝  ╚═══╝ ╚═════╝  ╚═════╝ 
`

func printBanner() {
	fmt.Print("\033[36m")
	fmt.Print(banner)
	fmt.Print("\033[0m")
	fmt.Printf("\033[33m[*] v%s | Ultra High-Performance Port Scanner\033[0m\n\n", version)
}

func main() {
	// CLI Flags - Basic
	targetFlag := flag.String("t", "", "Target IP/domain/CIDR")
	portRange := flag.String("p", "", "Port range (e.g: 1-1000, 80,443)")
	topPorts := flag.Int("top", 0, "Scan top N common ports")
	timeout := flag.Duration("timeout", 1*time.Second, "Connection timeout")
	threads := flag.Int("threads", 100, "Concurrent threads")

	// Features
	serviceDetect := flag.Bool("service", false, "Service detection")
	bannerGrab := flag.Bool("banner", false, "Banner grabbing")
	vulnCheck := flag.Bool("vuln", false, "Vulnerability check")
	sslInfo := flag.Bool("ssl", false, "SSL/TLS info")
	traceroute := flag.Bool("traceroute", false, "Traceroute")
	httpInfo := flag.Bool("http", false, "HTTP info & tech detection")

	// Profiles
	quickProfile := flag.Bool("quick", false, "Quick scan profile")
	aggressiveProfile := flag.Bool("aggressive", false, "Aggressive scan profile")
	stealthMode := flag.Bool("stealth", false, "Stealth scan with random delays")

	// Discovery
	discoverMode := flag.Bool("discover", false, "Host discovery mode")

	// Output
	outputFile := flag.String("o", "", "Output file (json/txt/html/xml/csv/md)")
	noColor := flag.Bool("no-color", false, "Disable colored output")
	showVersion := flag.Bool("v", false, "Show version")

	// Notifications
	discordWebhook := flag.String("discord", "", "Discord webhook URL")

	flag.Parse()

	if *showVersion {
		fmt.Printf("PortScanGO v%s\n", version)
		os.Exit(0)
	}

	printBanner()

	// Interactive mode
	if *targetFlag == "" {
		fmt.Println("\033[33m╔══════════════════════════════════════════════════════════════╗\033[0m")
		fmt.Println("\033[33m║              WELCOME TO PORTSCANGO v4.0!                     ║\033[0m")
		fmt.Println("\033[33m╚══════════════════════════════════════════════════════════════╝\033[0m")
		fmt.Println()
		fmt.Println("\033[36mNo target specified. Let's set up a scan!\033[0m")
		fmt.Println()

		fmt.Print("\033[32m[+] Enter target (IP or domain): \033[0m")
		var inputTarget string
		fmt.Scanln(&inputTarget)

		if inputTarget == "" {
			fmt.Println("\033[31m[!] Target is required. Exiting.\033[0m")
			os.Exit(1)
		}
		*targetFlag = inputTarget

		fmt.Println()
		fmt.Println("\033[36mSelect scan type:\033[0m")
		fmt.Println("  \033[33m[1]\033[0m Quick Scan   - Top 20 ports, fast")
		fmt.Println("  \033[33m[2]\033[0m Normal Scan  - Top 1000 ports")
		fmt.Println("  \033[33m[3]\033[0m Full Scan    - All 65535 ports")
		fmt.Println("  \033[33m[4]\033[0m Stealth Scan - Slow & undetectable")
		fmt.Println()
		fmt.Print("\033[32m[+] Enter choice (1-4) [default: 1]: \033[0m")
		var choice string
		fmt.Scanln(&choice)

		switch choice {
		case "2":
			*topPorts = 1000
			fmt.Println("\033[34m[*] Selected: Normal Scan\033[0m")
		case "3":
			*aggressiveProfile = true
			fmt.Println("\033[34m[*] Selected: Full Scan\033[0m")
		case "4":
			*stealthMode = true
			*quickProfile = true
			fmt.Println("\033[34m[*] Selected: Stealth Scan\033[0m")
		default:
			*quickProfile = true
			fmt.Println("\033[34m[*] Selected: Quick Scan\033[0m")
		}

		*serviceDetect = true
		*bannerGrab = true
		fmt.Println("\033[34m[*] Service detection: enabled\033[0m")
		fmt.Println()
	}

	// Discovery mode
	if *discoverMode {
		runDiscovery(*targetFlag, *threads, *timeout)
		return
	}

	// Setup stealth config
	stealthConfig := stealth.DefaultStealthConfig()
	if *stealthMode {
		stealthConfig.Enabled = true
		*threads = 10
		fmt.Println("\033[35m[*] Stealth mode: enabled (slow scan, random delays)\033[0m")
	}

	// Profile check
	var portList []int
	if *quickProfile {
		profile := scanner.QuickProfile
		portList = profile.Ports
		if !*stealthMode {
			*threads = profile.Threads
			*timeout = time.Duration(profile.TimeoutMs) * time.Millisecond
		}
		fmt.Printf("\033[35m[*] Profile: QUICK (%d ports)\033[0m\n", len(portList))
	} else if *aggressiveProfile {
		portList = scanner.GenerateFullPortRange()
		*threads = 500
		*timeout = 200 * time.Millisecond
		fmt.Println("\033[35m[*] Profile: AGGRESSIVE (65535 ports)\033[0m")
	} else if *topPorts > 0 {
		if *topPorts <= 25 {
			portList = ports.TopPorts[:min(*topPorts, len(ports.TopPorts))]
		} else if *topPorts <= 100 {
			portList = ports.Top100Ports[:min(*topPorts, len(ports.Top100Ports))]
		} else {
			for i := 1; i <= *topPorts; i++ {
				portList = append(portList, i)
			}
		}
		fmt.Printf("\033[34m[*] Scanning %d ports\033[0m\n", len(portList))
	} else if *portRange != "" {
		portList = parsePortRange(*portRange)
	} else {
		for i := 1; i <= 1024; i++ {
			portList = append(portList, i)
		}
	}

	// Stealth: randomize port order
	if *stealthMode {
		portList = stealth.ShuffleOrder(portList)
	}

	// Parse targets
	targets, err := target.ParseTargets(*targetFlag)
	if err != nil {
		fmt.Printf("\033[31m[!] Target parse error: %s\033[0m\n", err)
		os.Exit(1)
	}

	// Scan each target
	for _, tgt := range targets {
		results := scanTarget(tgt, portList, *timeout, *threads, *serviceDetect, *bannerGrab,
			*vulnCheck, *sslInfo, *traceroute, *httpInfo, stealthConfig, !*noColor)

		// Save results
		if *outputFile != "" {
			saveResults(*outputFile, tgt, portList, results, time.Duration(0))
		}

		// Discord notification
		if *discordWebhook != "" {
			fmt.Println("\033[34m[*] Sending Discord notification...\033[0m")
			if err := notify.SendDiscord(*discordWebhook, tgt, results, "completed"); err != nil {
				fmt.Printf("\033[33m[!] Discord notification failed: %s\033[0m\n", err)
			} else {
				fmt.Println("\033[32m[✓] Discord notification sent\033[0m")
			}
		}
	}
}

func runDiscovery(targetStr string, threads int, timeout time.Duration) {
	fmt.Println("\033[36m╔══════════════════════════════════════════════════════════════╗\033[0m")
	fmt.Println("\033[36m║              HOST DISCOVERY MODE                             ║\033[0m")
	fmt.Println("\033[36m╚══════════════════════════════════════════════════════════════╝\033[0m")
	fmt.Println()

	hosts, err := target.ParseTargets(targetStr)
	if err != nil {
		fmt.Printf("\033[31m[!] Error: %s\033[0m\n", err)
		return
	}

	fmt.Printf("\033[34m[*] Scanning %d hosts...\033[0m\n", len(hosts))
	results := discovery.PingSweep(hosts, threads, timeout)
	discovery.PrintDiscoveryResults(results)
}

func scanTarget(tgt string, portList []int, timeout time.Duration, threads int,
	serviceDetect, bannerGrab, vulnCheck, sslInfo, traceroute, httpInfo bool,
	stealthConfig *stealth.StealthConfig, useColor bool) []scanner.Result {

	fmt.Printf("\n\033[36m╔══════════════════════════════════════════════════════════════╗\033[0m\n")
	fmt.Printf("\033[36m║  TARGET: %-54s ║\033[0m\n", tgt)
	fmt.Printf("\033[36m╚══════════════════════════════════════════════════════════════╝\033[0m\n\n")

	// DNS resolution
	ips, err := net.LookupIP(tgt)
	if err != nil {
		fmt.Printf("\033[31m[!] DNS resolution error: %s\033[0m\n", err)
		return nil
	}
	fmt.Printf("\033[34m[*] IP: %s\033[0m\n", ips[0].String())
	fmt.Printf("\033[34m[*] Ports: %d | Threads: %d | Timeout: %s\033[0m\n", len(portList), threads, timeout)
	fmt.Println()

	// Traceroute
	if traceroute {
		fmt.Println("\033[34m[*] Running traceroute...\033[0m")
		hops, _ := network.Traceroute(tgt, 15)
		network.PrintTraceroute(hops)
	}

	// SSL Info
	if sslInfo {
		fmt.Println("\033[34m[*] Checking SSL certificate...\033[0m")
		info, err := network.GetSSLInfo(tgt, 443)
		if err == nil {
			network.PrintSSLInfo(info)
		}
	}

	// HTTP Info
	if httpInfo {
		fmt.Println("\033[34m[*] Detecting HTTP technologies...\033[0m")
		info, err := detect.GetHTTPInfo(tgt, 80, false)
		if err == nil {
			detect.PrintHTTPInfo(info)
		}
	}

	// Port scan
	startTime := time.Now()
	fmt.Printf("\033[33m[*] Starting port scan...\033[0m\n")

	s := scanner.NewScanner(tgt, portList, timeout, threads)
	if serviceDetect {
		s.WithServiceDetection()
	}
	if bannerGrab {
		s.WithBannerGrab()
	}

	// Progress
	progressChan := make(chan int, len(portList))
	done := make(chan bool)

	go func() {
		scanned := 0
		total := len(portList)
		for range progressChan {
			scanned++
			if stealthConfig.Enabled {
				stealthConfig.ApplyDelay()
			}
			percent := float64(scanned) / float64(total) * 100
			bar := strings.Repeat("█", int(percent/5)) + strings.Repeat("░", 20-int(percent/5))
			fmt.Printf("\r\033[34m[*] Scanning: [%s] %.1f%% (%d/%d)\033[0m", bar, percent, scanned, total)
		}
		fmt.Println()
		done <- true
	}()

	s.Scan(progressChan)
	<-done

	elapsed := time.Since(startTime)
	pps := float64(len(portList)) / elapsed.Seconds()

	fmt.Printf("\n\033[32m[✓] Scan completed! Duration: %s (%.0f ports/sec)\033[0m\n", elapsed.Round(time.Millisecond), pps)
	fmt.Printf("\033[32m[✓] Open ports: %d\033[0m\n", s.OpenPortCount())

	if s.OpenPortCount() > 0 {
		output.PrintTable(s.Results, useColor)
	} else {
		fmt.Println("\033[33m[!] No open ports found.\033[0m")
	}

	// Vulnerability Check
	if vulnCheck && s.OpenPortCount() > 0 {
		fmt.Println("\033[34m[*] Running vulnerability check...\033[0m")
		vc := detection.NewVulnChecker(s.Results)
		vulns := vc.Check()
		detection.PrintVulns(vulns)
	}

	return s.Results
}

func saveResults(outputFile, tgt string, portList []int, results []scanner.Result, elapsed time.Duration) {
	var saveErr error

	switch {
	case strings.HasSuffix(outputFile, ".json"):
		scanOutput := output.ScanOutput{
			Target:     tgt,
			TotalPorts: len(portList),
			OpenPorts:  len(results),
			ScanTime:   elapsed.String(),
			Results:    results,
		}
		saveErr = output.WriteJSON(outputFile, scanOutput)
	case strings.HasSuffix(outputFile, ".html"):
		report := output.HTMLReport{
			Target:     tgt,
			TotalPorts: len(portList),
			OpenPorts:  len(results),
			ScanTime:   elapsed.Round(time.Millisecond).String(),
			Results:    results,
		}
		saveErr = output.WriteHTML(outputFile, report)
	case strings.HasSuffix(outputFile, ".xml"):
		saveErr = output.WriteXML(outputFile, tgt, results, elapsed)
	case strings.HasSuffix(outputFile, ".csv"):
		saveErr = export.WriteCSV(outputFile, tgt, results)
	case strings.HasSuffix(outputFile, ".md"):
		report := export.MarkdownReport{
			Target:     tgt,
			TotalPorts: len(portList),
			OpenPorts:  len(results),
			ScanTime:   elapsed.Round(time.Millisecond).String(),
			Results:    results,
		}
		saveErr = export.WriteMarkdown(outputFile, report)
	default:
		scanOutput := output.ScanOutput{
			Target:     tgt,
			TotalPorts: len(portList),
			OpenPorts:  len(results),
			ScanTime:   elapsed.String(),
			Results:    results,
		}
		saveErr = output.WriteTXT(outputFile, scanOutput)
	}

	if saveErr != nil {
		fmt.Printf("\033[31m[!] File write error: %s\033[0m\n", saveErr)
	} else {
		fmt.Printf("\033[32m[✓] Results saved: %s\033[0m\n", outputFile)
	}
}

func parsePortRange(input string) []int {
	var result []int
	parts := strings.Split(input, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.Contains(part, "-") {
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) == 2 {
				start, _ := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
				end, _ := strconv.Atoi(strings.TrimSpace(rangeParts[1]))
				for i := start; i <= end; i++ {
					result = append(result, i)
				}
			}
		} else {
			port, _ := strconv.Atoi(part)
			if port > 0 {
				result = append(result, port)
			}
		}
	}
	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
