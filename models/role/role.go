package role

import (
	"time"

	"github.com/jbaikge/sparky/modules/pagination"
)

type Role struct {
	RoleId      int
	Name        string
	Permissions []string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (r Role) Validate() (errs map[string]string) {
	errs = make(map[string]string)

	if r.Name == "" {
		errs["Name"] = "Name is required"
	}

	return
}

type UserRoles struct {
	UserId  int
	RoleIds []int
}

type RoleUsers struct {
	RoleId  int
	UserIds []int
}

type RoleListParams struct {
	Pagination *pagination.Pagination
}
