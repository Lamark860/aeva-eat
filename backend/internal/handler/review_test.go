package handler

import (
	"testing"
)

func TestValidateRatings_Valid(t *testing.T) {
	if err := validateRatings(1, 5, 10); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidateRatings_TooLow(t *testing.T) {
	if err := validateRatings(0, 5, 5); err == nil {
		t.Fatal("expected error for rating 0")
	}
}

func TestValidateRatings_TooHigh(t *testing.T) {
	if err := validateRatings(5, 11, 5); err == nil {
		t.Fatal("expected error for rating 11")
	}
}

func TestValidateRatings_AllBoundaries(t *testing.T) {
	tests := []struct {
		name          string
		food, svc, vb int
		wantErr       bool
	}{
		{"all min", 1, 1, 1, false},
		{"all max", 10, 10, 10, false},
		{"food zero", 0, 5, 5, true},
		{"service negative", 5, -1, 5, true},
		{"vibe eleven", 5, 5, 11, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRatings(tt.food, tt.svc, tt.vb)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateRatings(%d,%d,%d) error=%v, wantErr=%v", tt.food, tt.svc, tt.vb, err, tt.wantErr)
			}
		})
	}
}
