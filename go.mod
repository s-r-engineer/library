module github.com/s-r-engineer/library

go 1.24.2

replace github.com/getsentry/sentry-go => github.com/s-r-engineer/sentry-go v0.0.0-20250510204047-972e22263f6a

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/getsentry/sentry-go v0.32.0
	github.com/influxdata/line-protocol/v2 v2.2.1
	github.com/stretchr/testify v1.10.0
	go.uber.org/multierr v1.11.0
	golang.org/x/crypto v0.37.0
	golang.org/x/sync v0.13.0
	golang.org/x/term v0.31.0
)

require (
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
