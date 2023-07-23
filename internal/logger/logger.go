package logger

import (
	"go.uber.org/zap"
)

// Log - logger object.
var Log = zap.NewNop()

// NewLogger the constructor creates a global variable Log.
func NewLogger(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}

	// используется для ведения журнала разработки.
	cfg := zap.NewProductionConfig()

	cfg.Level = lvl

	Log, err = cfg.Build()
	if err != nil {
		return err
	}

	return nil
}
