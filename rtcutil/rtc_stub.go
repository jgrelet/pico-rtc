//go:build !tinygo || (!pico && !pico_w && !rp2040 && !pico2 && !pico2_w && !rp2350)

package rtcutil

import "time"

type rtcStub struct{}

func newRTC() RTC { return &rtcStub{} }

func (s *rtcStub) Init1Hz(_ uint32)   {}
func (s *rtcStub) Set(_ time.Time)    {}
func (s *rtcStub) Now() time.Time     { return time.Unix(0, 0).UTC() }

