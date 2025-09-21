//go:build tinygo && (pico || pico_w || rp2040)

package rtcutil

import (
	drp "device/rp" // RP2040
	"time"
)

const (
	ctrlEnable = 1 << 0
	ctrlActive = 1 << 1
	ctrlLoad   = 1 << 4
)

type rtcRP2040 struct{}

func newRTC() RTC { return &rtcRP2040{} }

func (r *rtcRP2040) Init1Hz(clkRtcHz uint32) {
	if clkRtcHz == 0 {
		clkRtcHz = 46875 // XOSC 12 MHz / 256
	}
	divMinus1 := clkRtcHz - 1

	// Stop + attendre arrêt
	drp.RTC.CTRL.ClearBits(ctrlEnable)
	for (drp.RTC.CTRL.Get() & ctrlActive) != 0 {
	}

	// Régler le diviseur 1 Hz (N-1)
	drp.RTC.CLKDIV_M1.Set(uint32(divMinus1))
}

func (r *rtcRP2040) Set(t time.Time) {
	y, m, d := t.Date()
	hh, mm, ss := t.Clock()
	dotw := int(t.Weekday()) // 0=dimanche…6=samedi

	// Stop + attendre arrêt
	drp.RTC.CTRL.ClearBits(ctrlEnable)
	for (drp.RTC.CTRL.Get() & ctrlActive) != 0 {
	}

	// SETUP_0 : YEAR[23:12] | MONTH[11:8] | DAY[4:0]
	drp.RTC.SETUP_0.Set((uint32(y) << 12) | (uint32(m) << 8) | uint32(d))
	// SETUP_1 : DOTW[26:24] | HOUR[20:16] | MIN[13:8] | SEC[5:0]
	drp.RTC.SETUP_1.Set((uint32(dotw) << 24) | (uint32(hh) << 16) | (uint32(mm) << 8) | uint32(ss))

	// LOAD + ENABLE + attendre ACTIVE
	drp.RTC.CTRL.SetBits(ctrlLoad)
	drp.RTC.CTRL.SetBits(ctrlEnable)
	for (drp.RTC.CTRL.Get() & ctrlActive) == 0 {
	}
}

func (r *rtcRP2040) Now() time.Time {
	for {
		// Lire RTC_0 PUIS RTC_1 et relire RTC_0 pour éviter le passage de seconde
		u0a := drp.RTC.RTC_0.Get()
		u1a := drp.RTC.RTC_1.Get()
		u0b := drp.RTC.RTC_0.Get()
		if (u0a & 0x3F) != (u0b & 0x3F) {
			continue
		}

		sec  := int(u0b & 0x3F)
		min  := int((u0b >> 8) & 0x3F)
		hour := int((u0b >> 16) & 0x1F)
		day   := int(u1a & 0x1F)
		month := int((u1a >> 8) & 0x0F)
		year  := int((u1a >> 12) & 0x0FFF)

		return time.Date(year, time.Month(month), day, hour, min, sec, 0, time.UTC)
	}
}
