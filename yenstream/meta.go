package yenstream

type HaveMeta interface {
	GetMeta(key string) any
	SetMeta(key string, value any)
}

// GetMeta implements HaveMeta.
func (h *Metadata) GetMeta(key string) any {
	return h.StreamMeta[key]
}

// SetMeta implements HaveMeta.
func (h *Metadata) SetMeta(key string, value any) {
	if h.StreamMeta == nil {
		h.StreamMeta = map[string]any{}
	}
	h.StreamMeta[key] = value
}

type Metadata struct {
	StreamMeta map[string]any `json:"meta"`
}

func SetMeta(ctx *RunnerContext, key, value string) Pipeline {
	return NewMap(ctx, func(data any) (any, error) {
		mdata := data.(HaveMeta)
		mdata.SetMeta(key, value)
		return data, nil
	})
}
