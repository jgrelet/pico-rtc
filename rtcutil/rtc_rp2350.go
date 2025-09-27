//go:build tinygo && (pico2 || pico2_w || rp2350)

package rtcutil

import (
	"sync"
	"time"
)

// RP2350 : pas de RTC calendrier exposé → Simule via un temps monotone.
// RTC : Now() = base + (monotonic_now - monotonic_base).
// Avantages : pas de registres spécifiques, pas de dérive sur reset logiciel,
// et fonctionne tant que l’horloge monotone TinyGo tourne (toujours le cas).
// Inconvénient : pas de persistance sur coupure totale d’alimentation.
type rtcRP2350 struct {
	mu        sync.Mutex
	setBase   time.Time
	ticksBase uint64
	inited    bool
	freqHz    uint64 // ticks/s pour nowTicks() ; ici 1 tick = 1 ns (UnixNano)
}

// newRTC creates and returns a new instance of RTC configured for the RP2350 platform.
// The returned RTC uses a frequency of 1 GHz (1,000,000,000 Hz).
func newRTC() RTC { return &rtcRP2350{freqHz: 1_000_000_000} }

// Init1Hz initializes the RTC to generate a 1Hz signal.
// In simulation mode, this function is a no-op as there is no hardware divider to configure.
// The input parameter is ignored.
func (r *rtcRP2350) Init1Hz(_ uint32) {
	// No-op en mode simulation (pas de diviseur matériel à régler).
}

// Set initializes the RTC with the provided time value.
// It sets the base time to the given UTC time, records the current tick count,
// and marks the RTC as initialized. This method is thread-safe.
func (r *rtcRP2350) Set(t time.Time) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.setBase = t.UTC()
	r.ticksBase = r.nowTicks()
	r.inited = true
}

// Now returns the current time as maintained by the rtcRP2350 instance.
// If the RTC has not been initialized, it returns the Unix epoch (UTC).
// The method calculates the elapsed ticks since initialization, converts them
// to nanoseconds based on the RTC frequency, and adds the duration to the base time.
// This function is safe for concurrent use.
func (r *rtcRP2350) Now() time.Time {
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.inited {
		return time.Unix(0, 0).UTC()
	}
	cur := r.nowTicks()
	dt := cur - r.ticksBase
	ns := (dt * 1_000_000_000) / r.freqHz
	return r.setBase.Add(time.Duration(ns) * time.Nanosecond)
}

// Backend ticks : par défaut, horloge monotone (1 tick = 1 ns)
// nowTicks returns the current time in nanoseconds as a uint64 value.
// It uses time.Now().UnixNano() to obtain the current timestamp.
// This can be used for high-resolution time measurements.
func (r *rtcRP2350) nowTicks() uint64 {
	return uint64(time.Now().UnixNano())
}

