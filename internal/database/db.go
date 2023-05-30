package database

import (
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"

	log "github.com/sirupsen/logrus"
	errors2 "github.com/unbeman/ya-prac-go-first-grade/internal/apperrors"
	"github.com/unbeman/ya-prac-go-first-grade/internal/model"
)

type PG struct {
	conn *gorm.DB
}

func GetDatabase(dsn string) (*PG, error) {
	db := &PG{}
	if err := db.connect(dsn); err != nil {
		return nil, err
	}
	if err := db.migrate(); err != nil {
		return nil, err
	}
	return db, nil
}

func (db *PG) connect(dsn string) error {
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)}) //todo: use custom logger based on logrus
	if err != nil {
		return err
	}
	db.conn = conn
	return nil
}

func (db *PG) migrate() error {
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

func (db *PG) Ping() bool {
	pg, err := db.conn.DB()
	if err != nil {
		log.Error(err)
	}
	err = pg.Ping()
	if err != nil {
		log.Error(err)
		return false
	}
	return true
}

func (db *PG) CreateNewUser(user *model.User) error {
	result := db.conn.Create(user)
	if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		return fmt.Errorf("%w: user with login (%v)", errors2.ErrAlreadyExists, user.Login)
	}
	if result.Error != nil {
		return fmt.Errorf("%w: %v", errors2.ErrDB, result.Error)
	}
	return nil
}

func (db *PG) GetUserByLogin(login string) (user *model.User, err error) {
	user = &model.User{}
	result := db.conn.First(user, "login = ?", login)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors2.ErrInvalidUserCredentials
		return
	}
	if result.Error != nil {
		err = fmt.Errorf("%w: %v", errors2.ErrDB, result.Error)
		return
	}
	return
}

func (db *PG) CreateNewSession(session *model.Session) error {
	result := db.conn.Create(session)
	if result.Error != nil {
		return fmt.Errorf("%w: %v", errors2.ErrDB, result.Error)
	}
	return nil
}

func (db *PG) GetUserByToken(token string) (user *model.User, err error) {
	user = &model.User{}
	result := db.conn.Joins(
		"JOIN sessions ON users.id = sessions.user_id where token = ? AND expire_at > ?",
		token,
		time.Now(),
	).First(user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		err = errors2.ErrInvalidToken
		return
	}
	if result.Error != nil {
		err = fmt.Errorf("%w: %v", errors2.ErrDB, result.Error)
		return
	}
	return
}

func (db *PG) GetOrderByNumber(number string) (existingOrder *model.Order, err error) {
	existingOrder = &model.Order{}
	result := db.conn.First(existingOrder, "number = ?", number)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		err = errors2.ErrNoRecords
		return
	}
	if result.Error != nil {
		err = fmt.Errorf("%w: %v", errors2.ErrDB, result.Error)
		return
	}
	return
}

func (db *PG) CreateNewUserOrder(userID uint, number string) error {
	newOrder := &model.Order{UserID: userID, Status: model.StatusNew, Number: number}
	result := db.conn.Create(newOrder)
	if result.Error != nil {
		return fmt.Errorf("%w: %v", errors2.ErrDB, result.Error)
	}
	return nil
}

func (db *PG) UpdateUserBalanceAndOrder(order *model.Order) error {
	err := db.conn.Transaction(func(tx *gorm.DB) (txErr error) {
		result := tx.Save(order)
		if result.Error != nil {
			return result.Error
		}
		if order.Status == model.StatusProcessed {
			user := &model.User{}
			user.ID = order.UserID
			result = tx.Model(user).Update("current_balance", gorm.Expr("current_balance + ?", order.Accrual))
			if result.Error != nil {
				return result.Error
			}
		}
		return
	})

	if err != nil {
		return fmt.Errorf("%w: %v", errors2.ErrDB, err)
	}
	return nil
}

func (db *PG) GetUserOrders(userID uint) (orders []model.Order, err error) {
	result := db.conn.Find(&orders, "user_id = ?", userID).Order("created_at ASC")
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		err = errors2.ErrNoRecords
		return
	}
	if result.Error != nil {
		err = fmt.Errorf("%w: %v", errors2.ErrDB, result.Error)
		return
	}
	return
}

func (db *PG) CreateWithdraw(user *model.User, withdrawInfo model.WithdrawnInput) error {
	err := db.conn.Transaction(func(tx *gorm.DB) (txErr error) {
		result := tx.Where("id = ? and current_balance >= ?", user.ID, withdrawInfo.Sum).First(user)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors2.ErrNotEnoughPoints
		}
		user.CurrentBalance -= withdrawInfo.Sum
		user.Withdrawn += withdrawInfo.Sum
		if txErr = tx.Save(user).Error; txErr != nil {
			return
		}
		withdraw := model.Withdrawal{Order: withdrawInfo.OrderNumber, Sum: withdrawInfo.Sum, User: *user}
		if txErr = tx.Create(&withdraw).Error; txErr != nil {
			return
		}

		return
	})
	if errors.Is(err, errors2.ErrNotEnoughPoints) {
		return err
	}
	if err != nil {
		return fmt.Errorf("%w: %v", errors2.ErrDB, err)
	}
	return nil
}

func (db *PG) GetUserWithdrawals(userID uint) (withdrawals []model.Withdrawal, err error) {
	result := db.conn.Find(&withdrawals, "user_id = ?", userID).Order("created_at ASC")
	if result.Error != nil {
		err = fmt.Errorf("%w: %v", errors2.ErrDB, err)
		return
	}
	return
}
