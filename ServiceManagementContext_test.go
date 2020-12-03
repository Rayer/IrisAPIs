package IrisAPIs

import (
	"fmt"
	"testing"
)

func TestServiceManagementContext_CheckAllServerStatus(t *testing.T) {
	s := NewServiceManagementContext()
	fmt.Printf("%+v", s.CheckAllServerStatus())
}
