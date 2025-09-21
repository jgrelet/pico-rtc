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

func newRTC() RTC { return &rtcRP2350{freqHz: 1_000_000_000} }

func (r *rtcRP2350) Init1Hz(_ uint32) {
	// No-op en mode simulation (pas de diviseur matériel à régler).
}

func (r *rtcRP2350) Set(t time.Time) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.setBase = t.UTC()
	r.ticksBase = r.nowTicks()
	r.inited = true
}

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
func (r *rtcRP2350) nowTicks() uint64 {
	return uint64(time.Now().UnixNano())
}

