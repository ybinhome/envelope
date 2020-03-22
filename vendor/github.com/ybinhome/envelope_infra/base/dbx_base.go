package base

import (
	"context"
	"database/sql"
	"github.com/tietang/dbx"
	"log"
)

const TX = "tx"

//提供一个基本Dao基础功能的支持
type BaseDao struct {
	TX *sql.Tx
}

func (d *BaseDao) SetTx(tx *sql.Tx) {
	d.TX = tx

}

type txFunc func(*dbx.TxRunner) error

//事务执行帮助函数，简化代码
func Tx(fn func(*dbx.TxRunner) error) error {
	return TxContext(context.Background(), fn)
}

//事务执行帮助函数，简化代码，需要传入上下文
func TxContext(ctx context.Context, fn func(runner *dbx.TxRunner) error) error {
	return DbxDatabase().Tx(fn)
}

//将runner绑定到上下文，并创建一个新的WithValueContext
func WithValueContext(parent context.Context, runner *dbx.TxRunner) context.Context {
	return context.WithValue(parent, TX, runner)
}

//在事务上下文中执行事务逻辑
//传入绑定了runner的上下文，并执行事务函数代码
//函数只能用在绑定了runner的事务上下文中，也就是说用在事务函数内部，和WithValueContext配合一起来完成。
// 举例：
/*
func CreateXyz() error {
	err := Tx(func(runner *dbx.TxRunner) error {
		//将runner绑定到上下文
		ctx := WithValueContext(context.Background(), runner)
		//dao.insert xxx
		TxStepY(ctx)
		TxStepZ(ctx)
		return nil
	})

	return err
}

func TxStepY(ctx context.Context) error {
	//事务上下文中执行Y数据库操作
	return ExecuteContext(ctx, func(runner *dbx.TxRunner) error {
		//dao.insert yyy
		return nil
	})

}

func TxStepZ(ctx context.Context) error {
	//事务上下文中执行Z数据库操作
	return ExecuteContext(ctx, func(runner *dbx.TxRunner) error {
		//dao.update zzz
		return nil
	})

}
*/
func ExecuteContext(ctx context.Context, fn func(*dbx.TxRunner) error) error {
	tx, ok := ctx.Value(TX).(*dbx.TxRunner)
	if !ok || tx == nil {
		log.Panic("是否在事务函数块中使用？")
	}
	return fn(tx)
}
