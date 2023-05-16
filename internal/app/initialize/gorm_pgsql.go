package initialize

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"nft-collect/internal/app/global"
	"nft-collect/internal/app/initialize/internal"
)

// GormPgSql 初始化 Postgresql 数据库
func GormPgSql(Prefix string) *gorm.DB {
	p := global.CONFIG.Pgsql
	if p.Dbname == "" {
		return nil
	}
	pgsqlConfig := postgres.Config{
		DSN:                  p.Dsn(), // DSN data source name
		PreferSimpleProtocol: false,
	}
	if db, err := gorm.Open(postgres.New(pgsqlConfig), internal.Gorm.Config(Prefix)); err != nil {
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(p.MaxIdleConns)
		sqlDB.SetMaxOpenConns(p.MaxOpenConns)
		return db
	}
}
