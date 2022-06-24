package env

type system interface {
	GetHomeDir() string
	GetContextAppPath() string
	SpecificSetup()
}
