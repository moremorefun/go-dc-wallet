package xenv

import (
	"go-dc-wallet/eosclient"
	"go-dc-wallet/ethclient"
	"go-dc-wallet/omniclient"
	"math/rand"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/moremorefun/mcommon"
	"github.com/timest/env"
)

// DbCon 数据库链接
var DbCon *sqlx.DB

// 读取配置
type config struct {
	IsDebug bool `env:"IS-DEBUG" default:"false"`

	MySqlDataSourceName string `env:"MYSQL"`
	MySqlIsShowSQL      bool   `env:"MYSQL-IS-SHOW-SQL" default:"false"`

	Proxy string `env:"PROXY"`

	AESKey string `env:"AES-KEY"`

	BtcEnable      bool   `env:"BTC_ENABLE" default:"true"`
	BtcNetworkType string `env:"BTC-NETWORK-TYPE" default:"btc"`

	EthEnable bool   `env:"ETH_ENABLE" default:"true"`
	EthRPC    string `env:"ETH_RPC"`

	OmniRPCHost string `env:"OMNI_RPC_HOST"`
	OmniRPCUser string `env:"OMNI_RPC_USER"`
	OmniRPCPwd  string `env:"OMNI_RPC_PWD"`

	EosRPC    string `env:"EOS_RPC"`
	EosEnable bool   `env:"EOS_ENABLE" default:"true"`
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
			mcommon.Log.Fatalf("read env from .env err: [%T] %s", err, err.Error())
		} else {
			mcommon.Log.Warnf("no .env file specialty, set read from system or docker env")
		}
	}
	Cfg = new(config)
	env.IgnorePrefix()
	err = env.Fill(Cfg)
	if err != nil {
		mcommon.Log.Fatalf("read env config err: [%T] %s", err, err.Error())
	}
	if len(Cfg.MySqlDataSourceName) == 0 {
		mcommon.Log.Fatalf("mysql dataSourceName is empty")
	}

	// 初始化数据库
	DbCon = mcommon.DbCreate(Cfg.MySqlDataSourceName, Cfg.MySqlIsShowSQL)
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
