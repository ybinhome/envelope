package envelopes

import (
	"context"
	"github.com/tietang/dbx"
	"github.com/ybinhome/envelope/core/accounts"
	"github.com/ybinhome/envelope/infra/base"
	"github.com/ybinhome/envelope/services"
	"path"
)

// 发送红包业务领域代码
func (g *goodsDomain) SendOut(goods services.RedEnvelopeGoodsDTO) (activity *services.RedEnvelopeActivity, err error) {
	// 1. 创建红包商品
	g.Create(goods)

	// 2. 创建红包活动
	//    实例化红包活动对象
	activity = new(services.RedEnvelopeActivity)
	//    获取红包链接和域名配置
	link := base.GetEnvelopeActivityLink()
	domain := base.GetEnvelopeDomain()
	//    拼接活动完整路径
	activity.Link = path.Join(domain, link, g.EnvelopeNo)

	accountDomain := accounts.NewAccountDomain()
	// 3. 完成发红包
	err = base.Tx(func(runner *dbx.TxRunner) (err error) {
		ctx := base.WithValueContext(context.Background(), runner)
		// (1) 保存红包商品
		id, err := g.Save(ctx)
		if err != nil || id <= 0 {
			return err
		}

		// (2) 支付红包金额：此处需要将保存红包商品和支付红包金额放在同一个事务中完成，两件事情同时失败或者同时成功
		//    需要一个红包中间商的资金账户，此处定义在配置文件中，实现初始化到资金账户表中
		//    从红包发送人的资金账户中扣减红包金额
		//    将扣减的红包金额转入红包中间商的资金账户
		body := services.TradeParticipator{
			AccountNo: goods.AccountNo,
			UserId:    goods.UserId,
			Username:  goods.Username,
		}
		systemAccount := base.GetSystemAccount()
		target := services.TradeParticipator{
			AccountNo: systemAccount.AccountNo,
			UserId:    systemAccount.UserId,
			Username:  systemAccount.Username,
		}

		// 出账
		transfer := services.AccountTransferDTO{
			TradeBody:   body,
			TradeTarget: target,
			TradeNo:     g.EnvelopeNo,
			Amount:      g.Amount,
			ChangeType:  services.EnvelopeOutgoing,
			ChangeFlag:  services.FlagTransferOut,
			Decs:        "红包金额支付",
		}
		status, err := accountDomain.TransferWithContext(ctx, transfer)
		if status == services.TransferedStatusSuccess {
			return nil
		}

		// 入账
		transfer = services.AccountTransferDTO{
			TradeBody:   target,
			TradeTarget: body,
			TradeNo:     g.EnvelopeNo,
			Amount:      g.Amount,
			ChangeType:  services.EnvelopeIncoming,
			ChangeFlag:  services.FlagTransferIn,
			Decs:        "红包金额入账",
		}
		status, err = accountDomain.TransferWithContext(ctx, transfer)
		if status == services.TransferedStatusSuccess {
			return nil
		}

		return err
	})
	if err != nil {
		return nil, err
	}

	// 4. 扣减金额完成后，返回活动
	activity.RedEnvelopeGoodsDTO = *g.RedEnvelopeGoods.ToDTO()

	return activity, err
}
