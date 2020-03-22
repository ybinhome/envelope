package accounts

import (
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/ybinhome/envelope/services"
	"testing"
)

func TestAccountDomain_Create(t *testing.T) {
	dto := services.AccountDTO{
		UserId:   ksuid.New().Next().String(),
		Username: "测试用户",
		Balance:  decimal.NewFromFloat(10),
		Status:   1,
	}

	domain := new(accountDomain)
	Convey("账户创建", t, func() {
		rdto, err := domain.Create(dto)
		So(err, ShouldBeNil)
		So(rdto, ShouldNotBeNil)
		So(rdto.Balance.String(), ShouldEqual, dto.Balance.String())
		So(rdto.UserId, ShouldEqual, dto.UserId)
		So(rdto.Username, ShouldEqual, dto.Username)
		So(rdto.Status, ShouldEqual, dto.Status)
	})
}

func TestAccountDomain_Transfer(t *testing.T) {
	// 1. 创建两个账户，一个主体账户，一个目标账户，主体账户有余额
	adto1 := &services.AccountDTO{
		UserId:   ksuid.New().Next().String(),
		Username: "测试账户1",
		Balance:  decimal.NewFromFloat(100),
		Status:   1,
	}
	adto2 := &services.AccountDTO{
		UserId:   ksuid.New().Next().String(),
		Username: "测试账户2",
		Balance:  decimal.NewFromFloat(100),
		Status:   1,
	}

	domain := accountDomain{}
	Convey("转账测试", t, func() {
		// 创建账户1
		dto1, err := domain.Create(*adto1)
		So(err, ShouldBeNil)
		So(dto1, ShouldNotBeNil)
		So(dto1.Balance.String(), ShouldEqual, adto1.Balance.String())
		So(dto1.UserId, ShouldEqual, adto1.UserId)
		So(dto1.Username, ShouldEqual, adto1.Username)
		So(dto1.Status, ShouldEqual, adto1.Status)
		adto1 = dto1

		// 创建账户2
		dto2, err := domain.Create(*adto2)
		So(err, ShouldBeNil)
		So(dto2, ShouldNotBeNil)
		So(dto2.Balance.String(), ShouldEqual, adto2.Balance.String())
		So(dto2.UserId, ShouldEqual, adto2.UserId)
		So(dto2.Username, ShouldEqual, adto2.Username)
		So(dto2.Status, ShouldEqual, adto2.Status)
		adto2 = dto2

		// 2. 余额充足，资金转入其他账户
		Convey("余额充足，资金转入其他账户", func() {
			// 设置转账金额
			amount := decimal.NewFromFloat(1)
			// 创建转账主体
			body := services.TradeParticipator{
				AccountNo: adto1.AccountNo,
				UserId:    adto1.UserId,
				Username:  adto1.Username,
			}
			// 创建转账目标
			target := services.TradeParticipator{
				AccountNo: adto2.AccountNo,
				UserId:    adto2.UserId,
				Username:  adto2.Username,
			}
			// 创建账户转账 DTO
			dto := services.AccountTransferDTO{
				TradeBody:   body,
				TradeTarget: target,
				TradeNo:     ksuid.New().Next().String(),
				Amount:      amount,
				ChangeType:  services.ChangeType(-1),
				ChangeFlag:  services.FlagTransferOut,
				Decs:        "转账",
			}

			// 完成转账，并验证转账逻辑
			status, err := domain.Transfer(dto)
			So(err, ShouldBeNil)
			So(status, ShouldEqual, services.TransferedStatusSuccess)

			// 验证实际余额更新后的预期值
			a2 := domain.GetAccount(adto1.AccountNo)
			So(a2, ShouldNotBeNil)
			So(a2.Balance.String(), ShouldEqual, adto1.Balance.Sub(amount).String())
		})

		// 3. 余额不足，资金转入其他账户
		Convey("余额不足，资金转入其他账户", func() {
			// 设置转账金额
			amount := adto1.Balance
			amount = amount.Add(decimal.NewFromFloat(200))

			// 创建转账主体
			body := services.TradeParticipator{
				AccountNo: adto1.AccountNo,
				UserId:    adto1.UserId,
				Username:  adto1.Username,
			}
			// 创建转账目标
			target := services.TradeParticipator{
				AccountNo: adto2.AccountNo,
				UserId:    adto2.UserId,
				Username:  adto2.Username,
			}
			// 创建账户转账 DTO
			dto := services.AccountTransferDTO{
				TradeBody:   body,
				TradeTarget: target,
				TradeNo:     ksuid.New().Next().String(),
				Amount:      amount,
				ChangeType:  services.ChangeType(-1),
				ChangeFlag:  services.FlagTransferOut,
				Decs:        "转账",
			}

			// 完成转账，并验证转账逻辑
			status, err := domain.Transfer(dto)
			So(err, ShouldNotBeNil)
			So(status, ShouldEqual, services.TransferedStatusSufficientFunds)

			// 验证实际余额更新后的预期值
			a2 := domain.GetAccount(adto1.AccountNo)
			So(a2, ShouldNotBeNil)
			So(a2.Balance.String(), ShouldEqual, adto1.Balance.String())
		})

		// 4. 主体账户充值
		Convey("主体账户充值", func() {
			// 设置转账金额
			amount := decimal.NewFromFloat(11.1)

			// 创建转账主体
			body := services.TradeParticipator{
				AccountNo: adto1.AccountNo,
				UserId:    adto1.UserId,
				Username:  adto1.Username,
			}
			// 创建转账目标
			target := services.TradeParticipator{
				AccountNo: adto2.AccountNo,
				UserId:    adto2.UserId,
				Username:  adto2.Username,
			}
			// 创建账户转账 DTO
			dto := services.AccountTransferDTO{
				TradeBody:   body,
				TradeTarget: target,
				TradeNo:     ksuid.New().Next().String(),
				Amount:      amount,
				ChangeType:  services.AccountStoreValue,
				ChangeFlag:  services.FlagTransferIn,
				Decs:        "储值",
			}

			// 完成转账，并验证转账逻辑
			status, err := domain.Transfer(dto)
			So(err, ShouldBeNil)
			So(status, ShouldEqual, services.TransferedStatusSuccess)

			// 验证实际余额更新后的预期值
			a2 := domain.GetAccount(adto1.AccountNo)
			So(a2, ShouldNotBeNil)
			So(a2.Balance.String(), ShouldEqual, adto1.Balance.Add(amount).String())
		})
	})
}
