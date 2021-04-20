package IrisAPIs

import (
	"fmt"
	"testing"
)

func TestGetFromPbs(t *testing.T) {
	ret, err := FetchPbsFromServer()
	fmt.Println(ret)
	fmt.Println(err)
}

func TestPbsWriteDb(t *testing.T) {
	data, err := FetchPbsFromServer()
	if err != nil {
		t.Fatal(err)
	}
	db, err := NewTestDatabaseContext()
	if err != nil {
		t.Fatal(err)
	}
	err = UpdateDatabase(db.DbObject, data)
	if err != nil {
		t.Fatal(err)
	}
}
