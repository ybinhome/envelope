package envelopes

import (
	"context"
	"database/sql"
	"errors"
	"github.com/shopspring/decimal"
	"github.com/tietang/dbx"
	"github.com/ybinhome/envelope/core/accounts"
	"github.com/ybinhome/envelope/infra/algo"
	"github.com/ybinhome/envelope/infra/base"
	"github.com/ybinhome/envelope/services"
)

var multiple = decimal.NewFromFloat(100.0)

// 抢红包业务逻辑代码
func (g *goodsDomain) Receive(ctx context.Context, dto services.RedEnvelopeReceiveDTO) (item *services.RedEnvelopeItemDTO, err error) {
	// 1. 创建收红包的订单明细
	// -- 通过 preCreateItem 子方法完成
	g.preCreateItem(dto)

	// 2. 查出当前红包的剩余数量和剩余金额信息
	goods := g.Get(dto.EnvelopeNo)

	// 3. 校验剩余红包和剩余金额
	// - 如果没有剩余，直接返回无可用红包金额
	if goods.RemainQuantity <= 0 || goods.RemainAmount.Cmp(decimal.NewFromFloat(0)) <= 0 {
		return nil, errors.New("没有足够的红包和金额")
	}

	// 4. 使用红包算法计算红包金额
	nextAmount := g.nextAmount(goods)

	// 5. 在一个事务中，更新红包数量和剩余金额，保存订单明细数据，完成转账
	err = base.Tx(func(runner *dbx.TxRunner) error {
		// (1) 使用乐观锁更新语句，尝试更新剩余数量和剩余金额
		// - 如果更新成功，也就是返回 1，表示抢到红包
		// - 如果更新失败，也就是返回 0，表示无可用红包数量和金额，抢红包失败
		dao := RedEnvelopeGoodsDao{runner: runner}
		rows, err := dao.UpdateBalance(goods.EnvelopeNo, nextAmount)
		if rows <= 0 || err != nil {
			return errors.New("没有足够的红包和金额了")
		}

		// (2) 保存订单明细数据
		g.item.Quantity = 1
		g.item.PayStatus = int(services.Paying)
		g.item.AccountNo = dto.AccountNo
		g.item.RemainAmount = goods.RemainAmount.Sub(nextAmount)
		g.item.Amount = nextAmount

		txCtx := base.WithValueContext(ctx, runner)
		_, err = g.item.Save(txCtx)
		if err != nil {
			return err
		}

		// (3) 将抢到的红包金额从系统红包中间账户转入当前用户的资金账户
		// -- 通过 transfer 子方法完成
		status, err := g.transfer(txCtx, dto)
		if status == services.TransferedStatusSuccess {
			return nil
		}
		return err
	})

	return g.item.ToDTO(), err
}

// 预创建收红包的订单明细子方法
func (g *goodsDomain) preCreateItem(dto services.RedEnvelopeReceiveDTO) {
	g.item.AccountNo = dto.AccountNo
	g.item.EnvelopeNo = dto.EnvelopeNo
	g.item.RecvUsername = sql.NullString{String: dto.RecvUsername, Valid: true}
	g.item.RecvUserId = dto.RecvUserId
	g.item.createItemNo()
}

// 计算红包金额
func (g *goodsDomain) nextAmount(goods *RedEnvelopeGoods) (amount decimal.Decimal) {
	if goods.RemainQuantity == 1 {
		return goods.RemainAmount
	} else {
		if goods.EnvelopeType == services.GeneralEnvelopeType {
			return goods.AmountOne
		} else {
			cent := goods.RemainAmount.Mul(multiple).IntPart()
			next := algo.DoubleAverage(int64(g.RemainQuantity), cent)
			amount = decimal.NewFromFloat(float64(next)).Div(multiple)
		}
	}

	return amount
}

// 将抢到的红包金额从系统红包中间账户转入当前用户的资金账户
func (g *goodsDomain) transfer(ctx context.Context, dto services.RedEnvelopeReceiveDTO) (status services.TransferedStatus, err error) {
	systemAccount := base.GetSystemAccount()
	body := services.TradeParticipator{
		AccountNo: systemAccount.AccountNo,
		UserId:    systemAccount.UserId,
		Username:  systemAccount.Username,
	}
	target := services.TradeParticipator{
		AccountNo: dto.AccountNo,
		UserId:    dto.RecvUserId,
		Username:  dto.RecvUsername,
	}
	transfer := services.AccountTransferDTO{
		TradeNo:     dto.EnvelopeNo,
		TradeBody:   body,
		TradeTarget: target,
		Amount:      g.item.Amount,
		ChangeType:  services.EnvelopeIncoming,
		ChangeFlag:  services.FlagTransferIn,
		Decs:        "红包收入",
	}

	adomain := accounts.NewAccountDomain()
	return adomain.TransferWithContext(ctx, transfer)
}
