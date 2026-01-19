# PortScanGO 🔍

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux%20%7C%20macOS-blue)](https://github.com/portscango)

**Ultra High-Performance Port Scanner** written in Go. Fast, feature-rich, and professional.

![PortScanGO Banner](assets/banner.png)

## ⚡ Features

- 🚀 **Ultra Fast** - Concurrent scanning with customizable threads
- 🎯 **Multiple Targets** - CIDR, IP ranges, file input support
- 🔍 **Service Detection** - Identify running services
- 📡 **Banner Grabbing** - Extract service banners
- 🛡️ **Vulnerability Check** - Detect known vulnerabilities
- 🔐 **SSL/TLS Info** - Certificate details & expiry
- 🌐 **HTTP Detection** - Tech stack, WAF, CMS detection
- 🥷 **Stealth Mode** - Random delays & shuffled ports
- 🔎 **Host Discovery** - Ping sweep for live hosts
- 💬 **Notifications** - Discord webhook support
- 📊 **Multiple Exports** - JSON, HTML, XML, CSV, Markdown

## 📦 Installation

### From Source

```bash
git clone https://github.com/p0is0n3r404/portscango.git
cd portscango
go build -o portscango.exe .
```

### Download Binary

Download the latest release from [Releases](https://github.com/p0is0n3r404/portscango/releases).

## 🚀 Quick Start

```bash
# Interactive mode - just run without arguments
./portscango

# Quick scan
./portscango -t scanme.nmap.org --quick --service

# Full scan with all features
./portscango -t target.com --aggressive --vuln --ssl --http

# Stealth scan
./portscango -t target.com --stealth

# Network discovery
./portscango --discover 192.168.1.0/24

# Export to Markdown
./portscango -t target.com -o report.md
```

## 📖 Usage

```
Usage: portscango [options]

Target:
  -t string         Target IP/domain/CIDR
  -discover         Host discovery mode

Scan Options:
  -p string         Port range (e.g: 1-1000, 80,443)
  -top int          Scan top N common ports
  -threads int      Concurrent threads (default 100)
  -timeout duration Connection timeout (default 1s)

Profiles:
  -quick            Quick scan (top 20 ports)
  -aggressive       Full scan (65535 ports)
  -stealth          Stealth mode with random delays

Features:
  -service          Service detection
  -banner           Banner grabbing
  -vuln             Vulnerability check
  -ssl              SSL/TLS certificate info
  -http             HTTP info & technology detection
  -traceroute       Traceroute to target

Output:
  -o string         Output file (json/txt/html/xml/csv/md)
  -no-color         Disable colored output

Notifications:
  -discord string   Discord webhook URL

Other:
  -v                Show version
```

## 📸 Screenshots

### Interactive Mode

```
╔══════════════════════════════════════════════════════════════╗
║              WELCOME TO PORTSCANGO v4.0!                     ║
╚══════════════════════════════════════════════════════════════╝

[+] Enter target (IP or domain): scanme.nmap.org

Select scan type:
  [1] Quick Scan   - Top 20 ports, fast
  [2] Normal Scan  - Top 1000 ports
  [3] Full Scan    - All 65535 ports
  [4] Stealth Scan - Slow & undetectable
```

### Scan Results

```
╔══════════════════════════════════════════════════════════════╗
║  TARGET: scanme.nmap.org                                     ║
╚══════════════════════════════════════════════════════════════╝

[*] IP: 45.33.32.156
[*] Ports: 20 | Threads: 200 | Timeout: 500ms

[*] Scanning: [████████████████████] 100.0% (20/20)

[✓] Scan completed! Duration: 502ms (40 ports/sec)
[✓] Open ports: 2

PORT       STATE      SERVICE         BANNER
─────────────────────────────────────────────────────────────────
22         open       SSH             SSH-2.0-OpenSSH_6.6.1p1 Ubuntu...
80         open       HTTP            HTTP/1.1 200 OK
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## ⚠️ Disclaimer

This tool is for educational and authorized security testing purposes only. Users are responsible for obtaining proper authorization before scanning any systems they do not own.

## 🙏 Acknowledgments

- Inspired by Nmap, Masscan, and RustScan
- Built with ❤️ in Go

---

**Made with 💻 by [p0is0n3r404](https://github.com/p0is0n3r404)**
