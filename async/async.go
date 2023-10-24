package async

import "sync"

type Signal struct{}

var flag = Signal{}

func Run(f func()) <-chan Signal {
	future := make(chan Signal)

	go func() {
		defer close(future)
		f()
		future <- flag
	}()

	return future
}

func Wait[D any](futures []<-chan D) <-chan Signal {
	return Run(func() {
		for _, f := range futures {
			<-f
		}
	})
}

func Get[D any](f func() D) <-chan D {
	future := make(chan D)

	go func() {
		defer close(future)
		future <- f()
	}()

	return future
}

func Gather[D any](futures []<-chan D) []D {
	result := make([]D, len(futures))

	for i, f := range futures {
		result[i] = <-f
	}

	return result
}

func Collect[D any](futures <-chan D) []D {
	result := []D{}

	for f := range futures {
		result = append(result, f)
	}

	return result
}

func Execute[E any, D any](f func(E) D, e E) <-chan D {
	future := make(chan D)

	go func() {
		defer close(future)
		future <- f(e)
	}()

	return future
}

func Map[E any, D any](f func(E) D, es []E) <-chan []D {
	future := make(chan []D)

	go func() {
		defer close(future)

		result := make([]D, len(es))
		wait := sync.WaitGroup{}

		for i, e := range es {
			wait.Add(1)

			go func(j int, d E) {
				defer wait.Done()
				result[j] = f(d)
			}(i, e)
		}

		wait.Wait()
		future <- result
	}()

	return future
}

func Pipe[E any, D any](f func(E) D, es <-chan E) <-chan D {
	future := make(chan D)

	go func() {
		defer close(future)

		wait := sync.WaitGroup{}

		for e := range es {
			wait.Add(1)

			go func(d E) {
				defer wait.Done()
				future <- f(d)
			}(e)
		}

		wait.Wait()
	}()

	return future
}

func Reduce[E any, D any](f func(E, D) D, es <-chan E, d D) <-chan D {
	future := make(chan D)

	go func() {
		defer close(future)

		for e := range es {
			d = f(e, d)
		}

		future <- d
	}()

	return future
}
