package yenstream

type KeyedItem[D any] interface {
	Key() any
	Data() D
}

type keyedItemImpl[D any] struct {
	key  any
	data D
}

func NewKeyedItem[D any](key any, data D) *keyedItemImpl[D] {
	return &keyedItemImpl[D]{
		key:  key,
		data: data,
	}
}

// Data implements KeyedItem.
func (k *keyedItemImpl[D]) Data() D {
	return k.data
}

// Key implements KeyedItem.
func (k *keyedItemImpl[D]) Key() any {
	return k.key
}
