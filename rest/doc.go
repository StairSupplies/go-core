/*
Package rest provides a fluent REST client for making HTTP requests.

It offers a chainable API for constructing and sending HTTP requests,
with built-in support for JSON serialization/deserialization.

# Client Setup

Create a new client with default settings (30-second timeout):

    client := rest.NewClient()

Configure the client with a fluent interface:

    client := rest.NewClient().
        WithBaseURL("https://api.example.com").
        WithHeader("Authorization", "Bearer " + token).
        WithTimeout(5 * time.Second)

# Making Requests

The client supports all common HTTP methods:

    // GET request
    var users []User
    err := client.Get(context.Background(), "/users", &users)
    
    // POST request
    newUser := User{Name: "John", Email: "john@example.com"}
    var createdUser User
    err := client.Post(context.Background(), "/users", newUser, &createdUser)
    
    // PUT request
    updatedUser := User{ID: 1, Name: "John Updated"}
    var result User
    err := client.Put(context.Background(), "/users/1", updatedUser, &result)
    
    // PATCH request
    patch := map[string]string{"name": "John Updated"}
    var result User
    err := client.Patch(context.Background(), "/users/1", patch, &result)
    
    // DELETE request
    var result struct { Success bool `json:"success"` }
    err := client.Delete(context.Background(), "/users/1", &result)

# Advanced Usage

For more control, use the Request method directly:

    var response struct {
        ID   int    `json:"id"`
        Name string `json:"name"`
    }
    
    err := client.Request(
        context.Background(),
        http.MethodGet,
        "/users/1",
        nil,
        &response,
    )
*/
package rest