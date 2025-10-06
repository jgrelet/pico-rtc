//go:build debug

package logger


import (
	"log/slog"
	"machine"
    "time"
)

// Logger est visible depuis les autres packages.
var Logger *slog.Logger



func init() {

    machine.Serial.Configure(machine.UARTConfig{BaudRate: 115200})

    time.Sleep(1500 * time.Millisecond)

	// En mode DEBUG, on écrit les logs sur la sortie série du Pico.
	Logger = slog.New(slog.NewTextHandler(machine.Serial, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
    
	Logger.Debug("Logger initialized in DEBUG mode")
}
