package ssd1306x

import (
	"image/color"
	"machine"

	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/tinyfont"
)

type Config struct {
	I2C     machine.I2C // ex: machine.I2C0
	Address uint16      // 0x3C (courant) ou 0x3D
	SCL     machine.Pin // ex: machine.I2C0_SCL_PIN (GP5)
	SDA     machine.Pin // ex: machine.I2C0_SDA_PIN (GP4)
	Freq    uint32      // ex: 400_000
	Width   int16       // ex: 128
	Height  int16       // ex: 64 ou 32
}

type Display struct {
	dev    ssd1306.Device
	width  int16
	height int16
}

// NewI2C initializes and configures an SSD1306 display over I2C using the provided Config.
// It sets up the I2C interface with the specified frequency, SCL, and SDA pins,
// then creates and configures the SSD1306 display with the given address, width, and height.
// Returns a pointer to a Display instance representing the configured display.
func NewI2C(cfg Config) *Display {
	cfg.I2C.Configure(machine.I2CConfig{
		Frequency: cfg.Freq,
		SCL:       cfg.SCL,
		SDA:       cfg.SDA,
	})
	d := ssd1306.NewI2C(&cfg.I2C)
	d.Configure(ssd1306.Config{
		Address: cfg.Address,
		Width:   cfg.Width,
		Height:  cfg.Height,
	})
	return &Display{dev: *d, width: cfg.Width, height: cfg.Height}
}

// Begin initializes the display with the specified contrast level.
// It sets the display contrast, clears the display buffer, and updates the display.
// contrast: The contrast value to set for the display.
func (d *Display) Begin(contrast byte) {
	d.SetContrast(contrast)
	d.ClearDisplay()
	d.Display()
}

// --- Buffer / rendu ---

func (d *Display) ClearDisplay()      { d.dev.ClearDisplay() }
func (d *Display) ClearBuffer()       { d.dev.ClearBuffer() }
func (d *Display) Display() error     { return d.dev.Display() }
func (d *Display) SetContrast(v byte) { d.SetContrast(v) }

// Inversion via commandes INVERTDISPLAY / NORMALDISPLAY
func (d *Display) Invert(on bool) {
	if on {
		d.dev.Command(ssd1306.INVERTDISPLAY)
	} else {
		d.dev.Command(ssd1306.NORMALDISPLAY)
	}
}

// Alim logique : Sleep(false)=ON, Sleep(true)=OFF
func (d *Display) PowerOn()  { _ = d.dev.Sleep(false) }
func (d *Display) PowerOff() { _ = d.dev.Sleep(true) }

// --- Texte ---

// WriteLine écrit du texte à x,y (y = baseline de la police)
func (d *Display) WriteLine(font *tinyfont.Font, text string, x, y int16, col color.RGBA) {
	tinyfont.WriteLine(&d.dev, font, x, y, text, col)
}

// WriteCentered centre horizontalement autour de cx (baseline à y)
func (d *Display) WriteCentered(font *tinyfont.Font, text string, cx, y int16, col color.RGBA) {
	_, outW := tinyfont.LineWidth(font, text) // largeur bbox en pixels
	x := cx - int16(outW/2)
	tinyfont.WriteLine(&d.dev, font, x, y, text, col)
}

// --- Infos / accès bas-niveau ---

func (d *Display) Width() int16              { return d.width }
func (d *Display) Height() int16             { return d.height }
func (d *Display) Device() *ssd1306.Device   { return &d.dev }
