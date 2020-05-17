package IrisAPIs

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/xormplus/xorm"
)

type DatabaseContext struct {
	DbObject *xorm.Engine
}

var databaseContext *DatabaseContext

func GetDatabaseContext() *DatabaseContext {

	if databaseContext == nil {
		databaseContext = &DatabaseContext{}
		var err error
		databaseContext.DbObject, err = initDatabaseContext()
		if err != nil {
			//Do something panic
			panic("Fail to init Database Object!")
		}
	}
	return databaseContext
}

func initDatabaseContext() (engine *xorm.Engine, err error){
	engine, err = xorm.NewEngine("mysql", "acc:12qw34er@tcp(node.rayer.idv.tw:3306)/apps?charset=utf8&loc=Asia%2FTaipei&parseTime=true")
	if err != nil {
		return nil, err
	}
	return engine,nil
}
