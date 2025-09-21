//go:build tinygo && (pico2 || pico2_w || rp2350) 

package main

import (
	"time"
	"fmt"

	"github.com/jgrelet/pico-rtc/rtcutil"
	ntp "github.com/jgrelet/pico-rtc/ntputil"


)

func main() {
	// machine.Serial.Configure(machine.UARTConfig{BaudRate: 115200})
	time.Sleep(1 * time.Second)
	println("RTC simulé Pico2-W (RP2350)")
		time.Sleep(2 * time.Second) // attendre que tout soit prêt
	fmt.Println("NTP dhcp started ...")

	// Initialiser le Wi-Fi et la connexion NTP
	conn, err := ntp.NewNTPConn("DHCP-pico-w","", 10)
	if err != nil {
		fmt.Println("Erreur connexion Wi-Fi :", err)
		return
	}
	fmt.Println(conn.String())

	now, err := conn.GetNTPTime()
	if err != nil {
		fmt.Println("NTP error:", err)
	} else {
		// format HH:MM:SS MM/DD/YYYY
		// fmt.Printf("NTP time: %02d:%02d:%02d %02d/%02d/%04d\n",
		// 	now.Hour(), now.Minute(), now.Second(),
		// 	now.Day(), now.Month(), now.Year())
		fmt.Println("NTP time :", now.String())
	}

	// Mise à l'heure de référence
	rtcutil.Set(now)

	// Affiche l'heure chaque seconde
	for {
		time.Sleep(1 * time.Second)
		// Lire l'heure "RTC"
		now := rtcutil.Now()
		println(now.Format("15:04:05 02/01/2006"))
		
	}
}
