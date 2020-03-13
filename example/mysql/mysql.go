package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func main() {
	// 创建 sql.DB 对象
	dsName := "po:111111@tcp(127.0.0.1:3306)/po?charset=utf8&parseTime=true&loc=Local"
	db, err := sql.Open("mysql", dsName)
	if err != nil {
		fmt.Println(err)
	}

	// 设置 sql.DB 对象连接池属性
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(7 * time.Hour)

	// 使用数据库连接
	fmt.Println(db.Query("select now()"))

	// 使用 defer 函数关闭 sql.DB 资源
	defer db.Close()
}
