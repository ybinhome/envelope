package accounts

import (
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
)

// 在数据库访问层，每一个方法代表一个元子性操作，对于数据库的事务，通过外面控制，每一个事务可以使用多个元子操作组合，因此使用 TxRunner

type AccountDao struct {
	runner *dbx.TxRunner
}

// 查询数据库持久化对象的单实例，获取一行数据
func (dao *AccountDao) GetOne(accountNo string) *Account {
	// 将 accountNo 实例化到 Account 对象中，传递给 runner.GetOne 方法，GetOne 方法将返回结果存储在 Account 对象中
	a := &Account{AccountNo: accountNo}
	ok, err := dao.runner.GetOne(a)

	// 如果查询出错，返回空的 Account 对象，并且打印错误到日志输出
	if err != nil {
		logrus.Error(err)
		return nil
	}

	// 如果数据不存在，返回空的 Account 对象
	if !ok {
		return nil
	}

	// 查询成功，将 Account 对象返回出去
	return a
}

// 通过用户 id 和账户类型来查询账户信息
func (dao *AccountDao) GetByUserId(userId string, accountType int) *Account {
	// 定义一个空的 Account 对象，用于接受查询结果
	a := &Account{}

	// 定义查询的 sql 语句
	sql := "select * from account where user_id=? and account_type=?"

	// 使用 runner.Get 方法查询结果
	ok, err := dao.runner.Get(a, sql, userId, accountType)

	// 如果查询出错，返回空的 Account 对象，并且打印错误到日志输出
	if err != nil {
		logrus.Error(err)
		return nil
	}

	// 如果数据不存在，返回空的 Account 对象
	if !ok {
		return nil
	}

	// 查询成功，将 Account 对象返回出去
	return a
}

// 账户数据的插入
func (dao *AccountDao) Insert(a *Account) (id int64, err error) {
	// 直接调用 runner.Insert 方法将 Account 对象插入即可，dbx 的 ORM 会自动将结构体转换为 sql 语句中的表和字段
	rs, err := dao.runner.Insert(a)

	// 如果插入发生错误，将自增 id 设置为 0，然后将自增 id 和 err 信息返回出去
	if err != nil {
		return 0, err
	}

	// 如果插入成功，直接插入结果 sql.Result 的 LastInsertId 返回出去，LastInsertId 包含两个值，第一个是自增 id，第二个是错误信息
	return rs.LastInsertId()
}

// 账户余额的更新，返回更新影响行数和错误信息
func (dao *AccountDao) UpdateBalance(accountNo string, amount decimal.Decimal) (rows int64, err error) {
	// 通过 runner.Exec 方法更新账户余额，amount 为 decimal 类型，数据库语句中不能直接使用，因此我们在 Exec 的参数中，转换为 string 类型，
	//    然后在 sql 语句中，使用 CAST(? AS DECIMAL(30,6)，转换为 30 位长度 6 位精度的数值
	sql := "update account set balance=balance+CAST(? AS DECIMAL(30,6)) where account_no=? and balance>=-1*CAST(? AS DECIMAL(30,6))"
	rs, err := dao.runner.Exec(sql, amount.String(), accountNo, amount.String())

	// 如果发生错误，返回影响行数为 0 和错误信息
	if err != nil {
		return 0, err
	}

	// 如果插入成功，直接插入结果 sql.Result 的 RowsAffected 返回出去，RowsAffected 包含两个值，第一个是影响行数，第二个是错误信息
	return rs.RowsAffected()
}

// 账户状态更新，返回影响行数和错误信息
func (dao *AccountDao) UpdateStatus(accountNo string, status int) (rows int64, err error) {
	sql := "update account set status=? where account_no=?"
	rs, err := dao.runner.Exec(sql, status, accountNo)
	if err != nil {
		return 0, err
	}
	return rs.RowsAffected()
}
