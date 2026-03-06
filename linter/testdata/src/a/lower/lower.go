package lower

import (
	"go.uber.org/zap"
	"log/slog"
)

func someFun() {
	// Zap:

	logger := zap.NewExample()
	defer logger.Sync()

	logger.Info("Upper case letter") // want "log messages must begin with a lowercase letter."
	logger.Info("jUST message")
	logger.Info("Однозначно сообщение с БОЛЬШОЙ буквы") // want "log messages must begin with a lowercase letter."
	logger.Info("однозначно сообщение с МАЛЕНЬКОЙ буквы")

	logger.Warn("Upper case letter")                    // want "log messages must begin with a lowercase letter."
	logger.Warn("Однозначно сообщение с БОЛЬШОЙ буквы") // want "log messages must begin with a lowercase letter."

	logger.Debug("Upper case letter")                    // want "log messages must begin with a lowercase letter."
	logger.Debug("Однозначно сообщение с БОЛЬШОЙ буквы") // want "log messages must begin with a lowercase letter."

	logger.Fatal("Upper case letter")                    // want "log messages must begin with a lowercase letter."
	logger.Fatal("Однозначно сообщение с БОЛЬШОЙ буквы") // want "log messages must begin with a lowercase letter."

	logger.Error("Upper case letter")                    // want "log messages must begin with a lowercase letter."
	logger.Error("Однозначно сообщение с БОЛЬШОЙ буквы") // want "log messages must begin with a lowercase letter."

	// slog:

	slog.Info("Upper case letter") // want "log messages must begin with a lowercase letter."
	slog.Info("jUST message")
	slog.Info("Однозначно сообщение с БОЛЬШОЙ буквы") // want "log messages must begin with a lowercase letter."
	slog.Info("однозначно сообщение с МАЛЕНЬКОЙ буквы")

	slog.Warn("Upper case letter") // want "log messages must begin with a lowercase letter."
	slog.Warn("jUST message")
	slog.Warn("Однозначно сообщение с БОЛЬШОЙ буквы") // want "log messages must begin with a lowercase letter."
	slog.Warn("однозначно сообщение с МАЛЕНЬКОЙ буквы")

	slog.Debug("Upper case letter") // want "log messages must begin with a lowercase letter."
	slog.Debug("jUST message")
	slog.Debug("Однозначно сообщение с БОЛЬШОЙ буквы") // want "log messages must begin with a lowercase letter."
	slog.Debug("однозначно сообщение с МАЛЕНЬКОЙ буквы")

	slog.Debug("Upper case letter") // want "log messages must begin with a lowercase letter."
	slog.Debug("jUST message")
	slog.Debug("Однозначно сообщение с БОЛЬШОЙ буквы") // want "log messages must begin with a lowercase letter."
	slog.Debug("однозначно сообщение с МАЛЕНЬКОЙ буквы")

}
