package yenstream

type coGroupImpl struct{}

func CoGroupByKey(group map[string]Pipeline) *coGroupImpl {
	return &coGroupImpl{}
}
