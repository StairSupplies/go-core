package jsonutils_test

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/StairSupplies/go-core/jsonutils"
)

// ExamplePretty demonstrates pretty-printing a JSON object.
func ExamplePretty() {
	// Create a sample data structure
	data := map[string]interface{}{
		"name": "John Doe",
		"age":  30,
		"address": map[string]string{
			"street": "123 Main St",
			"city":   "Anytown",
			"state":  "CA",
			"zip":    "12345",
		},
		"hobbies": []string{"reading", "hiking", "photography"},
	}

	// Pretty-print the JSON
	prettyJSON, err := jsonutils.Pretty(data)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println(prettyJSON)
	// Output:
	// {
	//   "address": {
	//     "city": "Anytown",
	//     "state": "CA",
	//     "street": "123 Main St",
	//     "zip": "12345"
	//   },
	//   "age": 30,
	//   "hobbies": [
	//     "reading",
	//     "hiking",
	//     "photography"
	//   ],
	//   "name": "John Doe"
	// }
}

// ExampleEncode demonstrates encoding a data structure to JSON.
func ExampleEncode() {
	// Create a sample struct
	type Person struct {
		Name    string   `json:"name"`
		Age     int      `json:"age"`
		Email   string   `json:"email"`
		Hobbies []string `json:"hobbies,omitempty"`
	}

	// Create an instance of the struct
	person := Person{
		Name:    "Jane Smith",
		Age:     28,
		Email:   "jane.smith@example.com",
		Hobbies: []string{"coding", "gaming"},
	}

	// Create a buffer to hold the JSON output
	var buf bytes.Buffer

	// Encode the struct to JSON
	err := jsonutils.Encode(&buf, person)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Print the resulting JSON
	fmt.Println(buf.String())
	// Output:
	// {
	//   "name": "Jane Smith",
	//   "age": 28,
	//   "email": "jane.smith@example.com",
	//   "hobbies": [
	//     "coding",
	//     "gaming"
	//   ]
	// }
}

// ExampleDecode demonstrates decoding JSON data into a Go struct.
func ExampleDecode() {
	// Define a struct to hold the decoded data
	type Product struct {
		ID    int     `json:"id"`
		Name  string  `json:"name"`
		Price float64 `json:"price"`
		Tags  []string `json:"tags"`
	}

	// Sample JSON data
	jsonData := `{
		"id": 1001,
		"name": "Smartphone",
		"price": 599.99,
		"tags": ["electronics", "gadget", "mobile"]
	}`

	// Create a reader from the JSON string
	reader := strings.NewReader(jsonData)

	// Create an instance of the struct to decode into
	var product Product

	// Decode the JSON data
	err := jsonutils.Decode(reader, &product)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Use the decoded struct
	fmt.Printf("Product: %s (ID: %d)\n", product.Name, product.ID)
	fmt.Printf("Price: $%.2f\n", product.Price)
	fmt.Printf("Tags: %v\n", product.Tags)
	// Output:
	// Product: Smartphone (ID: 1001)
	// Price: $599.99
	// Tags: [electronics gadget mobile]
}

// ExampleDecode_errorHandling demonstrates the enhanced error handling when decoding invalid JSON.
func ExampleDecode_errorHandling() {
	// Define a struct to hold the decoded data
	type User struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	// Example of invalid JSON with syntax error
	jsonData := `{"id": 1, "name": "John", "email": "john@example.com"`

	// Try to decode the invalid JSON
	var user User
	err := jsonutils.Decode(strings.NewReader(jsonData), &user)
	
	// Print the error message
	fmt.Printf("Error: %v\n", err)
	// Output: Error: body contains badly-formed JSON
}