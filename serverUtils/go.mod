module ServerUtils

go 1.15

replace IrisAPIs => ./../

require (
	IrisAPIs v0.0.0-00010101000000-000000000000
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.0
)
