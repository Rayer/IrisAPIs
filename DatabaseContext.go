package IrisAPIs

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/xormplus/xorm"
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

func initDatabaseContext(connectionString string, showSql bool) (engine *xorm.Engine, err error) {
	//engine, err = xorm.NewEngine("mysql", "acc:12qw34er@tcp(node.rayer.idv.tw:3306)/apps?charset=utf8&loc=Asia%2FTaipei&parseTime=true")
	engine, err = xorm.NewEngine("mysql", connectionString)

	if err != nil {
		return nil, err
	}

	engine.ShowSQL(showSql)
	return engine, nil
}
