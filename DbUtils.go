package IrisAPIs

func NullBool() *bool {
	var ret *bool
	ret = nil
	return ret
}

func PBool(value bool) *bool {
	return &value
}