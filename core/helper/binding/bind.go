package binding

import (
	"github.com/ohwin/micro/core/errx"
	"github.com/ohwin/micro/core/rest/response"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type bindType string

const (
	JSON          bindType = "JSON"
	XML           bindType = "XML"
	Form          bindType = "Form"
	Query         bindType = "Query"
	FormPost      bindType = "FormPost"
	FormMultipart bindType = "FormMultipart"
	ProtoBuf      bindType = "ProtoBuf"
	MsgPack       bindType = "MsgPack"
	YAML          bindType = "YAML"
	URI           bindType = "Uri"
	Header        bindType = "Header"
	TOML          bindType = "TOML"
)

var bindMap = map[bindType]interface{}{
	JSON:          binding.JSON,
	XML:           binding.XML,
	Form:          binding.Form,
	Query:         binding.Query,
	FormPost:      binding.FormPost,
	FormMultipart: binding.FormMultipart,
	ProtoBuf:      binding.ProtoBuf,
	MsgPack:       binding.MsgPack,
	YAML:          binding.YAML,
	URI:           binding.Uri,
	Header:        binding.Header,
	TOML:          binding.TOML,
}

type Binding struct {
	err error
	ctx *gin.Context
	req any
}

func Auto(ctx *gin.Context, req any, opts ...bindType) error {
	b := new(Binding)
	b.ctx = ctx
	b.req = req
	if len(opts) > 0 {
		for _, opt := range opts {
			b.withBind(opt)
		}
	}
	if b.err != nil {
		response.Fail(b.ctx, errx.ErrBadRequest)
	}
	return b.err
}

func (b *Binding) withBind(types bindType) {
	if b.err != nil {
		return
	}
	bind := bindMap[types]
	switch bind.(type) {
	case binding.BindingBody:
		b.err = b.ctx.ShouldBindBodyWith(b.req, bind.(binding.BindingBody))
	case binding.BindingUri:
		b.err = b.ctx.ShouldBindUri(b.req)
	default:
		b.err = b.ctx.ShouldBindWith(b.req, bind.(binding.Binding))
	}
}
