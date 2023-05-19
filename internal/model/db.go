package model

import (
	"fmt"
	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PG struct {
	Conn *gorm.DB
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
	db.Conn = conn
	return nil
}

func (db *PG) migrate() error {
	tx := db.Conn.Exec(fmt.Sprintf(`
	DO $$ BEGIN
		CREATE TYPE order_status AS ENUM ('%v', '%v', '%v', '%v');
	EXCEPTION
		WHEN duplicate_object THEN null;
	END $$;`, StatusNew, StatusProcessing, StatusInvalid, StatusProcessed))
	if tx.Error != nil {
		return tx.Error
	}
	return db.Conn.AutoMigrate(
		&User{},
		&Session{},
		&Order{},
		&Withdrawal{},
	)
}

func (db PG) Ping() bool {
	return true
}
