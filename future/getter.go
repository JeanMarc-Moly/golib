package future

type Getter[T any] interface {
	Get() T
}

func Chain[T, U any](f Getter[T], g func(T) (U, error)) *Future[U] {
	return New(func() (U, error) {
		return g(f.Get())
	})
}
