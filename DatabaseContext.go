package IrisAPIs

import (
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"github.com/xormplus/xorm"
	"os"
)

type DatabaseContext struct {
	init             bool
	DbObject         *xorm.Engine
	ConnectionString string
	ShowSql          bool
}

func NewDatabaseContext(connectionString string, showSql bool) (*DatabaseContext, error) {
	engine, err := initDatabaseContext(connectionString, showSql)
	if err != nil {
		return nil, err
	}
	return &DatabaseContext{
		init:             true,
		DbObject:         engine,
		ConnectionString: connectionString,
	}, nil
}

// This is for test cases. It will try to fetch DB test connect strings with these order :
// 1. Environment Parameters
// 2. gtest pass argument
// 3. config file
func NewTestDatabaseContext() (*DatabaseContext, error) {
	log.SetLevel(log.DebugLevel)
	var connStr string
	log.Debug("Trying initializing Test DB with environment \"TEST_DB_CONN_STR\"...")
	connStr, exist := os.LookupEnv("TEST_DB_CONN_STR")
	if !exist || connStr == "" {
		log.Debug("Trying initializing Test DB with parameter \"test_db_conn_str\"...")
	} else {
		log.Debug("Initialized Test DB from environment")
	}

	//Fetch configuration file. It usually only exists in local test environment
	if connStr == "" {
		log.Debug("Trying initializing Test DB with configuration file...")
		connStr = NewConfiguration().TestConnectionString
	} else {
		log.Debug("Initialized DB from parameter")
	}

	if connStr == "" {
		log.Warn("Fail to initialize test database form any of source")
		return nil, nil
	} else {
		log.Debug("Initialized DB from configuration file")
	}

	return NewDatabaseContext(connStr, true)
}

func initDatabaseContext(connectionString string, showSql bool) (engine *xorm.Engine, err error) {
	engine, err = xorm.NewEngine("mysql", connectionString)

	if err != nil {
		return nil, err
	}

	engine.ShowSQL(showSql)
	return engine, nil
}
