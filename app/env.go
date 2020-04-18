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

// 读取配置
type config struct {
	IsDebug bool `env:"IS-DEBUG" default:"false"`

	MySqlDataSourceName string `env:"MYSQL"`
	MySqlIsShowSQL      bool   `env:"MYSQL-IS-SHOW-SQL" default:"false"`

	AESKey string `env:"AES-KEY"`

	EthRPC string `env:"ETH_RPC"`
}

// Cfg
var Cfg *config

// EnvCreate 初始化运行环境
func EnvCreate() {
	// 初始化随机种子
	rand.Seed(time.Now().UnixNano())

	err := godotenv.Load()
	if err != nil {
		hcommon.Log.Fatalf("read env from .env err: [%T] %s", err, err.Error())
	}
	Cfg = new(config)
	env.IgnorePrefix()
	err = env.Fill(Cfg)
	if err != nil {
		hcommon.Log.Fatalf("read env config err: [%T] %s", err, err.Error())
	}
	if len(Cfg.MySqlDataSourceName) == 0 {
		hcommon.Log.Fatalf("mysql dataSourceName is empty")
	}

	// 初始化数据库
	DbCon = hcommon.DbCreate(Cfg.MySqlDataSourceName, Cfg.MySqlIsShowSQL)
	// 初始化eth rpc
	hcommon.EthInitClient(Cfg.EthRPC)
}

// EnvDestroy 销毁运行环境
func EnvDestroy() {
	if DbCon != nil {
		_ = DbCon.Close()
	}
}
