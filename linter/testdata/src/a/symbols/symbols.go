package symbols

import (
	"log/slog"
)

func someFunc() {
	slog.Info("server started!��")                // want "log messages must not contain special symbols."
	slog.Error("connection failed!!!")            // want "log messages must not contain special symbols."
	slog.Warn("warning: something went wrong...") // want "log messages must not contain special symbols."
}
