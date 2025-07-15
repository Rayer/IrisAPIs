module ServerUtils

go 1.23.0

toolchain go1.24.1

replace IrisAPIs => ./../

require (
	IrisAPIs v0.0.0-00010101000000-000000000000
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.20.1
)
