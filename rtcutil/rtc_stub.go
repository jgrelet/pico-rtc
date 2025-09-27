//go:build !tinygo || (!pico && !pico_w && !rp2040 && !pico2 && !pico2_w && !rp2350)

package rtcutil

import "time"

type rtcStub struct{}

// newRTC creates and returns a new instance of RTC using the rtcStub implementation.
// This function is typically used for testing or as a placeholder when a real RTC is not available.
func newRTC() RTC { return &rtcStub{} }

// Init1Hz is a stub implementation that initializes a 1Hz timer or clock.
// The input parameter is ignored in this stub.
func (s *rtcStub) Init1Hz(_ uint32)   {}

// Set is a stub implementation that does not perform any action.
// It is intended to satisfy the interface requirements for setting the RTC time.
func (s *rtcStub) Set(_ time.Time)    {}

// Now returns a fixed UTC time representing the Unix epoch (January 1, 1970, 00:00:00 UTC).
// This is a stub implementation and does not provide the current time.
func (s *rtcStub) Now() time.Time     { return time.Unix(0, 0).UTC() }

