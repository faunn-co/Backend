package orm

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	defaultLog "log"
	"os"
	"time"
)

const (
	AFFILIATE_MANAGER_DB = "affiliate_manager_db"
	BOOKING_SLOTS_TABLE  = "booking_slots_table"
)

type SQLConnection struct {
	IP           string
	Port         string
	Username     string
	Password     string
	DatabaseName string
}

var sqlConnections = map[string]SQLConnection{
	AFFILIATE_MANAGER_DB: {
		IP:           "127.0.0.1",
		Port:         "3306",
		Username:     "root",
		Password:     "Xuanze94",
		DatabaseName: AFFILIATE_MANAGER_DB,
	},
}

var (
	db        *gorm.DB
	newLogger = logger.New(
		defaultLog.New(os.Stdout, "\r\n", defaultLog.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Warn, // Log level
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)
)

func DbInstance(ctx echo.Context) *gorm.DB {
	if db == nil {
		if err := ConnectMySQL(ctx, AFFILIATE_MANAGER_DB); err != nil {
		}
	}
	return db
}

func ConnectMySQL(ctx echo.Context, dbName string) error {
	var table SQLConnection
	var exists bool
	if table, exists = sqlConnections[dbName]; !exists {
		return errors.New(fmt.Sprintf("db not supported | %v", table))
	}
	URL := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", table.Username, table.Password, table.IP, table.Port, table.DatabaseName)

	openDb, err := gorm.Open(mysql.Open(URL), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		return err
	}

	db = openDb

	return nil
}
