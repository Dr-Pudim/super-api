package actions

import (
    "net/http"
    
	"github.com/gobuffalo/buffalo"
)

// SupersCreate default implementation.
func SupersCreate(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("supers/create.html"))
}

