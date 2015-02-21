package godo

// Handler is the interface which all task handlers eventually implement.
type Handler interface {
	Handle(*Context) error
}

// HandlerFunc is Handler adapter.
type HandlerFunc func() error

// Handle implements Handler.
func (f HandlerFunc) Handle(*Context) error {
	return f()
}

// VoidHandlerFunc is a Handler adapter.
type VoidHandlerFunc func()

// Handle implements Handler.
func (v VoidHandlerFunc) Handle(*Context) error {
	v()
	return nil
}

// ContextHandlerFunc is a Handler adapter.
type ContextHandlerFunc func(*Context) error

// Handle implements Handler.
func (c ContextHandlerFunc) Handle(ctx *Context) error {
	return c(ctx)
}

// VoidContextHandlerFunc is a Handler adapter.
type VoidContextHandlerFunc func(*Context)

// Handle implements Handler.
func (f VoidContextHandlerFunc) Handle(ctx *Context) error {
	f(ctx)
	return nil
}
