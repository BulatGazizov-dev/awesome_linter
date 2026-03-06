package sensitive

import (
	"go.uber.org/zap"
	"log/slog"
)

func someFun() {
	// Zap:

	logger := zap.NewExample()
	defer logger.Sync()

	var (
		token  string
		apiKey string
		secret string
	)

	logger.Info(secret)              // want "log messages should not contain potentially sensitive data."
	logger.Info("ApiKey" + apiKey)   // want "log messages should not contain potentially sensitive data."
	logger.Info("ApiKey" + (apiKey)) // want "log messages should not contain potentially sensitive data."

	slog.Info(secret)                       // want "log messages should not contain potentially sensitive data."
	slog.Info("ApiKey" + apiKey)            // want "log messages should not contain potentially sensitive data."
	slog.Info("ApiKey" + (apiKey))          // want "log messages should not contain potentially sensitive data."
	slog.Info("some sensitive data", token) // want "log messages should not contain potentially sensitive data."

}
