package IrisAPIs

import (
	"flag"
	_ "github.com/go-sql-driver/mysql"
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
	var connStr string
	connStr, exist := os.LookupEnv("TEST_DB_CONN_STR")
	if !exist || connStr == "" {
		connStr = *flag.String("db_conn_str", "", "test db password")
	}

	//Fetch configuration file. It usually only exists in local test environment
	if connStr == "" {
		connStr = NewConfiguration().TestConnectionString
	}

	if connStr == "" {
		return nil, nil
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
