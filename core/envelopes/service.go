package envelopes

import (
	"context"
	"errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/ybinhome/envelope/infra/base"
	"github.com/ybinhome/envelope/services"
	"sync"
)

var once sync.Once

func init() {
	once.Do(func() {
		services.IRedEnvelopeService = new(redEnvelopeService)
	})
}

type redEnvelopeService struct {
}

// 发红包逻辑代码
func (r *redEnvelopeService) SendOut(dto services.RedEnvelopeSendingDTO) (activity *services.RedEnvelopeActivity, err error) {
	// 1. 验证参数是否合法
	if err = base.ValidateStruct(&dto); err != nil {
		return activity, err
	}

	// 2. 获取发送红包人的资金账户信息
	account := services.GetAccountService().GetEnvelopeAccountByUserId(dto.UserId)
	if account == nil {
		return nil, errors.New("用户账户不存在" + dto.UserId)
	}
	goods := dto.ToGoods()
	goods.AccountNo = account.AccountNo

	if goods.Blessing == "" {
		goods.Blessing = services.DefaultBlessing
	}
	if goods.EnvelopeType == services.GeneralEnvelopeType {
		goods.AmountOne = goods.Amount
		goods.Amount = decimal.Decimal{}
	}

	// 3. 执行发送红包逻辑
	domain := new(goodsDomain)
	activity, err = domain.SendOut(*goods)
	if err != nil {
		logrus.Error(err)
	}

	return activity, err
}

// 收红包逻辑代码
func (r *redEnvelopeService) Receive(dto services.RedEnvelopeReceiveDTO) (item *services.RedEnvelopeItemDTO, err error) {
	// 校验参数
	if err = base.ValidateStruct(&dto); err != nil {
		return nil, err
	}

	// 获取当前收红包用户的账户信息
	account := services.GetAccountService().GetEnvelopeAccountByUserId(dto.RecvUserId)
	if account == nil {
		return nil, errors.New("红包资金账户不存在: " + dto.RecvUserId)
	}

	// 尝试收红包
	domain := goodsDomain{}
	item, err = domain.Receive(context.Background(), dto)
	return item, err
}

func (r *redEnvelopeService) Refund(envelopeNo string) (order *services.RedEnvelopeGoodsDTO) {
	return nil
}

func (r *redEnvelopeService) Get(envelopeNo string) (order *services.RedEnvelopeGoodsDTO) {
	return nil
}
