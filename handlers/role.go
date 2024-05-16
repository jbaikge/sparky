package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/jbaikge/sparky/models/role"
	"github.com/jbaikge/sparky/modules/page"
	"github.com/jbaikge/sparky/modules/pagination"
)

func roleList(w http.ResponseWriter, r *http.Request) {
	p := page.New(r.Context())
	p.Data["PageActiveRoles"] = true

	roleRepo := role.NewRoleRepository(p.Database())
	params := &role.RoleListParams{
		Pagination: pagination.NewPagination(r),
	}
	roles, err := roleRepo.ListRoles(r.Context(), params)
	if err != nil {
		slog.Error("ListRoles failed", "error", err)
	}
	p.Data["Roles"] = roles
	p.Data["Params"] = params
	p.Render(w, "role/list")
}

func roleForm(w http.ResponseWriter, r *http.Request) {
	tpl := "role/form"

	p := page.New(r.Context())
	p.Data["PageActiveRoles"] = true
	p.Data["FormAction"] = r.URL.Path

	urlId := r.PathValue("id")
	isEditing := urlId != ""

	roleRepo := role.NewRoleRepository(p.Database())

	formRole := &role.Role{}

	if isEditing {
		id, err := strconv.Atoi(urlId)
		if err != nil {
			slog.Error("invalid user id", "value", urlId)
			http.Redirect(w, r, "/admin/roles/list", http.StatusSeeOther)
			return
		}

		formRole, err = roleRepo.GetRoleById(r.Context(), id)
		if err != nil {
			slog.Error("role does not exist", "id", id)
			http.Redirect(w, r, "/admin/roles/list", http.StatusSeeOther)
			return
		}
	}
	p.Data["Role"] = formRole

	if r.Method == http.MethodGet {
		p.Render(w, tpl)
		return
	}

	for key, err := range formRole.Validate() {
		p.AddError(key, err)
	}
	if p.HasErrors() {
		p.Render(w, tpl)
		return
	}

	if err := roleRepo.UpsertRole(r.Context(), formRole); err != nil {
		slog.Error("failed to upsert role", "id", formRole.RoleId, "error", err)
		p.AddError("Database", err.Error())
		p.Render(w, tpl)
		return
	}

	http.Redirect(w, r, "/admin/roles/list", http.StatusSeeOther)
}
