package base

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
	"github.com/tietang/props/kvs"
	"github.com/ybinhome/envelope_infra"
)

// dbx 数据库实例
var database *dbx.Database

func DbxDatabase() *dbx.Database {
	return database
}

// dbx 数据库 starter，并设置为全局
type DbxDatabaseStarter struct {
	infra.BaseStarter
}

// 数据库连接的生命周期要晚于配置文件的加载，为了更好的符合此逻辑，将连接的初始化设置在 Setup 阶段
func (s *DbxDatabaseStarter) Setup(ctx infra.StarterContext) {
	// 获取配置
	config := ctx.Props()

	// 创建 dbx 配置对象
	settings := dbx.Settings{}

	// 利用 props 的 Unmarshal 功能，直接将 config 中的内容解析到结构体中
	err := kvs.Unmarshal(config, &settings, "mysql")
	if err != nil {
		panic(err)
	}

	// 为了更好地 debug，输出一些 log 信息
	logrus.Infof("%+v", settings)
	logrus.Info("mysql.conn url: ", settings.ShortDataSourceName())

	// 实例化 dbx 的 DB 对象
	dbxdb, err := dbx.Open(settings)
	if err != nil {
		panic(err)
	}

	// 输出数据库连接的将康状况
	logrus.Info(dbxdb.Ping())

	database = dbxdb
}
