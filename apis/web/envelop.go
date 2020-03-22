package web

import (
	"github.com/kataras/iris/v12"
	"github.com/ybinhome/envelope/infra"
	"github.com/ybinhome/envelope/infra/base"
	"github.com/ybinhome/envelope/services"
)

func init() {
	infra.RegisterApi(&EnvelopeApi{})
}

type EnvelopeApi struct {
	service services.RedEnvelopeService
}

func (e *EnvelopeApi) Init() {
	e.service = services.GetRedEnvelopeService()
	groupRouter := base.Iris().Party("/v1/envelope")
	groupRouter.Post("/sendout", e.sendOutHandler)
	groupRouter.Post("/receive", e.receiveHandler)
}

func (e *EnvelopeApi) sendOutHandler(ctx iris.Context) {
	r := base.Response{
		Code: base.ResCodeOk,
	}

	dto := services.RedEnvelopeSendingDTO{}
	err := ctx.ReadJSON(&dto)
	if err != nil {
		r.Code = base.ResCodeRequestParamsError
		r.Message = err.Error()
		ctx.JSON(&r)
		return
	}

	activity, err := e.service.SendOut(dto)
	if err != nil {
		r.Code = base.ResCodeRequestParamsError
		r.Message = err.Error()
		ctx.JSON(&r)
		return
	}

	r.Data = activity
	ctx.JSON(r)
}

func (e *EnvelopeApi) receiveHandler(ctx iris.Context) {
	r := base.Response{
		Code: base.ResCodeOk,
	}

	dto := services.RedEnvelopeReceiveDTO{}
	err := ctx.ReadJSON(&dto)
	if err != nil {
		r.Code = base.ResCodeRequestParamsError
		r.Message = err.Error()
		ctx.JSON(&r)
		return
	}

	item, err := e.service.Receive(dto)
	if err != nil {
		r.Code = base.ResCodeRequestParamsError
		r.Message = err.Error()
		ctx.JSON(&r)
		return
	}

	r.Data = item
	ctx.JSON(r)
}
