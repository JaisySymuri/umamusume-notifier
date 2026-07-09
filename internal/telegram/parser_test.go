package telegram

import "testing"

func TestParseUse(t *testing.T) {
	tests := []struct {
		name      string
		args      string
		wantID    string
		wantValue int
		wantErr   bool
	}{
		{
			name:      "valid",
			args:      "TP 20",
			wantID:    "TP",
			wantValue: 20,
		},
		{
			name:      "lowercase system",
			args:      "tp 15",
			wantID:    "TP",
			wantValue: 15,
		},
		{
			name:    "missing amount",
			args:    "TP",
			wantErr: true,
		},
		{
			name:    "invalid amount",
			args:    "TP abc",
			wantErr: true,
		},
		{
			name:    "too many arguments",
			args:    "TP 20 extra",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotID, gotValue, err := ParseUse(tt.args)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if gotID != tt.wantID {
				t.Fatalf("id = %q, want %q", gotID, tt.wantID)
			}

			if gotValue != tt.wantValue {
				t.Fatalf("value = %d, want %d", gotValue, tt.wantValue)
			}
		})
	}
}

func TestParseElapsed(t *testing.T) {
	tests := []struct {
		name      string
		args      string
		wantID    string
		wantValue int
		wantErr   bool
	}{
		{
			name:      "valid",
			args:      "TP 30",
			wantID:    "TP",
			wantValue: 30,
		},
		{
			name:      "lowercase system",
			args:      "rp 10",
			wantID:    "RP",
			wantValue: 10,
		},
		{
			name:    "missing minutes",
			args:    "TP",
			wantErr: true,
		},
		{
			name:    "invalid minutes",
			args:    "TP xyz",
			wantErr: true,
		},
		{
			name:    "too many arguments",
			args:    "TP 20 extra",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotID, gotValue, err := ParseElapsed(tt.args)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if gotID != tt.wantID {
				t.Fatalf("id = %q, want %q", gotID, tt.wantID)
			}

			if gotValue != tt.wantValue {
				t.Fatalf("value = %d, want %d", gotValue, tt.wantValue)
			}
		})
	}
}

func TestParseRegen(t *testing.T) {
	tests := []struct {
		name      string
		args      string
		wantID    string
		wantValue int
		wantErr   bool
	}{
		{
			name:      "valid",
			args:      "TP 30",
			wantID:    "TP",
			wantValue: 30,
		},
		{
			name:      "lowercase system",
			args:      "rp 10",
			wantID:    "RP",
			wantValue: 10,
		},
		{
			name:    "missing minutes",
			args:    "TP",
			wantErr: true,
		},
		{
			name:    "invalid minutes",
			args:    "TP xyz",
			wantErr: true,
		},
		{
			name:    "negative minutes",
			args:    "TP -1",
			wantErr: true,
		},
		{
			name:    "too many arguments",
			args:    "TP 20 extra",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotID, gotValue, err := ParseRegen(tt.args)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if gotID != tt.wantID {
				t.Fatalf("id = %q, want %q", gotID, tt.wantID)
			}

			if gotValue != tt.wantValue {
				t.Fatalf("value = %d, want %d", gotValue, tt.wantValue)
			}
		})
	}
}
