package accounts

import (
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/ybinhome/envelope/services"
	"testing"
)

// 应用服务层 - 账户创建测试用例
func TestAccountService_CreateAccount(t *testing.T) {
	dto := services.AccountCreateDTO{
		UserId:       ksuid.New().Next().String(),
		Username:     "测试用户",
		Amount:       "100",
		AccountName:  "测试账户",
		AccountType:  2,
		CurrencyCode: "CNY",
	}
	service := new(accountService)
	Convey("账户创建", t, func() {
		rdto, err := service.CreateAccount(dto)
		So(err, ShouldBeNil)
		So(rdto, ShouldNotBeNil)
		So(rdto.Balance.String(), ShouldEqual, dto.Amount)
		So(rdto.UserId, ShouldEqual, dto.UserId)
		So(rdto.Username, ShouldEqual, dto.Username)
		So(rdto.Status, ShouldEqual, 1)
	})
}

// 应用服务层 - 转账业务测试用例
func TestAccountService_Transfer(t *testing.T) {
	Convey("转账", t, func() {
		// 1. 创建两个账户，一个主体账户，一个目标账户
		a1 := services.AccountCreateDTO{
			UserId:       ksuid.New().Next().String(),
			Username:     "测试用户1",
			Amount:       "100",
			AccountName:  "测试账户1",
			AccountType:  2,
			CurrencyCode: "CNY",
		}
		a2 := services.AccountCreateDTO{
			UserId:       ksuid.New().Next().String(),
			Username:     "测试用户2",
			Amount:       "100",
			AccountName:  "测试账户2",
			AccountType:  2,
			CurrencyCode: "CNY",
		}
		service := new(accountService)
		adto1, err := service.CreateAccount(a1)
		So(err, ShouldBeNil)
		So(adto1, ShouldNotBeNil)
		adto2, err := service.CreateAccount(a2)
		So(err, ShouldBeNil)
		So(adto2, ShouldNotBeNil)

		// 2. 主体账户余额足够
		Convey("主体账户余额足够", func() {
			// 创建两个交易主体，以及交易 dto
			body := services.TradeParticipator{
				AccountNo: adto1.AccountNo,
				UserId:    adto1.UserId,
				Username:  adto1.Username,
			}
			target := services.TradeParticipator{
				AccountNo: adto2.AccountNo,
				UserId:    adto2.UserId,
				Username:  adto2.Username,
			}
			amount := decimal.NewFromFloat(10)
			dto := services.AccountTransferDTO{
				TradeNo:     ksuid.New().Next().String(),
				TradeBody:   body,
				TradeTarget: target,
				AmountStr:   amount.String(),
				ChangeType:  services.ChangeType(-1),
				ChangeFlag:  services.FlagTransferOut,
				Decs:        "转出",
			}

			// 验证交易逻辑
			status, err := service.Transfer(dto)
			So(err, ShouldBeNil)
			So(status, ShouldEqual, services.TransferedStatusSuccess)

			// 验证交易后的余额
			ra1 := service.GetAccount(adto1.AccountNo)
			So(ra1, ShouldNotBeNil)
			So(ra1.Balance.String(), ShouldEqual, adto1.Balance.Sub(amount).String())
		})

		// 3。 主体账户余额不足
		Convey("主体账户余额不足", func() {
			// 创建两个交易主体，以及交易 dto
			body := services.TradeParticipator{
				AccountNo: adto1.AccountNo,
				UserId:    adto1.UserId,
				Username:  adto1.Username,
			}
			target := services.TradeParticipator{
				AccountNo: adto2.AccountNo,
				UserId:    adto2.UserId,
				Username:  adto2.Username,
			}
			amount := adto1.Balance.Add(decimal.NewFromFloat(200))
			dto := services.AccountTransferDTO{
				TradeNo:     ksuid.New().Next().String(),
				TradeBody:   body,
				TradeTarget: target,
				AmountStr:   amount.String(),
				ChangeType:  services.ChangeType(-1),
				ChangeFlag:  services.FlagTransferOut,
				Decs:        "转出",
			}

			// 验证交易逻辑
			status, err := service.Transfer(dto)
			So(err, ShouldNotBeNil)
			So(status, ShouldEqual, services.TransferedStatusSufficientFunds)

			// 验证交易后的余额
			ra1 := service.GetAccount(adto1.AccountNo)
			So(ra1, ShouldNotBeNil)
			So(ra1.Balance.String(), ShouldEqual, adto1.Balance.String())
		})

		// 4. 主体账户储值
		Convey("主体账户储值", func() {
			// 创建两个交易主体，以及交易 dto
			body := services.TradeParticipator{
				AccountNo: adto1.AccountNo,
				UserId:    adto1.UserId,
				Username:  adto1.Username,
			}
			target := body
			amount := decimal.NewFromFloat(10)
			dto := services.AccountTransferDTO{
				TradeNo:     ksuid.New().Next().String(),
				TradeBody:   body,
				TradeTarget: target,
				AmountStr:   amount.String(),
				ChangeType:  services.AccountStoreValue,
				ChangeFlag:  services.FlagTransferIn,
				Decs:        "储值",
			}

			// 验证交易逻辑
			status, err := service.Transfer(dto)
			So(err, ShouldBeNil)
			So(status, ShouldEqual, services.TransferedStatusSuccess)

			// 验证交易后的余额
			ra1 := service.GetAccount(adto1.AccountNo)
			So(ra1, ShouldNotBeNil)
			So(ra1.Balance.String(), ShouldEqual, adto1.Balance.Add(amount).String())
		})
	})

}
