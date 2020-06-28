package heos

import (
	"context"
	"encoding/json"
	"fmt"
	"go-dc-wallet/app"
	"go-dc-wallet/app/model"
	"go-dc-wallet/eosclient"
	"go-dc-wallet/hcommon"
	"strings"
	"time"
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

// CheckBlockSeek 检测到账
func CheckBlockSeek() {
	lockKey := "EosCheckBlockSeek"
	app.LockWrap(lockKey, func() {
		// 获取状态 当前处理完成的最新的block number
		seekValue, err := app.SQLGetTAppStatusIntValueByK(
			context.Background(),
			app.DbCon,
			"eos_seek_num",
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		rpcChainInfo, err := eosclient.RpcChainGetInfo()
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		startI := seekValue + 1
		endI := rpcChainInfo.LastIrreversibleBlockNum + 1
		hcommon.Log.Debugf("eos check block: %d->%d", startI, endI)
		if startI < endI {
			// 获取冷钱包地址
			eosColdAddressValue, err := app.SQLGetTAppConfigStrValueByK(
				context.Background(),
				app.DbCon,
				"cold_wallet_address_eos",
			)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			// 遍历获取需要查询的block信息
			now := time.Now().Unix()
			for i := startI; i < endI; i++ {
				hcommon.Log.Debugf("eos check block: %d", i)
				rpcBlock, err := eosclient.RpcChainGetBlock(i)
				if err != nil {
					hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
				var memos []string
				type stAction struct {
					rpcTrx        eosclient.StTransactionTrx
					rpcActionData eosclient.StActionData
					actionIndex   int64
				}
				memosMap := make(map[string][]stAction)
				for _, rpcTransaction := range rpcBlock.Transactions {
					if rpcTransaction.Status != "executed" {
						continue
					}
					var rpcTrx eosclient.StTransactionTrx
					err := json.Unmarshal(rpcTransaction.Trx, &rpcTrx)
					if err != nil {
						_, ok := err.(*json.UnmarshalTypeError)
						if !ok {
							hcommon.Log.Debugf("err: [%T] %s", err, err.Error())
							return
						} else {
							continue
						}
					}
					if rpcTrx.Transaction.DelaySec != 0 {
						continue
					}
					for i, rpcAction := range rpcTrx.Transaction.Actions {
						if rpcAction.Account != "eosio.token" {
							continue
						}
						if rpcAction.Name != "transfer" {
							continue
						}
						var rpcActionData eosclient.StActionData
						err := json.Unmarshal(rpcAction.Data, &rpcActionData)
						if err != nil {
							_, ok := err.(*json.UnmarshalTypeError)
							if !ok {
								hcommon.Log.Debugf("err: [%T] %s", err, err.Error())
								return
							} else {
								continue
							}
						}
						if rpcActionData.Quantity != "" {
							if eosColdAddressValue != "" && rpcActionData.To == eosColdAddressValue {
								hcommon.Log.Debugf(
									"%s:%d %s->%s: memo:%s value:%s",
									rpcTrx.ID[:5],
									i,
									rpcActionData.From,
									rpcActionData.To,
									rpcActionData.Memo,
									rpcActionData.Quantity,
								)
								// 打款到冷钱包
								if !hcommon.IsStringInSlice(memos, rpcActionData.Memo) {
									memos = append(memos, rpcActionData.Memo)
								}
								quantitys := strings.Split(rpcActionData.Quantity, " ")
								if len(quantitys) != 2 {
									hcommon.Log.Errorf("error %s", rpcActionData.Quantity)
									return
								}
								if quantitys[1] != "EOS" {
									hcommon.Log.Errorf("error %s", rpcActionData.Quantity)
									return
								}
								rpcActionData.Quantity = quantitys[0]
								memosMap[rpcActionData.Memo] = append(
									memosMap[rpcActionData.Memo],
									stAction{
										rpcTrx:        rpcTrx,
										rpcActionData: rpcActionData,
										actionIndex:   int64(i),
									},
								)
							}
						}
					}
				}
				// 从db中查询这些地址是否是冲币地址中的地址
				dbAddressRows, err := app.SQLSelectTAddressKeyColByAddress(
					context.Background(),
					app.DbCon,
					[]string{
						model.DBColTAddressKeyAddress,
						model.DBColTAddressKeyUseTag,
					},
					memos,
				)
				if err != nil {
					hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
				var txRows []*model.DBTTxEos
				for _, dbAddressRow := range dbAddressRows {
					tActions := memosMap[dbAddressRow.Address]
					for _, tAction := range tActions {
						txRows = append(
							txRows,
							&model.DBTTxEos{
								ProductID:    dbAddressRow.UseTag,
								TxHash:       tAction.rpcTrx.ID,
								LogIndex:     tAction.actionIndex,
								FromAddress:  tAction.rpcActionData.From,
								ToAddress:    tAction.rpcActionData.To,
								Memo:         tAction.rpcActionData.Memo,
								BalanceReal:  tAction.rpcActionData.Quantity,
								CreateAt:     now,
								HandleStatus: 0,
								HandleMsg:    "",
								HandleAt:     now,
							},
						)
					}
				}
				_, err = model.SQLCreateIgnoreManyTTxEos(
					context.Background(),
					app.DbCon,
					txRows,
				)
				if err != nil {
					hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
				// 更新block num
				_, err = app.SQLUpdateTAppStatusIntByK(
					context.Background(),
					app.DbCon,
					&model.DBTAppStatusInt{
						K: "eos_seek_num",
						V: i,
					},
				)
				if err != nil {
					hcommon.Log.Errorf("SQLUpdateTAppStatusIntByK err: [%T] %s", err, err.Error())
					return
				}
			}
		}
	})
}
