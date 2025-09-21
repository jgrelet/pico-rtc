//go:build tinygo && (pico2 || pico2_w || rp2350)

package rtcutil

import (
	"sync"
	"time"
)

// Simule un RTC : Now() = base + (monotonic_now - monotonic_base).
// Avantages : pas de registres spécifiques, pas de dérive sur reset logiciel,
// et fonctionne tant que l’horloge monotone TinyGo tourne (toujours le cas).

var (
	mu          sync.Mutex
	baseUnix    int64 // secondes Unix au Set()
	baseMonoNS  int64 // nanosecondes monotones au Set()
	inited      bool
)

// Set fixe l'heure "RTC". Utilisez de préférence un temps en UTC.
func Set(t time.Time) {
	mu.Lock()
	defer mu.Unlock()
	nowMono := monotonicNS()
	baseUnix = t.Unix()
	baseMonoNS = nowMono
	inited = true
}

// Now renvoie l'heure "RTC". Si Set n'a jamais été appelé, retourne l'époque Unix 0.
func Now() time.Time {
	mu.Lock()
	defer mu.Unlock()
	if !inited {
		return time.Unix(0, 0).UTC()
	}
	elapsedNS := monotonicNS() - baseMonoNS
	secs := elapsedNS / 1_000_000_000
	nsecRema := int(elapsedNS % 1_000_000_000)
	return time.Unix(baseUnix+secs, int64(nsecRema)).UTC()
}

// monotonicNS retourne un temps monotone en nanosecondes depuis le boot.
// TinyGo implémente time.Now() avec une base monotone sur microcontrôleur.
func monotonicNS() int64 {
	return time.Now().UnixNano()
}
