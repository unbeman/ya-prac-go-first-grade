package database

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/unbeman/ya-prac-go-first-grade/internal/apperrors"
	"github.com/unbeman/ya-prac-go-first-grade/internal/config"
	"github.com/unbeman/ya-prac-go-first-grade/internal/model"
)

type pg struct {
	conn *gorm.DB
}

func getPG(cfg config.DatabaseConfig) (*pg, error) {
	db := &pg{}
	if err := db.connect(cfg.DatabaseURI); err != nil {
		return nil, err
	}
	if err := db.migrate(); err != nil {
		return nil, err
	}
	return db, nil
}

func (db *pg) connect(dsn string) error {
	conn, err := gorm.Open(postgres.Open(dsn))
	//todo: use custom logger based on logrus
	if err != nil {
		return err
	}
	db.conn = conn
	return nil
}

func (db *pg) migrate() error {
	tx := db.conn.Exec(fmt.Sprintf(`
	DO $$ BEGIN
		CREATE TYPE order_status AS ENUM ('%v', '%v', '%v', '%v');
	EXCEPTION
		WHEN duplicate_object THEN null;
	END $$;`, model.StatusNew, model.StatusProcessing, model.StatusInvalid, model.StatusProcessed))
	if tx.Error != nil {
		return tx.Error
	}
	return db.conn.AutoMigrate(
		&model.User{},
		&model.Session{},
		&model.Order{},
		&model.Withdrawal{},
	)
}

func (db *pg) CreateNewUser(ctx context.Context, user *model.User) error {
	result := db.conn.WithContext(ctx).Create(user)
	if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		return fmt.Errorf("%w: user with login (%v)", apperrors.ErrAlreadyExists, user.Login)
	}
	if result.Error != nil {
		return fmt.Errorf("%w: %v", apperrors.ErrDB, result.Error)
	}
	return nil
}

func (db *pg) GetUserByLogin(ctx context.Context, login string) (user *model.User, err error) {
	user = &model.User{}
	result := db.conn.WithContext(ctx).First(user, "login = ?", login)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		err = apperrors.ErrInvalidUserCredentials
		return
	}
	if result.Error != nil {
		err = fmt.Errorf("%w: %v", apperrors.ErrDB, result.Error)
		return
	}
	return
}

func (db *pg) GetUserByID(ctx context.Context, userID uint) (user *model.User, err error) {
	user = &model.User{}
	result := db.conn.WithContext(ctx).First(user, userID)
	if result.Error != nil {
		err = fmt.Errorf("%w: %v", apperrors.ErrDB, result.Error)
		return
	}
	return
}

func (db *pg) CreateNewSession(ctx context.Context, session *model.Session) error {
	result := db.conn.WithContext(ctx).Create(session)
	if result.Error != nil {
		return fmt.Errorf("%w: %v", apperrors.ErrDB, result.Error)
	}
	return nil
}

func (db *pg) GetUserByToken(ctx context.Context, token string) (user *model.User, err error) {
	user = &model.User{}
	result := db.conn.WithContext(ctx).Joins(
		"JOIN sessions ON users.id = sessions.user_id where token = ? AND expire_at > ?",
		token,
		time.Now(),
	).First(user)
	err = result.Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = apperrors.ErrInvalidToken
		return
	}
	if err != nil {
		err = fmt.Errorf("%w: %v", apperrors.ErrDB, err)
		return
	}
	return
}

func (db *pg) GetOrderByNumber(ctx context.Context, number string) (existingOrder *model.Order, err error) {
	existingOrder = &model.Order{}
	result := db.conn.WithContext(ctx).First(existingOrder, "number = ?", number)
	err = result.Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = apperrors.ErrNoRecords
		return
	}
	if err != nil {
		err = fmt.Errorf("%w: %v", apperrors.ErrDB, err)
		return
	}
	return
}

func (db *pg) CreateNewUserOrder(ctx context.Context, userID uint, number string) error {
	newOrder := &model.Order{UserID: userID, Status: model.StatusNew, Number: number}
	result := db.conn.WithContext(ctx).Create(newOrder)
	if result.Error != nil {
		return fmt.Errorf("%w: %v", apperrors.ErrDB, result.Error)
	}
	return nil
}

func (db *pg) UpdateUserBalanceAndOrder(order *model.Order, accrualInfo model.OrderAccrualInfo) error {
	err := db.conn.Transaction(func(tx *gorm.DB) error {
		result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(order, order.ID)
		if result.Error != nil {
			return result.Error
		}
		if order.Status == accrualInfo.Status {
			log.Infof("nothing for update for order (%v)", order)
			return nil
		}
		order.Status = accrualInfo.Status
		order.Accrual = accrualInfo.Accrual
		result = tx.Save(order)
		if result.Error != nil {
			return result.Error
		}
		if order.Status == model.StatusProcessed {
			user := &model.User{}
			result = tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(user, order.UserID)
			if result.Error != nil {
				return nil
			}
			user.CurrentBalance += order.Accrual
			result = tx.Save(user)
			if result.Error != nil {
				return result.Error
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("%w: %v", apperrors.ErrDB, err)
	}
	return nil
}

func (db *pg) GetUserOrders(ctx context.Context, userID uint) (orders []model.Order, err error) {
	result := db.conn.WithContext(ctx).Find(&orders, "user_id = ?", userID).Order("created_at ASC")
	err = result.Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = apperrors.ErrNoRecords
		return
	}
	if err != nil {
		err = fmt.Errorf("%w: %v", apperrors.ErrDB, err)
		return
	}
	return
}

func (db *pg) GetNotReadyOrders(ctx context.Context, userID uint) (orders []model.Order, err error) {
	result := db.conn.WithContext(ctx).Find(&orders,
		"user_id = ? AND status != ? AND status != ?",
		userID,
		model.StatusProcessed,
		model.StatusInvalid,
	)
	err = result.Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = apperrors.ErrNoRecords
		return
	}
	if err != nil {
		err = fmt.Errorf("%w: %v", apperrors.ErrDB, err)
		return
	}
	return
}

func (db *pg) CreateWithdraw(ctx context.Context, user *model.User, withdrawInfo model.WithdrawnInput) error {
	err := db.conn.WithContext(ctx).Transaction(func(tx *gorm.DB) (txErr error) {
		result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("current_balance >= ?", withdrawInfo.Sum).
			First(user, user.ID)
		txErr = result.Error
		if errors.Is(txErr, gorm.ErrRecordNotFound) {
			return apperrors.ErrNotEnoughPoints
		}
		if txErr != nil {
			return txErr
		}
		user.CurrentBalance -= withdrawInfo.Sum
		user.Withdrawn += withdrawInfo.Sum
		if txErr = tx.Save(user).Error; txErr != nil {
			return
		}
		withdraw := model.Withdrawal{
			Order:  withdrawInfo.OrderNumber,
			Sum:    withdrawInfo.Sum,
			UserID: user.ID,
		}
		if txErr = tx.Create(&withdraw).Error; txErr != nil {
			return
		}
		return
	})
	if errors.Is(err, apperrors.ErrNotEnoughPoints) {
		return err
	}
	if err != nil {
		return fmt.Errorf("%w: %v", apperrors.ErrDB, err)
	}
	return nil
}

func (db *pg) GetUserWithdrawals(ctx context.Context, userID uint) (withdrawals []model.Withdrawal, err error) {
	result := db.conn.WithContext(ctx).Find(&withdrawals, "user_id = ?", userID).Order("created_at ASC")
	err = result.Error
	if err != nil {
		err = fmt.Errorf("%w: %v", apperrors.ErrDB, err)
		return
	}
	return
}
