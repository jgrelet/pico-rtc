//go:build !debug && !info

package logger


import (
	"io"
	"log/slog"
)

var Logger *slog.Logger


func init() {
	// En mode release, on envoie tout dans /dev/null (rien affiché)
	Logger = slog.New(slog.NewTextHandler(io.Discard, nil))
}
