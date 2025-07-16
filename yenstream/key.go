package yenstream

import "encoding/json"

type KeyedItem[D any] interface {
	Key() any
	Data() D
}

type keyedItemImpl[D any] struct {
	Metadata
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

func (k *keyedItemImpl[D]) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"key":  k.Key(),
		"data": k.Data(), // Even though age is private
	})
}
