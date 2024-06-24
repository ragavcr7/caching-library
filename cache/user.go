// user.go
/*
package cache

import (
	"encoding/json"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// MarshalRedis serializes the User struct to a JSON byte slice
func (u User) MarshalRedis() ([]byte, error) {
	return json.Marshal(u)
}

// UnmarshalRedis deserializes a JSON byte slice into the User struct
func (u *User) UnmarshalRedis(data []byte) error {
	return json.Unmarshal(data, u)
}
*/
package cache

import (
	"encoding/json"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// MarshalBinary encodes the receiver into a binary form and returns the result.
func (u *User) MarshalBinary() ([]byte, error) {
	// Implementing MarshalBinary is not necessary for JSON serialization
	return json.Marshal(u)
}

// UnmarshalBinary decodes the receiver from binary form.
func (u *User) UnmarshalBinary(data []byte) error {
	// Implementing UnmarshalBinary is not necessary for JSON deserialization
	return json.Unmarshal(data, u)
}
