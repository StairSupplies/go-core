package str

import (
	"math"
	"strconv"
	"strings"
)

// FormatPriceToCurrency formats a float64 price to a currency string with commas and 2 decimal places
func FormatPriceToCurrency(price float64) string {
	clampedPrice := math.Round(price*100.0) / 100.0
	return FormatFloatWithCommas(clampedPrice)
}

// FormatFloatWithCommas formats a float64 with commas for thousands and 2 decimal places
func FormatFloatWithCommas(f float64) string {
	str := strconv.FormatFloat(f, 'f', 2, 64)
	parts := strings.Split(str, ".")

	reversed := ReverseString(parts[0])
	reversedWithCommas := ""

	for i, c := range reversed {
		if i > 0 && i%3 == 0 {
			reversedWithCommas += ","
		}
		reversedWithCommas += string(c)
	}

	result := ReverseString(reversedWithCommas)

	if len(parts) > 1 {
		result += "." + parts[1]
	}

	return result
}

// ReverseString reverses a string
func ReverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// TruncateString truncates a string to the given max length and adds an ellipsis if needed
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		if maxLen <= 0 {
			return ""
		}
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// IsEmpty checks if a string is empty or contains only whitespace
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}
