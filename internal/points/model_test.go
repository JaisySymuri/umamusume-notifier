package points

import (
	"errors"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name         string
		id           string
		systemName   string
		max          int
		regenMinutes int
		wantErr      error
	}{
		{
			name:         "valid",
			id:           "TP",
			systemName:   "Training Points",
			max:          100,
			regenMinutes: 10,
			wantErr:      nil,
		},
		{
			name:         "empty id",
			id:           "",
			systemName:   "Training Points",
			max:          100,
			regenMinutes: 10,
			wantErr:      ErrEmptyID,
		},
		{
			name:         "empty name",
			id:           "TP",
			systemName:   "",
			max:          100,
			regenMinutes: 10,
			wantErr:      ErrEmptyName,
		},
		{
			name:         "invalid max",
			id:           "TP",
			systemName:   "Training Points",
			max:          0,
			regenMinutes: 10,
			wantErr:      ErrInvalidMax,
		},
		{
			name:         "invalid regen",
			id:           "TP",
			systemName:   "Training Points",
			max:          100,
			regenMinutes: 0,
			wantErr:      ErrInvalidRegenTime,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := New(Definition{
				ID:           tt.id,
				Name:         tt.systemName,
				Max:          tt.max,
				RegenMinutes: tt.regenMinutes,
			})

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error = %v, want %v", err, tt.wantErr)
			}

			if tt.wantErr != nil {
				return
			}

			if p == nil {
				t.Fatal("New() returned nil")
			}

			if p.ID != tt.id {
				t.Errorf("ID = %q, want %q", p.ID, tt.id)
			}

			if p.Name != tt.systemName {
				t.Errorf("Name = %q, want %q", p.Name, tt.systemName)
			}

			if p.Max != tt.max {
				t.Errorf("Max = %d, want %d", p.Max, tt.max)
			}

			if p.Current != 0 {
				t.Errorf("Current = %d, want 0", p.Current)
			}

			if p.RegenMinutes != tt.regenMinutes {
				t.Errorf("RegenMinutes = %d, want %d", p.RegenMinutes, tt.regenMinutes)
			}

			if p.Elapsed != 0 {
				t.Errorf("Elapsed = %v, want 0", p.Elapsed)
			}

			if p.LastTick.IsZero() {
				t.Error("LastTick was not initialized")
			}
		})
	}
}
