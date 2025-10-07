//go:build tinygo && (pico || pico_w || rp2040 || pico2 || pico2_w || rp2350)

package main

import (
	"fmt"
	"machine"
	"time"

	font "github.com/Nondzu/ssd1306_font"
	ntp "github.com/jgrelet/pico-rtc/ntputil"
	"github.com/jgrelet/pico-rtc/rtcutil" // remplace par ton module si publié
	"github.com/jgrelet/pico-rtc/ssd1306x"
	"github.com/jgrelet/pico-rtc/logger"
)

func main() {

	// Initialisation du port série pour le debug
	time.Sleep(2 * time.Second)
	logger.Logger.Info("RTC unified (RP2040 / RP2350)")

	// --- OLED ---
	disp := ssd1306x.NewI2C(ssd1306x.Config{
		I2C:     *machine.I2C1,
		Address: 0x3C,
		SCL:     machine.I2C1_SCL_PIN, // Pico/Pico2: GP5
		SDA:     machine.I2C1_SDA_PIN, // Pico/Pico2: GP4
		Freq:    400 * machine.KHz,
		Width:   128,
		Height:  64,
	})
	logger.Logger.Info("OLED init ...")

	//font library init
	display := font.NewDisplay(*disp.Device())               //pass by value
	display.Configure(font.Config{FontType: font.FONT_7x10}) //set font here
	display.YPos = 0                                         // set position Y
	display.XPos = 0                                         // set position X

	logger.Logger.Info("NTP dhcp started ...")
	display.PrintText("Dhcp started...") // print text
	disp.ClearBuffer()

	// Initialiser le Wi-Fi et la connexion NTP
	conn, err := ntp.NewNTPConn("Pico2-w", "192.168.1.149", 10, logger.Logger)
	if err != nil {
		fmt.Println("Error connect Wi-Fi :", err)
		display.PrintText(fmt.Sprintf("Error Wi-Fi:", err))
		disp.Display()
		return
	}
	logger.Logger.Info(conn.String())

	now, err := conn.GetNTPTime()
	if err != nil {
		fmt.Println("NTP error:", err)
		display.PrintText(fmt.Sprintf("NTP error:", err))
		disp.Display()
	} else {
		logger.Logger.Info("NTP time :", now.String())
	}
	display.YPos = 12
	display.PrintText("NTP OK")
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
		//dev.ClearDisplay()
		//println(now.Format("15:04:05 02/01/2006"))
		display.YPos = 0
		display.PrintText(now.Format("15:04:05 02/01/06"))
		disp.Display()
		disp.ClearBuffer()
	}

}
