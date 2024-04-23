package page

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log/slog"

	"github.com/jbaikge/sparky/models/user"
	"github.com/jbaikge/sparky/modules/database"
	"github.com/jbaikge/sparky/modules/middleware"
)

type Page struct {
	Data  map[string]interface{}
	admin *user.User
	db    database.Database
	errs  map[string]string
	tpl   *template.Template
}

func New(ctx context.Context) *Page {
	p := &Page{
		Data: make(map[string]interface{}),
		errs: make(map[string]string),
	}

	p.db = ctx.Value(middleware.ContextDatabase).(database.Database)
	p.tpl = ctx.Value(middleware.ContextTemplate).(*template.Template)
	if admin := ctx.Value(middleware.ContextAdminUser); admin != nil {
		p.admin = admin.(*user.User)
	}
	// TODO capture regular user?

	p.Data["Admin"] = p.admin

	return p
}

func (p *Page) AddError(key, err string) {
	p.errs[key] = err
}

func (p *Page) Admin() *user.User {
	return p.admin
}

func (p *Page) Database() database.Database {
	return p.db
}

func (p *Page) HasErrors() bool {
	return len(p.errs) > 0
}

func (p *Page) Render(w io.Writer, name string) {
	p.Data["Errors"] = p.errs
	p.Data["RenderedTemplate"] = name
	if err := p.tpl.ExecuteTemplate(w, name, p.Data); err != nil {
		fmt.Fprintf(w, "Error during execution: %v", err)
		slog.Error("template execution failed", "error", err)
	}
}
