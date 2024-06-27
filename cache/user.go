package cache

import (
	"encoding/json"
	"time"
)

type User struct {
	ID        int       `json:"id"` //when struct is serialized to json these names are used
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// MarshalBinary encodes the receiver into a binary form and returns the result.
func (u *User) MarshalBinary() ([]byte, error) {
	// Implementing MarshalBinary is not necessary for JSON serialization
	return json.Marshal(u) //struct to json
}

// UnmarshalBinary decodes the receiver from binary form.
func (u *User) UnmarshalBinary(data []byte) error {
	// Implementing UnmarshalBinary is not necessary for JSON deserialization
	return json.Unmarshal(data, u)
}

//redis or memcache storage systems uses json type so we habe to convert binary to json
