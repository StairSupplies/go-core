/*
Package jsonutils provides enhanced JSON utilities for encoding and decoding.

It offers better error handling and additional features beyond the standard
encoding/json package.

# Pretty Printing

Pretty prints JSON for debugging and logging:

    user := User{ID: 1, Name: "John"}
    prettyJSON, err := jsonutils.Pretty(user)
    fmt.Println(prettyJSON)

# Encoding

Writes JSON data to a writer with nice formatting:

    user := User{ID: 1, Name: "John"}
    var buf bytes.Buffer
    err := jsonutils.Encode(&buf, user)
    fmt.Println(buf.String())

# Decoding with Enhanced Error Handling

Reads JSON from a reader into a target struct with detailed error messages:

    var user User
    err := jsonutils.Decode(r.Body, &user)
    if err != nil {
        // Detailed error message available
        fmt.Println(err)
    }

Decoding provides detailed error messages for common JSON parsing issues:

- Syntax errors with position information
- Type mismatch errors with field names
- Empty body detection
- Unknown field identification
- Multiple JSON value detection
*/
package jsonutils