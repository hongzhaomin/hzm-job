package sdk

const (
	SucCode    = 200
	SucResult  = true
	FailCode   = -1
	FailResult = false
)

type Result[T any] struct {
	Code    *int    `json:"code"`
	Msg     *string `json:"msg"`
	Data    *T      `json:"data"`
	Success bool    `json:"success"`
}

func Ok() Result[any] {
	rb := new(Result[any])
	rb.ok()
	return *rb
}

func Ok2[T any](data T) Result[T] {
	rb := new(Result[T])
	rb.ok()
	rb.Data = &data
	return *rb
}

func Fail[T any](msg string) Result[T] {
	rb := new(Result[T])
	rb.fail(msg)
	return *rb
}

func Fail2[T any](code int, msg string) Result[T] {
	rb := new(Result[T])
	rb.fail(msg)
	rb.Code = &code
	return *rb
}

func (rb *Result[T]) ok() *Result[T] {
	code := SucCode
	rb.Code = &code
	rb.Success = SucResult
	return rb
}

func (rb *Result[T]) fail(msg string) *Result[T] {
	code := FailCode
	rb.Code = &code
	rb.Success = FailResult
	rb.Msg = &msg
	return rb
}
