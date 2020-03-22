package accounts

import (
	"context"
	"errors"
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
	"github.com/ybinhome/envelope/infra/base"
	"github.com/ybinhome/envelope/services"
)

// 领域模型是有状态的，每次使用时都要实例化
type accountDomain struct {
	account    Account
	accountLog AccountLog
}

// 用于发红包时红包支付环节
func NewAccountDomain() *accountDomain {
	return new(accountDomain)
}

// 创建 logNo 的逻辑
func (domain *accountDomain) createAccountLogNo() {
	// 全局唯一的 id，暂时采用 ksuid 的 id 生成策略创建 logNo，后期会优化成可读性较好的分布式 id
	domain.accountLog.LogNo = ksuid.New().Next().String()
}

// 创建 accountNo 的逻辑
func (domain *accountDomain) createAccountNo() {
	domain.account.AccountNo = ksuid.New().Next().String()
}

// 创建流水记录
func (domain *accountDomain) createAccountLog() {
	// 通过 account 来创建流水，需要注意的是创建账户逻辑在前
	domain.accountLog = AccountLog{}
	domain.createAccountLogNo()
	domain.accountLog.TradeNo = domain.accountLog.LogNo

	// 流水中的交易主题信息
	domain.accountLog.AccountNo = domain.account.AccountNo
	domain.accountLog.UserId = domain.account.UserId
	domain.accountLog.Username = domain.account.Username.String

	// 交易对象信息
	domain.accountLog.TargetAccountNo = domain.account.AccountNo
	domain.accountLog.TargetUserId = domain.account.UserId
	domain.accountLog.TargetUsername = domain.account.Username.String

	// 交易金额
	domain.accountLog.Amount = domain.account.Balance
	domain.accountLog.Balance = domain.account.Balance

	// 交易变化属性
	domain.accountLog.Decs = "账户创建"
	domain.accountLog.ChangeType = services.AccountCreated
	domain.accountLog.ChangeFlag = services.FlagAccountCreated
}

// 账户创建的业务逻辑
func (domain *accountDomain) Create(dto services.AccountDTO) (*services.AccountDTO, error) {
	// 创建账户持久化对象
	domain.account = Account{}
	domain.account.FromDTO(&dto)
	domain.createAccountNo()
	domain.account.Username.Valid = true

	// 创建账户流水持久化对象
	domain.createAccountLog()

	// 数据库持久化操作对象，dao 对象是有状态的，每次使用的连接和事务都是不同的，因此需要先实例化
	accountDao := AccountDao{}
	accountLogDao := AccountLogDao{}

	// 整个持久化账户对象和持久化账户流水对象在同一个事务中
	var rdto *services.AccountDTO
	// base.Tx 是一个快捷的事务操作函数，函数中所有的操作被认为是一个事务，如果函数返回一个非空的 error，整个事务就会回滚
	err := base.Tx(func(runner *dbx.TxRunner) error {
		accountDao.runner = runner
		accountLogDao.runner = runner

		// 插入账户数据
		id, err := accountDao.Insert(&domain.account)
		if err != nil {
			return err
		}
		if id <= 0 {
			return errors.New("账户创建失败")
		}

		// 如果插入成功就插入账户流水数据
		id, err = accountLogDao.Insert(&domain.accountLog)
		if err != nil {
			return err
		}
		if id <= 0 {
			return errors.New("账户流水创建失败")
		}

		domain.account = *accountDao.GetOne(domain.account.AccountNo)
		rdto = domain.account.ToDTO()
		return nil
	})
	return rdto, err
}

// 转账业务逻辑
//func (a *accountDomain) Transfer(dto services.AccountTransferDTO) (status services.TransferedStatus, err error) {
//	// 如果是支出交易，修正 amount
//	amount := dto.Amount
//	if dto.ChangeFlag == services.FlagTransferOut {
//		amount = amount.Mul(decimal.NewFromFloat(-1))
//	}
//
//	// 创建账户流水记录
//	a.accountLog = AccountLog{}
//	a.accountLog.FromTransferDTO(&dto)
//	a.createAccountLogNo()
//
//	// 更新余额，通过乐观锁的方式来验证余额是否足够，即在 sql 语句中更新余额的同时，验证余额是否足够
//	// 更新余额成功后，写入流水记录
//	err = base.Tx(func(runner *dbx.TxRunner) error {
//		accountDao := AccountDao{runner: runner}
//		accountLogDao := AccountLogDao{runner: runner}
//
//		// 更新余额
//		rows, err := accountDao.UpdateBalance(dto.TradeBody.AccountNo, amount)
//		if err != nil {
//			status = services.TransferedStatusFailure
//			return err
//		}
//		if rows <= 0 && dto.ChangeFlag == services.FlagTransferOut {
//			status = services.TransferedStatusSufficientFunds
//			return errors.New("余额不足")
//		}
//
//		// 更新账户流水
//		account := accountDao.GetOne(dto.TradeBody.AccountNo)
//		if account == nil {
//			return errors.New("账户出错")
//		}
//		a.account = *account
//		a.accountLog.Balance = a.account.Balance
//		id, err := accountLogDao.Insert(&a.accountLog)
//		if err != nil || id <= 0 {
//			status = services.TransferedStatusFailure
//			return errors.New("账户流水创建失败")
//		}
//
//		return nil
//	})
//
//	if err != nil {
//		logrus.Error(err)
//	} else {
//		status = services.TransferedStatusSuccess
//	}
//	return status, err
//}
// 改造原来的 Transfer 函数，使其支持夸方法的事务，从而支持发红包时红包支付部分的事务需求
func (a *accountDomain) Transfer(dto services.AccountTransferDTO) (status services.TransferedStatus, err error) {
	err = base.Tx(func(runner *dbx.TxRunner) error {
		ctx := base.WithValueContext(context.Background(), runner)
		status, err = a.TransferWithContext(ctx, dto)
		return err
	})
	return status, err
}

// TransferWithContext 方法必须在 base.Tx 事务中运行，不支持单独运行
func (a *accountDomain) TransferWithContext(ctx context.Context, dto services.AccountTransferDTO) (status services.TransferedStatus, err error) {
	// 如果是支出交易，修正 amount
	amount := dto.Amount
	if dto.ChangeFlag == services.FlagTransferOut {
		amount = amount.Mul(decimal.NewFromFloat(-1))
	}

	// 创建账户流水记录
	a.accountLog = AccountLog{}
	a.accountLog.FromTransferDTO(&dto)
	a.createAccountLogNo()

	// 更新余额，通过乐观锁的方式来验证余额是否足够，即在 sql 语句中更新余额的同时，验证余额是否足够
	// 更新余额成功后，写入流水记录
	err = base.ExecuteContext(ctx, func(runner *dbx.TxRunner) error {
		accountDao := AccountDao{runner: runner}
		accountLogDao := AccountLogDao{runner: runner}

		// 更新余额
		rows, err := accountDao.UpdateBalance(dto.TradeBody.AccountNo, amount)
		if err != nil {
			status = services.TransferedStatusFailure
			return err
		}
		if rows <= 0 && dto.ChangeFlag == services.FlagTransferOut {
			status = services.TransferedStatusSufficientFunds
			return errors.New("余额不足")
		}

		// 更新账户流水
		account := accountDao.GetOne(dto.TradeBody.AccountNo)
		if account == nil {
			return errors.New("账户出错")
		}
		a.account = *account
		a.accountLog.Balance = a.account.Balance
		id, err := accountLogDao.Insert(&a.accountLog)
		if err != nil || id <= 0 {
			status = services.TransferedStatusFailure
			return errors.New("账户流水创建失败")
		}

		return nil
	})

	if err != nil {
		logrus.Error(err)
	} else {
		status = services.TransferedStatusSuccess
	}
	return status, err
}

// 根据账户编号来查询账户信息
func (a *accountDomain) GetAccount(accountNo string) *services.AccountDTO {
	accountDao := AccountDao{}
	var account *Account

	err := base.Tx(func(runner *dbx.TxRunner) error {
		accountDao.runner = runner
		account = accountDao.GetOne(accountNo)
		return nil
	})

	if err != nil {
		return nil
	}

	if account == nil {
		return nil
	}

	return account.ToDTO()
}

// 根据用户 id 来查询红包账户信息
func (a *accountDomain) GetEnvelopeAccountByUserId(userId string) *services.AccountDTO {
	accountDao := AccountDao{}
	var account *Account

	err := base.Tx(func(runner *dbx.TxRunner) error {
		accountDao.runner = runner
		account = accountDao.GetByUserId(userId, int(services.EnvelopeAccountType))
		return nil
	})

	if err != nil {
		return nil
	}

	if account == nil {
		return nil
	}

	return account.ToDTO()
}

// 根据流水 id 来查询账户流水
func (a *accountDomain) GetAccountLog(logNo string) *services.AccountLogDTO {
	dao := AccountLogDao{}
	var log *AccountLog

	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao.runner = runner
		log = dao.GetOne(logNo)
		return nil
	})

	if err != nil {
		logrus.Error(err)
		return nil
	}

	if log == nil {
		return nil
	}

	return log.ToDTO()
}

// 根据交易编号来查询账户流水
func (a *accountDomain) GetAccountLogByTrade(tradeNo string) *services.AccountLogDTO {
	dao := AccountLogDao{}
	var log *AccountLog

	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao.runner = runner
		log = dao.GetByTradeNo(tradeNo)
		return nil
	})

	if err != nil {
		logrus.Error(err)
		return nil
	}

	if log == nil {
		return nil
	}

	return log.ToDTO()
}
