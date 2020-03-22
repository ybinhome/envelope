package base

// 常见的 web 响应码
type ResCode int

const (
	ResCodeOk                 ResCode = 1000
	ResCodeValidationError    ResCode = 2000
	ResCodeRequestParamsError ResCode = 2100
	ResCodeInterServerError   ResCode = 5000
	ResCodeBizError           ResCode = 6000
)

// http 响应结构体
type Response struct {
	Code    ResCode     `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
