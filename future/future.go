package future

type signal struct{}

type Future[T any] struct {
	ready chan signal
	value T
	err   error
}

var flag = signal{}

func New[T any](f func() (T, error)) *Future[T] {
	future := &Future[T]{ready: make(chan signal)}

	go func() {
		defer close(future.ready)
		future.value, future.err = f()
		future.ready <- flag
	}()

	return future
}

func (f *Future[T]) Await() (T, error) {
	<-f.ready
	return f.value, f.err
}

func (f *Future[T]) Get() T {
	v, _ := f.Await()
	return v
}

func (f *Future[T]) Error() error {
	_, err := f.Await()
	return err
}
