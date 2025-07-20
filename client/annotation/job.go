package annotation

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/hongzhaomin/hzm-job/client/internal/anno"
	"github.com/hongzhaomin/hzm-job/core/tools"
	"reflect"
	"strings"
)

type HzmJob[Param any] struct {
	anno.HzmJobIface[Param]
}

// DoHandle 任务处理逻辑
func (my HzmJob[Param]) DoHandle(ctx context.Context, param *Param) error {
	panic("implement me")
}

func (my HzmJob[Param]) ParseParam(args *string) *Param {
	if args == nil || *args == "" {
		return nil
	}
	s := *args

	var param Param
	rt := reflect.TypeOf(param)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	isJsonStr := (strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]")) ||
		(strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}"))
	canParseJsonType := rt.Kind() == reflect.Struct || rt.Kind() == reflect.Map || rt.Kind() == reflect.Slice || rt.Kind() == reflect.Array
	if isJsonStr && canParseJsonType {
		if err := json.Unmarshal([]byte(s), &param); err != nil {
			panic(errors.New("jobParams parse error: " + err.Error()))
		}
		return &param
	}

	paramPtr := &param
	rvArgs, err := tools.ReflectConvert4Str(rt, s)
	if err != nil {
		panic(errors.New("jobParams parse error: " + err.Error()))
	}
	reflect.ValueOf(paramPtr).Elem().Set(rvArgs.Elem())
	return paramPtr
}

// JobFunc 任务处理函数
type JobFunc func(ctx context.Context, param *string) error

func (my JobFunc) DoHandle(ctx context.Context, param *string) error {
	return my(ctx, param)
}

func (my JobFunc) ParseParam(args *string) *string {
	return args
}

func (my JobFunc) MarkHzmJob() {}
