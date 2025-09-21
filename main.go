//go:build tinygo && (pico || pico_w || rp2040 || pico2 || pico2_w || rp2350)

package main

import (
	"machine"
	"time"
	"fmt"

	//"github.com/jgrelet/pico-rtc/rtcutil"
	"pico-rtc/rtcutil" // remplace par ton module si publié
	// ntp "github.com/jgrelet/pico-rtc/ntputil" 
	ntp "pico-rtc/ntputil" 
)

func main() {
	machine.Serial.Configure(machine.UARTConfig{BaudRate: 115200})
	time.Sleep(2 * time.Second)
	fmt.Println("RTC unifié (RP2040 / RP2350)")
	fmt.Println("NTP dhcp started ...")

	// Initialiser le Wi-Fi et la connexion NTP
	conn, err := ntp.NewNTPConn("DHCP-pico-w","192.168.1.150", 10)
	if err != nil {
		fmt.Println("Erreur connexion Wi-Fi :", err)
		return
	}
	fmt.Println(conn.String())

	now, err := conn.GetNTPTime()
	if err != nil {
		fmt.Println("NTP error:", err)
	} else {

		fmt.Println("NTP time :", now.String())
	}

	rtc := rtcutil.NewRTC()

	// RP2040: calibre 1 Hz (0 => fréquence par défaut 46875 Hz)
	// RP2350: no-op (simulation monotone)
	rtc.Init1Hz(0)

	// Mise à l'heure de référence
	rtc.Set(now)

	// Affiche l'heure chaque seconde
	for {
		time.Sleep(1 * time.Second)
		// Lire l'heure "RTC"
		now := rtc.Now()
		println(now.Format("15:04:05 02/01/2006"))	
	}

	
	

	
}
