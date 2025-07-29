package anno

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

type HzmJobIface[Param any] interface {
	Job

	// DoHandle 任务处理逻辑
	DoHandle(ctx context.Context, param *Param) error

	// ParseParam 转换参数
	ParseParam(args *string) *Param
}

// Job 任务接口，为了能统一存储泛型类型 job，故而设计此接口
type Job interface {
	MarkHzmJob()
}

// CallDoHandle 反射调用 DoHandle 方法
// 由于go泛型不支持协变（Covariance），即A[string]不能自动转换为A[any]，即使string满足any的约束‌
// 只能使用反射的方式调用
func CallDoHandle(job Job, ctx context.Context, params *string) (err error) {
	defer func() {
		if ex := recover(); ex != nil {
			switch e := ex.(type) {
			case error:
				err = e
			case string:
				err = errors.New(e)
			default:
				err = fmt.Errorf("%v", e)
			}
		}
	}()
	if job == nil {
		return errors.New(fmt.Sprintf("job is not exist"))
	}
	rvJob := reflect.ValueOf(job)
	rvParam := rvJob.MethodByName("ParseParam").Call([]reflect.Value{reflect.ValueOf(params)})[0]
	rvResult := rvJob.MethodByName("DoHandle").Call([]reflect.Value{reflect.ValueOf(ctx), rvParam})
	errResult := rvResult[0].Interface()
	if errResult != nil {
		err = errResult.(error)
	}
	return
}
