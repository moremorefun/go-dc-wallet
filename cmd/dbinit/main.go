package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"go-dc-wallet/app"
	"go-dc-wallet/eosclient"
	"go-dc-wallet/ethclient"
	"go-dc-wallet/hbtc"
	"go-dc-wallet/heth"
	"go-dc-wallet/model"
	"go-dc-wallet/omniclient"
	"go-dc-wallet/xenv"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/moremorefun/mcommon"
	"github.com/parnurzeal/gorequest"
)

func main() {
	xenv.EnvCreate()
	defer xenv.EnvDestroy()

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
	_, err := model.SQLCreateManyTAppConfigInt(
		context.Background(),
		xenv.DbCon,
		configIntRows,
		true,
	)
	if err != nil {
		mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}

	// 2. 初始化 t_app_config_str
	// 获取可用地址
	ethAddressRows, err := app.SQLSelectTAddressKeyColByTagAndSymbol(
		context.Background(),
		xenv.DbCon,
		[]string{
			model.DBColTAddressKeyAddress,
		},
		-1,
		heth.CoinSymbol,
	)
	if err != nil {
		mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
	var ethAddresses []string
	for _, ethAddressRow := range ethAddressRows {
		ethAddresses = append(ethAddresses, ethAddressRow.Address)
	}
	if len(ethAddresses) < 10 {
		// 创建可用地址
		ethAddresses, err = heth.CreateHotAddress(50)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
	}
	// 获取可用地址
	btcAddressRows, err := app.SQLSelectTAddressKeyColByTagAndSymbol(
		context.Background(),
		xenv.DbCon,
		[]string{
			model.DBColTAddressKeyAddress,
		},
		-1,
		hbtc.CoinSymbol,
	)
	if err != nil {
		mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
	var btcAddresses []string
	for _, btcAddressRow := range btcAddressRows {
		btcAddresses = append(btcAddresses, btcAddressRow.Address)
	}
	if len(btcAddresses) < 10 {
		btcAddresses, err = hbtc.CreateHotAddress(50)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
	}
	if ethAddresses == nil {
		mcommon.Log.Errorf("ethAddresses nil")
		return
	}
	if btcAddresses == nil {
		mcommon.Log.Errorf("btcAddresses nil")
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
			V: "",
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
		{
			// eos 冷钱包地址
			K: "cold_wallet_address_eos",
			V: "",
		},
		{
			// eos 热钱包地址
			K: "hot_wallet_address_eos",
			V: "",
		},
		{
			// eos 热钱包加密私钥
			K: "hot_wallet_key_eos",
			V: "",
		},
	}
	_, err = model.SQLCreateManyTAppConfigStr(
		context.Background(),
		xenv.DbCon,
		configStrRows,
		true,
	)
	if err != nil {
		mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
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
	_, err = model.SQLCreateManyTAppConfigToken(
		context.Background(),
		xenv.DbCon,
		configTokenRows,
		true,
	)
	if err != nil {
		mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
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
	_, err = model.SQLCreateManyTAppConfigTokenBtc(
		context.Background(),
		xenv.DbCon,
		configTokenBtcRows,
		true,
	)
	if err != nil {
		mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}

	// 5. 初始化 t_app_status_int
	ethRpcBlockNum, err := ethclient.RpcBlockNumber(context.Background())
	if err != nil {
		mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
	btcRpcBlockNum, err := omniclient.RpcGetBlockCount()
	if err != nil {
		mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
	rpcChainInfo, err := eosclient.RpcChainGetInfo()
	if err != nil {
		mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
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
		mcommon.Log.Errorf("err: [%T] %s", errs[0], errs[0].Error())
		return
	}
	if gresp.StatusCode != http.StatusOK {
		// 状态错误
		mcommon.Log.Errorf("req status error: %d", gresp.StatusCode)
		return
	}
	var resp StRespGasPrice
	err = json.Unmarshal([]byte(body), &resp)
	if err != nil {
		mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
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
		mcommon.Log.Errorf("err: [%T] %s", errs[0], errs[0].Error())
		return
	}
	if gresp.StatusCode != http.StatusOK {
		// 状态错误
		mcommon.Log.Errorf("req status error: %d", gresp.StatusCode)
		return
	}
	var respBtc BtcStRespGasPrice
	err = json.Unmarshal([]byte(body), &respBtc)
	if err != nil {
		mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
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
			// eos blocknum
			K: "eos_seek_num",
			V: rpcChainInfo.LastIrreversibleBlockNum,
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
	_, err = model.SQLCreateManyTAppStatusInt(
		context.Background(),
		xenv.DbCon,
		appStatusIntRows,
		true,
	)
	if err != nil {
		mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}

	// 6. 更新 t_app_config_str
	feeAddressValue, err := app.SQLGetTAppConfigStrValueByK(
		context.Background(),
		xenv.DbCon,
		"fee_wallet_address",
	)
	if err != nil && err != sql.ErrNoRows {
		mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
	feeAddressValue = strings.TrimSpace(feeAddressValue)
	feeAddressListValue, err := app.SQLGetTAppConfigStrValueByK(
		context.Background(),
		xenv.DbCon,
		"fee_wallet_address_list",
	)
	if err != nil && err != sql.ErrNoRows {
		mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
	feeAddressListValue = strings.TrimSpace(feeAddressListValue)
	if feeAddressValue != "" && !strings.Contains(feeAddressListValue, feeAddressValue) {
		if feeAddressListValue == "" {
			feeAddressListValue = feeAddressValue
		} else {
			feeAddressListValue += fmt.Sprintf(",%s", feeAddressValue)
		}
	}
	_, err = app.SQLUpdateTAppConfigStrByK(
		context.Background(),
		xenv.DbCon,
		&model.DBTAppConfigStr{
			K: "fee_wallet_address_list",
			V: feeAddressListValue,
		},
	)
	if err != nil && err != sql.ErrNoRows {
		mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
}
