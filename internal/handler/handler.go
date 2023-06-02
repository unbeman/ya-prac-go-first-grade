package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	logger "github.com/chi-middleware/logrus-logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"

	"github.com/unbeman/ya-prac-go-first-grade/internal/apperrors"
	"github.com/unbeman/ya-prac-go-first-grade/internal/controller"
	"github.com/unbeman/ya-prac-go-first-grade/internal/model"
	"github.com/unbeman/ya-prac-go-first-grade/internal/utils"
)

type AppHandler struct {
	*chi.Mux
	authControl   *controller.AuthController
	pointsControl *controller.PointsController
}

func GetAppHandler(authControl *controller.AuthController, pointsControl *controller.PointsController) *AppHandler {
	h := &AppHandler{
		Mux:           chi.NewMux(),
		authControl:   authControl,
		pointsControl: pointsControl,
	}
	h.Use(middleware.RequestID)
	h.Use(middleware.RealIP)
	h.Use(logger.Logger("router", log.New()))
	h.Use(middleware.Recoverer)

	h.Route("/api/user", func(router chi.Router) {

		router.Post("/register", h.Register)
		router.Post("/login", h.Login)

		router.Group(func(ra chi.Router) {
			ra.Use(h.authorized)
			ra.Post("/orders", h.AddOrder)
			ra.Get("/withdrawals", h.GetWithdrawals)
			ra.Group(func(ru chi.Router) {
				ru.Use(h.updOrdersInfo)
				ru.Get("/orders", h.GetOrders)
				ru.Get("/balance", h.GetBalance)
				ru.Post("/balance/withdraw", h.WithdrawPoints)
			})
		})
	})
	return h
}

func (h AppHandler) Register(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	writer.Header().Set("Content-Type", "application/json")

	var userInfo model.UserInput
	err := render.DecodeJSON(request.Body, &userInfo) //todo: use another decoder
	if err != nil {
		utils.WriteJSONError(writer, request, err, http.StatusBadRequest)
		return
	}

	user, err := h.authControl.CreateUser(ctx, userInfo)
	if errors.Is(err, apperrors.ErrAlreadyExists) {
		utils.WriteJSONError(writer, request, err, http.StatusConflict)
		return
	}
	if err != nil {
		utils.WriteJSONError(writer, request, err, http.StatusInternalServerError)
		return
	}
	session, err := h.authControl.CreateSession(ctx, user)
	if err != nil {
		utils.WriteJSONError(writer, request, err, http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Authorization", session.Token)
	writer.WriteHeader(http.StatusOK)
}

func (h AppHandler) Login(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	writer.Header().Set("Content-Type", "application/json")

	var userInfo model.UserInput
	err := render.DecodeJSON(request.Body, &userInfo) //todo: use another decoder
	if err != nil {
		utils.WriteJSONError(writer, request, err, http.StatusBadRequest)
		return
	}

	user, err := h.authControl.GetUser(ctx, userInfo)
	if errors.Is(err, apperrors.ErrInvalidUserCredentials) {
		utils.WriteJSONError(writer, request, err, http.StatusUnauthorized)
		return
	}
	if err != nil {
		utils.WriteJSONError(writer, request, err, http.StatusInternalServerError)
		return
	}
	log.Debug(user)

	session, err := h.authControl.CreateSession(ctx, user)
	if err != nil {
		utils.WriteJSONError(writer, request, err, http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Authorization", session.Token) //todo: wrap
	writer.WriteHeader(http.StatusOK)
}
func (h AppHandler) AddOrder(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	writer.Header().Set("Content-Type", "text/plain")

	user := h.getUserFromContext(request.Context())

	orderNumber, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	isNewOrder, err := h.pointsControl.AddUserOrder(ctx, user, string(orderNumber))
	if errors.Is(err, apperrors.ErrInvalidOrderNumberFormat) {
		http.Error(writer, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	if errors.Is(err, apperrors.ErrAlreadyExists) {
		http.Error(writer, err.Error(), http.StatusConflict)
	}
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	if isNewOrder {
		writer.WriteHeader(http.StatusAccepted)
		return
	}
	writer.WriteHeader(http.StatusOK)
}

func (h AppHandler) GetOrders(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	writer.Header().Set("Content-Type", "application/json")

	user := h.getUserFromContext(request.Context())
	orders, err := h.pointsControl.GetUserOrders(ctx, user)
	if errors.Is(err, apperrors.ErrNoRecords) {
		writer.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil {
		utils.WriteJSONError(writer, request, err, http.StatusInternalServerError)
		return
	}
	jsonOrders, err := json.Marshal(orders)
	if err != nil {
		utils.WriteJSONError(writer, request, err, http.StatusInternalServerError)
		return
	}
	log.Info("GetOrders result", string(jsonOrders))
	writer.Write(jsonOrders)
	writer.WriteHeader(http.StatusOK)

}

func (h AppHandler) GetBalance(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	writer.Header().Set("Content-Type", "application/json")

	user := h.getUserFromContext(request.Context())
	userBalance, err := h.pointsControl.GetUserBalance(ctx, user)
	if err != nil {
		utils.WriteJSONError(writer, request, err, http.StatusInternalServerError)
		return
	}
	jsonUserBalance, err := json.Marshal(userBalance)
	if err != nil {
		utils.WriteJSONError(writer, request, err, http.StatusInternalServerError)
		return
	}
	writer.Write(jsonUserBalance)
	writer.WriteHeader(http.StatusOK)
}

func (h AppHandler) WithdrawPoints(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	writer.Header().Set("Content-Type", "application/json")

	user := h.getUserFromContext(request.Context())

	var withdrawInfo model.WithdrawnInput
	err := render.DecodeJSON(request.Body, &withdrawInfo) //todo: use another decoder
	if err != nil {
		utils.WriteJSONError(writer, request, err, http.StatusBadRequest)
		return
	}

	err = h.pointsControl.CreateWithdraw(ctx, user, withdrawInfo)
	if errors.Is(err, apperrors.ErrNotEnoughPoints) {
		utils.WriteJSONError(writer, request, err, http.StatusPaymentRequired)
		return
	}
	if errors.Is(err, apperrors.ErrInvalidOrderNumberFormat) {
		utils.WriteJSONError(writer, request, err, http.StatusUnprocessableEntity)
		return
	}
	if err != nil {
		utils.WriteJSONError(writer, request, err, http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)

}

func (h AppHandler) GetWithdrawals(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	writer.Header().Set("Content-Type", "application/json")

	user := h.getUserFromContext(request.Context())

	withdrawals, err := h.pointsControl.GetUserWithdrawals(ctx, user)
	if err != nil {
		utils.WriteJSONError(writer, request, err, http.StatusInternalServerError)
		return
	}
	if len(withdrawals) == 0 { //todo: wrap
		writer.WriteHeader(http.StatusNoContent)
		return
	}
	jsonWithdrawals, err := json.Marshal(withdrawals)
	if err != nil {
		utils.WriteJSONError(writer, request, err, http.StatusInternalServerError)
		return
	}
	writer.Write(jsonWithdrawals)
	writer.WriteHeader(http.StatusOK)

}

func (h AppHandler) getUserFromContext(ctx context.Context) *model.User {
	return ctx.Value(UserContextKey).(*model.User) //todo: check context not nil and get value is ok
}
