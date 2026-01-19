package output

import (
	"fmt"
	"os"
	"time"

	"portscango/internal/scanner"
)

// HTMLReport HTML rapor yapısı
type HTMLReport struct {
	Target     string
	IP         string
	ScanTime   string
	TotalPorts int
	OpenPorts  int
	Results    []scanner.Result
	OSInfo     string
	Vulns      int
}

// WriteHTML HTML rapor oluşturur
func WriteHTML(filename string, report HTMLReport) error {
	html := generateHTML(report)
	return os.WriteFile(filename, []byte(html), 0644)
}

func generateHTML(r HTMLReport) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="tr">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>PortScanGO - Scan Report</title>
    <style>
        :root {
            --bg-primary: #0a0a0f;
            --bg-secondary: #12121a;
            --bg-card: #1a1a24;
            --text-primary: #e4e4e7;
            --text-secondary: #a1a1aa;
            --accent: #00d4ff;
            --accent-green: #10b981;
            --accent-red: #ef4444;
            --accent-yellow: #f59e0b;
            --accent-purple: #8b5cf6;
        }
        
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: 'Segoe UI', system-ui, sans-serif;
            background: var(--bg-primary);
            color: var(--text-primary);
            line-height: 1.6;
            min-height: 100vh;
        }
        
        .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 40px 20px;
        }
        
        .header {
            text-align: center;
            margin-bottom: 40px;
            padding: 40px;
            background: linear-gradient(135deg, var(--bg-secondary) 0%%, var(--bg-card) 100%%);
            border-radius: 16px;
            border: 1px solid rgba(255,255,255,0.1);
        }
        
        .logo {
            font-size: 48px;
            font-weight: 800;
            background: linear-gradient(135deg, var(--accent) 0%%, var(--accent-purple) 100%%);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            margin-bottom: 10px;
        }
        
        .subtitle {
            color: var(--text-secondary);
            font-size: 14px;
        }
        
        .stats-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 20px;
            margin-bottom: 40px;
        }
        
        .stat-card {
            background: var(--bg-card);
            padding: 24px;
            border-radius: 12px;
            border: 1px solid rgba(255,255,255,0.05);
            text-align: center;
            transition: transform 0.2s, box-shadow 0.2s;
        }
        
        .stat-card:hover {
            transform: translateY(-4px);
            box-shadow: 0 10px 40px rgba(0,212,255,0.1);
        }
        
        .stat-value {
            font-size: 36px;
            font-weight: 700;
            color: var(--accent);
        }
        
        .stat-value.green { color: var(--accent-green); }
        .stat-value.red { color: var(--accent-red); }
        .stat-value.yellow { color: var(--accent-yellow); }
        
        .stat-label {
            color: var(--text-secondary);
            font-size: 14px;
            text-transform: uppercase;
            letter-spacing: 1px;
            margin-top: 8px;
        }
        
        .section {
            background: var(--bg-card);
            border-radius: 12px;
            padding: 24px;
            margin-bottom: 24px;
            border: 1px solid rgba(255,255,255,0.05);
        }
        
        .section-title {
            font-size: 18px;
            font-weight: 600;
            margin-bottom: 20px;
            display: flex;
            align-items: center;
            gap: 10px;
        }
        
        .section-title::before {
            content: '';
            width: 4px;
            height: 20px;
            background: var(--accent);
            border-radius: 2px;
        }
        
        table {
            width: 100%%;
            border-collapse: collapse;
        }
        
        th, td {
            padding: 14px 16px;
            text-align: left;
            border-bottom: 1px solid rgba(255,255,255,0.05);
        }
        
        th {
            color: var(--text-secondary);
            font-weight: 500;
            font-size: 13px;
            text-transform: uppercase;
            letter-spacing: 1px;
        }
        
        tr:hover {
            background: rgba(255,255,255,0.02);
        }
        
        .port-open {
            display: inline-flex;
            align-items: center;
            gap: 6px;
            color: var(--accent-green);
            font-weight: 500;
        }
        
        .port-open::before {
            content: '';
            width: 8px;
            height: 8px;
            background: var(--accent-green);
            border-radius: 50%%;
            animation: pulse 2s infinite;
        }
        
        @keyframes pulse {
            0%%, 100%% { opacity: 1; }
            50%% { opacity: 0.5; }
        }
        
        .banner {
            font-family: 'Consolas', monospace;
            font-size: 12px;
            color: var(--text-secondary);
            background: var(--bg-secondary);
            padding: 4px 8px;
            border-radius: 4px;
            max-width: 300px;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
        }
        
        .footer {
            text-align: center;
            padding: 40px;
            color: var(--text-secondary);
            font-size: 13px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="logo">⚡ PortScanGO</div>
            <div class="subtitle">High-Performance Port Scanner Report</div>
        </div>
        
        <div class="stats-grid">
            <div class="stat-card">
                <div class="stat-value">%s</div>
                <div class="stat-label">Hedef</div>
            </div>
            <div class="stat-card">
                <div class="stat-value">%d</div>
                <div class="stat-label">Taranan Port</div>
            </div>
            <div class="stat-card">
                <div class="stat-value green">%d</div>
                <div class="stat-label">Açık Port</div>
            </div>
            <div class="stat-card">
                <div class="stat-value yellow">%s</div>
                <div class="stat-label">Tarama Süresi</div>
            </div>
        </div>
        
        <div class="section">
            <div class="section-title">Açık Portlar</div>
            <table>
                <thead>
                    <tr>
                        <th>Port</th>
                        <th>Durum</th>
                        <th>Servis</th>
                        <th>Banner</th>
                    </tr>
                </thead>
                <tbody>
                    %s
                </tbody>
            </table>
        </div>
        
        <div class="footer">
            <p>Rapor oluşturma tarihi: %s</p>
            <p>PortScanGO v3.0 - https://github.com/portscango</p>
        </div>
    </div>
</body>
</html>`, r.Target, r.TotalPorts, r.OpenPorts, r.ScanTime, generateTableRows(r.Results), time.Now().Format("2006-01-02 15:04:05"))
}

func generateTableRows(results []scanner.Result) string {
	var rows string
	for _, r := range results {
		banner := r.Banner
		if banner == "" {
			banner = "-"
		}
		rows += fmt.Sprintf(`
                    <tr>
                        <td><strong>%d</strong></td>
                        <td><span class="port-open">OPEN</span></td>
                        <td>%s</td>
                        <td><span class="banner">%s</span></td>
                    </tr>`, r.Port, r.Service, banner)
	}
	return rows
}
