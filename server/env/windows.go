//go:build windows

package env

type windows struct {
}

func (sys *windows) GetHomeDir() string {
	//TODO
}

func (sys *windows) GetContextAppPath() string {
	return ""
	//TODO
}

func (sys *windows) SpecificSetup() {
	//TODO
}

func GetOsManager() system {
	return new(windows)
}
