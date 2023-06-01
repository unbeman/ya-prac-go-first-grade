package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/unbeman/ya-prac-go-first-grade/internal/apperrors"
)

type ContextKey string

const UserContextKey ContextKey = "user"

func (h AppHandler) authorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		inputToken := request.Header.Get("Authorization")
		user, err := h.authControl.GetUserByToken(inputToken)
		if errors.Is(err, apperrors.ErrInvalidToken) {
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

func (h AppHandler) updOrdersInfo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		user := h.getUserFromContext(request.Context())
		err := h.pointsControl.UpdateUserOrders(user)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		next.ServeHTTP(writer, request)

	})
}
