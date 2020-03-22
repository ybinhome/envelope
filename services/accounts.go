package services

import (
	"github.com/shopspring/decimal"
	"github.com/ybinhome/envelope/infra/base"
	"time"
)

// 对外暴漏 AccountService 接口的实例化方法，AccountService 仅需要被实例化一次，为了防止 AccountService 对象被多次实例化，
//    此处在 AccountService 接口实现中通过 init() 函数结合 sync 包 Once 的 Do 方法实例化。
var IAccountService AccountService

func GetAccountService() AccountService {
	base.Check(IAccountService)
	return IAccountService
}

// 创建账户服务接口
type AccountService interface {
	CreateAccount(dto AccountCreateDTO) (*AccountDTO, error)
	Transfer(dto AccountTransferDTO) (TransferedStatus, error)
	StoreValue(dto AccountTransferDTO) (TransferedStatus, error)
	GetEnvelopeAccountByUserId(userId string) *AccountDTO
	GetAccount(accountNo string) *AccountDTO
}

// 定义 账户服务接口 - 账户创建 的传参
type AccountCreateDTO struct {
	UserId       string
	Username     string
	AccountName  string
	AccountType  int
	CurrencyCode string
	Amount       string // 定义为 float34 和 float64 计算时会丢失精度
}

// 定义 账户服务接口 - 账户创建 返回值
type AccountDTO struct {
	AccountNo    string          //账户编号,账户唯一标识
	AccountName  string          //账户名称,用来说明账户的简短描述,账户对应的名称或者命名，比如xxx积分、xxx零钱
	AccountType  int             //账户类型，用来区分不同类型的账户：积分账户、会员卡账户、钱包账户、红包账户
	CurrencyCode string          //货币类型编码：CNY人民币，EUR欧元，USD美元 。。。
	UserId       string          //用户编号, 账户所属用户
	Username     string          //用户名称
	Balance      decimal.Decimal //账户可用余额
	Status       int             //账户状态，账户状态：0账户初始化，1启用，2停用
	CreatedAt    time.Time       //创建时间
	UpdatedAt    time.Time       //更新时间
}

// 定义 账户服务接口 - 转账 - 交易参与者
type TradeParticipator struct {
	AccountNo string
	UserId    string
	Username  string
}

// 定义 账户服务接口 - 转账 的传参
type AccountTransferDTO struct {
	TradeNo     string            `validate:"required"`         //交易单号 全局不重复字符或数字，唯一性标识
	TradeBody   TradeParticipator `validate:"required"`         //交易主体
	TradeTarget TradeParticipator `validate:"required"`         //交易对象
	AmountStr   string            `validate:"required,numeric"` //交易金额,该交易涉及的金额
	Amount      decimal.Decimal   ``                            //交易金额,该交易涉及的金额
	ChangeType  ChangeType        `validate:"required,numeric"` //流水交易类型，0 创建账户，>0 为收入类型，<0 为支出类型，自定义
	ChangeFlag  ChangeFlag        `validate:"required,numeric"` //交易变化标识：-1 出账 1为进账，枚举
	Decs        string            ``                            //交易描述
}

type AccountLogDTO struct {
	LogNo           string          //流水编号 全局不重复字符或数字，唯一性标识
	TradeNo         string          //交易单号 全局不重复字符或数字，唯一性标识
	AccountNo       string          //账户编号 账户ID
	TargetAccountNo string          //账户编号 账户ID
	UserId          string          //用户编号
	Username        string          //用户名称
	TargetUserId    string          //目标用户编号
	TargetUsername  string          //目标用户名称
	Amount          decimal.Decimal //交易金额,该交易涉及的金额
	Balance         decimal.Decimal //交易后余额,该交易后的余额
	ChangeType      ChangeType      //流水交易类型，0 创建账户，>0 为收入类型，<0 为支出类型，自定义
	ChangeFlag      ChangeFlag      //交易变化标识：-1 出账 1为进账，枚举
	Status          int             //交易状态：
	Decs            string          //交易描述
	CreatedAt       time.Time       //创建时间
}
