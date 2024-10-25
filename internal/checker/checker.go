package checker

type Checker interface {
	InitConnection() error
	CheckAvailability() bool
	Name() string
}
