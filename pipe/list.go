package pipe

// Turn a slice into a chan.
func Stream[T any](input []T) (output chan T) {
	if input == nil {
		close(output)
	} else {
		go func() {
			defer func() { _ = recover() }()
			defer close(output)

			for _, value := range input {
				output <- value
			}
		}()
	}

	return
}

// Turn a chan into a slice.
func Gather[T any](input <-chan T) (list []T) {
	if input != nil {
		for value := range input {
			list = append(list, value)
		}
	}

	return
}

// Group items from a chan into chunks of a given size.
func Cluster[T any](input chan T, chunkSize int) (output chan []T) {
	if input == nil {
		close(output)
	} else {
		go func() {
			defer func() { _ = recover() }()
			defer close(output)
			defer close(input)

			chunk := make([]T, 0, chunkSize)

			for value := range input {
				chunk = append(chunk, value)

				if len(chunk) == chunkSize {
					output <- chunk
					chunk = make([]T, 0, chunkSize)
				}
			}

			if len(chunk) > 0 {
				output <- chunk
			}
		}()
	}

	return
}

// Open all lists in the input and send each item individually.
func Flatten[T any](input chan []T) (output chan T) {
	if input == nil {
		close(output)
	} else {
		go func() {
			defer func() { _ = recover() }()
			defer close(output)
			defer close(input)

			for chunk := range input {
				for _, value := range chunk {
					output <- value
				}
			}
		}()
	}

	return
}
