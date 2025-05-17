package jsonutils

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

type TestStruct struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

func TestDecode(t *testing.T) {
	t.Run("valid JSON", func(t *testing.T) {
		jsonStr := `{"name":"John","age":30,"email":"john@example.com"}`
		reader := strings.NewReader(jsonStr)
		
		var result TestStruct
		err := Decode(reader, &result)
		
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		
		if result.Name != "John" {
			t.Errorf("Expected Name to be 'John', got '%s'", result.Name)
		}
		if result.Age != 30 {
			t.Errorf("Expected Age to be 30, got %d", result.Age)
		}
		if result.Email != "john@example.com" {
			t.Errorf("Expected Email to be 'john@example.com', got '%s'", result.Email)
		}
	})
	
	t.Run("invalid JSON syntax", func(t *testing.T) {
		jsonStr := `{"name":"John","age":30,"email":"john@example.com`
		reader := strings.NewReader(jsonStr)
		
		var result TestStruct
		err := Decode(reader, &result)
		
		if err == nil {
			t.Fatal("Expected error for invalid JSON syntax, got nil")
		}
		
		if !strings.Contains(err.Error(), "badly-formed JSON") {
			t.Errorf("Expected error message to contain 'badly-formed JSON', got '%s'", err.Error())
		}
	})
	
	t.Run("wrong type", func(t *testing.T) {
		jsonStr := `{"name":"John","age":"thirty","email":"john@example.com"}`
		reader := strings.NewReader(jsonStr)
		
		var result TestStruct
		err := Decode(reader, &result)
		
		if err == nil {
			t.Fatal("Expected error for wrong type, got nil")
		}
		
		if !strings.Contains(err.Error(), "incorrect JSON type") {
			t.Errorf("Expected error message to contain 'incorrect JSON type', got '%s'", err.Error())
		}
	})
	
	t.Run("empty body", func(t *testing.T) {
		jsonStr := ``
		reader := strings.NewReader(jsonStr)
		
		var result TestStruct
		err := Decode(reader, &result)
		
		if err == nil {
			t.Fatal("Expected error for empty body, got nil")
		}
		
		if err.Error() != "body must not be empty" {
			t.Errorf("Expected error message to be 'body must not be empty', got '%s'", err.Error())
		}
	})
	
	t.Run("unknown field", func(t *testing.T) {
		jsonStr := `{"name":"John","age":30,"email":"john@example.com","unknown_field":"value"}`
		reader := strings.NewReader(jsonStr)
		
		var result TestStruct
		err := Decode(reader, &result)
		
		if err == nil {
			t.Fatal("Expected error for unknown field, got nil")
		}
		
		if !strings.Contains(err.Error(), "body contains unknown field") {
			t.Errorf("Expected error message to contain 'body contains unknown field', got '%s'", err.Error())
		}
	})
	
	t.Run("multiple JSON values", func(t *testing.T) {
		jsonStr := `{"name":"John","age":30,"email":"john@example.com"}{"name":"Jane"}`
		reader := strings.NewReader(jsonStr)
		
		var result TestStruct
		err := Decode(reader, &result)
		
		if err == nil {
			t.Fatal("Expected error for multiple JSON values, got nil")
		}
		
		if err.Error() != "body must only contain a single JSON value" {
			t.Errorf("Expected error message to be 'body must only contain a single JSON value', got '%s'", err.Error())
		}
	})
}

func TestPretty(t *testing.T) {
	data := TestStruct{
		Name:  "John",
		Age:   30,
		Email: "john@example.com",
	}
	
	pretty, err := Pretty(data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	expected := `{
  "name": "John",
  "age": 30,
  "email": "john@example.com"
}`
	
	if pretty != expected {
		t.Errorf("Pretty output not as expected. \nGot: %s\nWant: %s", pretty, expected)
	}
	
	// Test with a value that can't be marshaled
	_, err = Pretty(make(chan int))
	if err == nil {
		t.Fatal("Expected error for unmarshallable value, got nil")
	}
}

func TestEncode(t *testing.T) {
	data := TestStruct{
		Name:  "John",
		Age:   30,
		Email: "john@example.com",
	}
	
	var buf bytes.Buffer
	err := Encode(&buf, data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	expected := `{
  "name": "John",
  "age": 30,
  "email": "john@example.com"
}
`
	
	if buf.String() != expected {
		t.Errorf("Encode output not as expected. \nGot: %s\nWant: %s", buf.String(), expected)
	}
	
	// Test with a value that can't be marshaled
	err = Encode(&buf, make(chan int))
	if err == nil {
		t.Fatal("Expected error for unmarshallable value, got nil")
	}
}

// Helper type for testing error case in the decode function
type ErrorReader struct{}

func (r ErrorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}

func TestDecodeWithReaderError(t *testing.T) {
	var result TestStruct
	err := Decode(ErrorReader{}, &result)
	
	if err == nil {
		t.Fatal("Expected error for reader error, got nil")
	}
}

func TestDecodeWithNilValue(t *testing.T) {
	jsonStr := `{"name":"John","age":30,"email":"john@example.com"}`
	reader := strings.NewReader(jsonStr)
	
	err := Decode(reader, nil)
	
	if err == nil {
		t.Fatal("Expected error for nil destination, got nil")
	}
}

func TestPrettyWithComplexTypes(t *testing.T) {
	// Test with nested structures
	data := map[string]interface{}{
		"person": TestStruct{
			Name:  "John",
			Age:   30,
			Email: "john@example.com",
		},
		"tags":    []string{"tag1", "tag2", "tag3"},
		"numbers": []int{1, 2, 3, 4, 5},
		"nested": map[string]interface{}{
			"key1": "value1",
			"key2": 42,
			"key3": true,
		},
	}
	
	pretty, err := Pretty(data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	// Just check if it contains expected keys
	for _, key := range []string{"person", "tags", "numbers", "nested"} {
		if !strings.Contains(pretty, key) {
			t.Errorf("Expected pretty output to contain key '%s', but it doesn't", key)
		}
	}
}