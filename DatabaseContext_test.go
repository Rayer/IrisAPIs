package IrisAPIs

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	db := GetDatabaseContext()

	result, _ := db.DbObject.QueryString("select * from mcds_tw_members")
	output, _ := json.MarshalIndent(result, "", "\t")
	fmt.Println(string(output))
}
