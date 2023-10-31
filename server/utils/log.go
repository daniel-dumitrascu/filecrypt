package utils

type Log interface {
	Info(v ...any)
	Error(v ...any)
	Fatal(v ...any)
}
