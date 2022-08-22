package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type User struct {
	Name        string `json:"full_name"`
	Age         int    `json:"age,omitempty"`
	Active      bool   `json:"-"`
	LastLoginAt time.Time
}

func (u *User) MarshalJSON() ([]byte, error) {
	type Alias User
	return json.Marshal(&struct {
		*Alias
		LastLoginAt int64 `json:"last_login_at"`
	}{
		Alias:       (*Alias)(u),
		LastLoginAt: u.LastLoginAt.Unix(),
	})
}

func main() {

	u, err := json.Marshal(User{Name: "Ball", Age: 35, Active: true, LastLoginAt: time.Now()})
	if err != nil {
		panic(err)
	}
	fmt.Println(string(u)) // {"full_name":"Ball","age":35}
}
