package tools

import "reflect"

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

func FindAnnotationValueByType(target, annotationBean any, tagName string) (string, bool) {
	targetType := reflect.TypeOf(target)
	if targetType.Kind() == reflect.Ptr {
		targetType = targetType.Elem()
	}

	if targetType.Kind() != reflect.Struct {
		return "", false
	}

	annotationType := reflect.TypeOf(annotationBean)
	if annotationType.Kind() == reflect.Ptr {
		annotationType = annotationType.Elem()
	}

	for i := 0; i < targetType.NumField(); i++ {
		structField := targetType.Field(i)
		if structField.Type == annotationType || structField.Type.Implements(annotationType) {
			return structField.Tag.Lookup(tagName)
		}
	}
	return "", false
}
