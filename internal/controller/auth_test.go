package controller

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/unbeman/ya-prac-go-first-grade/internal/apperrors"
	"github.com/unbeman/ya-prac-go-first-grade/internal/config"
	mock_database "github.com/unbeman/ya-prac-go-first-grade/internal/database/mock"
	"github.com/unbeman/ya-prac-go-first-grade/internal/model"
	"github.com/unbeman/ya-prac-go-first-grade/internal/utils"
)

func randomUser(t *testing.T) (user model.User, password string) {
	password = "12345"
	hashedPassword, err := utils.HashPassword(password)
	require.NoError(t, err)

	user = model.User{
		Login:        "Vasya",
		HashPassword: hashedPassword,
	}
	return
}

func TestAuthController_CreateUser(t *testing.T) {
	user, password := randomUser(t)
	type args struct {
		ctx       context.Context
		userInput model.UserInput
	}
	ctx := context.TODO()
	tests := []struct {
		name    string
		args    args
		setup   func(*mock_database.MockDatabase)
		wantErr bool
	}{
		{
			name: "register new user",
			args: args{
				ctx:       ctx,
				userInput: model.UserInput{Login: user.Login, Password: password}},
			setup: func(mDB *mock_database.MockDatabase) {
				mDB.EXPECT().
					CreateNewUser(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "register existed user",
			args: args{
				ctx:       ctx,
				userInput: model.UserInput{Login: user.Login, Password: password}},
			setup: func(mDB *mock_database.MockDatabase) {
				mDB.EXPECT().
					CreateNewUser(gomock.Any(), gomock.Any()).
					Return(apperrors.ErrAlreadyExists)
			},
			wantErr: true,
		},
		{
			name: "db error",
			args: args{
				ctx:       ctx,
				userInput: model.UserInput{Login: user.Login, Password: password}},
			setup: func(mDB *mock_database.MockDatabase) {
				mDB.EXPECT().
					CreateNewUser(gomock.Any(), gomock.Any()).
					Return(apperrors.ErrDB)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockDatabase := mock_database.NewMockDatabase(ctrl)
			tt.setup(mockDatabase)
			c := AuthController{
				db:            mockDatabase,
				tokenLifeTime: config.TokenLifeTimeDefault,
			}
			_, err := c.CreateUser(context.TODO(), tt.args.userInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestAuthController_GetUser(t *testing.T) {
	user, password := randomUser(t)
	type args struct {
		ctx       context.Context
		userInput model.UserInput
	}
	ctx := context.TODO()
	tests := []struct {
		name    string
		args    args
		setup   func(*mock_database.MockDatabase)
		wantErr bool
	}{
		{
			name: "get user",
			args: args{
				ctx:       ctx,
				userInput: model.UserInput{Login: user.Login, Password: password}},
			setup: func(mDB *mock_database.MockDatabase) {
				mDB.EXPECT().
					GetUserByLogin(gomock.Any(), gomock.Any()).
					Return(&user, nil)
			},
			wantErr: false,
		},
		{
			name: "invalid login",
			args: args{
				ctx:       ctx,
				userInput: model.UserInput{Login: user.Login, Password: password}},
			setup: func(mDB *mock_database.MockDatabase) {
				mDB.EXPECT().
					GetUserByLogin(gomock.Any(), gomock.Any()).
					Return(nil, apperrors.ErrInvalidUserCredentials)
			},
			wantErr: true,
		},
		{
			name: "invalid password",
			args: args{
				ctx:       ctx,
				userInput: model.UserInput{Login: user.Login, Password: password}},
			setup: func(mDB *mock_database.MockDatabase) {
				mDB.EXPECT().
					GetUserByLogin(gomock.Any(), gomock.Any()).
					Return(nil, apperrors.ErrInvalidUserCredentials)
			},
			wantErr: true,
		},
		{
			name: "db error",
			args: args{
				ctx:       ctx,
				userInput: model.UserInput{Login: user.Login, Password: password}},
			setup: func(mDB *mock_database.MockDatabase) {
				mDB.EXPECT().
					GetUserByLogin(gomock.Any(), gomock.Any()).
					Return(nil, apperrors.ErrDB)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockDatabase := mock_database.NewMockDatabase(ctrl)
			tt.setup(mockDatabase)
			c := AuthController{
				db:            mockDatabase,
				tokenLifeTime: config.TokenLifeTimeDefault,
			}
			_, err := c.GetUser(context.TODO(), tt.args.userInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestAuthController_CreateSession(t *testing.T) {
	user, _ := randomUser(t)
	type args struct {
		ctx  context.Context
		user model.User
	}
	ctx := context.TODO()
	tests := []struct {
		name    string
		args    args
		setup   func(*mock_database.MockDatabase)
		wantErr bool
	}{
		{
			name: "get user",
			args: args{
				ctx:  ctx,
				user: user,
			},
			setup: func(mDB *mock_database.MockDatabase) {
				mDB.EXPECT().
					CreateNewSession(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "db error",
			args: args{
				ctx:  ctx,
				user: user,
			},
			setup: func(mDB *mock_database.MockDatabase) {
				mDB.EXPECT().
					CreateNewSession(gomock.Any(), gomock.Any()).
					Return(apperrors.ErrDB)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockDatabase := mock_database.NewMockDatabase(ctrl)
			tt.setup(mockDatabase)
			c := AuthController{
				db:            mockDatabase,
				tokenLifeTime: config.TokenLifeTimeDefault,
			}
			_, err := c.CreateSession(context.TODO(), &tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
