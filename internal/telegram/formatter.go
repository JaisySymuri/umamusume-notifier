package telegram

import (
	"fmt"
	"strings"
	"time"

	"umamusume-notifier/internal/app"
)

// FormatStatus formats all point systems for display.
func FormatStatus(statuses []app.Status) string {
	if len(statuses) == 0 {
		return "No point systems configured."
	}

	var b strings.Builder

	b.WriteString("📊 Point Status\n\n")

	for i, status := range statuses {
		fmt.Fprintf(
			&b,
			"%s (%s)\n",
			status.Name,
			status.ID,
		)

		fmt.Fprintf(
			&b,
			"  %d/%d\n",
			status.Current,
			status.Max,
		)

		if status.Full {
			b.WriteString("  FULL\n")
		} else {
			fmt.Fprintf(
				&b,
				"  Full in: %s\n",
				formatDuration(status.TimeUntilFull),
			)
		}

		if i != len(statuses)-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

func formatDuration(d time.Duration) string {
	totalMinutes := int(d.Minutes())

	hours := totalMinutes / 60
	minutes := totalMinutes % 60

	if hours == 0 {
		return fmt.Sprintf("%dm", minutes)
	}

	return fmt.Sprintf("%dh %dm", hours, minutes)
}

// FormatHelp returns the help text shown by the /help command.
func FormatHelp() string {
	var b strings.Builder

	b.WriteString("Available commands:\n\n")
	b.WriteString("/status - Show all point systems.\n")
	b.WriteString("/help - Show this help message.\n")
	b.WriteString("/use <SYSTEM> <AMOUNT> - Consume or add points.\n")
	b.WriteString("/elapsed <SYSTEM> <MINUTES> - Set elapsed regeneration time.")

	return b.String()
}
