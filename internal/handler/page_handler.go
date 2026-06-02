package handler

import (
	"net/http"
	"voidsounds/internal/components"
)

type PageHandler struct{}

func NewPageHandler() *PageHandler {
	return &PageHandler{}
}

func (h *PageHandler) Artists(w http.ResponseWriter, r *http.Request) {
	components.ArtistsPage().Render(r.Context(), w)
}

func (h *PageHandler) ForOrganizers(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") == "true" {
		components.ForOrganizersContent().Render(r.Context(), w)
	} else {
		components.ForOrganizersPage().Render(r.Context(), w)
	}
}
