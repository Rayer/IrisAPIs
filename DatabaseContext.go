package IrisAPIs

import (
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"github.com/xormplus/xorm"
	"sync"
)

type DatabaseContext struct {
	DbObject *xorm.Engine
}

var databaseContext *DatabaseContext = nil

func GetDatabaseContext() *DatabaseContext {

	log.Debugf("Trying fetching DB : %+v", databaseContext)

	var mutex = &sync.Mutex{}
	mutex.Lock()
	if databaseContext == nil {
		log.Debugln("Database Initializing")
		databaseContext = &DatabaseContext{}
		var err error
		databaseContext.DbObject, err = initDatabaseContext()
		if err != nil {
			//Do something panic
			panic("Fail to init Database Object, error is : " + err.Error())
		}
	}
	log.Debugf("Out with database object : %+v", databaseContext)
	mutex.Unlock()
	return databaseContext
}

func initDatabaseContext() (engine *xorm.Engine, err error) {
	//engine, err = xorm.NewEngine("mysql", "acc:12qw34er@tcp(node.rayer.idv.tw:3306)/apps?charset=utf8&loc=Asia%2FTaipei&parseTime=true")
	engine, err = xorm.NewEngine("mysql", "acc:12qw34er@tcp(node.rayer.idv.tw:3306)/apps?charset=utf8&loc=Local&parseTime=true")

	if err != nil {
		return nil, err
	}

	engine.ShowSQL(true)
	//engine.SetLogger(log.Logger{})
	if err != nil {
		return nil, err
	}
	return engine, nil
}
