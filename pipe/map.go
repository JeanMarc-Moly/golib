package pipe

type Item[K, V any] struct {
	Key   K
	Value V
}

func StreamItems[K comparable, V any](input map[K]V) chan *Item[K, V] {
	items := make(chan *Item[K, V])

	if input != nil {
		go func() {
			defer func() { _ = recover() }()
			defer close(items)

			for key, value := range input {
				items <- &Item[K, V]{key, value}
			}
		}()
	} else {
		close(items)
	}

	return items
}

func GetItemKey[K, V any](item *Item[K, V]) K   { return item.Key }
func GetItemValue[K, V any](item *Item[K, V]) V { return item.Value }
