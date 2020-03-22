package envelopes

import (
	"context"
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
	"github.com/ybinhome/envelope/infra/base"
	"github.com/ybinhome/envelope/services"
	"time"
)

type goodsDomain struct {
	RedEnvelopeGoods
	item itemDomain
}

// 发红包 - 创建一个红包商品对象
func (g *goodsDomain) createEnvelopeNo() {
	g.EnvelopeNo = ksuid.New().Next().String()
}

// 发红包 - 生成一个红包编号
func (g *goodsDomain) Create(goods services.RedEnvelopeGoodsDTO) {
	g.RedEnvelopeGoods.FromDTO(&goods)
	g.RemainQuantity = goods.Quantity
	g.Username.Valid = true
	g.Blessing.Valid = true
	if g.EnvelopeType == services.GeneralEnvelopeType {
		g.Amount = goods.AmountOne.Mul(decimal.NewFromFloat(float64(goods.Quantity)))
	} else {
		g.AmountOne = decimal.NewFromFloat(0)
	}
	g.RemainAmount = g.Amount
	g.ExpiredAt = time.Now().Add(24 * time.Hour)
	g.Status = services.OrderCreate
	g.createEnvelopeNo()
}

// 发红包 - 保存到红包商品表 (上下文用于红包支付事务的实现)
func (g *goodsDomain) Save(ctx context.Context) (id int64, err error) {
	err = base.ExecuteContext(ctx, func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeGoodsDao{runner: runner}
		id, err = dao.Insert(&g.RedEnvelopeGoods)
		return nil
	})
	return id, err
}

// 组合红包创建和保存方法 (上下文用于红包支付事务的实现)
func (g *goodsDomain) CreateAndSave(ctx context.Context, goods services.RedEnvelopeGoodsDTO) (id int64, err error) {
	// 创建红包商品
	g.Create(goods)
	// 保存红包商品
	return g.Save(ctx)
}

// 查询红包商品信息
func (g *goodsDomain) Get(envelopeNo string) (goods *RedEnvelopeGoods) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeGoodsDao{runner: runner}
		goods = dao.GetOne(envelopeNo)
		return nil
	})
	if err != nil {
		logrus.Error(err)
	}
	return goods
}
