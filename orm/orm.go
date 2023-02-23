package orm

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	defaultLog "log"
	"os"
	"time"
)

const (
	AFFILIATE_MANAGER_DB      = "affiliate_manager_db"
	AFFILIATE_MANAGER_TEST_DB = "affiliate_manager_test_db"
	USER_TABLE                = "user_table"
	REFERRAL_TABLE            = "referral_table"
	BOOKING_DETAILS_TABLE     = "booking_details_table"
	AFFILIATE_DETAILS_TABLE   = "affiliate_details_table"
	BOOKING_SLOTS_TABLE       = "booking_slots_table"
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
		IP:           os.Getenv("PROD_DB_HOST"),
		Port:         os.Getenv("PROD_DB_PORT"),
		Username:     os.Getenv("PROD_DB_USERNAME"),
		Password:     os.Getenv("PROD_DB_PASS"),
		DatabaseName: AFFILIATE_MANAGER_DB,
	},
	AFFILIATE_MANAGER_TEST_DB: {
		IP:           os.Getenv("TEST_DB_HOST"),
		Port:         os.Getenv("TEST_DB_PORT"),
		Username:     os.Getenv("TEST_DB_USERNAME"),
		Password:     os.Getenv("TEST_DB_PASS"),
		DatabaseName: AFFILIATE_MANAGER_TEST_DB,
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
	ENV = "PROD"
)

func DbInstance(ctx echo.Context) *gorm.DB {
	if db == nil {
		switch ENV {
		case "PROD":
			if err := ConnectMySQL(ctx, AFFILIATE_MANAGER_DB); err != nil {
				log.Error(err)
			}
			break
		case "TEST":
			if err := ConnectMySQL(ctx, AFFILIATE_MANAGER_TEST_DB); err != nil {
				log.Error(err)
			}
			break
		}
	}
	return db
}

func ConnectMySQL(ctx echo.Context, dbName string) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
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
