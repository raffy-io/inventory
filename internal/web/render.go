package web

import (
	"bytes"
	"net/http"

	"github.com/a-h/templ"
	"github.com/alexedwards/scs/v2"
	"github.com/justinas/nosurf"
	"github.com/raffy-io/inventory/internal/shared"
)

type Renderer struct {
	SessionManager *scs.SessionManager
}

func NewRenderer(sm *scs.SessionManager) *Renderer {
	return &Renderer{SessionManager: sm}
}

func (rd *Renderer) Render(w http.ResponseWriter, r *http.Request, title string, component templ.Component) error {
	csrfToken := nosurf.Token(r)

	var buf bytes.Buffer
	if err := shared.Base(title, csrfToken, component).Render(r.Context(), &buf); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err := buf.WriteTo(w)
	return err
}
