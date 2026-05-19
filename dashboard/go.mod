module monitoring.local/dashboard

go 1.25.5

require (
	github.com/charmbracelet/bubbletea v1.3.4
	github.com/charmbracelet/lipgloss v1.1.0
	monitoring.local/api v0.0.0
)

replace monitoring.local/api => ../api