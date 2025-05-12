/*
Package timeutils provides date and time utility functions.

It offers common time format constants, date range calculations,
and human-readable duration formatting.

# Time Constants

The package provides constants for common time formats:

	// ISO format with UTC timezone: "2006-01-02T15:04:05Z"
	ISO8601UTC

	// ISO format with timezone offset: "2006-01-02T15:04:05-07:00"
	ISO8601

	// Date only: "2006-01-02"
	DateOnly

	// Time only: "15:04:05"
	TimeOnly

	// Date and time without timezone: "2006-01-02 15:04:05"
	DateTime

	// Human-readable date: "January 2, 2006"
	HumanReadableDate

	// Human-readable date and time: "January 2, 2006 3:04 PM"
	HumanReadableDateTime

# Parsing

Parse time strings in multiple formats:

	// Tries multiple formats automatically
	t, err := timeutils.ParseAny("2023-05-15T14:30:00Z")
	t, err := timeutils.ParseAny("2023-05-15")
	t, err := timeutils.ParseAny("May 15, 2023")

# Formatting

Format times in standard formats:

	// Returns "2023-05-15T14:30:00Z"
	formatted := timeutils.FormatISO8601(time.Now())

	// Returns "2023-05-15T14:30:00-07:00"
	formatted := timeutils.FormatISO8601WithTZ(time.Now())

	// Returns "2023-05-15"
	formatted := timeutils.FormatDate(time.Now())

# Time Ranges

Get the start and end of day or month:

	// Returns 2023-05-15 00:00:00
	start := timeutils.StartOfDay(time.Now())

	// Returns 2023-05-15 23:59:59.999999999
	end := timeutils.EndOfDay(time.Now())

	// Returns 2023-05-01 00:00:00
	start := timeutils.StartOfMonth(time.Now())

	// Returns 2023-05-31 23:59:59.999999999
	end := timeutils.EndOfMonth(time.Now())

# Duration Formatting

Format durations in a human-readable way:

	// Returns "2 days, 3 hours, 45 minutes and 10 seconds"
	formatted := timeutils.FormatDuration(
	    51 * time.Hour + 45 * time.Minute + 10 * time.Second
	)

# Business Logic

Work with business days and day comparisons:

	// Check if a date is a weekend
	if timeutils.IsWeekend(time.Now()) {
	    // It's the weekend!
	}

	// Add business days (skipping weekends)
	deadline := timeutils.AddBusinessDays(time.Now(), 5)

	// Check if two times are on the same day
	if timeutils.IsSameDay(user.CreatedAt, time.Now()) {
	    // User was created today
	}

# UTC Helper

Get the current time in UTC:

	now := timeutils.Now()
*/
package timeutils
