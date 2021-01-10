package IrisAPIs

import (
	"google.golang.org/protobuf/types/known/wrapperspb"
	"time"
)

func NullBool() *bool {
	var ret *bool
	ret = nil
	return ret
}

func PBool(value bool) *bool {
	return &value
}

func PInt(value int) *int {
	return &value
}

func PString(value string) *string {
	return &value
}

func PTime(value time.Time) *time.Time {
	return &value
}

func PValue(value interface{}) *interface{} {
	return &value
}

func PTimestamp(value *time.Time) *int64 {
	if value == nil {
		return nil
	}
	ret := value.Unix()
	return &ret
}

func PGTimestamp(value *time.Time) *wrapperspb.Int64Value {
	if value == nil {
		return nil
	}
	return &wrapperspb.Int64Value{Value: value.Unix()}
}
