package handler

import (
	"testing"
)

func TestValidateRatings_Valid(t *testing.T) {
	if err := validateRatings(0, 5.5, 10); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidateRatings_TooLow(t *testing.T) {
	if err := validateRatings(-0.1, 5, 5); err == nil {
		t.Fatal("expected error for rating -0.1")
	}
}

func TestValidateRatings_TooHigh(t *testing.T) {
	if err := validateRatings(5, 10.1, 5); err == nil {
		t.Fatal("expected error for rating 10.1")
	}
}

func TestValidateRatings_AllBoundaries(t *testing.T) {
	tests := []struct {
		name          string
		food, svc, vb float64
		wantErr       bool
	}{
		{"all zero", 0, 0, 0, false},
		{"all max", 10, 10, 10, false},
		{"decimal values", 7.5, 8.3, 9.1, false},
		{"small decimal", 0.1, 0.2, 0.3, false},
		{"just under max", 9.9, 9.9, 9.9, false},
		{"food negative", -0.1, 5, 5, true},
		{"service negative", 5, -1, 5, true},
		{"vibe over ten", 5, 5, 10.1, true},
		{"large negative", -5.0, 5, 5, true},
		{"large over", 5, 5, 15.5, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRatings(tt.food, tt.svc, tt.vb)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateRatings(%v,%v,%v) error=%v, wantErr=%v", tt.food, tt.svc, tt.vb, err, tt.wantErr)
			}
		})
	}
}
