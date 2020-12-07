package IrisAPIs

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
)

type DatabaseContextTest struct {
	suite.Suite
	db *DatabaseContext
}

func (d *DatabaseContextTest) SetupSuite() {
	db, err := NewTestDatabaseContext()
	if err != nil {
		fmt.Println(err.Error())
		d.db = nil
	}
	d.db = db
}

func (d *DatabaseContextTest) SetupTest() {
	if d.db == nil {
		d.T().Skip("Test case skipped due to no database available.")
	}
}

func (d *DatabaseContextTest) TearDownSuite() {

}

func TestDatabaseContextTest(t *testing.T) {
	suite.Run(t, new(DatabaseContextTest))
}

func (d *DatabaseContextTest) TestDbConnection() {
	result, _ := d.db.DbObject.QueryString("select * from mcds_tw_members")
	output, _ := json.MarshalIndent(result, "", "\t")
	fmt.Println(string(output))
}
