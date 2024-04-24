package user

import (
	"fmt"
	"strings"
	"time"
)

type User struct {
	UserId    int
	FirstName string
	LastName  string
	Email     string
	Password  string
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

func (u User) Validate() (errs map[string]string) {
	errs = make(map[string]string)

	if u.FirstName == "" {
		errs["FirstName"] = "First name is required"
	}
	if u.LastName == "" {
		errs["LastName"] = "Last name is required"
	}
	if u.Email == "" {
		errs["Email"] = "Email is required"
	}
	if !strings.ContainsRune(u.Email, '@') {
		errs["Email"] = "Invalid email address"
	}
	if u.Password == "" {
		errs["Password"] = "Password is required"
	}

	return
}
