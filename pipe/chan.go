package pipe

import "sync"

// Send all items from the input channels to one output channel.
func Merge[T any](buffer uint, channels ...<-chan T) <-chan T {
	out := make(chan T, buffer)
	queue := sync.WaitGroup{}
	queue.Add(len(channels))

	output := func(c <-chan T) {
		defer func() { _ = recover() }()
		defer queue.Done()

		for n := range c {
			out <- n
		}
	}

	for _, c := range channels {
		go output(c)
	}

	go func() {
		defer close(out)

		queue.Wait()
	}()

	return out
}

// Send all items from one input channel to multiple output channels.
func Clone[T any](input chan T, count uint) (output []chan T) {
	if input != nil && count > 0 {
		go func() {
			defer func() { _ = recover() }()
			defer close(input)

			output = make([]chan T, count)

			for i := range output {
				output[i] = make(chan T, cap(input))
			}

			var o chan<- T

			for item := range input {
				for _, o = range output {
					o <- item
				}
			}

			for _, o = range output {
				close(o)
			}
		}()
	}

	return
}

// Modify all items from the input channel and send each result to the output channel.
func Apply[I, O any](input chan I, callback func(I) O, buffer int) chan O {
	output := make(chan O, buffer)

	go func() {
		defer func() { _ = recover() }()
		defer close(input)
		defer close(output)

		for i := range input {
			output <- callback(i)
		}
	}()

	return output
}

// Modify all items from the input channel and send each result to the output channel,
// unless an error is returned.
func ApplyUntilError[I, O any](input chan I, callback func(I) (O, error), buffer int) (chan O, *error) {
	var err error

	output := make(chan O, buffer)

	if input == nil {
		close(output)
	} else {
		go func() {
			defer func() { _ = recover() }()
			defer close(input)
			defer close(output)

			var tmp O

			for i := range input {
				tmp, err = callback(i)
				if err != nil {
					break
				}

				output <- tmp
			}
		}()
	}

	return output, &err
}

// Keep only the items from the input channel that match the callback.
func Filter[T comparable](input chan T, callback func(T) bool, buffer uint) chan T {
	output := make(chan T, buffer)

	if input == nil {
		close(output)
	} else {
		go func() {
			defer func() { _ = recover() }()
			defer close(input)
			defer close(output)

			for item := range input {
				if callback(item) {
					output <- item
				}
			}
		}()
	}

	return output
}

// Apply a callback to each each item, and pass the result to the next item.
func Aggregate[T comparable, A any](input <-chan T, callback func(T, A), accumulator A) {
	if input != nil {
		for item := range input {
			callback(item, accumulator)
		}
	}
}
