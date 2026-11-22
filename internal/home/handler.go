package home

import (
	"net/http"

	"github.com/raffy-io/inventory/internal/web"
)

type homeHandler struct {
	render *web.Renderer
}

func newHomeHandler(render *web.Renderer)*homeHandler{
	return &homeHandler{
		render: render,
	}
}

func(h *homeHandler) serveHomepage(w http.ResponseWriter, r *http.Request){
	component := homepage()
	title := "Welcome"

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.render.Render(w,r,title,component); err != nil {
		http.Error(w,"Failed to render template",http.StatusInternalServerError)
		return
	}
}