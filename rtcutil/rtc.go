package rtcutil

import "time"

type rtcStub struct{}

// RTC est l’interface commune.
type RTC interface {
	// Init1Hz calibre le “tick” interne à ~1 Hz (selon la plateforme).
	// Passer 0 utilise la valeur par défaut de la carte (ex: 46875 Hz).
	Init1Hz(clkRtcHz uint32)
	// Set met à l’heure (UTC conseillé).
	Set(t time.Time)
	// Now lit l’heure courante (UTC).
	Now() time.Time
}

// NewRTC retourne l’implémentation courante (RP2040 ou RP2350) selon les build tags.
// La définition concrète de NewRTC se trouve dans rtc_rp2040.go / rtc_rp2350.go / rtc_stub.go.
func NewRTC() RTC { return newRTC() }

