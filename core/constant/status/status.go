package status

const (
	Success           = 200 // 成功
	InternalServerErr = 500 // 服务器内部错误
	Unauthorized      = 401 // 未登录
	BadRequest        = 402 // 错误请求
	Forbidden         = 403 // 没有权限
	NotFound          = 404 // 数据未找到
	UnknownErr        = 405 // 未知错误
	VerifyErr         = 406 // 验证错误
	TooManyRequests   = 429 // 请求频繁
)

var msgMap = map[int]string{
	Success:           "Success",
	Unauthorized:      "Unauthorized",
	InternalServerErr: "Internal Server",
	NotFound:          "Not Found",
	UnknownErr:        "Unknown Error",
	BadRequest:        "Bad Request",
	Forbidden:         "Forbidden",
	VerifyErr:         "Verify Error",
	TooManyRequests:   "Too Many Requests",
}

func Msg(code int) string {
	return msgMap[code]
}
