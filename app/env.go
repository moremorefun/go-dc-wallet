package app

import (
	"go-dc-wallet/eosclient"
	"go-dc-wallet/ethclient"
	"go-dc-wallet/hcommon"
	"go-dc-wallet/omniclient"
	"math/rand"
	"strings"
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

	BtcNetworkType string `env:"BTC-NETWORK-TYPE" default:"btc"`

	EthRPC string `env:"ETH_RPC"`

	OmniRPCHost string `env:"OMNI_RPC_HOST"`
	OmniRPCUser string `env:"OMNI_RPC_USER"`
	OmniRPCPwd  string `env:"OMNI_RPC_PWD"`

	EosRPC string `env:"EOS_RPC"`
}

// Cfg
var Cfg *config

// EnvCreate 初始化运行环境
func EnvCreate() {
	// 初始化随机种子
	rand.Seed(time.Now().UnixNano())

	err := godotenv.Load()
	if err != nil {
		if !strings.Contains(err.Error(), "no such file or directory") {
			hcommon.Log.Fatalf("read env from .env err: [%T] %s", err, err.Error())
		} else {
			hcommon.Log.Warnf("no .env file specialty, set read from system or docker env")
		}
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
	ethclient.InitClient(Cfg.EthRPC)
	// 初始化omni rpc
	omniclient.InitClient(Cfg.OmniRPCHost, Cfg.OmniRPCUser, Cfg.OmniRPCPwd)
	// 初始化eos rpc
	eosclient.InitClient(Cfg.EosRPC)
}

// EnvDestroy 销毁运行环境
func EnvDestroy() {
	if DbCon != nil {
		_ = DbCon.Close()
	}
}
