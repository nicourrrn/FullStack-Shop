package subj

import (
	"math/rand"
	"time"
)

type User struct {
	ID           int
	OrderIDs     []int
	Login, Email string
	HomeAddress  string
	PassHash     string
	AccessToken  string
	LastVisit    time.Time
}

func NewUser(Email, PassHath string) *User {
	return &User{
		ID:       rand.Int(),
		Email:    Email,
		PassHash: PassHath,
	}
}
