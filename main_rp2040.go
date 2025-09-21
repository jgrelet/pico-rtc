//go:build tinygo && (pico || pico_w || rp2040)

package main

import (
	"fmt"
	"time"
	"runtime"

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

	now, err := conn.GetNTPTime()
	if err != nil {
		fmt.Println("NTP error:", err)
	} else {
		// format HH:MM:SS MM/DD/YYYY
		fmt.Printf("%02d:%02d:%02d %02d/%02d/%04d\n",
			now.Hour(), now.Minute(), now.Second(),
			now.Day(), now.Month(), now.Year())
			fmt.Println("NTP time =", now.String())
	}

	// Régler le RTC
	//rtcutil.Set(t)

	fmt.Println("RTC set.")

	/* actualTime, _ := time.Parse(time.RFC3339, fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%04dZ",
			t.Year(), t.Month(), t.Day(),
			t.Hour(), t.Minute(), t.Second())) */
			t := time.Now()
			fmt.Println("System time =", t.String())
	offset := now.Sub(t)

	// adjust internal clock by adding the offset to the internal clock
	runtime.AdjustTimeOffset(int64(offset))

	// display time from RTC every second
	for {
		/* now := rtcutil.Now()
		println(now.Format("15:04:05 02/01/2006")) */
		println(time.Now().Format(time.RFC3339))
		time.Sleep(1 * time.Second)
	}
}
