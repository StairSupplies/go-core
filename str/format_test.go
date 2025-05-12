package str

import "testing"

func TestFormatPriceToCurrency(t *testing.T) {
	tests := []struct {
		name     string
		price    float64
		expected string
	}{
		{"zero", 0, "0.00"},
		{"integer", 1000, "1,000.00"},
		{"decimal", 1234.56, "1,234.56"},
		{"negative", -1234.56, "-1,234.56"},
		{"rounding up", 1234.567, "1,234.57"},
		{"rounding down", 1234.563, "1,234.56"},
		{"small decimal", 0.12, "0.12"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := FormatPriceToCurrency(tc.price)
			if result != tc.expected {
				t.Errorf("FormatPriceToCurrency(%f) = %s; expected %s", tc.price, result, tc.expected)
			}
		})
	}
}

func TestFormatFloatWithCommas(t *testing.T) {
	tests := []struct {
		name     string
		f        float64
		expected string
	}{
		{"zero", 0, "0.00"},
		{"integer", 1000, "1,000.00"},
		{"decimal", 1234.56, "1,234.56"},
		{"large number", 1234567.89, "1,234,567.89"},
		{"negative", -1234.56, "-1,234.56"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := FormatFloatWithCommas(tc.f)
			if result != tc.expected {
				t.Errorf("FormatFloatWithCommas(%f) = %s; expected %s", tc.f, result, tc.expected)
			}
		})
	}
}

func TestReverseString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"a", "a"},
		{"abc", "cba"},
		{"hello", "olleh"},
		{"racecar", "racecar"}, // Palindrome
		{"áéíóú", "úóíéá"},     // Unicode characters
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result := ReverseString(tc.input)
			if result != tc.expected {
				t.Errorf("ReverseString(%s) = %s; expected %s", tc.input, result, tc.expected)
			}
		})
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		maxLen   int
		expected string
	}{
		{"empty string", "", 10, ""},
		{"short string", "hello", 10, "hello"},
		{"exact length", "hello", 5, "hello"},
		{"truncate", "hello world", 8, "hello..."},
		{"maxLen less than 3", "hello", 2, "he"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := TruncateString(tc.s, tc.maxLen)
			if result != tc.expected {
				t.Errorf("TruncateString(%s, %d) = %s; expected %s", tc.s, tc.maxLen, result, tc.expected)
			}
		})
	}
}

func TestPadLeft(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		padChar  rune
		length   int
		expected string
	}{
		{"empty string", "", ' ', 5, "     "},
		{"already at length", "hello", ' ', 5, "hello"},
		{"longer than length", "hello", ' ', 3, "hello"},
		{"needs padding", "hello", ' ', 10, "     hello"},
		{"different pad char", "123", '0', 5, "00123"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := PadLeft(tc.s, tc.padChar, tc.length)
			if result != tc.expected {
				t.Errorf("PadLeft(%s, %c, %d) = %s; expected %s", tc.s, tc.padChar, tc.length, result, tc.expected)
			}
		})
	}
}

func TestPadRight(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		padChar  rune
		length   int
		expected string
	}{
		{"empty string", "", ' ', 5, "     "},
		{"already at length", "hello", ' ', 5, "hello"},
		{"longer than length", "hello", ' ', 3, "hello"},
		{"needs padding", "hello", ' ', 10, "hello     "},
		{"different pad char", "123", '0', 5, "12300"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := PadRight(tc.s, tc.padChar, tc.length)
			if result != tc.expected {
				t.Errorf("PadRight(%s, %c, %d) = %s; expected %s", tc.s, tc.padChar, tc.length, result, tc.expected)
			}
		})
	}
}

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		expected bool
	}{
		{"empty string", "", true},
		{"whitespace", "   ", true},
		{"tabs and newlines", "\t\n", true},
		{"non-empty", "hello", false},
		{"whitespace with text", "  hello  ", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := IsEmpty(tc.s)
			if result != tc.expected {
				t.Errorf("IsEmpty(%s) = %v; expected %v", tc.s, result, tc.expected)
			}
		})
	}
}