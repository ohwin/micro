package errx

import (
	"errors"
	. "github.com/ohwin/micro/core/constant/status"
)

type ErrX struct {
	Code int
	Err  error
}

func (e ErrX) Error() string {
	return e.Err.Error()
}

func New(code int, msg string) ErrX {

	return ErrX{
		Code: code,
		Err:  errors.New(msg),
	}
}

func errX(code int) ErrX {
	return ErrX{
		Code: code,
		Err:  errors.New(Msg(code)),
	}
}

var (
	ErrNotFound        = errX(NotFound)          // 数据未找到
	ErrInternalServer  = errX(InternalServerErr) // 服务器内部错误
	ErrUnknown         = errX(UnknownErr)        // 未知错误
	ErrBadRequest      = errX(BadRequest)        // 错误请求
	ErrVerify          = errX(VerifyErr)         // 验证错误
	ErrTooManyRequests = errX(TooManyRequests)   // 请求频繁
)
