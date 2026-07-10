# Files Changed
 internal/app/elapsed.go          | 51 ++++++++++++++++++++++++++++-
 internal/app/elapsed_test.go     | 61 +++++++++++++++++++++++++++++++++-
 internal/telegram/commands.go    | 29 +++++++++++++++++
 internal/telegram/formatter.go   |  1 +
 internal/telegram/parser.go      | 28 +++++++++++++++-
 internal/telegram/parser_test.go | 70 +++++++++++++++++++++++++++++++++++++++-
 internal/telegram/service.go     |  6 ++++
 7 files changed, 242 insertions(+), 4 deletions(-)

# Detailed Changes
diff --git a/internal/app/elapsed.go b/internal/app/elapsed.go
index 9cb8859..25c12b4 100644
--- a/internal/app/elapsed.go
+++ b/internal/app/elapsed.go
@@ -46,4 +46,53 @@ func (m *Manager) SetElapsed(
 	}
 
 	return nil
-}
\ No newline at end of file
+}
+
+func (m *Manager) SetRegen(
+	ctx context.Context,
+	systemID string,
+	minutesLeft int,
+) error {
+	m.mu.Lock()
+
+	system, reminder, ok := m.system(systemID)
+	if !ok {
+		m.mu.Unlock()
+		return fmt.Errorf("unknown point system %q", systemID)
+	}
+
+	regenDuration := time.Duration(system.RegenMinutes) * time.Minute
+	elapsed := regenDuration - time.Duration(minutesLeft)*time.Minute
+
+	if minutesLeft <= 0 {
+		elapsed = regenDuration
+	} else if minutesLeft >= system.RegenMinutes {
+		elapsed = 0
+	}
+
+	system.SetElapsed(elapsed)
+
+	reminder.AlertSent = false
+	reminder.FullSent = false
+
+	systemToSave := system
+	reminderToSave := reminder
+
+	m.mu.Unlock()
+
+	if err := m.store.SavePointSystems(
+		ctx,
+		[]*points.PointSystem{systemToSave},
+	); err != nil {
+		return fmt.Errorf("save point system: %w", err)
+	}
+
+	if err := m.store.SaveReminderState(
+		ctx,
+		reminderToSave,
+	); err != nil {
+		return fmt.Errorf("save reminder state: %w", err)
+	}
+
+	return nil
+}
diff --git a/internal/app/elapsed_test.go b/internal/app/elapsed_test.go
index f731e52..c1a8da1 100644
--- a/internal/app/elapsed_test.go
+++ b/internal/app/elapsed_test.go
@@ -86,4 +86,63 @@ func TestManagerSetElapsed_UnknownSystem(t *testing.T) {
 	); err == nil {
 		t.Fatal("expected error")
 	}
-}
\ No newline at end of file
+}
+
+func TestManagerSetRegen(t *testing.T) {
+	store := &mockStore{}
+
+	manager := &Manager{
+		store: store,
+		pointSystems: map[string]*points.PointSystem{
+			"TP": {
+				Definition: points.Definition{
+					ID:           "TP",
+					Name:         "Training Points",
+					Max:          100,
+					RegenMinutes: 10,
+				},
+				Current: 50,
+				Elapsed: 4 * time.Minute,
+			},
+		},
+		reminders: map[string]*points.ReminderState{
+			"TP": {
+				SystemID:  "TP",
+				AlertSent: true,
+				FullSent:  true,
+			},
+		},
+	}
+
+	if err := manager.SetRegen(
+		context.Background(),
+		"TP",
+		6,
+	); err != nil {
+		t.Fatalf("SetRegen() error = %v", err)
+	}
+
+	system := manager.pointSystems["TP"]
+
+	if system.Elapsed != 4*time.Minute {
+		t.Fatalf("Elapsed = %v, want %v", system.Elapsed, 4*time.Minute)
+	}
+
+	reminder := manager.reminders["TP"]
+
+	if reminder.AlertSent {
+		t.Fatal("AlertSent should be reset")
+	}
+
+	if reminder.FullSent {
+		t.Fatal("FullSent should be reset")
+	}
+
+	if !store.savePointSystemsCalled {
+		t.Fatal("SavePointSystems was not called")
+	}
+
+	if !store.saveReminderCalled {
+		t.Fatal("SaveReminderState was not called")
+	}
+}
diff --git a/internal/telegram/commands.go b/internal/telegram/commands.go
index 019824e..6981dd1 100644
--- a/internal/telegram/commands.go
+++ b/internal/telegram/commands.go
@@ -23,6 +23,9 @@ func (b *Bot) handleCommand(msg *tgbotapi.Message) {
 	case "elapsed":
 		b.handleElapsed(msg)
 
+	case "regen":
+		b.handleRegen(msg)
+
 	default:
 		b.handleUnknownCommand(msg)
 	}
@@ -96,6 +99,32 @@ func (b *Bot) handleElapsed(msg *tgbotapi.Message) {
 	)
 }
 
+func (b *Bot) handleRegen(msg *tgbotapi.Message) {
+	systemID, minutesLeft, err := ParseRegen(msg.CommandArguments())
+	if err != nil {
+		b.SendText(msg.Chat.ID, err.Error())
+		return
+	}
+
+	if err := b.service.SetRegen(
+		context.Background(),
+		systemID,
+		minutesLeft,
+	); err != nil {
+		b.SendText(msg.Chat.ID, err.Error())
+		return
+	}
+
+	b.SendText(
+		msg.Chat.ID,
+		fmt.Sprintf(
+			"Updated %s: %d minute(s) left until the next point.",
+			systemID,
+			minutesLeft,
+		),
+	)
+}
+
 func (b *Bot) handleReply(msg *tgbotapi.Message) {
 	amount, err := strconv.Atoi(strings.TrimSpace(msg.Text))
 	if err != nil {
diff --git a/internal/telegram/formatter.go b/internal/telegram/formatter.go
index 3bf4d65..a1016ec 100644
--- a/internal/telegram/formatter.go
+++ b/internal/telegram/formatter.go
@@ -73,6 +73,7 @@ func FormatHelp() string {
 	b.WriteString("/help - Show this help message.\n")
 	b.WriteString("/use <SYSTEM> <AMOUNT> - Consume or add points.\n")
 	b.WriteString("/elapsed <SYSTEM> <MINUTES> - Set elapsed regeneration time.")
+	b.WriteString("\n/regen <SYSTEM> <MINUTES_LEFT> - Set minutes left until the next point.")
 
 	return b.String()
 }
diff --git a/internal/telegram/parser.go b/internal/telegram/parser.go
index 1019d24..9cbd7d9 100644
--- a/internal/telegram/parser.go
+++ b/internal/telegram/parser.go
@@ -48,4 +48,30 @@ func ParseElapsed(args string) (systemID string, minutes int, err error) {
 	}
 
 	return strings.ToUpper(fields[0]), minutes, nil
-}
\ No newline at end of file
+}
+
+// ParseRegen parses:
+//
+//	/regen TP 15
+//
+// returning:
+//
+//	systemID = "TP"
+//	minutesLeft = 15
+func ParseRegen(args string) (systemID string, minutesLeft int, err error) {
+	fields := strings.Fields(args)
+	if len(fields) != 2 {
+		return "", 0, fmt.Errorf("usage: /regen <SYSTEM> <MINUTES_LEFT>")
+	}
+
+	minutesLeft, err = strconv.Atoi(fields[1])
+	if err != nil {
+		return "", 0, fmt.Errorf("minutes must be an integer")
+	}
+
+	if minutesLeft < 0 {
+		return "", 0, fmt.Errorf("minutes must be non-negative")
+	}
+
+	return strings.ToUpper(fields[0]), minutesLeft, nil
+}
diff --git a/internal/telegram/parser_test.go b/internal/telegram/parser_test.go
index 025553c..c496bc9 100644
--- a/internal/telegram/parser_test.go
+++ b/internal/telegram/parser_test.go
@@ -126,4 +126,72 @@ func TestParseElapsed(t *testing.T) {
 			}
 		})
 	}
-}
\ No newline at end of file
+}
+
+func TestParseRegen(t *testing.T) {
+	tests := []struct {
+		name      string
+		args      string
+		wantID    string
+		wantValue int
+		wantErr   bool
+	}{
+		{
+			name:      "valid",
+			args:      "TP 30",
+			wantID:    "TP",
+			wantValue: 30,
+		},
+		{
+			name:      "lowercase system",
+			args:      "rp 10",
+			wantID:    "RP",
+			wantValue: 10,
+		},
+		{
+			name:    "missing minutes",
+			args:    "TP",
+			wantErr: true,
+		},
+		{
+			name:    "invalid minutes",
+			args:    "TP xyz",
+			wantErr: true,
+		},
+		{
+			name:    "negative minutes",
+			args:    "TP -1",
+			wantErr: true,
+		},
+		{
+			name:    "too many arguments",
+			args:    "TP 20 extra",
+			wantErr: true,
+		},
+	}
+
+	for _, tt := range tests {
+		t.Run(tt.name, func(t *testing.T) {
+			gotID, gotValue, err := ParseRegen(tt.args)
+
+			if tt.wantErr {
+				if err == nil {
+					t.Fatal("expected error")
+				}
+				return
+			}
+
+			if err != nil {
+				t.Fatalf("unexpected error: %v", err)
+			}
+
+			if gotID != tt.wantID {
+				t.Fatalf("id = %q, want %q", gotID, tt.wantID)
+			}
+
+			if gotValue != tt.wantValue {
+				t.Fatalf("value = %d, want %d", gotValue, tt.wantValue)
+			}
+		})
+	}
+}
diff --git a/internal/telegram/service.go b/internal/telegram/service.go
index 30c868a..9702003 100644
--- a/internal/telegram/service.go
+++ b/internal/telegram/service.go
@@ -21,6 +21,12 @@ type Service interface {
 		minutes int,
 	) error
 
+	SetRegen(
+		ctx context.Context,
+		systemID string,
+		minutesLeft int,
+	) error
+
 	ConsumeReply(
 		ctx context.Context,
 		messageID int,
