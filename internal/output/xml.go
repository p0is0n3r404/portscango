package output

import (
	"encoding/xml"
	"fmt"
	"os"
	"time"

	"portscango/internal/scanner"
)

// NmapRun Nmap uyumlu XML root elementi
type NmapRun struct {
	XMLName          xml.Name `xml:"nmaprun"`
	Scanner          string   `xml:"scanner,attr"`
	Args             string   `xml:"args,attr"`
	Start            int64    `xml:"start,attr"`
	StartStr         string   `xml:"startstr,attr"`
	Version          string   `xml:"version,attr"`
	XMLOutputVersion string   `xml:"xmloutputversion,attr"`
	ScanInfo         ScanInfo `xml:"scaninfo"`
	Host             Host     `xml:"host"`
	RunStats         RunStats `xml:"runstats"`
}

// ScanInfo tarama bilgisi
type ScanInfo struct {
	Type     string `xml:"type,attr"`
	Protocol string `xml:"protocol,attr"`
}

// Host hedef bilgisi
type Host struct {
	StartTime int64     `xml:"starttime,attr"`
	EndTime   int64     `xml:"endtime,attr"`
	Status    Status    `xml:"status"`
	Address   Address   `xml:"address"`
	Hostnames Hostnames `xml:"hostnames"`
	Ports     Ports     `xml:"ports"`
}

// Status durum
type Status struct {
	State  string `xml:"state,attr"`
	Reason string `xml:"reason,attr"`
}

// Address adres
type Address struct {
	Addr     string `xml:"addr,attr"`
	AddrType string `xml:"addrtype,attr"`
}

// Hostnames hostname listesi
type Hostnames struct {
	Hostname []Hostname `xml:"hostname"`
}

// Hostname tek hostname
type Hostname struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
}

// Ports port listesi
type Ports struct {
	Port []Port `xml:"port"`
}

// Port tek port
type Port struct {
	Protocol string  `xml:"protocol,attr"`
	PortID   int     `xml:"portid,attr"`
	State    State   `xml:"state"`
	Service  Service `xml:"service"`
}

// State port durumu
type State struct {
	State  string `xml:"state,attr"`
	Reason string `xml:"reason,attr"`
}

// Service servis bilgisi
type Service struct {
	Name    string `xml:"name,attr"`
	Product string `xml:"product,attr,omitempty"`
	Version string `xml:"version,attr,omitempty"`
}

// RunStats tarama istatistikleri
type RunStats struct {
	Finished Finished `xml:"finished"`
	Hosts    Hosts    `xml:"hosts"`
}

// Finished bitiş bilgisi
type Finished struct {
	Time    int64  `xml:"time,attr"`
	TimeStr string `xml:"timestr,attr"`
	Elapsed string `xml:"elapsed,attr"`
}

// Hosts host sayıları
type Hosts struct {
	Up    int `xml:"up,attr"`
	Down  int `xml:"down,attr"`
	Total int `xml:"total,attr"`
}

// WriteXML Nmap uyumlu XML dosyası oluşturur
func WriteXML(filename string, target string, results []scanner.Result, elapsed time.Duration) error {
	startTime := time.Now().Add(-elapsed)

	var ports []Port
	for _, r := range results {
		ports = append(ports, Port{
			Protocol: "tcp",
			PortID:   r.Port,
			State:    State{State: r.State, Reason: "syn-ack"},
			Service:  Service{Name: r.Service, Product: r.Banner},
		})
	}

	nmapRun := NmapRun{
		Scanner:          "portscango",
		Args:             fmt.Sprintf("portscango -t %s", target),
		Start:            startTime.Unix(),
		StartStr:         startTime.Format("Mon Jan 2 15:04:05 2006"),
		Version:          "3.0",
		XMLOutputVersion: "1.05",
		ScanInfo: ScanInfo{
			Type:     "connect",
			Protocol: "tcp",
		},
		Host: Host{
			StartTime: startTime.Unix(),
			EndTime:   time.Now().Unix(),
			Status:    Status{State: "up", Reason: "conn-refused"},
			Address:   Address{Addr: target, AddrType: "ipv4"},
			Hostnames: Hostnames{
				Hostname: []Hostname{{Name: target, Type: "user"}},
			},
			Ports: Ports{Port: ports},
		},
		RunStats: RunStats{
			Finished: Finished{
				Time:    time.Now().Unix(),
				TimeStr: time.Now().Format("Mon Jan 2 15:04:05 2006"),
				Elapsed: fmt.Sprintf("%.2f", elapsed.Seconds()),
			},
			Hosts: Hosts{Up: 1, Down: 0, Total: 1},
		},
	}

	output, err := xml.MarshalIndent(nmapRun, "", "  ")
	if err != nil {
		return err
	}

	xmlHeader := []byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	return os.WriteFile(filename, append(xmlHeader, output...), 0644)
}
