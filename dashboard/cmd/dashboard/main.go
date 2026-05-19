package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"monitoring.local/api"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	refreshInterval = 2 * time.Second
	historyLimit    = 30
)

var (
	blue   = lipgloss.Color("39")
	green  = lipgloss.Color("42")
	red    = lipgloss.Color("203")
	yellow = lipgloss.Color("220")
	muted  = lipgloss.Color("245")

	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(blue)
	mutedStyle = lipgloss.NewStyle().Foreground(muted)
	errorStyle = lipgloss.NewStyle().Foreground(red).Bold(true)

	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(1, 2)

	selectedRowStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("24")).
				Bold(true)
)

type DashboardClient struct {
	baseURL string
	client  *http.Client
}

type model struct {
	client          *DashboardClient
	preferredHostID string

	devices       []api.Device
	selectedIndex int
	latest        api.Metrics
	history       []api.Metrics

	width   int
	height  int
	loading bool
	err     error
}

type tickMsg time.Time

type dataMsg struct {
	devices       []api.Device
	selectedIndex int
	latest        api.Metrics
	history       []api.Metrics
}

type errMsg struct{ err error }

func main() {
	baseURL := os.Getenv("API_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	m := model{
		client: &DashboardClient{
			baseURL: strings.TrimRight(baseURL, "/"),
			client:  &http.Client{Timeout: 5 * time.Second},
		},
		preferredHostID: os.Getenv("HOST_ID"),
		loading:         true,
	}

	program := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		fmt.Println("dashboard error:", err)
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(fetchDataCmd(m.client, m.preferredHostID, m.selectedIndex), tickCmd())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "r":
			m.loading = true
			return m, fetchDataCmd(m.client, "", m.selectedIndex)
		case "down", "j", "n", "tab":
			if len(m.devices) > 0 {
				m.selectedIndex = (m.selectedIndex + 1) % len(m.devices)
				m.loading = true
				return m, fetchDataCmd(m.client, "", m.selectedIndex)
			}
		case "up", "k", "p", "shift+tab":
			if len(m.devices) > 0 {
				m.selectedIndex--
				if m.selectedIndex < 0 {
					m.selectedIndex = len(m.devices) - 1
				}
				m.loading = true
				return m, fetchDataCmd(m.client, "", m.selectedIndex)
			}
		}

	case tickMsg:
		return m, tea.Batch(fetchDataCmd(m.client, "", m.selectedIndex), tickCmd())

	case dataMsg:
		m.devices = msg.devices
		m.selectedIndex = msg.selectedIndex
		m.latest = msg.latest
		m.history = msg.history
		m.loading = false
		m.err = nil
		m.preferredHostID = ""
		return m, nil

	case errMsg:
		m.err = msg.err
		m.loading = false
		return m, nil
	}

	return m, nil
}

func (m model) View() string {
	if m.width == 0 {
		return "loading..."
	}

	contentWidth := m.width - 4
	if contentWidth < 72 {
		contentWidth = 72
	}

	sections := []string{
		renderHeader(m, contentWidth),
		renderDevices(m, contentWidth),
		renderLatest(m, contentWidth),
		renderHistory(m, contentWidth),
		renderFooter(m),
	}

	return appStyle.Width(contentWidth).Render(strings.Join(sections, "\n"))
}

func fetchDataCmd(client *DashboardClient, preferredHostID string, selectedIndex int) tea.Cmd {
	return func() tea.Msg {
		devices, err := client.GetDevices()
		if err != nil {
			return errMsg{err: err}
		}

		if len(devices) == 0 {
			return dataMsg{devices: devices, selectedIndex: 0}
		}

		selectedIndex = resolveSelectedIndex(devices, preferredHostID, selectedIndex)
		selectedHostID := devices[selectedIndex].HostID

		latest, err := client.GetLatestMetrics(selectedHostID)
		if err != nil {
			return errMsg{err: err}
		}

		history, err := client.GetMetricsHistory(selectedHostID, historyLimit)
		if err != nil {
			return errMsg{err: err}
		}

		return dataMsg{devices: devices, selectedIndex: selectedIndex, latest: latest, history: history}
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(refreshInterval, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func resolveSelectedIndex(devices []api.Device, preferredHostID string, selectedIndex int) int {
	if preferredHostID != "" {
		for i, d := range devices {
			if d.HostID == preferredHostID {
				return i
			}
		}
	}
	if selectedIndex >= 0 && selectedIndex < len(devices) {
		return selectedIndex
	}
	for i, d := range devices {
		if d.Online {
			return i
		}
	}
	return 0
}

func (d *DashboardClient) GetDevices() ([]api.Device, error) {
	resp, err := d.client.Get(d.baseURL + "/api/v1/devices")
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
	endpoint := fmt.Sprintf("%s/api/v1/devices/metrics?host_id=%s&limit=%d", d.baseURL, url.QueryEscape(hostID), limit)
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

func renderHeader(m model, width int) string {
	state := greenText("READY")
	if m.loading {
		state = yellowText("REFRESHING")
	}
	if m.err != nil {
		state = redText("ERROR")
	}

	online, offline := deviceCounts(m.devices)
	line1 := titleStyle.Render("Linux Nodes Monitoring")
	line2 := mutedStyle.Render(fmt.Sprintf("API: %s | refresh: %s | %s", m.client.baseURL, refreshInterval, state))
	line3 := fmt.Sprintf("devices: %d | online: %s | offline: %s", len(m.devices), greenText(fmt.Sprint(online)), redText(fmt.Sprint(offline)))

	if m.err != nil {
		line3 += " | " + errorStyle.Render(m.err.Error())
	}

	return panelStyle.Width(width).Render(strings.Join([]string{line1, line2, line3}, "\n"))
}

func renderDevices(m model, width int) string {
	lines := []string{titleStyle.Render("Devices")}

	if len(m.devices) == 0 {
		lines = append(lines, mutedStyle.Render("No devices found"))
		return panelStyle.Width(width).Render(strings.Join(lines, "\n"))
	}

	for i, d := range m.devices {
		marker := " "
		if i == m.selectedIndex {
			marker = ">"
		}

		status := redText("offline")
		if d.Online {
			status = greenText("online")
		}

		line := fmt.Sprintf("%s %-34s %-8s last_seen=%s", marker, shortID(d.HostID, 34), status, formatTime(d.LastSeenAt))
		if i == m.selectedIndex {
			line = selectedRowStyle.Render(line)
		}
		lines = append(lines, line)
	}

	return panelStyle.Width(width).Render(strings.Join(lines, "\n"))
}

func renderLatest(m model, width int) string {
	lines := []string{titleStyle.Render("Latest metrics")}

	if len(m.devices) == 0 {
		lines = append(lines, mutedStyle.Render("No selected device"))
		return panelStyle.Width(width).Render(strings.Join(lines, "\n"))
	}

	selected := m.devices[m.selectedIndex]
	barWidth := width - 36
	if barWidth < 16 {
		barWidth = 16
	}
	if barWidth > 60 {
		barWidth = 60
	}

	lines = append(lines,
		fmt.Sprintf("host_id: %s", selected.HostID),
		fmt.Sprintf("last_seen: %s", formatDateTime(selected.LastSeenAt)),
		"",
		fmt.Sprintf("CPU  %s %7.2f%%", progressBar(m.latest.CPU.TotalPct, barWidth, blue), m.latest.CPU.TotalPct),
		fmt.Sprintf("MEM  %s %7.2f%%", progressBar(m.latest.Mem.UsedPct, barWidth, green), m.latest.Mem.UsedPct),
		fmt.Sprintf("RX   %-14s", formatBps(m.latest.Network.RxBpsTotal)),
		fmt.Sprintf("TX   %-14s", formatBps(m.latest.Network.TxBpsTotal)),
	)

	return panelStyle.Width(width).Render(strings.Join(lines, "\n"))
}

func renderHistory(m model, width int) string {
	lines := []string{titleStyle.Render("History")}

	if len(m.history) == 0 {
		lines = append(lines, mutedStyle.Render("No history"))
		return panelStyle.Width(width).Render(strings.Join(lines, "\n"))
	}

	sparkWidth := width - 16
	if sparkWidth < 20 {
		sparkWidth = 20
	}
	if sparkWidth > 80 {
		sparkWidth = 80
	}

	lines = append(lines,
		"CPU "+blueText(sparkline(collectCPU(m.history), sparkWidth)),
		"MEM "+greenText(sparkline(collectMEM(m.history), sparkWidth)),
		"",
		mutedStyle.Render("time      cpu       mem       rx          tx"),
	)

	for _, row := range recentRows(m.history, 8) {
		lines = append(lines, row)
	}

	return panelStyle.Width(width).Render(strings.Join(lines, "\n"))
}

func renderFooter(m model) string {
	return mutedStyle.Render("↑/k previous | ↓/j next | r refresh | q quit")
}

func progressBar(value float64, width int, color lipgloss.Color) string {
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
	fill := lipgloss.NewStyle().Foreground(color).Render(strings.Repeat("█", filled))
	empty := lipgloss.NewStyle().Foreground(lipgloss.Color("238")).Render(strings.Repeat("░", width-filled))
	return fill + empty
}

func sparkline(values []float64, width int) string {
	if len(values) == 0 {
		return "no history"
	}

	if len(values) > width {
		values = values[len(values)-width:]
	}

	minValue := values[0]
	maxValue := values[0]
	for _, value := range values {
		if value < minValue {
			minValue = value
		}
		if value > maxValue {
			maxValue = value
		}
	}

	var b strings.Builder

	if maxValue-minValue < 0.01 {
		for range values {
			b.WriteString("▃")
		}
		return b.String()
	}

	for _, value := range values {
		normalized := (value - minValue) / (maxValue - minValue) * 100
		b.WriteString(sparkChar(normalized))
	}

	return b.String()
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

func recentRows(history []api.Metrics, limit int) []string {
	rows := make([]string, 0, limit)
	count := min(limit, len(history))
	for i := 0; i < count; i++ {
		m := history[i]
		rows = append(rows, fmt.Sprintf(
			"%s  %7.2f%%  %7.2f%%  %-10s  %-10s",
			formatTime(m.Timestamp),
			m.CPU.TotalPct,
			m.Mem.UsedPct,
			formatBps(m.Network.RxBpsTotal),
			formatBps(m.Network.TxBpsTotal),
		))
	}
	return rows
}

func collectCPU(history []api.Metrics) []float64 {
	values := make([]float64, 0, len(history))
	for i := len(history) - 1; i >= 0; i-- {
		values = append(values, history[i].CPU.TotalPct)
	}
	return values
}

func collectMEM(history []api.Metrics) []float64 {
	values := make([]float64, 0, len(history))
	for i := len(history) - 1; i >= 0; i-- {
		values = append(values, history[i].Mem.UsedPct)
	}
	return values
}

func deviceCounts(devices []api.Device) (int, int) {
	online := 0
	for _, d := range devices {
		if d.Online {
			online++
		}
	}
	return online, len(devices) - online
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

func formatDateTime(t time.Time) string {
	if t.IsZero() {
		return "—"
	}
	return t.Format("2006-01-02 15:04:05")
}

func shortID(value string, maxLen int) string {
	if value == "" {
		return "none"
	}
	if len(value) <= maxLen {
		return value
	}
	if maxLen <= 3 {
		return value[:maxLen]
	}
	return value[:maxLen-3] + "..."
}

func greenText(s string) string  { return lipgloss.NewStyle().Foreground(green).Bold(true).Render(s) }
func redText(s string) string    { return lipgloss.NewStyle().Foreground(red).Bold(true).Render(s) }
func yellowText(s string) string { return lipgloss.NewStyle().Foreground(yellow).Bold(true).Render(s) }
func blueText(s string) string   { return lipgloss.NewStyle().Foreground(blue).Bold(true).Render(s) }

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
