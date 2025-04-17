package streampipe

type Publisher[T any] interface {
	Chan() chan T
	Start() error
	Cancel() error
}
