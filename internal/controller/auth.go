package controller

import (
	"time"

	"github.com/unbeman/ya-prac-go-first-grade/internal/apperrors"
	"github.com/unbeman/ya-prac-go-first-grade/internal/config"
	"github.com/unbeman/ya-prac-go-first-grade/internal/database"
	"github.com/unbeman/ya-prac-go-first-grade/internal/model"
	"github.com/unbeman/ya-prac-go-first-grade/internal/utils"
)

type AuthController struct {
	db            *database.PG
	tokenLifeTime time.Duration
}

func GetAuthController(db *database.PG, cfg config.AuthConfig) *AuthController {
	return &AuthController{db: db, tokenLifeTime: cfg.TokenLifeTime}
}

func (c AuthController) CreateUser(userInput model.UserInput) (user *model.User, err error) {
	hashPassword, err := utils.HashPassword(userInput.Password)
	if err != nil {
		return
	}
	user = &model.User{
		Login:          userInput.Login,
		HashPassword:   hashPassword,
		CurrentBalance: 0,
		Withdrawn:      0,
		CreatedAt:      time.Now(),
	}
	err = c.db.CreateNewUser(user)
	return
}

func (c AuthController) GetUser(userInput model.UserInput) (user *model.User, err error) {
	user, err = c.db.GetUserByLogin(userInput.Login)
	if err != nil {
		return
	}
	err = utils.CheckPassword(user.HashPassword, userInput.Password)
	if err != nil {
		err = apperrors.ErrInvalidUserCredentials
		return
	}
	return
}

func (c AuthController) CreateSession(user *model.User) (session *model.Session, err error) {
	token := utils.GenerateToken()
	created := time.Now()
	expired := created.Add(c.tokenLifeTime)
	session = &model.Session{
		User:      *user,
		Token:     token,
		CreatedAt: created,
		ExpireAt:  expired,
	}
	err = c.db.CreateNewSession(session)
	return
}

func (c AuthController) GetUserByToken(token string) (user *model.User, err error) {
	return c.db.GetUserByToken(token)
}
