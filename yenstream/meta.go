package yenstream

type HaveMeta interface {
	GetMeta(key string) any
	SetMeta(key string, value any)
}

// GetMeta implements HaveMeta.
func (h *Metadata) GetMeta(key string) any {
	return h.meta[key]
}

// SetMeta implements HaveMeta.
func (h *Metadata) SetMeta(key string, value any) {
	if h.meta == nil {
		h.meta = map[string]any{}
	}
	h.meta[key] = value
}

type Metadata struct {
	meta map[string]any
}

func SetMeta(ctx *RunnerContext, key, value string) Pipeline {
	return NewMap(ctx, func(data any) (any, error) {
		mdata := data.(HaveMeta)
		mdata.SetMeta(key, value)
		return data, nil
	})
}
