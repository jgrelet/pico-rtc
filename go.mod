module ntp-client/new

go 1.25.0

replace rtcutil => ./rtcutil

require (
	github.com/soypat/cyw43439 v0.0.0-20250505012923-830110c8f4af
	github.com/soypat/seqs v0.0.0-20250630134107-01c3f05666ba
	rtcutil v0.0.0-00010101000000-000000000000
)

require (
	github.com/tinygo-org/pio v0.2.0 // indirect
	golang.org/x/exp v0.0.0-20240808152545-0cdaa3abc0fa // indirect
)
