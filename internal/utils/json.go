package utils

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/unbeman/ya-prac-go-first-grade/internal/model"
)

func WriteJSONError(writer http.ResponseWriter, request *http.Request, err error, httpCode int) {
	errMsg := model.ErrorOutput{Message: err.Error(), HTTPStatusCode: httpCode}
	render.Render(writer, request, errMsg)
}
