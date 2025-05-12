package timeutils

import (
	"strings"
	"testing"
	"time"
)

func TestParseAny(t *testing.T) {
	tests := []struct {
		name      string
		timeStr   string
		expectErr bool
		validate  func(t *testing.T, tm time.Time)
	}{
		{
			name:      "RFC3339",
			timeStr:   "2023-07-15T14:30:20Z",
			expectErr: false,
			validate: func(t *testing.T, tm time.Time) {
				if tm.Year() != 2023 || tm.Month() != 7 || tm.Day() != 15 || tm.Hour() != 14 || tm.Minute() != 30 || tm.Second() != 20 {
					t.Errorf("Parsed time is incorrect: %v", tm)
				}
			},
		},
		{
			name:      "ISO8601UTC",
			timeStr:   "2023-07-15T14:30:05Z",
			expectErr: false,
			validate: func(t *testing.T, tm time.Time) {
				if tm.Year() != 2023 || tm.Month() != 7 || tm.Day() != 15 {
					t.Errorf("Parsed time is incorrect: %v", tm)
				}
			},
		},
		{
			name:      "ISO8601",
			timeStr:   "2023-07-15T14:30:05-07:00",
			expectErr: false,
			validate: func(t *testing.T, tm time.Time) {
				if tm.Year() != 2023 || tm.Month() != 7 || tm.Day() != 15 {
					t.Errorf("Parsed time is incorrect: %v", tm)
				}
			},
		},
		{
			name:      "DateOnly",
			timeStr:   "2023-07-15",
			expectErr: false,
			validate: func(t *testing.T, tm time.Time) {
				if tm.Year() != 2023 || tm.Month() != 7 || tm.Day() != 15 {
					t.Errorf("Parsed time is incorrect: %v", tm)
				}
			},
		},
		{
			name:      "DateTime",
			timeStr:   "2023-07-15 14:30:05",
			expectErr: false,
			validate: func(t *testing.T, tm time.Time) {
				if tm.Year() != 2023 || tm.Month() != 7 || tm.Day() != 15 || tm.Hour() != 14 || tm.Minute() != 30 || tm.Second() != 5 {
					t.Errorf("Parsed time is incorrect: %v", tm)
				}
			},
		},
		{
			name:      "DD/MM/YYYY",
			timeStr:   "15/07/2023",
			expectErr: false,
			validate: func(t *testing.T, tm time.Time) {
				if tm.Year() != 2023 || tm.Month() != 7 || tm.Day() != 15 {
					t.Errorf("Parsed time is incorrect: %v", tm)
				}
			},
		},
		{
			name:      "Human readable date",
			timeStr:   "July 15, 2023",
			expectErr: false,
			validate: func(t *testing.T, tm time.Time) {
				if tm.Year() != 2023 || tm.Month() != 7 || tm.Day() != 15 {
					t.Errorf("Parsed time is incorrect: %v", tm)
				}
			},
		},
		{
			name:      "Invalid format",
			timeStr:   "not a time",
			expectErr: true,
			validate:  func(t *testing.T, tm time.Time) {},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tm, err := ParseAny(tc.timeStr)
			
			if tc.expectErr && err == nil {
				t.Error("Expected error but got nil")
			}
			
			if !tc.expectErr && err != nil {
				t.Errorf("Did not expect error but got: %v", err)
			}
			
			if !tc.expectErr {
				tc.validate(t, tm)
			}
		})
	}
}

func TestFormatISO8601(t *testing.T) {
	// Create a time in a specific timezone
	loc, _ := time.LoadLocation("America/New_York")
	tm := time.Date(2023, 7, 15, 14, 30, 5, 0, loc)
	
	// Format to ISO8601 UTC
	formatted := FormatISO8601(tm)
	
	// Expected format: "2023-07-15T18:30:05Z" (NYC is UTC-4 in summer)
	expected := "2023-07-15T18:30:05Z"
	
	if formatted != expected {
		t.Errorf("FormatISO8601() = %s; expected %s", formatted, expected)
	}
}

func TestFormatISO8601WithTZ(t *testing.T) {
	// Create a time in a specific timezone
	loc, _ := time.LoadLocation("America/New_York")
	tm := time.Date(2023, 7, 15, 14, 30, 5, 0, loc)
	
	// Format with timezone
	formatted := FormatISO8601WithTZ(tm)
	
	// Should contain the negative offset for New York
	if !strings.Contains(formatted, "-04:00") {
		t.Errorf("FormatISO8601WithTZ() = %s; expected to contain -04:00", formatted)
	}
}

func TestFormatDate(t *testing.T) {
	tm := time.Date(2023, 7, 15, 14, 30, 5, 0, time.UTC)
	
	formatted := FormatDate(tm)
	expected := "2023-07-15"
	
	if formatted != expected {
		t.Errorf("FormatDate() = %s; expected %s", formatted, expected)
	}
}

func TestNow(t *testing.T) {
	// This test checks that Now() returns a time in UTC
	now := Now()
	
	if now.Location() != time.UTC {
		t.Errorf("Now() returned time in %v location; expected UTC", now.Location())
	}
}

func TestStartOfDay(t *testing.T) {
	// Create a time with non-zero hours, minutes, seconds
	tm := time.Date(2023, 7, 15, 14, 30, 5, 123456789, time.UTC)
	
	// Get start of day
	start := StartOfDay(tm)
	
	// Check that date components are preserved
	if start.Year() != 2023 || start.Month() != 7 || start.Day() != 15 {
		t.Errorf("StartOfDay() changed the date: %v", start)
	}
	
	// Check that time components are zeroed
	if start.Hour() != 0 || start.Minute() != 0 || start.Second() != 0 || start.Nanosecond() != 0 {
		t.Errorf("StartOfDay() did not zero the time components: %v", start)
	}
	
	// Check that location is preserved
	if start.Location() != time.UTC {
		t.Errorf("StartOfDay() changed the location from %v to %v", time.UTC, start.Location())
	}
}

func TestEndOfDay(t *testing.T) {
	// Create a time with non-max hours, minutes, seconds
	tm := time.Date(2023, 7, 15, 14, 30, 5, 123456789, time.UTC)
	
	// Get end of day
	end := EndOfDay(tm)
	
	// Check that date components are preserved
	if end.Year() != 2023 || end.Month() != 7 || end.Day() != 15 {
		t.Errorf("EndOfDay() changed the date: %v", end)
	}
	
	// Check that time components are set to end of day
	if end.Hour() != 23 || end.Minute() != 59 || end.Second() != 59 || end.Nanosecond() != 999999999 {
		t.Errorf("EndOfDay() did not set time to end of day: %v", end)
	}
	
	// Check that location is preserved
	if end.Location() != time.UTC {
		t.Errorf("EndOfDay() changed the location from %v to %v", time.UTC, end.Location())
	}
}

func TestStartOfMonth(t *testing.T) {
	// Create a time with non-first day
	tm := time.Date(2023, 7, 15, 14, 30, 5, 123456789, time.UTC)
	
	// Get start of month
	start := StartOfMonth(tm)
	
	// Check that year and month are preserved
	if start.Year() != 2023 || start.Month() != 7 {
		t.Errorf("StartOfMonth() changed the year or month: %v", start)
	}
	
	// Check that day is set to 1
	if start.Day() != 1 {
		t.Errorf("StartOfMonth() did not set day to 1: %v", start)
	}
	
	// Check that time components are zeroed
	if start.Hour() != 0 || start.Minute() != 0 || start.Second() != 0 || start.Nanosecond() != 0 {
		t.Errorf("StartOfMonth() did not zero the time components: %v", start)
	}
}

func TestEndOfMonth(t *testing.T) {
	tests := []struct {
		year     int
		month    time.Month
		expected int // Last day of month
	}{
		{2023, time.January, 31},
		{2023, time.February, 28},
		{2023, time.April, 30},
		{2024, time.February, 29}, // Leap year
	}
	
	for _, tc := range tests {
		t.Run(tc.month.String(), func(t *testing.T) {
			tm := time.Date(tc.year, tc.month, 15, 14, 30, 5, 123456789, time.UTC)
			
			// Get end of month
			end := EndOfMonth(tm)
			
			// Check that year and month are preserved
			if end.Year() != tc.year || end.Month() != tc.month {
				t.Errorf("EndOfMonth() changed the year or month: %v", end)
			}
			
			// Check that day is the last day of the month
			if end.Day() != tc.expected {
				t.Errorf("EndOfMonth() did not set day to %d: %v", tc.expected, end)
			}
			
			// Check that time components are set to end of day
			if end.Hour() != 23 || end.Minute() != 59 || end.Second() != 59 || end.Nanosecond() != 999999999 {
				t.Errorf("EndOfMonth() did not set time to end of day: %v", end)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{"seconds only", 45 * time.Second, "45 seconds"},
		{"one second", 1 * time.Second, "1 second"},
		{"minutes and seconds", 2*time.Minute + 30*time.Second, "2 minutes and 30 seconds"},
		{"one minute", 1 * time.Minute, "1 minute"},
		{"hours, minutes and seconds", 1*time.Hour + 15*time.Minute + 30*time.Second, "1 hour, 15 minutes and 30 seconds"},
		{"one hour", 1 * time.Hour, "1 hour"},
		{"days", 2*24*time.Hour + 5*time.Hour + 15*time.Minute, "2 days, 5 hours and 15 minutes"},
		{"one day", 24 * time.Hour, "1 day"},
		{"zero", 0, "0 seconds"},
	}
	
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := FormatDuration(tc.duration)
			if result != tc.expected {
				t.Errorf("FormatDuration(%v) = %s; expected %s", tc.duration, result, tc.expected)
			}
		})
	}
}

func TestIsWeekend(t *testing.T) {
	tests := []struct {
		date     string
		expected bool
	}{
		{"2023-07-15", true},  // Saturday
		{"2023-07-16", true},  // Sunday
		{"2023-07-17", false}, // Monday
		{"2023-07-21", false}, // Friday
	}
	
	for _, tc := range tests {
		t.Run(tc.date, func(t *testing.T) {
			date, _ := time.Parse("2006-01-02", tc.date)
			result := IsWeekend(date)
			if result != tc.expected {
				t.Errorf("IsWeekend(%s) = %v; expected %v", tc.date, result, tc.expected)
			}
		})
	}
}

func TestAddBusinessDays(t *testing.T) {
	// Start on a Monday (2023-07-17)
	start, _ := time.Parse("2006-01-02", "2023-07-17")
	
	tests := []struct {
		name      string
		days      int
		expected  string // YYYY-MM-DD
		dayOfWeek time.Weekday
	}{
		{"add 1 day", 1, "2023-07-18", time.Tuesday},
		{"add 5 days", 5, "2023-07-24", time.Monday},     // 5 business days = 7 calendar days (including weekend)
		{"add 10 days", 10, "2023-07-31", time.Monday},   // 10 business days = 14 calendar days (including 2 weekends)
		{"add 0 days", 0, "2023-07-17", time.Monday},     // No change
		{"add 7 days", 7, "2023-07-26", time.Wednesday},  // 7 business days = 9 calendar days (including 1 weekend)
	}
	
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := AddBusinessDays(start, tc.days)
			
			// Format to YYYY-MM-DD for comparison
			resultStr := result.Format("2006-01-02")
			if resultStr != tc.expected {
				t.Errorf("AddBusinessDays(%s, %d) = %s; expected %s", start.Format("2006-01-02"), tc.days, resultStr, tc.expected)
			}
			
			// Check weekday
			if result.Weekday() != tc.dayOfWeek {
				t.Errorf("AddBusinessDays(%s, %d) returned %s which is a %s; expected %s", 
					start.Format("2006-01-02"), tc.days, resultStr, result.Weekday(), tc.dayOfWeek)
			}
			
			// Check that result is not a weekend
			if IsWeekend(result) {
				t.Errorf("AddBusinessDays(%s, %d) returned a weekend: %s", start.Format("2006-01-02"), tc.days, resultStr)
			}
		})
	}
}

func TestIsSameDay(t *testing.T) {
	tests := []struct {
		name     string
		t1       time.Time
		t2       time.Time
		expected bool
	}{
		{
			name:     "same day different time",
			t1:       time.Date(2023, 7, 15, 10, 30, 0, 0, time.UTC),
			t2:       time.Date(2023, 7, 15, 22, 45, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "different day",
			t1:       time.Date(2023, 7, 15, 10, 30, 0, 0, time.UTC),
			t2:       time.Date(2023, 7, 16, 10, 30, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "different month",
			t1:       time.Date(2023, 7, 15, 10, 30, 0, 0, time.UTC),
			t2:       time.Date(2023, 8, 15, 10, 30, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "different year",
			t1:       time.Date(2023, 7, 15, 10, 30, 0, 0, time.UTC),
			t2:       time.Date(2024, 7, 15, 10, 30, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "same day different timezone",
			t1:       time.Date(2023, 7, 15, 10, 30, 0, 0, time.UTC),
			t2:       time.Date(2023, 7, 15, 5, 30, 0, 0, time.FixedZone("EST", -5*60*60)),
			expected: true,
		},
	}
	
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := IsSameDay(tc.t1, tc.t2)
			if result != tc.expected {
				t.Errorf("IsSameDay(%v, %v) = %v; expected %v", tc.t1, tc.t2, result, tc.expected)
			}
		})
	}
}