package godo

type mustPanic struct {
	// err is the original error that caused the panic
	err error
}
