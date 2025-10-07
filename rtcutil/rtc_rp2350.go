//go:build tinygo && (pico2 || pico2_w || rp2350)

package rtcutil

import (
	"sync"
	"time"
)

// RP2350: Simulated RTC via monotonic clock.
// We anchor the “wall" time (UTC) to a time.Now() (monotonic included).
// Now() = anchorWall + time.Since(anchorMono)
type rtcRP2350 struct {
	mu         sync.Mutex
	anchorMono time.Time // instant monotone au Set()
	anchorWall time.Time // heure "mur" voulue au Set() (UTC)
	inited     bool
}

// newRTC creates and returns a new instance of the RTC implementation for the RP2350 platform.
// It returns an RTC interface that can be used to interact with the real-time clock hardware.
func newRTC() RTC { return &rtcRP2350{} }

// Init1Hz initializes the RTC to generate a 1Hz signal.
// For the simulated version, this function is a no-op.
func (r *rtcRP2350) Init1Hz(_ uint32) {
	// no-op pour la version simulée
}

// Set initializes the rtcRP2350 instance with the provided wall clock time.
// It locks the mutex to ensure thread safety, sets the current monotonic time as an anchor,
// stores the provided time in UTC as the wall clock reference, and marks the RTC as initialized.
func (r *rtcRP2350) Set(t time.Time) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.anchorMono = time.Now()    // porte l'horloge monotone
	r.anchorWall = t.UTC()       // stocke l'heure mur de référence
	r.inited = true
}

// Now returns the current wall-clock time as tracked by the rtcRP2350 instance.
// If the RTC has not been initialized, it returns the Unix epoch (UTC).
// The method calculates the elapsed time since the last anchor point using a monotonic clock,
// and adds this duration to the anchored wall time. If the elapsed time is negative or exceeds
// 24 hours (as a safeguard against anomalies), it is reset to zero.
func (r *rtcRP2350) Now() time.Time {
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.inited {
		return time.Unix(0, 0).UTC()
	}
	// durée écoulée selon l’horloge monotone
	elapsed := time.Since(r.anchorMono)
	// garde-fou simple en cas d’anomalie
	if elapsed < 0 || elapsed > 24*time.Hour {
		elapsed = 0
	}
	return r.anchorWall.Add(elapsed)
}
