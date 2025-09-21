//go:build tinygo && (pico || pico_w || rp2040)

package rtcutil

import (
	"device/rp"
	"time"
)

// Bits CTRL (datasheet RP2040)
const (
	ctrlEnable = 1 << 0 // RTC_ENABLE
	ctrlActive = 1 << 1 // RTC_ACTIVE (RO)
	ctrlLoad   = 1 << 4 // LOAD
)

// GetClkRTC lit la source de clk_rtc, calcule clk_rtc en Hz via le diviseur Q24.8.
func GetClkRTC() uint32 {
	// 1) source de clk_rtc
	src := rp.CLOCKS.CLK_RTC_CTRL.Get() & 0x3

	var srcHz uint32
	switch src {
		case 0:
    		srcHz = 12_000_000 // clk_ref from XOSC
		case 1:
    		srcHz = 125_000_000 // clk_sys (default ~125 MHz)
		case 2:
    		srcHz = 6_000_000   // rosc (very approximate)
		case 3:
    		srcHz = 12_000_000  // xosc direct
	}

	// 2) diviseur Q24.8 (N = raw/256)
	divRaw := rp.CLOCKS.CLK_RTC_DIV.Get()
	if divRaw == 0 {
		return 0
	}
	// clk_rtc = srcHz * 256 / divRaw  (arithmétique entière)
	return (srcHz * 256) / divRaw
}

// CalcDiv retourne la valeur à écrire dans CLKDIV_M1 pour obtenir ~1 Hz
func CalcDiv(clkRtcHz uint32) uint16 {
	if clkRtcHz == 0 {
		return 0
	}
	// DividerMinus1 = clkRtcHz - 1
	return uint16(clkRtcHz - 1)
}

// initRTC1Hz configure le diviseur du RTC pour obtenir un tick de 1 Hz.
// Hypothèse standard Pico/Pico W : clk_rtc ≈ 46_875 Hz (XOSC/256) -> CLKDIV_M1 = 46874.
func InitRTC1Hz() {
	// On peut stopper le RTC avant de toucher aux réglages
	rp.RTC.CTRL.ClearBits(ctrlEnable)
	for (rp.RTC.CTRL.Get() & ctrlActive) != 0 {
	}

	// Diviseur "N-1" (46875 - 1) pour produire 1 Hz
	rp.RTC.CLKDIV_M1.Set(46874)
}

// Set règle la date/heure du calendrier RTC (UTC conseillé).
func Set(t time.Time) {
	year, month, day := t.Date()
	hour, min, sec := t.Clock()
	dotw := int(t.Weekday()) // 0=dimanche…6=samedi

	// Désactiver et attendre l'arrêt
	rp.RTC.CTRL.ClearBits(ctrlEnable)
	for (rp.RTC.CTRL.Get() & ctrlActive) != 0 {
	}

	// SETUP_0 : YEAR[23:12] | MONTH[11:8] | DAY[4:0]
	setup0 := (uint32(year) << 12) | (uint32(month) << 8) | uint32(day)
	rp.RTC.SETUP_0.Set(setup0)

	// SETUP_1 : DOTW[26:24] | HOUR[20:16] | MIN[13:8] | SEC[5:0]
	setup1 := (uint32(dotw) << 24) | (uint32(hour) << 16) | (uint32(min) << 8) | uint32(sec)
	rp.RTC.SETUP_1.Set(setup1)

	// LOAD puis ENABLE, attendre ACTIVE
	rp.RTC.CTRL.SetBits(ctrlLoad)
	rp.RTC.CTRL.SetBits(ctrlEnable)
	for (rp.RTC.CTRL.Get() & ctrlActive) == 0 {
	}
}

// Now lit RTC_0 puis RTC_1 et reconstitue un time.Time (UTC).
func Now() time.Time {
	// Lire RTC_0 AVANT RTC_1 pour latching cohérent
	u0 := rp.RTC.RTC_0.Get()
	u1 := rp.RTC.RTC_1.Get()

	sec  := int(u0 & 0x3F)          // [5:0]
	min  := int((u0 >> 8) & 0x3F)   // [13:8]
	hour := int((u0 >> 16) & 0x1F)  // [20:16]
	day   := int(u1 & 0x1F)         // [4:0]
	month := int((u1 >> 8) & 0x0F)  // [11:8]
	year  := int((u1 >> 12) & 0xFFF)// [23:12]

	return time.Date(year, time.Month(month), day, hour, min, sec, 0, time.UTC)
}
