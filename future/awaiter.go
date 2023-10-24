package future

type Awaiter interface {
	Await()
}

func Await(waiters ...Awaiter) {
	AwaitAll(waiters)
}

func AwaitAll(waiters []Awaiter) {
	for _, w := range waiters {
		w.Await()
	}
}
