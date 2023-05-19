package handler

import (
	"context"
	"net/http"
)

func (h AppHandler) authorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		ctx := context.WithValue(request.Context(), "user_id", 10) //todo: rewrite
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}
