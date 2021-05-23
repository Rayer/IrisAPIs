package IrisAPIs

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestGetFromPbs(t *testing.T) {
	db, err := NewTestDatabaseContext(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	s := NewPbsTrafficDataService(db)
	ret, err := s.FetchPbsFromServer(context.TODO())
	fmt.Println(ret)
	fmt.Println(err)
}

func TestPbsWriteDb(t *testing.T) {
	db, err := NewTestDatabaseContext(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	s := NewPbsTrafficDataService(db)
	data, err := s.FetchPbsFromServer(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	err = s.UpdateDatabase(context.TODO(), data[:10], nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestJoinedPbsData(t *testing.T) {
	db, err := NewTestDatabaseContext(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	s := NewPbsTrafficDataService(db)
	fmt.Println(s.(*PbsTrafficDataServiceImpl).GetHistory(context.TODO(), 2*time.Hour))
}
