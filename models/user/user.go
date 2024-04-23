package user

import (
	"fmt"
	"time"
)

type User struct {
	UserId    int
	FirstName string
	LastName  string
	Email     string
	Password  string
	StartDate time.Time
	EndDate   time.Time
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u User) Initials() string {
	return u.FirstName[0:1] + u.LastName[0:1]
}

func (u User) Name() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}
