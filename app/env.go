package app

import (
	"go-dc-wallet/hcommon"
	"math/rand"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/timest/env"
)

// DbCon 数据库链接
var DbCon *sqlx.DB

// 加密密钥
var AESKey string

// EnvCreate 初始化运行环境
func EnvCreate() {
	// 初始化随机种子
	rand.Seed(time.Now().UnixNano())
	// 读取配置
	type config struct {
		IsDebug bool `env:"IS-DEBUG" default:"false"`

		MySqlDataSourceName string `env:"MYSQL"`
		MySqlIsShowSQL      bool   `env:"MYSQL-IS-SHOW-SQL" default:"false"`

		AESKey string `env:"AES-KEY"`
	}
	err := godotenv.Load()
	if err != nil {
		hcommon.Log.Fatalf("read env from .env err: [%T] %s", err, err.Error())
	}
	cfg := new(config)
	env.IgnorePrefix()
	err = env.Fill(cfg)
	if err != nil {
		hcommon.Log.Fatalf("read env config err: [%T] %s", err, err.Error())
	}
	if len(cfg.MySqlDataSourceName) == 0 {
		hcommon.Log.Fatalf("mysql dataSourceName is empty")
	}
	AESKey = cfg.AESKey

	// 初始化数据库
	DbCon = hcommon.DbCreate(cfg.MySqlDataSourceName, cfg.MySqlIsShowSQL)
}

// EnvDestroy 销毁运行环境
func EnvDestroy() {
	if DbCon != nil {
		_ = DbCon.Close()
	}
}
