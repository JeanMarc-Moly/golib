package pipe

import "sync"

func Wait(functions ...func()) {
	wait := sync.WaitGroup{}

	for _, function := range functions {
		wait.Add(1)

		go func(f func()) {
			defer wait.Done()
			f()
		}(function)
	}

	wait.Wait()
}
