package controller

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/unbeman/ya-prac-go-first-grade/internal/app-errors"
	"github.com/unbeman/ya-prac-go-first-grade/internal/model"
	"github.com/unbeman/ya-prac-go-first-grade/internal/utils"
)

type AuthController struct {
	db *model.PG
}

func GetAuthController(db *model.PG) *AuthController {
	return &AuthController{db: db}
}

func (c AuthController) CreateUser(userInput model.UserInput) (user model.User, err error) {
	hashPassword, err := utils.HashPassword(userInput.Password)
	if err != nil {
		return
	}

	user = model.User{
		Login:          userInput.Login,
		HashPassword:   hashPassword,
		CurrentBalance: 0,
		Withdrawn:      0,
		CreatedAt:      time.Now(),
	}

	result := c.db.Conn.Create(&user)
	if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		err = fmt.Errorf("%w: user with login (%v)", app_errors.ErrAlreadyExists, user.Login)
		return
	} else if result.Error != nil {
		err = fmt.Errorf("%w: %v", app_errors.ErrDb, result.Error)
		return
	}
	return
}

func (c AuthController) GetUser(userInput model.UserInput) (user model.User, err error) {
	result := c.db.Conn.First(&user, "login = ?", userInput.Login)
	if result.Error != nil {
		return user, app_errors.ErrInvalidUserCredentials
	}
	err = utils.CheckPassword(user.HashPassword, userInput.Password)
	if err != nil {
		err = app_errors.ErrInvalidUserCredentials
		return
	}
	return
}

func (c AuthController) CreateSession(user model.User) (session model.Session, err error) {
	token := utils.GenerateToken()
	created := time.Now()
	expired := created.Add(1 * time.Hour)
	session = model.Session{
		User:      user,
		Token:     token,
		CreatedAt: created,
		ExpireAt:  expired,
	}
	result := c.db.Conn.Create(&session)
	if result.Error != nil {
		err = fmt.Errorf("%w: %v", app_errors.ErrDb, result.Error)
		return
	}
	return
}

func (c AuthController) GetUserByToken(token string) (user model.User, err error) {
	result := c.db.Conn.Joins(
		"JOIN sessions ON users.id = sessions.user_id where token = ? AND expire_at > ?",
		token,
		time.Now(),
	).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		err = app_errors.ErrInvalidToken
		return
	}
	if result.Error != nil {
		err = fmt.Errorf("%w: %v", app_errors.ErrDb, result.Error)
		return
	}
	log.Debug(user)
	return
}
