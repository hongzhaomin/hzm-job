package tool

import (
	"crypto/md5"
	"encoding/hex"
	"reflect"
)

func BeanConv[A any, B any](as []*A, convFunc func(a *A) (*B, bool)) []*B {
	if len(as) <= 0 {
		return nil
	}

	var targetSlice []*B
	for _, a := range as {
		b, ok := convFunc(a)
		if ok {
			targetSlice = append(targetSlice, b)
		}
	}
	return targetSlice
}

func BeanConv4Basic[A ~int64 | ~int | ~string | ~bool, B any](as []A, convFunc func(a A) (*B, bool)) []*B {
	if len(as) <= 0 {
		return nil
	}

	var targetSlice []*B
	for _, a := range as {
		b, ok := convFunc(a)
		if ok {
			targetSlice = append(targetSlice, b)
		}
	}
	return targetSlice
}

func GetOrDefault[T any](obj any, defaultVal T, fn func() T) T {
	if isNil(obj) {
		return defaultVal
	}
	return fn()
}

func isNil(a any) bool {
	// 空接口判断是否为nil，必须类型和值都为nil才会返回真
	// 大多数情况是：type不为nil, value为nil
	// 那就需要继续判断了，使用反射判断
	if a == nil {
		return true
	}
	rv := reflect.ValueOf(a)
	if !rv.IsValid() {
		return true
	}
	return rv.IsZero()
}

// MD5 将字符串进行md5加密处理
func MD5(data string) string {
	resultArr := md5.Sum([]byte(data))
	return hex.EncodeToString(resultArr[:])
}
