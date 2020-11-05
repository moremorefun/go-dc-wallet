// 定时处理检测任务
package main

import (
	"go-dc-wallet/app"
	"go-dc-wallet/hbtc"
	"go-dc-wallet/heos"
	"go-dc-wallet/heth"
	"go-dc-wallet/xenv"

	"github.com/moremorefun/mcommon"
	"github.com/robfig/cron/v3"
)

func main() {
	xenv.EnvCreate()
	defer xenv.EnvDestroy()

	c := cron.New(
		cron.WithSeconds(),
		cron.WithChain(
			cron.Recover(cron.DefaultLogger),
		),
	)
	var err error
	// --- common --
	// 检测 通知发送
	_, err = c.AddFunc("@every 1m", app.CheckDoNotify)
	if err != nil {
		mcommon.Log.Errorf("cron add func error: %#v", err)
	}
	// --- eth ---
	// 检测 eth 生成地址
	_, err = c.AddFunc("@every 1m", heth.CheckAddressFree)
	if err != nil {
		mcommon.Log.Errorf("cron add func error: %#v", err)
	}
	// 检测 eth 冲币
	_, err = c.AddFunc("@every 5s", heth.CheckBlockSeek)
	if err != nil {
		mcommon.Log.Errorf("cron add func error: %#v", err)
	}
	// 检测 eth 零钱整理
	_, err = c.AddFunc("@every 10m", heth.CheckAddressOrg)
	if err != nil {
		mcommon.Log.Errorf("cron add func error: %#v", err)
	}
	// 检测 eth 提币
	_, err = c.AddFunc("@every 3m", heth.CheckWithdraw)
	if err != nil {
		mcommon.Log.Errorf("cron add func error: %#v", err)
	}
	// 检测 eth 发送交易
	_, err = c.AddFunc("@every 1m", heth.CheckRawTxSend)
	if err != nil {
		mcommon.Log.Errorf("cron add func error: %#v", err)
	}
	// 检测 eth 交易上链
	_, err = c.AddFunc("@every 5s", heth.CheckRawTxConfirm)
	if err != nil {
		mcommon.Log.Errorf("cron add func error: %#v", err)
	}
	// 检测 eth 通知到账
	_, err = c.AddFunc("@every 5s", heth.CheckTxNotify)
	if err != nil {
		mcommon.Log.Errorf("cron add func error: %#v", err)
	}
	// 检测 eth gas price
	_, err = c.AddFunc("@every 2m", heth.CheckGasPrice)
	if err != nil {
		mcommon.Log.Errorf("cron add func error: %#v", err)
	}

	// --- erc20 ---
	// 检测 erc20 冲币
	_, err = c.AddFunc("@every 5s", heth.CheckErc20BlockSeek)
	if err != nil {
		mcommon.Log.Errorf("cron add func error: %#v", err)
	}
	// 检测 erc20 通知到账
	_, err = c.AddFunc("@every 5s", heth.CheckErc20TxNotify)
	if err != nil {
		mcommon.Log.Errorf("cron add func error: %#v", err)
	}
	// 检测 erc20 零钱整理
	_, err = c.AddFunc("@every 10m", heth.CheckErc20TxOrg)
	if err != nil {
		mcommon.Log.Errorf("cron add func error: %#v", err)
	}
	// 检测 erc20 提币
	_, err = c.AddFunc("@every 3m", heth.CheckErc20Withdraw)
	if err != nil {
		mcommon.Log.Errorf("cron add func error: %#v", err)
	}

	// --- btc ---
	// 检测 btc 生成地址
	_, err = c.AddFunc("@every 1m", hbtc.CheckAddressFree)
	if err != nil {
		mcommon.Log.Errorf("cron add func error: %#v", err)
	}
	// 检测 btc 冲币
	_, err = c.AddFunc("@every 5m", hbtc.CheckBlockSeek)
	if err != nil {
		mcommon.Log.Errorf("cron add func error: %#v", err)
	}
	// 检测 btc hot and fee uxto
	_, err = c.AddFunc("@every 5m", hbtc.CheckBlockSeekHotAndFee)
	if err != nil {
		mcommon.Log.Errorf("cron add func error: %#v", err)
	}
	// 检测 btc 零钱整理
	_, err = c.AddFunc("@every 10m", hbtc.CheckTxOrg)
	if err != nil {
		mcommon.Log.Errorf("cron add func error: %#v", err)
	}
	// 检测 btc 提币
	_, err = c.AddFunc("@every 3m", hbtc.CheckWithdraw)
	if err != nil {
		mcommon.Log.Errorf("cron add func error: %#v", err)
	}
	// 检测 btc 发送交易
	_, err = c.AddFunc("@every 1m", hbtc.CheckRawTxSend)
	if err != nil {
		mcommon.Log.Errorf("cron add func error: %#v", err)
	}
	// 检测 btc 交易上链
	_, err = c.AddFunc("@every 5m", hbtc.CheckRawTxConfirm)
	if err != nil {
		mcommon.Log.Errorf("cron add func error: %#v", err)
	}
	// 检测 btc 通知到账
	_, err = c.AddFunc("@every 5s", hbtc.CheckTxNotify)
	if err != nil {
		mcommon.Log.Errorf("cron add func error: %#v", err)
	}
	// 检测 btc 手续费
	_, err = c.AddFunc("@every 5m", hbtc.CheckGasPrice)
	if err != nil {
		mcommon.Log.Errorf("cron add func error: %#v", err)
	}

	// --- omni ---
	// 检测 omni 冲币
	_, err = c.AddFunc("@every 5m", hbtc.OmniCheckBlockSeek)
	if err != nil {
		mcommon.Log.Errorf("cron add func error: %#v", err)
	}
	// 检测 omni 零钱整理
	_, err = c.AddFunc("@every 10m", hbtc.OmniCheckTxOrg)
	if err != nil {
		mcommon.Log.Errorf("cron add func error: %#v", err)
	}
	// 检测 omni 提币
	_, err = c.AddFunc("@every 3m", hbtc.OmniCheckWithdraw)
	if err != nil {
		mcommon.Log.Errorf("cron add func error: %#v", err)
	}
	// 检测 omni 通知到账
	_, err = c.AddFunc("@every 5s", hbtc.OmniCheckTxNotify)
	if err != nil {
		mcommon.Log.Errorf("cron add func error: %#v", err)
	}

	// --- eos ---
	// 检测 eos 生成地址
	_, err = c.AddFunc("@every 1m", heos.CheckAddressFree)
	if err != nil {
		mcommon.Log.Errorf("cron add func error: %#v", err)
	}
	// 检测 eos 冲币
	_, err = c.AddFunc("@every 3s", heos.CheckBlockSeek)
	if err != nil {
		mcommon.Log.Errorf("cron add func error: %#v", err)
	}
	// 检测 eos 提币
	_, err = c.AddFunc("@every 3m", heos.CheckWithdraw)
	if err != nil {
		mcommon.Log.Errorf("cron add func error: %#v", err)
	}
	// 检测 eos 发送交易
	_, err = c.AddFunc("@every 1s", heos.CheckRawTxSend)
	if err != nil {
		mcommon.Log.Errorf("cron add func error: %#v", err)
	}
	// 检测 eos 交易上链
	_, err = c.AddFunc("@every 3s", heos.CheckRawTxConfirm)
	if err != nil {
		mcommon.Log.Errorf("cron add func error: %#v", err)
	}
	// 检测 eos 通知到账
	_, err = c.AddFunc("@every 3s", heos.CheckTxNotify)
	if err != nil {
		mcommon.Log.Errorf("cron add func error: %#v", err)
	}

	c.Start()
	select {}
}
