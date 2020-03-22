package accounts

import (
	"errors"
	"github.com/shopspring/decimal"
	"github.com/ybinhome/envelope/infra/base"
	"github.com/ybinhome/envelope/services"
	"sync"
)

// AccountService 接口实例化，AccountService 全局只需要实例化一次即可，为了防止多次引用 accounts 包多次实例化，此处使用 sync 包中的 Once 的 Do 方法实例化 AccountService，
//    这样即便 accounts 包被多次引入，AccountService 结构体仅被实例化一次。
var once sync.Once

func init() {
	once.Do(func() {
		services.IAccountService = new(accountService)
	})
}

type accountService struct {
}

func (a *accountService) CreateAccount(dto services.AccountCreateDTO) (*services.AccountDTO, error) {
	domain := accountDomain{}

	// 1. 验证输入参数是否合法
	if err := base.ValidateStruct(&dto); err != nil {
		return nil, err
	}

	// 2. 执行账户创建的业务逻辑
	//    将 dto.Amount 从 string 类型转换为 decimal 类型
	amount, err := decimal.NewFromString(dto.Amount)
	if err != nil {
		return nil, err
	}
	//    将 AccountCreateDTO 转换为 AccountDTO
	account := services.AccountDTO{
		UserId:       dto.UserId,
		Username:     dto.Username,
		AccountType:  dto.AccountType,
		AccountName:  dto.AccountName,
		CurrencyCode: dto.CurrencyCode,
		Status:       1,
		Balance:      amount,
	}
	//    创建账户
	rdto, err := domain.Create(account)
	return rdto, err
}

func (a *accountService) Transfer(dto services.AccountTransferDTO) (services.TransferedStatus, error) {
	domain := accountDomain{}

	// 1. 验证输入参数是否合法
	if err := base.ValidateStruct(&dto); err != nil {
		return services.TransferedStatusFailure, err
	}
	//    将 dto.Amount 从 string 类型转换为 decimal 类型
	amount, err := decimal.NewFromString(dto.AmountStr)
	if err != nil {
		return services.TransferedStatusFailure, err
	}
	dto.Amount = amount
	//    验证 ChangeType 合法性
	if dto.ChangeFlag == services.FlagTransferOut {
		if dto.ChangeType > 0 {
			return services.TransferedStatusFailure, errors.New("如果 changeFlag 为支出，那么 changeType 必须小于 0")
		}
	} else {
		if dto.ChangeType < 0 {
			return services.TransferedStatusFailure, errors.New("如果 changeFlag 为收入，那么 changeType 必须大于 0")
		}
	}

	// 2. 执行转账逻辑
	status, err := domain.Transfer(dto)
	return status, err
}

func (a *accountService) StoreValue(dto services.AccountTransferDTO) (services.TransferedStatus, error) {
	// 储值是特殊类型的转账，修改 3 个差异变量：TradeTarget ChangeFlag ChangeType
	dto.TradeTarget = dto.TradeBody
	dto.ChangeFlag = services.FlagTransferIn
	dto.ChangeType = services.AccountStoreValue

	// 调用转账方法来实现存储
	return a.Transfer(dto)
}

func (a *accountService) GetEnvelopeAccountByUserId(userId string) *services.AccountDTO {
	domain := accountDomain{}
	account := domain.GetEnvelopeAccountByUserId(userId)
	return account
}
func (a *accountService) GetAccount(accountNo string) *services.AccountDTO {
	domain := accountDomain{}
	return domain.GetAccount(accountNo)
}
