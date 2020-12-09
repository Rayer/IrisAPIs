package IrisAPIs

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

func PValue(value interface{}) *interface{} {
	return &value
}
