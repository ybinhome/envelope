package web

import (
	"github.com/kataras/iris/v12"
	"github.com/sirupsen/logrus"
	"github.com/ybinhome/envelope/infra"
	"github.com/ybinhome/envelope/infra/base"
	"github.com/ybinhome/envelope/services"
)

// 定义 web api 的时候，对于每个子业务定义统一的前缀
//    资金账户定义为：/account
//    通常 api 前缀前都要加上一个版本号，方便后续迭代，则最终账户模块的 api 前缀为：/v1/account

func init() {
	infra.RegisterApi(new(AccountApi))
}

type AccountApi struct {
}

func (a *AccountApi) Init() {
	// 初始化一个 iris 分组
	groupRouter := base.Iris().Party("/v1/account")
	// 注册 create 接口
	groupRouter.Post("/create", createHandler)
	// 注册 transfer 接口
	groupRouter.Post("/transfer", transferHandler)
}

// 1. 创建用户接口：/v1/account/create
//    使用 POST 方法，body 部分为 json 格式
func createHandler(ctx iris.Context) {
	// 获取请求参数，通过 iris 上下文中的 ReadJSON 方法从请求的 body 中读取 json 格式的数据，并解析到 AccountCreateDTO 结构体中
	account := services.AccountCreateDTO{}
	err := ctx.ReadJSON(&account)
	r := base.Response{
		Code: base.ResCodeOk,
	}
	if err != nil {
		r.Code = base.ResCodeRequestParamsError
		r.Message = err.Error()
		// 通过 iris 上下文的 JSON 方法将自定义的 Response 结构体以 json 格式解析到 body 中
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}

	// 执行创建账户逻辑
	service := services.GetAccountService()
	dto, err := service.CreateAccount(account)
	if err != nil {
		r.Code = base.ResCodeInterServerError
		r.Message = err.Error()
		logrus.Error(err)
	}
	r.Data = dto
	ctx.JSON(&r)
}

// 2. 转账接口：/v1/account/transfer
func transferHandler(ctx iris.Context) {
	// 获取请求参数，通过 iris 上下文中的 ReadJSON 方法从请求的 body 中读取 json 格式的数据，并解析到 AccountCreateDTO 结构体中
	account := services.AccountTransferDTO{}
	err := ctx.ReadJSON(&account)
	r := base.Response{
		Code: base.ResCodeOk,
	}
	if err != nil {
		r.Code = base.ResCodeRequestParamsError
		r.Message = err.Error()
		// 通过 iris 上下文的 JSON 方法将自定义的 Response 结构体以 json 格式解析到 body 中
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}

	// 执行转账逻辑
	service := services.GetAccountService()
	status, err := service.Transfer(account)
	if err != nil {
		r.Code = base.ResCodeInterServerError
		r.Message = err.Error()
		logrus.Error(err)
	}
	r.Data = status
	if status != services.TransferedStatusSuccess {
		r.Code = base.ResCodeBizError
	}
	ctx.JSON(&r)
}

// 3. 查询红包账户接口：/v1/account/envelope/get

// 4. 查询账户信息接口：/v1/account/get
