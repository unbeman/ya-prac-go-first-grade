package model

import (
	"net/http"

	"github.com/go-chi/render"
)

type ErrorOutput struct {
	Message        string `json:"message,omitempty"`
	HTTPStatusCode int    `json:"-"`
}

func (e ErrorOutput) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}
