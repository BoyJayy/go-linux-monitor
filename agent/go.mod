module agent

go 1.25.5

require (
	github.com/google/uuid v1.6.0
	golang.org/x/sys v0.42.0
	monitoring/api v0.0.0
)

replace monitoring/api => ../api
