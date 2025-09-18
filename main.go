package main

import (
	"fmt"
	"time"

	"github.com/jgrelet/pico-rtc/rtcutil"
	ntp "github.com/jgrelet/pico-rtc/ntputil"
)


func main() {

	time.Sleep(2 * time.Second) // attendre que tout soit prêt
	fmt.Println("NTP dhcp started ...")

	// Lire la fréquence effective de clk_rtc
	freq := rtcutil.GetClkRTC()
	fmt.Println("clk_rtc =", freq, "Hz")

	// Forcer 1 Hz (valeur standard Pico/Pico W)
	rtcutil.InitRTC1Hz()

	// 2) Initialiser le Wi-Fi et la connexion NTP
	conn, err := ntp.NewNTPConn("DHCP-pico-w", "192.168.1.150", 10)
	if err != nil {
		fmt.Println("Erreur connexion Wi-Fi :", err)
		return
	}
	println(conn.String())

	t, err := conn.GetNTPTime()
	if err != nil {
		fmt.Println("NTP error:", err)
	} else {
		// format HH:MM:SS MM/DD/YYYY
		fmt.Printf("%02d:%02d:%02d %02d/%02d/%04d\n",
			t.Hour(), t.Minute(), t.Second(),
			t.Day(), t.Month(), t.Year())
	}

	// Régler le RTC
	rtcutil.Set(t)

	fmt.Println("RTC set.")

	// display time from RTC every second
	for {
		now := rtcutil.Now()
		println(now.Format("15:04:05 02/01/2006"))
		time.Sleep(1 * time.Second)
	}
}
