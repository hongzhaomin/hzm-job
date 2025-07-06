package tools

// GetIds4Slice 切片获取ids
func GetIds4Slice[S ~[]E, E any](list S, fn func(E) int64) []int64 {
	if len(list) <= 0 {
		return nil
	}

	var ids []int64
	for _, ele := range list {
		ids = append(ids, fn(ele))
	}
	return ids
}
