package telegram

import (
	"fmt"
	"strings"
	"time"

	"umamusume-notifier/internal/app"
)

type statusGroup struct {
	Header string
	IDs    []string
}

var statusGroups = []statusGroup{
	{
		Header: "🦄 Umamusume",
		IDs:    []string{"CP", "RP", "TP"},
	},
	{
		Header: "🎤 Holodori",
		IDs:    []string{"1P", "2P", "LP", "MP"},
	},
}

// FormatStatus formats all point systems for display.
func FormatStatus(statuses []app.Status) string {
	if len(statuses) == 0 {
		return "No point systems configured."
	}

	statusByID := make(map[string]app.Status, len(statuses))
	for _, status := range statuses {
		statusByID[status.ID] = status
	}

	used := make(map[string]struct{}, len(statuses))
	sections := make([]string, 0, len(statusGroups)+1)
	now := time.Now()

	for _, group := range statusGroups {
		lines := make([]string, 0, len(group.IDs)+1)
		for _, id := range group.IDs {
			status, ok := statusByID[id]
			if !ok {
				continue
			}

			lines = append(lines, renderStatus(status, now))
			used[id] = struct{}{}
		}

		if len(lines) == 0 {
			continue
		}

		sections = append(sections, group.Header+"\n\n"+strings.Join(lines, "\n\n"))
	}

	var other []string
	for _, status := range statuses {
		if _, ok := used[status.ID]; ok {
			continue
		}
		other = append(other, renderStatus(status, now))
	}

	if len(other) > 0 {
		sections = append(sections, strings.Join(other, "\n\n"))
	}

	return strings.Join(sections, "\n\n")
}

func renderStatus(status app.Status, now time.Time) string {
	var b strings.Builder

	fmt.Fprintf(&b, "%s (%s)\n", status.Name, status.ID)
	fmt.Fprintf(&b, "  %d/%d\n", status.Current, status.Max)

	if status.Full {
		b.WriteString("  FULL")
		return b.String()
	}

	fmt.Fprintf(
		&b,
		"  Full in: %s (%s)",
		formatDuration(status.TimeUntilFull),
		formatFullTime(now, status.TimeUntilFull),
	)

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
	b.WriteString("/status, /s - Show all point systems.\n")
	b.WriteString("/help - Show this help message.\n")
	b.WriteString("/use <SYSTEM> <AMOUNT> - Consume or add points.\n")
	b.WriteString("/set <SYSTEM> <AMOUNT> - Set current points directly.\n")
	b.WriteString("/elapsed <SYSTEM> <MINUTES> - Set elapsed regeneration time.")
	b.WriteString("\n/regen <SYSTEM> <MINUTES_LEFT> - Set minutes left until the next point.")

	return b.String()
}

func FormatServiceOnline() string {
	return "✅ Service is online"
}
