package models

import (
	"testing"
	"time"
)

func TestCalculateAge(t *testing.T) {
	now := time.Now().UTC()

	tests := []struct {
		name string
		dob  time.Time
		want int
	}{
		{
			name: "birthday already passed this year",
			dob:  time.Date(now.Year()-30, 1, 1, 0, 0, 0, 0, time.UTC),
			want: 30,
		},
		{
			name: "exactly today (birthday today)",
			dob:  time.Date(now.Year()-25, now.Month(), now.Day(), 0, 0, 0, 0, time.UTC),
			want: 25,
		},
		{
			name: "born today (age 0)",
			dob:  time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC),
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateAge(tt.dob)
			if got != tt.want {
				t.Errorf("CalculateAge(%v) = %d; want %d", tt.dob, got, tt.want)
			}
		})
	}

	// Birthday not yet occurred this year -> should be one less than naive diff.
	future := now.AddDate(0, 1, 0) // a month ahead of "now" in month/day terms
	dobFuture := time.Date(now.Year()-20, future.Month(), future.Day(), 0, 0, 0, 0, time.UTC)
	got := CalculateAge(dobFuture)
	if got != 19 {
		t.Errorf("CalculateAge(birthday not yet occurred) = %d; want 19", got)
	}
}
