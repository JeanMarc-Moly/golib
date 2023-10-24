package pipe

// func DropDuplicates[T comparable](input chan T, buffer, uniques int) chan T {
// 	output := make(chan T, buffer)

// 	go func() {
// 		defer func() { _ = recover() }()
// 		defer close(input)
// 		defer close(output)

// 		seen := make(set.Set[T], uniques)

// 		for item := range input {
// 			if !seen.Contains(item) {
// 				seen.Add(item)
// 				output <- item
// 			}
// 		}
// 	}()

// 	return output
// }

// func Exclude[T comparable](input chan T, forbidden set.Set[T]) chan T {
// 	return Filter(input, func(item T) bool { return !forbidden.Contains(item) })
// }
