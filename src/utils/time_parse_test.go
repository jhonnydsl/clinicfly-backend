package utils

import "testing"

func TestParseDate(t *testing.T) {
	tests := []struct {
		name      string
		date      string
		wantError bool
	}{
		{
			name:      "valid date",
			date:      "2003-02-03",
			wantError: false,
		},
		{
			name:      "invalid format",
			date:      "03-02-2003",
			wantError: true,
		},
		{
			name:      "invalid date",
			date:      "2003-99-99",
			wantError: true,
		},
		{
			name:      "empty date",
			date:      "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseDate(tt.date)

			hasError := err != nil

			if hasError != tt.wantError {
				t.Errorf("case test %q failed", tt.name)
			}
		})
	}
}

func TestParseTime(t *testing.T) {
	tests := []struct {
		name      string
		time      string
		wantError bool
	}{
		{
			name:      "valid time",
			time:      "14:30",
			wantError: false,
		},
		{
			name:      "midnight",
			time:      "00:00",
			wantError: false,
		},
		{
			name:      "last valid minute",
			time:      "23:59",
			wantError: false,
		},
		{
			name:      "invalid format",
			time:      "14-30",
			wantError: true,
		},
		{
			name:      "invalid hour",
			time:      "25:00",
			wantError: true,
		},
		{
			name:      "invalid minute",
			time:      "14:60",
			wantError: true,
		},
		{
			name:      "empty time",
			time:      "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseTime(tt.time)

			hasError := err != nil

			if hasError != tt.wantError {
				t.Errorf("case test %q failed", tt.name)
			}
		})
	}
}

func TestParseDateTime(t *testing.T) {
	tests := []struct {
		name      string
		date      string
		time      string
		wantError bool
	}{
		{
			name:      "valid date and time",
			date:      "2003-02-03",
			time:      "14:30",
			wantError: false,
		},
		{
			name:      "midnight",
			date:      "2003-02-03",
			time:      "00:00",
			wantError: false,
		},
		{
			name:      "last valid minute",
			date:      "2003-02-03",
			time:      "23:59",
			wantError: false,
		},
		{
			name:      "invalid date",
			date:      "2003-99-99",
			time:      "14:30",
			wantError: true,
		},
		{
			name:      "invalid time",
			date:      "2003-02-03",
			time:      "25:00",
			wantError: true,
		},
		{
			name:      "invalid date format",
			date:      "03-02-2003",
			time:      "14:30",
			wantError: true,
		},
		{
			name:      "invalid time format",
			date:      "2003-02-03",
			time:      "14-30",
			wantError: true,
		},
		{
			name:      "empty date",
			date:      "",
			time:      "14:30",
			wantError: true,
		},
		{
			name:      "empty time",
			date:      "2003-02-03",
			time:      "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseDateTime(tt.date, tt.time)

			hasError := err != nil

			if hasError != tt.wantError {
				t.Errorf("case test %q failed", tt.name)
			}
		})
	}
}