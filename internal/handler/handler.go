package handler

import (
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	errors2 "github.com/unbeman/ya-prac-go-first-grade/internal/app-errors"
	"github.com/unbeman/ya-prac-go-first-grade/internal/controller"
	"github.com/unbeman/ya-prac-go-first-grade/internal/model"
	"github.com/unbeman/ya-prac-go-first-grade/internal/utils"

	logger "github.com/chi-middleware/logrus-logger"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
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
	h.Get("/ping", h.Ping())
	h.Route("/api/user", func(router chi.Router) {
		router.Post("/register", h.Register())
		router.Post("/login", h.Login())
		h.Group(func(r chi.Router) {
			//h.Use(h.authorized)
			router.Post("/orders", h.AddOrder())
			router.Get("/orders", h.GetOrders())
			router.Get("/balance", h.GetBalance())
			router.Post("/balance/withdraw", h.WithdrawPoints())
			router.Get("/withdrawals", h.GetWithdrawals())
		})

	})
	return h
}

func (h AppHandler) Ping() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if h.pointsControl.Ping() {
			writer.WriteHeader(http.StatusOK)
		} else {
			writer.WriteHeader(http.StatusInternalServerError)
		}

	}
}

func (h AppHandler) Register() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")

		var userInfo model.UserInput
		err := render.DecodeJSON(request.Body, &userInfo) //todo: use another decoder
		if err != nil {
			utils.WriteJsonError(writer, request, err, http.StatusBadRequest)
			return
		}

		user, err := h.authControl.CreateUser(userInfo)
		if errors.Is(err, errors2.ErrAlreadyExists) {
			utils.WriteJsonError(writer, request, err, http.StatusConflict)
			return
		}
		if err != nil {
			utils.WriteJsonError(writer, request, err, http.StatusInternalServerError)
			return
		}
		session, err := h.authControl.CreateSession(user)
		if err != nil {
			utils.WriteJsonError(writer, request, err, http.StatusInternalServerError)
			return
		}
		writer.Header().Set("Authorization", session.Token)
		writer.WriteHeader(http.StatusOK)
	}
}

func (h AppHandler) Login() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")

		var userInfo model.UserInput
		err := render.DecodeJSON(request.Body, &userInfo) //todo: use another decoder
		if err != nil {
			utils.WriteJsonError(writer, request, err, http.StatusBadRequest)
			return
		}

		user, err := h.authControl.GetUser(userInfo)
		if errors.Is(err, errors2.ErrInvalidUserCredentials) {
			utils.WriteJsonError(writer, request, err, http.StatusUnauthorized)
			return
		}
		if err != nil {
			utils.WriteJsonError(writer, request, err, http.StatusInternalServerError)
			return
		}
		log.Debug(user)

		session, err := h.authControl.CreateSession(user)
		if err != nil {
			utils.WriteJsonError(writer, request, err, http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Authorization", session.Token) //todo: wrap
		writer.WriteHeader(http.StatusOK)
	}
}
func (h AppHandler) AddOrder() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/plain")
		contentType := request.Header.Get("Content-Type") //todo: wrap
		if contentType != "text/plain" {
			http.Error(writer, errors2.ErrInvalidContentType.Error(), http.StatusBadRequest) // http.StatusUnsupportedMediaType
			return
		}
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

		orderNumber, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		isNewOrder, err := h.pointsControl.AddUserOrder(user, string(orderNumber))
		if errors.Is(err, errors2.ErrInvalidOrderNumberFormat) {
			http.Error(writer, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		if errors.Is(err, errors2.ErrAlreadyExists) {
			http.Error(writer, err.Error(), http.StatusConflict)
		}
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		//todo: prolong session
		if isNewOrder {
			writer.WriteHeader(http.StatusAccepted)
			return
		}
		writer.WriteHeader(http.StatusOK)
	}
}

func (h AppHandler) GetOrders() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

	}
}

func (h AppHandler) GetBalance() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

	}
}

func (h AppHandler) WithdrawPoints() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

	}
}

func (h AppHandler) GetWithdrawals() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

	}
}
