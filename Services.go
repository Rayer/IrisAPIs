package IrisAPIs

type TeardownableServices interface {
	Teardown() error
}
