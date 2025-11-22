package lib

import (
	"fmt"
	"log/slog"
	"os"
	"time"
)

func NewLogger(name string, model *slog.Logger) {
	// Create logs directory
	if err := os.MkdirAll("logs/", 0755); err != nil {
		fmt.Println("[-] logs.mkdir:", err)
	}

	// Create filename and open file
	filename := fmt.Sprintf("logs/%s-%s.log", name, time.Now().Format(time.DateOnly))
	open, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("[-] logs.openfile:", err)
	}

	// New logger
	*model = *slog.New(slog.NewTextHandler(
		open, &slog.HandlerOptions{Level: slog.LevelInfo},
	))
}
