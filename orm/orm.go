package orm

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	defaultLog "log"
	"os"
	"time"
)

const (
	DATA_HELPER_V2_DB              = "data_helper_v2_db"
	DATA_INJECTION_AUDIT_LOG_TABLE = "data_injection_audit_log_table"
	PROJECT_TABLE                  = "project_table"
	ACCOUNT_LIST_TABLE             = "account_list_table"
)

type SQLConnection struct {
	IP                     string
	Port                   string
	Username               string
	Password               string
	DatabaseName           string
	DataTypeColumnPosition int
}

var sqlConnections = map[string]SQLConnection{
	"tts": {
		IP:                     "10.231.16.14",
		Port:                   "9967",
		Username:               "i_boe_ecom",
		Password:               "i_boe_ecom123",
		DatabaseName:           "i_boe_ecom",
		DataTypeColumnPosition: 2,
	},
	"fans": {
		IP:                     "10.231.14.130",
		Port:                   "3306",
		Username:               "ecom_fans_dw_w",
		Password:               "dwpmYISCgUzk0JF_mSXl3XTkxi9Mgtzs",
		DatabaseName:           "ecom_fans_dw",
		DataTypeColumnPosition: 1,
	},
	"doris": {
		IP:                     "10.231.31.1",
		Port:                   "9030",
		Username:               "root",
		Password:               "",
		DatabaseName:           "ecom",
		DataTypeColumnPosition: 1,
	},
	"navigator": {
		IP:                     "10.231.17.195",
		Port:                   "3306",
		Username:               "navigator_boe_w",
		Password:               "Sb3ZMVVYv97lNsD_c9lDEtjHbRMVbWWA",
		DatabaseName:           "navigator_boe_mock",
		DataTypeColumnPosition: 1,
	},
	"invoker": {
		IP:                     "10.231.20.70",
		Port:                   "3306",
		Username:               "invoker_mock__w",
		Password:               "mNKQ7absI9FViKa_fLDREGhrzAVu0FIj",
		DatabaseName:           "invoker_mock_table",
		DataTypeColumnPosition: 1,
	},
	"DataHelper": {
		IP:                     "10.231.117.196",
		Port:                   "3306",
		Username:               "data_helper_v_w",
		Password:               "iKi91jOKGubpac0_Z8bHD9onQaqnk4oM",
		DatabaseName:           "data_helper_v2_db",
		DataTypeColumnPosition: 1,
	},
}

var (
	db           *gorm.DB
	dataHelperDb *gorm.DB
	newLogger    = logger.New(
		defaultLog.New(os.Stdout, "\r\n", defaultLog.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Warn, // Log level
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)
)

func DbInstance(ctx context.Context) *gorm.DB {
	if db == nil {
		if err := ConnectMySQL(ctx, "doris"); err != nil {
			hlog.CtxErrorf(ctx, "Error while establishing doris DB Connection: %v", err.Error())
		}
	}
	return db
}

func DataHelperDbInstance(ctx context.Context) *gorm.DB {
	if dataHelperDb == nil {
		if err := ConnectMySQL(ctx, "DataHelper"); err != nil {
			hlog.CtxErrorf(ctx, "Error while establishing DataHelper DB Connection: %v", err.Error())
		}
	}
	return dataHelperDb
}

func ConnectMySQL(ctx context.Context, dbName string) error {
	var table SQLConnection
	var exists bool
	if table, exists = sqlConnections[dbName]; !exists {
		return errors.New(fmt.Sprintf("db not supported | %v", table))
	}
	URL := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", table.Username, table.Password, table.IP, table.Port, table.DatabaseName)

	hlog.CtxInfof(ctx, "Connecting to %v", URL)
	openDb, err := gorm.Open(mysql.Open(URL), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		hlog.CtxErrorf(ctx, err.Error())
		return err
	}
	hlog.CtxInfof(ctx, "Database connection established | %v", table.DatabaseName)

	if dbName == "DataHelper" {
		dataHelperDb = openDb
	} else {
		db = openDb
	}
	return nil
}
