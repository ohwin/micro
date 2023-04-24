package response

import (
	"github.com/gin-gonic/gin"
	. "github.com/ohwin/micro/core/constant/status"
	"github.com/ohwin/micro/core/errx"
	"github.com/ohwin/micro/core/rest/req"
	"math"
)

type PageInfo struct {
	Page      int   `json:"page"`
	Total     int   `json:"total"`
	PageSize  int   `json:"pageSize"`
	TotalSize int64 `json:"totalSize"`
}

type Response struct {
	Code      int         `json:"code"`
	Msg       string      `json:"msg"`
	Data      interface{} `json:"data"`
	*PageInfo `json:",omitempty"`
}

//
//func OK(ctx *gin.Context, data interface{}) {
//	resp := Response{
//		Code: Success,
//		Msg:  Msg(Success),
//		Data: data,
//	}
//	ctx.JSON(200, resp)
//}

type Option func(res *Response)

func Ok(ctx *gin.Context, data interface{}, opt ...Option) {
	resp := &Response{
		Code: Success,
		Msg:  Msg(Success),
		Data: data,
	}
	for _, option := range opt {
		option(resp)
	}
	ctx.JSON(200, resp)
}

// Fail 失败响应
func Fail(ctx *gin.Context, err error, opt ...Option) {
	var e errx.ErrX

	switch err.(type) {
	case errx.ErrX:
		e = err.(errx.ErrX)
	default:
		e = errx.ErrUnknown
	}

	resp := &Response{
		Code: e.Code,
		Msg:  e.Error(),
		Data: nil,
	}

	for _, option := range opt {
		option(resp)
	}

	ctx.JSON(200, resp)
}

func WithData(data interface{}) Option {
	return func(resp *Response) {
		resp.Data = data
	}
}

func WithMsg(msg string) Option {
	return func(resp *Response) {
		resp.Msg = msg
	}
}

func WithPage(totalSize int64, info req.PageInfo) Option {
	return func(resp *Response) {
		page := info.Page
		pageSize := info.PageSize
		total := int(math.Ceil(float64(totalSize) / float64(pageSize)))

		resp.PageInfo = &PageInfo{
			Page:      page,
			Total:     total,
			PageSize:  pageSize,
			TotalSize: totalSize,
		}
	}
}
