package pgsql

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gorm.io/driver/postgres"
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

func R() *gorm.DB {
	if std == nil {
		panic("pgsql not init")
	}
	return std
}

func Migrate(models ...any) error { return std.AutoMigrate(models...) }

func connect(cfg *db.Config) (err error) {
	newLogger := gLogger.Default.LogMode(gLogger.Info)
	if cfg.Debug {
		newLogger = gLogger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), gLogger.Config{
			SlowThreshold:             200 * time.Millisecond,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		})
	} else {
		f, err := os.OpenFile(filepath.Join("logs", "pgsql.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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
	dsnList := []string{}
	if cfg.Host != "" {
		dsnList = append(dsnList, "host="+cfg.Host)
	}
	if cfg.Port != 0 {
		dsnList = append(dsnList, "port="+strconv.Itoa(cfg.Port))
	}
	if cfg.User != "" {
		dsnList = append(dsnList, "user="+cfg.User)
	}
	if cfg.Pass != "" {
		dsnList = append(dsnList, "password="+cfg.Pass)
	}
	if cfg.Name != "" {
		dsnList = append(dsnList, "dbname="+cfg.Name)
	}
	if len(dsnList) > 0 {
		dsnList = append(dsnList, "sslmode=disable", "TimeZone=Asia/Shanghai")
	}
	dsn := strings.Join(dsnList, " ")
	if cfg.Debug {
		vlog.Info("Connecting to postgres", "dsn", dsn)
	}
	std, err = gorm.Open(postgres.New(postgres.Config{DSN: dsn, PreferSimpleProtocol: true}), gormCfg)
	if err != nil {
		vlog.Error("Failed to connect postgres", "err", err.Error())
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
