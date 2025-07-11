package yenstream

type KeyedItem[K comparable, D any] interface {
	Key() K
	Data() D
}

type keyedItemImpl[K comparable, D any] struct {
	key  K
	data D
}

func NewKeyedItem[K comparable, D any](key K, data D) *keyedItemImpl[K, D] {
	return &keyedItemImpl[K, D]{
		key:  key,
		data: data,
	}
}

// Data implements KeyedItem.
func (k *keyedItemImpl[K, D]) Data() D {
	return k.data
}

// Key implements KeyedItem.
func (k *keyedItemImpl[K, D]) Key() K {
	return k.key
}
