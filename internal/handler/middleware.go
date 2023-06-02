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
		ctx := request.Context()
		inputToken := request.Header.Get("Authorization")
		user, err := h.authControl.GetUserByToken(ctx, inputToken)
		if errors.Is(err, apperrors.ErrInvalidToken) {
			http.Error(writer, err.Error(), http.StatusUnauthorized)
			return
		}
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		ctx2 := context.WithValue(ctx, UserContextKey, user)
		next.ServeHTTP(writer, request.WithContext(ctx2))
	})
}

func (h AppHandler) updOrdersInfo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		user := h.getUserFromContext(request.Context())
		err := h.pointsControl.UpdateUserOrders(ctx, user)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		next.ServeHTTP(writer, request)

	})
}
