package yenstream

type Row[R any] struct {
	Metadata
	Data R `json:"data"`
}
