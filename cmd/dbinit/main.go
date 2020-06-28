package main

import (
	"context"
	"encoding/json"
	"go-dc-wallet/app"
	"go-dc-wallet/app/model"
	"go-dc-wallet/ethclient"
	"go-dc-wallet/hbtc"
	"go-dc-wallet/hcommon"
	"go-dc-wallet/heth"
	"go-dc-wallet/omniclient"
	"math"
	"net/http"
	"time"

	"github.com/parnurzeal/gorequest"
)

func main() {
	app.EnvCreate()
	defer app.EnvDestroy()

	// 1. 初始化 t_app_config_int
	configIntRows := []*model.DBTAppConfigInt{
		{
			// 最小可用剩余地址数
			K: "min_free_address",
			V: 1000,
		},
		{
			// eth 确认延迟数
			K: "block_confirm_num",
			V: 15,
		},
		{
			// erc20 默认转账 gas
			K: "erc20_gas_use",
			V: 90000,
		},
		{
			// btc 确认延迟数
			K: "btc_block_confirm_num",
			V: 2,
		},
	}
	_, err := model.SQLCreateIgnoreManyTAppConfigInt(
		context.Background(),
		app.DbCon,
		configIntRows,
	)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}

	// 2. 初始化 t_app_config_str
	// 创建可用地址
	ethAddresses, err := heth.CreateHotAddress(50)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
	btcAddresses, err := hbtc.CreateHotAddress(50)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
	configStrRows := []*model.DBTAppConfigStr{
		{
			// eth 冷钱包地址
			K: "cold_wallet_address",
			V: "",
		},
		{
			// eth 热钱包地址
			K: "hot_wallet_address",
			V: ethAddresses[0],
		},
		{
			// erc20 零钱整理手续费 热钱包地址
			K: "fee_wallet_address",
			V: ethAddresses[1],
		},
		{
			// erc20 零钱整理手续费 热钱包地址 列表
			K: "fee_wallet_address_list",
			V: ethAddresses[1],
		},
		{
			// btc 冷钱包地址
			K: "cold_wallet_address_btc",
			V: "",
		},
		{
			// btc 热钱包地址
			K: "hot_wallet_address_btc",
			V: btcAddresses[0],
		},
	}
	_, err = model.SQLCreateIgnoreManyTAppConfigStr(
		context.Background(),
		app.DbCon,
		configStrRows,
	)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}

	// 3. 初始化 t_app_config_token
	now := time.Now().Unix()
	configTokenRows := []*model.DBTAppConfigToken{
		{
			// erc20 token配置
			TokenAddress:  "0xdac17f958d2ee523a2206206994597c13d831ec7",
			TokenDecimals: 6,
			TokenSymbol:   "erc20_usdt",
			ColdAddress:   "",
			HotAddress:    ethAddresses[2],
			OrgMinBalance: "0.0",
			CreateTime:    now,
		},
	}
	_, err = model.SQLCreateIgnoreManyTAppConfigToken(
		context.Background(),
		app.DbCon,
		configTokenRows,
	)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}

	// 4. 初始化 t_app_config_token_btc
	configTokenBtcRows := []*model.DBTAppConfigTokenBtc{
		{
			// Omni token 配置
			TokenIndex:      31,
			TokenSymbol:     "omni_usdt",
			ColdAddress:     "",
			HotAddress:      btcAddresses[1],
			FeeAddress:      btcAddresses[2],
			TxOrgMinBalance: "0.0",
			CreateAt:        now,
		},
	}
	_, err = model.SQLCreateIgnoreManyTAppConfigTokenBtc(
		context.Background(),
		app.DbCon,
		configTokenBtcRows,
	)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}

	// 5. 初始化 t_app_status_int
	ethRpcBlockNum, err := ethclient.RpcBlockNumber(context.Background())
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
	btcRpcBlockNum, err := omniclient.RpcGetBlockCount()
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
	type StRespGasPrice struct {
		Fast        int64   `json:"fast"`
		Fastest     int64   `json:"fastest"`
		SafeLow     int64   `json:"safeLow"`
		Average     int64   `json:"average"`
		BlockTime   float64 `json:"block_time"`
		BlockNum    int64   `json:"blockNum"`
		Speed       float64 `json:"speed"`
		SafeLowWait float64 `json:"safeLowWait"`
		AvgWait     float64 `json:"avgWait"`
		FastWait    float64 `json:"fastWait"`
		FastestWait float64 `json:"fastestWait"`
	}
	gresp, body, errs := gorequest.New().
		Get("https://ethgasstation.info/api/ethgasAPI.json").
		Timeout(time.Second * 120).
		End()
	if errs != nil {
		hcommon.Log.Errorf("err: [%T] %s", errs[0], errs[0].Error())
		return
	}
	if gresp.StatusCode != http.StatusOK {
		// 状态错误
		hcommon.Log.Errorf("req status error: %d", gresp.StatusCode)
		return
	}
	var resp StRespGasPrice
	err = json.Unmarshal([]byte(body), &resp)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
	ethToUserGasPrice := resp.Fast * int64(math.Pow10(8))
	ethToColdGasPrice := resp.Average * int64(math.Pow10(8))
	type BtcStRespGasPrice struct {
		FastestFee  int64 `json:"fastestFee"`
		HalfHourFee int64 `json:"halfHourFee"`
		HourFee     int64 `json:"hourFee"`
	}
	gresp, body, errs = gorequest.New().
		Get("https://bitcoinfees.earn.com/api/v1/fees/recommended").
		Timeout(time.Second * 120).
		End()
	if errs != nil {
		hcommon.Log.Errorf("err: [%T] %s", errs[0], errs[0].Error())
		return
	}
	if gresp.StatusCode != http.StatusOK {
		// 状态错误
		hcommon.Log.Errorf("req status error: %d", gresp.StatusCode)
		return
	}
	var respBtc BtcStRespGasPrice
	err = json.Unmarshal([]byte(body), &respBtc)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
	btcToUserGasPrice := respBtc.FastestFee
	btcToColdGasPrice := respBtc.HalfHourFee

	appStatusIntRows := []*model.DBTAppStatusInt{
		{
			// eth blocknum
			K: "seek_num",
			V: ethRpcBlockNum,
		},
		{
			// eth blocknum
			K: "erc20_seek_num",
			V: ethRpcBlockNum,
		},
		{
			// btc blocknum
			K: "btc_seek_num",
			V: btcRpcBlockNum,
		},
		{
			// omni blocknum
			K: "omni_seek_num",
			V: btcRpcBlockNum,
		},
		{
			// eth 到冷钱包手续费
			K: "to_cold_gas_price",
			V: ethToColdGasPrice,
		},
		{
			// eth 到用户手续费
			K: "to_user_gas_price",
			V: ethToUserGasPrice,
		},
		{
			// btc 到冷钱包手续费
			K: "to_cold_gas_price_btc",
			V: btcToColdGasPrice,
		},
		{
			// btc 到用户手续费
			K: "to_user_gas_price_btc",
			V: btcToUserGasPrice,
		},
	}
	_, err = model.SQLCreateIgnoreManyTAppStatusInt(
		context.Background(),
		app.DbCon,
		appStatusIntRows,
	)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
}
