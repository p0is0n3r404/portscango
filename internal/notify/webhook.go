package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"portscango/internal/scanner"
)

// WebhookPayload generic webhook payload
type WebhookPayload struct {
	Target    string           `json:"target"`
	OpenPorts int              `json:"open_ports"`
	ScanTime  string           `json:"scan_time"`
	Timestamp string           `json:"timestamp"`
	Results   []scanner.Result `json:"results"`
}

// DiscordEmbed Discord embed structure
type DiscordEmbed struct {
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Color       int            `json:"color"`
	Fields      []DiscordField `json:"fields"`
	Footer      DiscordFooter  `json:"footer"`
	Timestamp   string         `json:"timestamp"`
}

// DiscordField embed field
type DiscordField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

// DiscordFooter embed footer
type DiscordFooter struct {
	Text string `json:"text"`
}

// DiscordWebhook Discord webhook message
type DiscordWebhook struct {
	Username  string         `json:"username"`
	AvatarURL string         `json:"avatar_url"`
	Embeds    []DiscordEmbed `json:"embeds"`
}

// SendDiscord sends scan results to Discord webhook
func SendDiscord(webhookURL string, target string, results []scanner.Result, scanTime string) error {
	// Build port list
	var portList string
	for i, r := range results {
		if i > 0 {
			portList += ", "
		}
		portList += fmt.Sprintf("`%d` (%s)", r.Port, r.Service)
		if i >= 9 {
			portList += fmt.Sprintf(" +%d more", len(results)-10)
			break
		}
	}

	if portList == "" {
		portList = "No open ports found"
	}

	// Color based on results
	color := 0x00FF00 // Green
	if len(results) == 0 {
		color = 0xFF0000 // Red
	} else if len(results) > 10 {
		color = 0xFFA500 // Orange
	}

	webhook := DiscordWebhook{
		Username:  "PortScanGO",
		AvatarURL: "https://raw.githubusercontent.com/portscango/assets/main/logo.png",
		Embeds: []DiscordEmbed{
			{
				Title:       "🔍 Port Scan Complete",
				Description: fmt.Sprintf("Scan results for `%s`", target),
				Color:       color,
				Fields: []DiscordField{
					{Name: "🎯 Target", Value: target, Inline: true},
					{Name: "🔓 Open Ports", Value: fmt.Sprintf("%d", len(results)), Inline: true},
					{Name: "⏱️ Duration", Value: scanTime, Inline: true},
					{Name: "📋 Ports", Value: portList, Inline: false},
				},
				Footer: DiscordFooter{
					Text: "PortScanGO v4.0",
				},
				Timestamp: time.Now().Format(time.RFC3339),
			},
		},
	}

	return sendWebhook(webhookURL, webhook)
}

// SendSlack sends scan results to Slack webhook
func SendSlack(webhookURL string, target string, results []scanner.Result, scanTime string) error {
	payload := map[string]interface{}{
		"text": fmt.Sprintf("🔍 *Port Scan Complete*\nTarget: `%s`\nOpen Ports: %d\nDuration: %s",
			target, len(results), scanTime),
	}

	return sendWebhook(webhookURL, payload)
}

// SendCustomWebhook sends results to custom webhook
func SendCustomWebhook(webhookURL string, target string, results []scanner.Result, scanTime string) error {
	payload := WebhookPayload{
		Target:    target,
		OpenPorts: len(results),
		ScanTime:  scanTime,
		Timestamp: time.Now().Format(time.RFC3339),
		Results:   results,
	}

	return sendWebhook(webhookURL, payload)
}

// sendWebhook sends HTTP POST request
func sendWebhook(url string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	return nil
}
