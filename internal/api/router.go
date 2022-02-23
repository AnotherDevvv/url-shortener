package api

type Router interface {
	Start(errc chan error)
	Close() error
}