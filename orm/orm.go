package orm

import (
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
	ENV         = "PROD"
	DB_HOST     = ""
	DB_PORT     = ""
	DB_USERNAME = ""
	DB_PASS     = ""
	DB_NAME     = ""
)

func DbInstance(ctx echo.Context) *gorm.DB {
	curDir, wdErr := os.Getwd()
	if wdErr != nil {
		log.Fatal(wdErr)
	}
	err := godotenv.Load(curDir + "/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	if db == nil {
		switch ENV {
		case "PROD":
			DB_HOST = os.Getenv("PROD_DB_HOST")
			DB_PORT = os.Getenv("PROD_DB_PORT")
			DB_USERNAME = os.Getenv("PROD_DB_USERNAME")
			DB_PASS = os.Getenv("PROD_DB_PASS")
			DB_NAME = AFFILIATE_MANAGER_DB
			if err := ConnectMySQL(ctx); err != nil {
				log.Error(err)
			}
			break
		case "TEST":
			DB_HOST = os.Getenv("TEST_DB_HOST")
			DB_PORT = os.Getenv("TEST_DB_PORT")
			DB_USERNAME = os.Getenv("TEST_DB_USERNAME")
			DB_PASS = os.Getenv("TEST_DB_PASS")
			DB_NAME = AFFILIATE_MANAGER_TEST_DB
			if err := ConnectMySQL(ctx); err != nil {
				log.Error(err)
			}
			break
		}
	}
	return db
}

func ConnectMySQL(ctx echo.Context) error {
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
