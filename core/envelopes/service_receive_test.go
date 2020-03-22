package envelopes

import (
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"github.com/smartystreets/goconvey/convey"
	"github.com/ybinhome/envelope/services"
	"strconv"
	"testing"
)

func TestRedEnvelopeService_Receive(t *testing.T) {
	accountService := services.GetAccountService()
	convey.Convey("收红包测试用例", t, func() {
		accounts := make([]*services.AccountDTO, 0)
		size := 10
		// 1. 准备多个账户，用于发红包和收红包
		for i := 0; i < size; i++ {
			account := services.AccountCreateDTO{
				UserId:       ksuid.New().Next().String(),
				Username:     "测试用户" + strconv.Itoa(i+1),
				AccountName:  "测试用户" + strconv.Itoa(i+1),
				AccountType:  int(services.EnvelopeAccountType),
				CurrencyCode: "CNY",
				Amount:       "2000",
			}

			acDto, err := accountService.CreateAccount(account)
			convey.So(err, convey.ShouldBeNil)
			convey.So(acDto, convey.ShouldNotBeNil)
			accounts = append(accounts, acDto)
		}

		// 2. 使用其中一个账户来发送红包
		acDto := accounts[0]
		convey.So(len(accounts), convey.ShouldEqual, size)

		re := services.GetRedEnvelopeService()
		goods := services.RedEnvelopeSendingDTO{
			EnvelopeType: services.GeneralEnvelopeType,
			Username:     acDto.Username,
			UserId:       acDto.UserId,
			Blessing:     services.DefaultBlessing,
			Amount:       decimal.NewFromFloat(1.88),
			Quantity:     size,
		}
		at, err := re.SendOut(goods)
		convey.So(err, convey.ShouldBeNil)
		convey.So(at, convey.ShouldNotBeNil)
		convey.So(at.Link, convey.ShouldNotBeEmpty)
		convey.So(at.RedEnvelopeGoodsDTO, convey.ShouldNotBeNil)
		// 验证每一个属性
		dto := at.RedEnvelopeGoodsDTO
		convey.So(dto.Username, convey.ShouldEqual, goods.Username)
		convey.So(dto.UserId, convey.ShouldEqual, dto.UserId)
		convey.So(dto.Quantity, convey.ShouldEqual, goods.Quantity)
		convey.So(dto.Amount.String(), convey.ShouldEqual, goods.Amount.Mul(decimal.NewFromFloat(float64(dto.Quantity))).String())

		// 3. 使用红包数量的人来收红包
		remainAmount := at.Amount
		convey.Convey("收普通红包", func() {
			for _, account := range accounts {
				rcv := services.RedEnvelopeReceiveDTO{
					EnvelopeNo:   at.EnvelopeNo,
					RecvUsername: account.Username,
					RecvUserId:   account.UserId,
					AccountNo:    account.AccountNo,
				}
				item, err := re.Receive(rcv)
				convey.So(err, convey.ShouldBeNil)
				convey.So(item, convey.ShouldNotBeNil)
				convey.So(item.Amount, convey.ShouldEqual, at.AmountOne)
				remainAmount = remainAmount.Sub(at.AmountOne)
				convey.So(item.RemainAmount.String(), convey.ShouldEqual, remainAmount.String())
			}
		})
	})
}
