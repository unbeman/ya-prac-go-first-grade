package handler

import (
	"context"
	"errors"
	"net/http"

	errors2 "github.com/unbeman/ya-prac-go-first-grade/internal/app-errors"
)

const UserContextKey = "user"

func (h AppHandler) authorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		inputToken := request.Header.Get("Authorization")
		user, err := h.authControl.GetUserByToken(inputToken)
		if errors.Is(err, errors2.ErrInvalidToken) {
			http.Error(writer, err.Error(), http.StatusUnauthorized)
			return
		}
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		ctx := context.WithValue(request.Context(), UserContextKey, user)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}
