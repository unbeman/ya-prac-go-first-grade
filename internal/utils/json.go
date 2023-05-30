package utils

import (
	"net/http"

	"github.com/unbeman/ya-prac-go-first-grade/internal/model"

	"github.com/go-chi/render"
)

func WriteJSONError(writer http.ResponseWriter, request *http.Request, err error, httpCode int) {
	errMsg := model.ErrorOutput{Message: err.Error(), HTTPStatusCode: httpCode}
	render.Render(writer, request, errMsg)
}
