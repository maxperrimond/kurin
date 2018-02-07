package domain

import (
	"encoding/hex"
	"math/rand"
	"time"
)

type (
	User struct {
		Id        string
		Username  string
		Email     string
		CreatedAt time.Time
	}
)

func (user *User) GenerateId() {
	id := make([]byte, 16)
	rand.Read(id)
	user.Id = hex.EncodeToString(id)
}
