package accounts

import (
	"database/sql"
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/smartystreets/assertions/should"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tietang/dbx"
	"github.com/ybinhome/envelope/infra/base"
	_ "github.com/ybinhome/envelope/test_brun"
	"testing"
)

func TestAccountDao_GetOne(t *testing.T) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := &AccountDao{
			runner: runner,
		}
		Convey("通过编号查询账户数据", t, func() {
			a := &Account{
				Balance:     decimal.NewFromFloat(100),
				Status:      1,
				AccountNo:   ksuid.New().Next().String(),
				AccountName: "测试资金账户",
				UserId:      ksuid.New().Next().String(),
				Username:    sql.NullString{String: "测试用户", Valid: true},
			}
			id, err := dao.Insert(a)
			So(err, ShouldBeNil)
			So(id, should.BeGreaterThan, 0)
			na := dao.GetOne(a.AccountNo)
			So(na, ShouldNotBeNil)
			So(na.Balance.String(), ShouldEqual, a.Balance.String())
			So(na.CreatedAt, ShouldNotBeNil)
			So(na.UpdatedAt, ShouldNotBeNil)
		})
		return nil
	})

	if err != nil {
		logrus.Error(err)
	}
}

func TestAccountDao_GetByUserId(t *testing.T) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := &AccountDao{
			runner: runner,
		}
		Convey("通过用户 id 和账户类型查询账户数据", t, func() {
			a := &Account{
				Balance:     decimal.NewFromFloat(100),
				Status:      1,
				AccountNo:   ksuid.New().Next().String(),
				AccountName: "测试资金账户",
				UserId:      ksuid.New().Next().String(),
				Username:    sql.NullString{String: "测试用户", Valid: true},
				AccountType: 2,
			}
			id, err := dao.Insert(a)
			So(err, ShouldBeNil)
			So(id, should.BeGreaterThan, 0)
			na := dao.GetByUserId(a.UserId, a.AccountType)
			So(na, ShouldNotBeNil)
			So(na.Balance.String(), ShouldEqual, a.Balance.String())
			So(na.CreatedAt, ShouldNotBeNil)
			So(na.UpdatedAt, ShouldNotBeNil)
		})
		return nil
	})

	if err != nil {
		logrus.Error(err)
	}
}

func TestAccountDao_UpdateBalance(t *testing.T) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := &AccountDao{
			runner: runner,
		}
		balance := decimal.NewFromFloat(100)
		Convey("更新账户余额", t, func() {
			a := &Account{
				Balance:     balance,
				Status:      1,
				AccountNo:   ksuid.New().Next().String(),
				AccountName: "测试资金账户",
				UserId:      ksuid.New().Next().String(),
				Username:    sql.NullString{String: "测试用户", Valid: true},
			}
			id, err := dao.Insert(a)
			So(err, ShouldBeNil)
			So(id, should.BeGreaterThan, 0)

			// 1. 增加余额
			Convey("增加余额", func() {
				amount := decimal.NewFromFloat(10)
				rows, err := dao.UpdateBalance(a.AccountNo, amount)
				So(err, ShouldBeNil)
				So(rows, ShouldEqual, 1)
				na := dao.GetOne(a.AccountNo)
				newBalance := balance.Add(amount)
				So(na, ShouldNotBeNil)
				So(na.Balance.String(), ShouldEqual, newBalance.String())
			})

			// 2. 扣减余额，余额足够
			Convey("扣减余额，余额足够", func() {
				a1 := dao.GetOne(a.AccountNo)
				So(a1, ShouldNotBeNil)
				amount := decimal.NewFromFloat(-50)
				rows, err := dao.UpdateBalance(a.AccountNo, amount)
				So(err, ShouldBeNil)
				So(rows, ShouldEqual, 1)
				a2 := dao.GetOne(a.AccountNo)
				So(a2, ShouldNotBeNil)
				So(a1.Balance.Add(amount).String(), ShouldEqual, a2.Balance.String())
			})

			// 3. 扣减余额，余额不够
			Convey("扣减余额，余额不够", func() {
				a1 := dao.GetOne(a.AccountNo)
				So(a1, ShouldNotBeNil)
				amount := decimal.NewFromFloat(-300)
				rows, err := dao.UpdateBalance(a.AccountNo, amount)
				So(err, ShouldBeNil)
				So(rows, ShouldEqual, 0)
				a2 := dao.GetOne(a.AccountNo)
				So(a2, ShouldNotBeNil)
				So(a1.Balance.String(), ShouldEqual, a2.Balance.String())
			})
		})
		return nil
	})

	if err != nil {
		logrus.Error(err)
	}
}
