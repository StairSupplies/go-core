/*
Package str provides string formatting and manipulation utilities.

It offers functions for currency formatting, string truncation, padding,
and other common string operations.

# Currency Formatting

Format a number as currency with commas and 2 decimal places:

    // Returns "1,234.50"
    formattedPrice := str.FormatPriceToCurrency(1234.5)

Format any float with commas for thousands:

    // Returns "1,234,567.89"
    formatted := str.FormatFloatWithCommas(1234567.89)

# String Manipulation

Reverse a string:

    // Returns "olleh"
    reversed := str.ReverseString("hello")

Truncate a string with ellipsis if it exceeds a length:

    // Returns "This is a lo..."
    truncated := str.TruncateString("This is a long string", 13)

# Padding

Pad a string on the left or right:

    // Returns "00042"
    paddedLeft := str.PadLeft("42", '0', 5)
    
    // Returns "Hello     "
    paddedRight := str.PadRight("Hello", ' ', 10)

# String Testing

Check if a string is empty or contains only whitespace:

    if str.IsEmpty(input) {
        // Handle empty input
    }
*/
package str