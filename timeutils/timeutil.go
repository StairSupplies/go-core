package timeutils

import (
	"fmt"
	"time"
)

// Standard time formats
const (
	// ISO8601UTC represents time in ISO 8601 format with UTC timezone: "2006-01-02T15:04:05Z"
	ISO8601UTC = "2006-01-02T15:04:05Z"

	// ISO8601 represents time in ISO 8601 format with timezone offset: "2006-01-02T15:04:05-07:00"
	ISO8601 = "2006-01-02T15:04:05-07:00"

	// DateOnly represents date in YYYY-MM-DD format: "2006-01-02"
	DateOnly = "2006-01-02"

	// TimeOnly represents time in HH:MM:SS format: "15:04:05"
	TimeOnly = "15:04:05"

	// DateTime represents date and time without timezone: "2006-01-02 15:04:05"
	DateTime = "2006-01-02 15:04:05"

	// HumanReadableDate represents date in human-readable format: "January 2, 2006"
	HumanReadableDate = "January 2, 2006"

	// HumanReadableDateTime represents date and time in human-readable format: "January 2, 2006 3:04 PM"
	HumanReadableDateTime = "January 2, 2006 3:04 PM"
)

// ParseAny attempts to parse a time string using multiple common formats
func ParseAny(timeString string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		ISO8601UTC,
		ISO8601,
		DateOnly,
		DateTime,
		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05.999999999Z07:00",
		"2006-01-02T15:04:05.999999999",
		"2006-01-02 15:04:05.999999999",
		"2006-01-02 15:04",
		"2006/01/02",
		"01/02/2006",
		"01-02-2006",
		"02/01/2006", // DD/MM/YYYY
		"02-01-2006", // DD-MM-YYYY
		"January 2, 2006",
		"Jan 2, 2006",
		"Jan 2 2006",
		"2 Jan 2006",
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339Nano,
		time.Kitchen,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timeString); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse time string: %s", timeString)
}

// FormatISO8601 formats time as ISO 8601 UTC
func FormatISO8601(t time.Time) string {
	return t.UTC().Format(ISO8601UTC)
}

// FormatISO8601WithTZ formats time as ISO 8601 with timezone
func FormatISO8601WithTZ(t time.Time) string {
	return t.Format(ISO8601)
}

// FormatDate formats time as YYYY-MM-DD
func FormatDate(t time.Time) string {
	return t.Format(DateOnly)
}

// Now returns the current time in UTC
func Now() time.Time {
	return time.Now().UTC()
}

// StartOfDay returns the start of day for a given time
func StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// EndOfDay returns the end of day for a given time
func EndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
}

// StartOfMonth returns the start of month for a given time
func StartOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// EndOfMonth returns the end of month for a given time
func EndOfMonth(t time.Time) time.Time {
	return StartOfMonth(t).AddDate(0, 1, 0).Add(-time.Nanosecond)
}

// FormatDuration formats a duration in a human-readable format
func FormatDuration(d time.Duration) string {
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	parts := make([]string, 0, 4)
	if days > 0 {
		if days == 1 {
			parts = append(parts, "1 day")
		} else {
			parts = append(parts, fmt.Sprintf("%d days", days))
		}
	}
	if hours > 0 {
		if hours == 1 {
			parts = append(parts, "1 hour")
		} else {
			parts = append(parts, fmt.Sprintf("%d hours", hours))
		}
	}
	if minutes > 0 {
		if minutes == 1 {
			parts = append(parts, "1 minute")
		} else {
			parts = append(parts, fmt.Sprintf("%d minutes", minutes))
		}
	}
	if seconds > 0 || len(parts) == 0 {
		if seconds == 1 {
			parts = append(parts, "1 second")
		} else {
			parts = append(parts, fmt.Sprintf("%d seconds", seconds))
		}
	}

	if len(parts) == 1 {
		return parts[0]
	}

	result := ""
	for i, part := range parts {
		if i == 0 {
			result = part
		} else if i == len(parts)-1 {
			result = result + " and " + part
		} else {
			result = result + ", " + part
		}
	}

	return result
}

// IsWeekend returns true if the time is a weekend (Saturday or Sunday)
func IsWeekend(t time.Time) bool {
	day := t.Weekday()
	return day == time.Saturday || day == time.Sunday
}

// AddBusinessDays adds the specified number of business days to a time
// (skipping weekends but not holidays)
func AddBusinessDays(t time.Time, days int) time.Time {
	result := t
	for days > 0 {
		result = result.AddDate(0, 0, 1)
		if !IsWeekend(result) {
			days--
		}
	}
	return result
}

// IsSameDay returns true if two times are on the same day, ignoring time parts
func IsSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
