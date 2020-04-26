// 定时处理检测任务
package main

import (
	"go-dc-wallet/app"
	"go-dc-wallet/eth"
	"go-dc-wallet/hcommon"

	"github.com/robfig/cron/v3"
)

func main() {
	app.EnvCreate()
	defer app.EnvDestroy()

	c := cron.New(
		cron.WithSeconds(),
		cron.WithChain(
			cron.Recover(cron.DefaultLogger),
		),
	)
	var err error
	// 检测 eth 生成地址
	_, err = c.AddFunc("@every 1m", eth.CheckAddressFree)
	if err != nil {
		hcommon.Log.Errorf("cron add func error: %#v", err)
	}
	// 检测 eth 冲币
	_, err = c.AddFunc("@every 5s", eth.CheckBlockSeek)
	if err != nil {
		hcommon.Log.Errorf("cron add func error: %#v", err)
	}
	// 检测 eth 零钱整理
	_, err = c.AddFunc("@every 10m", eth.CheckAddressOrg)
	if err != nil {
		hcommon.Log.Errorf("cron add func error: %#v", err)
	}
	// 检测 eth 提币
	_, err = c.AddFunc("@every 3m", eth.CheckWithdraw)
	if err != nil {
		hcommon.Log.Errorf("cron add func error: %#v", err)
	}
	// 检测 eth 发送交易
	_, err = c.AddFunc("@every 1m", eth.CheckRawTxSend)
	if err != nil {
		hcommon.Log.Errorf("cron add func error: %#v", err)
	}
	// 检测 eth 交易上链
	_, err = c.AddFunc("@every 5s", eth.CheckRawTxConfirm)
	if err != nil {
		hcommon.Log.Errorf("cron add func error: %#v", err)
	}
	// 检测 eth 通知到账
	_, err = c.AddFunc("@every 5s", eth.CheckTxNotify)
	if err != nil {
		hcommon.Log.Errorf("cron add func error: %#v", err)
	}
	// 检测 eth 通知发送
	_, err = c.AddFunc("@every 1m", app.CheckDoNotify)
	if err != nil {
		hcommon.Log.Errorf("cron add func error: %#v", err)
	}
	// 检测 erc20 冲币
	_, err = c.AddFunc("@every 5s", eth.CheckErc20BlockSeek)
	if err != nil {
		hcommon.Log.Errorf("cron add func error: %#v", err)
	}
	// 检测 erc20 通知到账
	_, err = c.AddFunc("@every 5s", eth.CheckErc20TxNotify)
	if err != nil {
		hcommon.Log.Errorf("cron add func error: %#v", err)
	}
	// 检测 erc20 通知到账
	_, err = c.AddFunc("@every 5s", eth.CheckErc20TxNotify)
	if err != nil {
		hcommon.Log.Errorf("cron add func error: %#v", err)
	}
	// 检测 erc20 提币
	_, err = c.AddFunc("@every 3m", eth.CheckErc20Withdraw)
	if err != nil {
		hcommon.Log.Errorf("cron add func error: %#v", err)
	}

	c.Start()
	select {}
}
