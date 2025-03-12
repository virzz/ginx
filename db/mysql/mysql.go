package mysql

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gLogger "gorm.io/gorm/logger"

	"github.com/virzz/utils/once"
	"github.com/virzz/vlog"

	"github.com/virzz/ginx/db"
)

var (
	std      *gorm.DB
	oncePlus once.OncePlus
)

func R() *gorm.DB { return std }

func Migrate(models []any) error { return std.AutoMigrate(models...) }

func connect(cfg *db.Config) (err error) {
	newLogger := gLogger.Default.LogMode(gLogger.Info)
	if cfg.Debug {
		newLogger = gLogger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), gLogger.Config{
			SlowThreshold:             200 * time.Millisecond,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		})
	} else {
		f, err := os.OpenFile(filepath.Join("logs", "mysql.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			vlog.Warn("Failed to open gorm log file", "err", err.Error())
		} else {
			newLogger = gLogger.New(log.New(f, "\r\n", log.LstdFlags),
				gLogger.Config{LogLevel: gLogger.Warn, IgnoreRecordNotFoundError: true},
			)
		}
	}
	gormCfg := &gorm.Config{Logger: newLogger,
		QueryFields:                              true,
		DisableForeignKeyConstraintWhenMigrating: true,
		IgnoreRelationshipsWhenMigrating:         true,
	}
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Name,
	)
	std, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,  // string 类型字段的默认长度
		DontSupportRenameIndex:    true, // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		SkipInitializeWithVersion: true, // 根据当前 MySQL 版本自动配置
	}), gormCfg)
	if err != nil {
		vlog.Error("Failed to connect mysql", "err", err.Error())
		return err
	}
	sqlDB, err := std.DB()
	if err != nil {
		vlog.Warn("Failed to get sql.db", "err", err.Error())
	} else {
		sqlDB.SetMaxIdleConns(cfg.Conn.Idle)                                     // 最大空闲连接
		sqlDB.SetMaxOpenConns(cfg.Conn.Open)                                     // 最大连接数
		sqlDB.SetConnMaxLifetime(time.Duration(cfg.Conn.Lifetime) * time.Second) // 最大可复用时间
	}
	if cfg.Debug {
		std = std.Debug()
	}
	return nil
}

func Connect(cfg *db.Config, force ...bool) error {
	if len(force) > 0 && force[0] {
		return connect(cfg)
	}
	return oncePlus.Do(func() (err error) {
		return connect(cfg)
	})
}

func Init(cfg *db.Config, force ...bool) error {
	return Connect(cfg, force...)
}
