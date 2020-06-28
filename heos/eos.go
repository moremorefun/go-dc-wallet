package heos

import (
	"context"
	"fmt"
	"go-dc-wallet/app"
	"go-dc-wallet/app/model"
	"go-dc-wallet/hcommon"
)

// CheckAddressFree 检测剩余地址数
func CheckAddressFree() {
	lockKey := "EosCheckAddressFree"
	app.LockWrap(lockKey, func() {
		// 获取配置 允许的最小剩余地址数
		minFreeValue, err := app.SQLGetTAppConfigIntValueByK(
			context.Background(),
			app.DbCon,
			"min_free_address",
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		// 获取当前剩余可用地址数
		freeCount, err := app.SQLGetTAddressKeyFreeCount(
			context.Background(),
			app.DbCon,
			CoinSymbol,
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		// 如果数据库中剩余可用地址小于最小允许可用地址
		if freeCount < minFreeValue {
			var rows []*model.DBTAddressKey
			// 获取最大值
			maxAddress, err := app.SQLGetTAddressMaxIntOfEos(
				context.Background(),
				app.DbCon,
			)
			if maxAddress < MiniAddress {
				maxAddress = MiniAddress
			}
			maxAddress++
			// 遍历差值次数
			for i := int64(0); i < minFreeValue-freeCount; i++ {
				// 存入待添加队列
				rows = append(rows, &model.DBTAddressKey{
					Symbol:  CoinSymbol,
					Address: fmt.Sprintf("%d", maxAddress+i),
					Pwd:     "",
					UseTag:  0,
				})
			}
			// 一次性将生成的地址存入数据库
			_, err = model.SQLCreateIgnoreManyTAddressKey(
				context.Background(),
				app.DbCon,
				rows,
			)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
		}
	})
}
