package main

import (
	"encoding/json"
	"fmt"
	"log"
	"monitoring/api"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type DashboardClient struct {
	baseURL string
	client  *http.Client
}

func main() {
	baseURL := os.Getenv("API_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	selectedHostID := os.Getenv("HOST_ID")
	dashboard := &DashboardClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
	for {
		devices, err := dashboard.GetDevices()
		if err != nil {
			clearScreen()
			fmt.Println("Linux Nodes Monitoring")
			fmt.Println()
			fmt.Println("failed to load devices:", err)
			time.Sleep(2 * time.Second)
			continue
		}
		if selectedHostID == "" && len(devices) > 0 {
			selectedHostID = chooseDefaultDevice(devices)
		}
		var latest api.Metrics
		var history []api.Metrics
		if selectedHostID != "" {
			latest, err = dashboard.GetLatestMetrics(selectedHostID)
			if err != nil {
				latest = api.Metrics{}
			}
			history, err = dashboard.GetMetricsHistory(selectedHostID, 30)
			if err != nil {
				history = nil
			}
		}
		renderDashboard(devices, selectedHostID, latest, history)
		time.Sleep(2 * time.Second)
	}
}

func (d *DashboardClient) GetDevices() ([]api.Device, error) {
	endpoint := d.baseURL + "/api/v1/devices"
	resp, err := d.client.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("get devices: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get devices: unexpected status %d", resp.StatusCode)
	}
	var devices []api.Device
	if err := json.NewDecoder(resp.Body).Decode(&devices); err != nil {
		return nil, fmt.Errorf("decode devices: %w", err)
	}
	return devices, nil
}

func (d *DashboardClient) GetLatestMetrics(hostID string) (api.Metrics, error) {
	endpoint := d.baseURL + "/api/v1/devices/latest?host_id=" + url.QueryEscape(hostID)
	resp, err := d.client.Get(endpoint)
	if err != nil {
		return api.Metrics{}, fmt.Errorf("get latest metrics: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return api.Metrics{}, fmt.Errorf("get latest metrics: unexpected status %d", resp.StatusCode)
	}
	var metrics api.Metrics
	if err := json.NewDecoder(resp.Body).Decode(&metrics); err != nil {
		return api.Metrics{}, fmt.Errorf("decode latest metrics: %w", err)
	}
	return metrics, nil
}

func (d *DashboardClient) GetMetricsHistory(hostID string, limit int) ([]api.Metrics, error) {
	endpoint := fmt.Sprintf(
		"%s/api/v1/devices/metrics?host_id=%s&limit=%d",
		d.baseURL,
		url.QueryEscape(hostID),
		limit,
	)
	resp, err := d.client.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("get metrics history: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get metrics history: unexpected status %d", resp.StatusCode)
	}
	var history []api.Metrics
	if err := json.NewDecoder(resp.Body).Decode(&history); err != nil {
		return nil, fmt.Errorf("decode metrics history: %w", err)
	}
	return history, nil
}

func renderDashboard(devices []api.Device, selectedHostID string, latest api.Metrics, history []api.Metrics) {
	clearScreen()
	fmt.Println("┌──────────────────────────── Linux Nodes Monitoring ────────────────────────────┐")
	fmt.Printf("│ API dashboard | refresh: 2s | devices: %-38d │\n", len(devices))
	fmt.Println("├────────────────────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ DEVICES                                                                        │")
	if len(devices) == 0 {
		fmt.Println("│   no devices found                                                              │")
	} else {
		for _, d := range devices {
			prefix := " "
			if d.HostID == selectedHostID {
				prefix = ">"
			}
			status := "offline"
			if d.Online {
				status = "online "
			}
			line := fmt.Sprintf(
				"│ %s %-36s %-8s last seen: %-22s │",
				prefix,
				shortHostID(d.HostID),
				status,
				formatTime(d.LastSeenAt),
			)
			fmt.Println(trimToWidth(line, 82))
		}
	}
	fmt.Println("├────────────────────────────────────────────────────────────────────────────────┤")
	fmt.Printf("│ SELECTED NODE: %-64s │\n", shortHostID(selectedHostID))
	fmt.Println("├────────────────────────────────────────────────────────────────────────────────┤")
	fmt.Printf("│ CPU     %s %8s                                                   │\n",
		bar(latest.CPU.TotalPct, 30),
		formatPercent(latest.CPU.TotalPct),
	)
	fmt.Printf("│ MEMORY  %s %8s                                                   │\n",
		bar(latest.Mem.UsedPct, 30),
		formatPercent(latest.Mem.UsedPct),
	)
	fmt.Printf("│ RX      %-20s                                                        │\n",
		formatBps(latest.Network.RxBpsTotal),
	)
	fmt.Printf("│ TX      %-20s                                                        │\n",
		formatBps(latest.Network.TxBpsTotal),
	)
	fmt.Println("├────────────────────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ CPU HISTORY                                                                     │")
	fmt.Printf("│ %s │\n", padRight(cpuSparkline(history), 78))
	fmt.Println("└────────────────────────────────────────────────────────────────────────────────┘")
}

func chooseDefaultDevice(devices []api.Device) string {
	for _, d := range devices {
		if d.Online {
			return d.HostID
		}
	}
	return devices[0].HostID
}

func cpuSparkline(history []api.Metrics) string {
	if len(history) == 0 {
		return "no history"
	}
	values := make([]float64, 0, len(history))
	for i := len(history) - 1; i >= 0; i-- {
		values = append(values, history[i].CPU.TotalPct)
	}
	var builder strings.Builder

	for _, value := range values {
		builder.WriteString(sparkChar(value))
	}
	return builder.String()
}

func sparkChar(value float64) string {
	switch {
	case value < 10:
		return "▁"
	case value < 25:
		return "▂"
	case value < 40:
		return "▃"
	case value < 55:
		return "▄"
	case value < 70:
		return "▅"
	case value < 85:
		return "▆"
	default:
		return "█"
	}
}

func bar(value float64, width int) string {
	if value < 0 {
		value = 0
	}
	if value > 100 {
		value = 100
	}
	filled := int(value / 100 * float64(width))
	if filled > width {
		filled = width
	}
	return "[" + strings.Repeat("█", filled) + strings.Repeat("░", width-filled) + "]"
}

func formatPercent(value float64) string {
	return fmt.Sprintf("%.2f%%", value)
}

func formatBps(value float64) string {
	if value >= 1024*1024 {
		return fmt.Sprintf("%.2f MB/s", value/1024/1024)
	}
	if value >= 1024 {
		return fmt.Sprintf("%.2f KB/s", value/1024)
	}
	return fmt.Sprintf("%.0f B/s", value)
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return "—"
	}
	return t.Format("15:04:05")
}

func shortHostID(hostID string) string {
	if hostID == "" {
		return "none"
	}
	if len(hostID) <= 36 {
		return hostID
	}
	return hostID[:36]
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func padRight(value string, width int) string {
	if len([]rune(value)) >= width {
		return value
	}
	return value + strings.Repeat(" ", width-len([]rune(value)))
}

func trimToWidth(value string, width int) string {
	runes := []rune(value)
	if len(runes) <= width {
		return value
	}
	return string(runes[:width])
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
