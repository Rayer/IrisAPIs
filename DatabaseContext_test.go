package IrisAPIs

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestDbConnection(t *testing.T) {
	db, _ := NewDatabaseContext("acc:12qw34er@tcp(node.rayer.idv.tw:3306)/apps?charset=utf8&loc=Asia%2FTaipei&parseTime=true", true)
	result, _ := db.DbObject.QueryString("select * from mcds_tw_members")
	output, _ := json.MarshalIndent(result, "", "\t")
	fmt.Println(string(output))
}
