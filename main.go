package main

import (
	"fmt"
	"time"

	"github.com/jgrelet/pico-rtc/rtcutil"
	"github.com/jgrelet/pico-rtc/ntputil"
)


func main() {

	time.Sleep(2 * time.Second) // attendre que tout soit prêt
	fmt.Println("NTP dhcp started ...")

	// 1) Configurer le tick 1 Hz (à faire une fois au boot)
	//initRTC1Hz()
	// 1) Config 1 Hz (valeur standard Pico/Pico W)
	//rtcutil.Init1Hz(46874)

	// Lire la fréquence effective de clk_rtc
	freq := rtcutil.GetClkRTC()
	fmt.Println("clk_rtc =", freq, "Hz")

	/* // Supposons clk_rtc = 46875 Hz (cas standard Pico/Pico W : XOSC/12MHz ÷ 256)
	div := rtcutil.CalcDiv(freq) */
	rtcutil.InitRTC1Hz()

	// 2) Initialiser le Wi-Fi et la connexion NTP
	conn, err := newNTPConn("DHCP-pico-w", "192.168.1.150", 10)
	if err != nil {
		fmt.Println("Erreur connexion Wi-Fi :", err)
		return
	}
	println(conn.String())

	t, err := conn.getNTPTime()
	if err != nil {
		fmt.Println("NTP error:", err)
	} else {
		// format HH:MM:SS MM/DD/YYYY
		fmt.Printf("%02d:%02d:%02d %02d/%02d/%04d\n",
			t.Hour(), t.Minute(), t.Second(),
			t.Day(), t.Month(), t.Year())
	}

	// Convertir en Europe/Paris (CET/CEST)
	//tLocal := toEuropeParis(tUTC)

	// Régler le RTC
	//setRTC(t)
	rtcutil.Set(t)

	fmt.Println("RTC set.")

	for {
		now := rtcutil.Now()
		println(now.Format("15:04:05 02/01/2006"))
		time.Sleep(1 * time.Second)
	}
}
