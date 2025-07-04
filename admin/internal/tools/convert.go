package tools

import (
	"errors"
	"reflect"
	"strconv"
)

// ReflectConvert4Any 将any(接口)类型的val转化为rt类型的反射值rv
// 注意：返回的是rt的指针类型的reflect.Value
func ReflectConvert4Any(rt reflect.Type, val any) (reflect.Value, error) {
	str, err := ConvertStr4Any(val)
	if err != nil {
		return reflect.Value{}, err
	}
	return ReflectConvert4Str(rt, str)
}

// ReflectConvert4Str 将string类型的val转化为rt类型的反射值rv
// 注意：返回的是rt的指针类型的reflect.Value
func ReflectConvert4Str(rt reflect.Type, p string) (reflect.Value, error) {
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	switch rt.Kind() {
	case reflect.String:
		return reflect.ValueOf(&p), nil
	case reflect.Int:
		pInt, err := Str2Int(p)
		return reflect.ValueOf(&pInt), err
	case reflect.Int8:
		pInt64, err := Str2Int64(p)
		pInt8 := int8(pInt64)
		return reflect.ValueOf(&pInt8), err
	case reflect.Int16:
		pInt64, err := Str2Int64(p)
		pInt16 := int16(pInt64)
		return reflect.ValueOf(&pInt16), err
	case reflect.Int32:
		pInt64, err := Str2Int64(p)
		pInt32 := int32(pInt64)
		return reflect.ValueOf(&pInt32), err
	case reflect.Int64:
		pInt64, err := Str2Int64(p)
		return reflect.ValueOf(&pInt64), err
	case reflect.Uint:
		pUint64, err := Str2Uint64(p)
		pUint := uint(pUint64)
		return reflect.ValueOf(&pUint), err
	case reflect.Uint8:
		pUint64, err := Str2Uint64(p)
		pUint8 := uint8(pUint64)
		return reflect.ValueOf(&pUint8), err
	case reflect.Uint16:
		pUint64, err := Str2Uint64(p)
		pUint16 := uint16(pUint64)
		return reflect.ValueOf(&pUint16), err
	case reflect.Uint32:
		pUint64, err := Str2Uint64(p)
		pUint32 := uint32(pUint64)
		return reflect.ValueOf(&pUint32), err
	case reflect.Uint64:
		pUint64, err := Str2Uint64(p)
		return reflect.ValueOf(&pUint64), err
	case reflect.Float32:
		pFloat64, err := Str2Float64(p)
		pFloat32 := float32(pFloat64)
		return reflect.ValueOf(&pFloat32), err
	case reflect.Float64:
		pFloat64, err := Str2Float64(p)
		return reflect.ValueOf(&pFloat64), err
	case reflect.Bool:
		var pBool bool
		pBool, err := strconv.ParseBool(p)
		return reflect.ValueOf(&pBool), err
	default:
		return reflect.Value{}, errors.New("无法转换的数据类型")
	}
}

// ConvertStr4Any 接口类型转string
func ConvertStr4Any(val any) (string, error) {
	var valStr string
	switch act := val.(type) {
	case int:
		valStr = strconv.FormatInt(int64(act), 10)
	case int8:
		valStr = strconv.FormatInt(int64(act), 10)
	case int16:
		valStr = strconv.FormatInt(int64(act), 10)
	case int32:
		valStr = strconv.FormatInt(int64(act), 10)
	case int64:
		valStr = strconv.FormatInt(act, 10)
	case uint:
		valStr = strconv.FormatInt(int64(act), 10)
	case uint8:
		valStr = strconv.FormatInt(int64(act), 10)
	case uint16:
		valStr = strconv.FormatInt(int64(act), 10)
	case uint32:
		valStr = strconv.FormatInt(int64(act), 10)
	case uint64:
		valStr = strconv.FormatInt(int64(act), 10)
	case float64:
		valStr = strconv.FormatFloat(act, 'f', -1, 64)
	case float32:
		valStr = strconv.FormatFloat(float64(act), 'f', -1, 32)
	case bool:
		valStr = strconv.FormatBool(act)
	case string:
		valStr = act
	default:
		return "", errors.New("无法转换的数据类型")
	}
	return valStr, nil
}

func Str2Int(p string) (int, error) {
	var pInt int
	if p == "" {
		return pInt, nil
	}
	return strconv.Atoi(p)
}

func Str2Int64(p string) (int64, error) {
	var pInt64 int64
	if p == "" {
		return pInt64, nil
	}
	return strconv.ParseInt(p, 10, 64)
}

func Str2Uint64(p string) (uint64, error) {
	var pUint64 uint64
	if p == "" {
		return pUint64, nil
	}
	return strconv.ParseUint(p, 10, 64)
}

func Str2Float64(p string) (float64, error) {
	var pFloat64 float64
	if p == "" {
		return pFloat64, nil
	}
	return strconv.ParseFloat(p, 64)
}
