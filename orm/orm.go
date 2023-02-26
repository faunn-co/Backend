package orm

import (
	"context"
	"fmt"
	"github.com/aaronangxz/AffiliateManager/logger"
	"github.com/aaronangxz/AffiliateManager/root"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	defaultLog "log"
	"os"
	"time"
)

const (
	AFFILIATE_MANAGER_DB      = "affiliate_manager_db"
	AFFILIATE_MANAGER_TEST_DB = "affiliate_manager_test_db"
	USER_TABLE                = "user_table"
	USER_AUTH_TABLE           = "user_auth_table"
	REFERRAL_TABLE            = "referral_table"
	BOOKING_DETAILS_TABLE     = "booking_details_table"
	AFFILIATE_DETAILS_TABLE   = "affiliate_details_table"
	BOOKING_SLOTS_TABLE       = "booking_slots_table"
)

var (
	db        *gorm.DB
	newLogger = gormLogger.New(
		defaultLog.New(os.Stdout, "\r\n", defaultLog.LstdFlags), // io writer
		gormLogger.Config{
			SlowThreshold:             time.Second,     // Slow SQL threshold
			LogLevel:                  gormLogger.Warn, // Log level
			IgnoreRecordNotFoundError: false,           // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,            // Disable color
		},
	)
	ENV         = "PROD"
	DB_HOST     = ""
	DB_PORT     = ""
	DB_USERNAME = ""
	DB_PASS     = ""
	DB_NAME     = ""
)

func getEnvDir() string {
	switch ENV {
	case "PROD":
		curDir, wdErr := os.Getwd()
		if wdErr != nil {
			log.Error(wdErr)
		}
		return curDir + "/.env"
	case "TEST":
		return root.Root + "/.env"
	}
	return ""
}

func DbInstance(ctx context.Context) *gorm.DB {
	if db == nil {
		err := godotenv.Load(getEnvDir())
		if err != nil {
			logger.Warn(ctx, "Error loading .env file")
		}
		ENV = os.Getenv("ENV")
		switch ENV {
		case "PROD":
			DB_HOST = os.Getenv("PROD_DB_HOST")
			DB_PORT = os.Getenv("PROD_DB_PORT")
			DB_USERNAME = os.Getenv("PROD_DB_USERNAME")
			DB_PASS = os.Getenv("PROD_DB_PASS")
			DB_NAME = AFFILIATE_MANAGER_DB
			logger.Info(ctx, "Connecting to PROD DB")
			break
		case "TEST":
			fallthrough
		case "STAGING":
			DB_HOST = os.Getenv("TEST_DB_HOST")
			DB_PORT = os.Getenv("TEST_DB_PORT")
			DB_USERNAME = os.Getenv("TEST_DB_USERNAME")
			DB_PASS = os.Getenv("TEST_DB_PASS")
			DB_NAME = AFFILIATE_MANAGER_TEST_DB
			logger.Info(ctx, "Connecting to TEST DB")
			break
		case "LOCAL":
			DB_HOST = "127.0.0.1"
			DB_PORT = "3306"
			DB_USERNAME = "root"
			DB_PASS = "Xuanze94"
			DB_NAME = AFFILIATE_MANAGER_DB
			logger.Info(ctx, "Connecting to LOCAL DB")
			break
		}
		if err := ConnectMySQL(ctx); err != nil {
			logger.Error(ctx, err)
		}
	}
	return db
}

func ConnectMySQL(ctx context.Context) error {
	URL := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", DB_USERNAME, DB_PASS, DB_HOST, DB_PORT, DB_NAME)
	openDb, err := gorm.Open(mysql.Open(URL), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		return err
	}

	db = openDb
	return nil
}
