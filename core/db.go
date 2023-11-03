package core

import (
	"authentik-go/config"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	DB                *gorm.DB
	Session           gorm.Session
	ErrRecordNotFound = gorm.ErrRecordNotFound
	Expr              = gorm.Expr
)

// InDB 加载数据库
func InDB(str string) (*gorm.DB, error) {
	sp := strings.Split(str, "://")
	dbType := sp[0]
	dbPath := str
	if len(sp) == 2 {
		dbType = strings.ToLower(sp[0])
		dbPath = sp[1]
	}
	// 配置项
	dbConfig := &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",   // 表名前缀
			SingularTable: true, // true:单数 false:复数
		},
	}
	// 数据库类型
	if dbType == "mysql" {
		return gorm.Open(mysql.Open(dbPath), dbConfig)
	} else if dbType == "postgres" {
		return gorm.Open(postgres.Open(str), dbConfig)
	} else {
		return gorm.Open(sqlite.Open(dbPath), dbConfig)
	}
}

// CloseDB 关闭数据库
func CloseDB(db *gorm.DB) {
	if sqlDB, err := db.DB(); err == nil {
		_ = sqlDB.Close()
	}
}

// InitDB 初始化数据库
func InitDB() error {
	db, err := InDB(config.CONF.System.Dsn)
	if err != nil {
		return err
	}
	// defer CloseDB(db)
	DB = db
	return nil
}

// DBTableName 获取表名
func DBTableName(model interface{}) string {
	stmt := &gorm.Statement{DB: DB}
	stmt.Parse(model)
	return stmt.Schema.Table
}
