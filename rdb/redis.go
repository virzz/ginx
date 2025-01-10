package rdb

import (
	"context"
	"fmt"
	"net"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/pflag"

	"github.com/virzz/utils/once"
	"github.com/virzz/vlog"
)

func FlagSet(name string) *pflag.FlagSet {
	fs := pflag.NewFlagSet("rdb", pflag.ContinueOnError)
	fs.Bool("redis.debug", false, "Database Debug Mode")
	fs.String("redis.host", "127.0.0.1", "Database Host")
	fs.Int("redis.port", 6379, "Database Port")
	fs.Int("redis.db", 0, "Database User")
	fs.String("redis.pass", "", "Database Password")
	return fs
}

type Config struct {
	Debug bool   `json:"debug" yaml:"debug"`
	Host  string `json:"host" yaml:"host"`
	Port  int    `json:"port" yaml:"port"`
	DB    int    `json:"db" yaml:"db"`
	Pass  string `json:"pass" yaml:"pass"`
}

type DebugHook struct{}

// 当创建网络连接时调用
func (DebugHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return next(ctx, network, addr)
	}
}

// 执行命令时调用
func (DebugHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		vlog.Debug(cmd.String())
		next(ctx, cmd)
		return nil
	}
}

// 执行管道命令时调用
func (DebugHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		return next(ctx, cmds)
	}
}

var (
	rdb      *redis.Client
	oncePlus once.OncePlus
	Nil      = redis.Nil
)

func connect(cfg *Config) error {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Pass,
		DB:       cfg.DB,
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			vlog.Info("Redis is connected")
			return nil
		},
	})
	if cfg.Debug {
		rdb.AddHook(DebugHook{})
	}
	return rdb.Ping(context.Background()).Err()
}

func Connect(cfg *Config, force ...bool) error {
	if len(force) > 0 && force[0] {
		return connect(cfg)
	}
	return oncePlus.Do(func() (err error) {
		return connect(cfg)
	})
}

func Init(cfg *Config, force ...bool) error {
	return Connect(cfg, force...)
}

func R() *redis.Client { return rdb }
