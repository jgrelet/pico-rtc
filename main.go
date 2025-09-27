//go:build tinygo && (pico || pico_w || rp2040 || pico2 || pico2_w || rp2350)

package main

import (
	"machine"
	"time"
	"fmt"
	//"image/color"

	font "github.com/Nondzu/ssd1306_font"
	// "tinygo.org/x/drivers/ssd1306"

	//"tinygo.org/x/tinyfont/freemono"

	//"github.com/jgrelet/pico-rtc/rtcutil"
	"pico-rtc/rtcutil" // remplace par ton module si publié
	// ntp "github.com/jgrelet/pico-rtc/ntputil" 
	ntp "pico-rtc/ntputil" 
	"pico-rtc/ssd1306x"  
)

func main() {

	/* var (
		with int16 = 128
		height int16 = 32
		//height int16 = 64
	) */
	machine.Serial.Configure(machine.UARTConfig{BaudRate: 115200})
	time.Sleep(2 * time.Second)
	println("RTC unified (RP2040 / RP2350)")

	// The default I2C1 pins are GP3 and GP4, so we use those here.
	/* machine.I2C1.Configure(machine.I2CConfig{
		//Frequency: 400000,
		Frequency: 400 * machine.KHz,
		// SCL: machine.I2C1_SCL_PIN,
		// SDA: machine.I2C1_SDA_PIN,
	})
	// Initialiser l'écran OLED SSD1306
	dev := ssd1306x.NewI2C(machine.I2C1)
	dev.Configure(ssd1306.Config{Width: with, Height: height, Address: 0x3C, VccState: ssd1306.SWITCHCAPVCC})
	dev.ClearBuffer()
	dev.ClearDisplay() */

	// --- OLED ---
	disp := ssd1306x.NewI2C(ssd1306x.Config{
		I2C:     *machine.I2C1,
		Address: 0x3C,
		SCL:     machine.I2C1_SCL_PIN, // Pico/Pico2: GP5
		SDA:     machine.I2C1_SDA_PIN, // Pico/Pico2: GP4
		Freq:    400 * machine.KHz,
		Width:   128,
		Height:  32,
	})
	println("Je passe 1")

	//disp.Begin(0x7F)
	//white := color.RGBA{255, 255, 255, 255}
	println("Je passe 2")
	
	//font library init
	display := font.NewDisplay(*disp.Device()) //pass by value
	display.Configure(font.Config{FontType: font.FONT_7x10}) //set font here
	display.YPos = 0                 // set position Y
	display.XPos = 0                 // set position X
	
	println("NTP dhcp started ...")
	display.PrintText("Dhcp started...") // print text
	disp.ClearBuffer()
	//disp.WriteLine(&freemono.Regular9pt7b,"Dhcp started...", 0, 0, white) // print text

	// Initialiser le Wi-Fi et la connexion NTP
	conn, err := ntp.NewNTPConn("DHCP-pico-w","192.168.1.150", 10)
	if err != nil {
		fmt.Println("Error connect Wi-Fi :", err)
		display.PrintText(fmt.Sprintf("Error Wi-Fi:", err))
		//disp.WriteLine(&freemono.Regular9pt7b,fmt.Sprintf("Error Wi-Fi:", err), 0, 0, white) // print text
		return
	}
	println(conn.String())

	now, err := conn.GetNTPTime()
	if err != nil {
		fmt.Println("NTP error:", err)
		display.PrintText(fmt.Sprintf("NTP error:", err))
		//disp.WriteLine(&freemono.Regular9pt7b,fmt.Sprintf("NTP error:", err), 0, 0, white) // print text
	} else {

		println("NTP time :", now.String())
	}
	display.YPos = 12
	display.PrintText("NTP OK")
	//disp.WriteLine(&freemono.Regular9pt7b,"NTP OK", 0, 12, white) // print text
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
		println(now.Format("15:04:05 02/01/2006"))		
		display.YPos = 0
		display.PrintText(now.Format("15:04:05 02/01/06"))
		//disp.WriteLine(&freemono.Regular9pt7b, now.Format("15:04:05 02/01/06"), 0, 24, white) // print text
		disp.Display()
		disp.ClearBuffer()	
	}

	
	

	
}
