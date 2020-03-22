package envelopes

import (
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
	"github.com/ybinhome/envelope/services"
	"time"
)

// 定义红包商品 Dao 结构体
type RedEnvelopeGoodsDao struct {
	runner *dbx.TxRunner
}

// 1. 发红包 - 将红包商品信息插入数据库
func (dao *RedEnvelopeGoodsDao) Insert(po *RedEnvelopeGoods) (int64, error) {
	rs, err := dao.runner.Insert(po)
	if err != nil {
		return 0, err
	}
	return rs.LastInsertId()
}

// 2. 抢红包 - 更新红包剩余数量和剩余金额
//    使用乐观锁而非事务行锁来更新红包剩余数量和剩余金额，通过在 sql 语句中的 where 部分实现，
//    除了能够避免出现负库存的问题，还能够减少实际数据更新操作，过滤掉部分无效的更新，从而提高总体的性能。
func (dao *RedEnvelopeGoodsDao) UpdateBalance(envelopeNo string, amount decimal.Decimal) (int64, error) {
	sql := "update red_envelope_goods set remain_amount=remain_amount-CAST(? AS DECIMAL(30,6)), remain_quantity=remain_quantity-1 where envelope_no=? " +
		// 乐观锁的实现
		"and remain_quantity>0 and remain_amount >= CAST(? AS DECIMAL(30,6))"
	rs, err := dao.runner.Exec(sql, amount.String(), envelopeNo, amount.String())
	if err != nil {
		return 0, err
	}
	return rs.RowsAffected()
}

// 3. 查询红包 - 根据红包编号查询红包的状态
func (dao *RedEnvelopeGoodsDao) GetOne(envelopeNo string) *RedEnvelopeGoods {
	po := &RedEnvelopeGoods{EnvelopeNo: envelopeNo}
	ok, err := dao.runner.GetOne(po)
	if err != nil || !ok {
		logrus.Error(err)
		return nil
	}
	return po
}

// 4. 红包过期退款 - 更新红包订单状态
func (dao *RedEnvelopeGoodsDao) UpdateOrderStatus(envelopeNo string, status services.OrderStatus) (int64, error) {
	sql := "update red_envelope_goods set order_status=? where envelope_no=?"
	rs, err := dao.runner.Exec(sql, status, envelopeNo)
	if err != nil {
		return 0, err
	}
	return rs.RowsAffected()
}

// 5. 红包过期退款 - 包过期的所有红包都查询出来，过期红包数量过多时，需要使用分页功能，使用数据库的 limit 实现，通过 offset 和 size 来设定
func (dao *RedEnvelopeGoodsDao) FindExpired(offset, size int) []RedEnvelopeGoods {
	var goods []RedEnvelopeGoods
	now := time.Now()
	sql := "select * from red_envelope_goods where expired_at>? limit ?,?"
	err := dao.runner.Find(&goods, sql, now, offset, size)
	if err != nil {
		logrus.Error(err)
	}
	return goods
}
