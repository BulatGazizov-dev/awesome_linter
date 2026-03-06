package english

import (
	"go.uber.org/zap"
	"log/slog"
)

func someFun() {
	logger := zap.NewExample()
	defer logger.Sync()

	logger.Info("usual english message")

	// Zap
	logger.Info("сообщение на русском")      // want "log messages must be in English only."
	logger.Info("log messagе with surprise") // want "log messages must be in English only."

	logger.Warn("сообщение на русском")      // want "log messages must be in English only."
	logger.Warn("log messagе with surprise") // want "log messages must be in English only."

	logger.Debug("сообщение на русском")      // want "log messages must be in English only."
	logger.Debug("log messagе with surprise") // want "log messages must be in English only."

	logger.Fatal("сообщение на русском")      // want "log messages must be in English only."
	logger.Fatal("log messagе with surprise") // want "log messages must be in English only."

	logger.Panic("сообщение на русском")      // want "log messages must be in English only."
	logger.Panic("log messagе with surprise") // want "log messages must be in English only."

	// slog
	slog.Info("сообщение на русском")      // want "log messages must be in English only."
	slog.Info("log messagе with surprise") // want "log messages must be in English only."

	slog.Error("сообщение на русском") // want "log messages must be in English only."
	slog.Debug("сообщение на русском") // want "log messages must be in English only."
	slog.Warn("сообщение на русском")  // want "log messages must be in English only."
}
