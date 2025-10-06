//go:build tinygo && ds3231

package rtcutil

import (
	"machine"
	"sync"
	"time"
)

//
type ds3231 struct {
	mu        sync.Mutex
	setBase   time.Time
	ticksBase uint64
	inited    bool
	freqHz    uint64 // ticks/s pour nowTicks() ; ici 1 tick = 1 ns (UnixNano)
}

// mustRetry attempts to execute the provided function f up to n times with a delay between each attempt.
// If f returns nil, mustRetry returns immediately. If all attempts fail, it prints the error and panics.
// Parameters:
//   n     - number of retry attempts
//   delay - duration to wait between attempts
//   f     - function to execute, which returns an error
func mustRetry(n int, delay time.Duration, f func() error) {
	for i := 0; i < n; i++ {
		if err := f(); err == nil {
			return
		} else if i == n-1 {
			println("ERR:", err.Error())
			panic(err)
		}
		time.Sleep(delay)
	}
}

func newRTC() RTC {

	// RTC DS3231
	// I2C0 sur GPIO4 (SDA) / GPIO5 (SCL) en 400kHz
	machine.I2C1.Configure(machine.I2CConfig{
		SCL:       machine.I2C1_SCL_PIN,
		SDA:       machine.I2C1_SDA_PIN,
		Frequency: 400 * machine.KHz,
	})

	rtcds := ds3231.New(machine.I2C1)
	ok := rtcds.Configure()
	if !ok {
		println("DS3231 not detected (addr 0x68) ?")
		for {
			time.Sleep(2 * time.Second)
		}
	}

	// Si l'oscillateur n'est pas en marche, on le démarre.
	if !rtcds.IsRunning() {
		mustRetry(5, 200*time.Millisecond, func() error { return rtc.SetRunning(true) })
	}
	return &ds3231{}
}
