package hbtc

import (
	"context"
	"go-dc-wallet/app"
	"go-dc-wallet/app/model"
	"go-dc-wallet/hcommon"
	"go-dc-wallet/omniclient"
	"time"

	"github.com/shopspring/decimal"
)

const (
	CoinSymbol = "btc"
)

func CheckAddressFree() {
	lockKey := "BtcCheckAddressFree"
	app.LockWrap(lockKey, func() {
		// 获取配置 允许的最小剩余地址数
		minFreeRow, err := app.SQLGetTAppConfigIntByK(
			context.Background(),
			app.DbCon,
			"min_free_address",
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		if minFreeRow == nil {
			hcommon.Log.Errorf("no config int of min_free_address")
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
		if freeCount < minFreeRow.V {
			var rows []*model.DBTAddressKey
			// 遍历差值次数
			for i := int64(0); i < minFreeRow.V-freeCount; i++ {
				// 生成私钥
				wif, err := GetNetwork(app.Cfg.BtcNetworkType).CreatePrivateKey()
				if err != nil {
					hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
				// 加密密钥
				wifStrEn := hcommon.AesEncrypt(wif.String(), app.Cfg.AESKey)
				// 获取地址
				address, err := GetNetwork(app.Cfg.BtcNetworkType).GetAddress(wif)
				if err != nil {
					hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
				// 存入待添加队列
				rows = append(rows, &model.DBTAddressKey{
					Symbol:  CoinSymbol,
					Address: address.EncodeAddress(),
					Pwd:     wifStrEn,
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

func CheckBlockSeek() {
	lockKey := "BtcCheckBlockSeek"
	app.LockWrap(lockKey, func() {
		// 获取配置 延迟确认数
		confirmRow, err := app.SQLGetTAppConfigIntByK(
			context.Background(),
			app.DbCon,
			"btc_block_confirm_num",
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		if confirmRow == nil {
			hcommon.Log.Errorf("no config int of btc_block_confirm_num")
			return
		}
		// 获取状态 当前处理完成的最新的block number
		seekRow, err := app.SQLGetTAppStatusIntByK(
			context.Background(),
			app.DbCon,
			"btc_seek_num",
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		if seekRow == nil {
			hcommon.Log.Errorf("no config int of btc_seek_num")
			return
		}
		rpcBlockNum, err := omniclient.RpcGetBlockCount()
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		vinTxMap := make(map[string]*omniclient.StTxResult)
		startI := seekRow.V + 1
		endI := rpcBlockNum - confirmRow.V
		hcommon.Log.Debugf("btc block seek %d->%d", startI, endI)
		if startI < endI {
			// 遍历获取需要查询的block信息
			for i := startI; i < endI; i++ {
				blockHash, err := omniclient.RpcGetBlockHash(i)
				if err != nil {
					hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
				// 一个block
				rpcBlock, err := omniclient.RpcGetBlockVerbose(blockHash)
				if err != nil {
					hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
				var toAddresses []string
				type StTxWithIndex struct {
					RpcTx *omniclient.StTxResult
					Index int64
				}
				toAddressTxMap := make(map[string][]*StTxWithIndex)
				// 所有tx
				for _, rpcTx := range rpcBlock.Tx {
					for _, vout := range rpcTx.Vout {
						if len(vout.ScriptPubKey.Addresses) == 1 {
							toAddress := vout.ScriptPubKey.Addresses[0]
							if !hcommon.IsStringInSlice(toAddresses, toAddress) {
								toAddresses = append(toAddresses, toAddress)
							}
							toAddressTxMap[toAddress] = append(toAddressTxMap[toAddress], &StTxWithIndex{
								RpcTx: rpcTx,
								Index: vout.N,
							})
						}
					}
				}
				hcommon.Log.Debugf("rpc get block: %d to addresses: %d", i, len(toAddresses))
				// 从db中查询这些地址是否是冲币地址中的地址
				dbAddressRows, err := app.SQLSelectTAddressKeyColByAddress(
					context.Background(),
					app.DbCon,
					[]string{
						model.DBColTAddressKeyAddress,
						model.DBColTAddressKeyUseTag,
					},
					toAddresses,
				)
				if err != nil {
					hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
				// 待插入数据
				var txBtcRows []*model.DBTTxBtc
				var txBtcUxtoRows []*model.DBTTxBtcUxto
				now := time.Now().Unix()
				// 遍历数据库地址
				for _, dbAddressRow := range dbAddressRows {
					toAddress := dbAddressRow.Address
					rpcTxWithIndexes := toAddressTxMap[toAddress]
					for _, rpcTxWithIndex := range rpcTxWithIndexes {
						rpcTx := rpcTxWithIndex.RpcTx
						voutIndex := rpcTxWithIndex.Index
						checkVout := rpcTx.Vout[voutIndex]

						voutAddress := checkVout.ScriptPubKey.Addresses[0]
						voutScript := checkVout.ScriptPubKey.Hex
						isVoutAddressInVin := false
						for _, vin := range rpcTx.Vin {
							rpcVinTx, ok := vinTxMap[vin.Txid]
							if !ok {
								rpcVinTx, err = omniclient.RpcGetRawTransactionVerbose(vin.Txid)
								if err != nil {
									hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
									return
								}
								vinTxMap[vin.Txid] = rpcVinTx
								hcommon.Log.Debugf("get tx: %s", vin.Txid)
							}
							if int64(len(rpcVinTx.Vout)) > vin.Vout {
								if len(rpcVinTx.Vout[vin.Vout].ScriptPubKey.Addresses) > 0 {
									if rpcVinTx.Vout[vin.Vout].ScriptPubKey.Addresses[0] == voutAddress {
										isVoutAddressInVin = true
										break
									}
								}
							}
						}
						if !isVoutAddressInVin {
							// 记录数据
							value := decimal.NewFromFloat(checkVout.Value).String()
							txBtcRows = append(
								txBtcRows,
								&model.DBTTxBtc{
									BlockHash:    rpcBlock.Hash,
									TxID:         rpcTx.Txid,
									VoutN:        voutIndex,
									VoutAddress:  voutAddress,
									VoutValue:    value,
									CreateTime:   now,
									HandleStatus: 0,
									HandleMsg:    "",
									HandleTime:   now,
								},
							)
							txBtcUxtoRows = append(
								txBtcUxtoRows,
								&model.DBTTxBtcUxto{
									BlockHash:    rpcBlock.Hash,
									TxID:         rpcTx.Txid,
									VoutN:        voutIndex,
									VoutAddress:  voutAddress,
									VoutValue:    value,
									VoutScript:   voutScript,
									CreateTime:   now,
									SpendTxID:    "",
									SpendN:       0,
									HandleStatus: 0,
									HandleMsg:    "",
									HandleTime:   now,
								},
							)
						}
					}
				}

				// 插入数据库
				_, err = model.SQLCreateIgnoreManyTTxBtc(
					context.Background(),
					app.DbCon,
					txBtcRows,
				)
				if err != nil {
					hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
				_, err = model.SQLCreateIgnoreManyTTxBtcUxto(
					context.Background(),
					app.DbCon,
					txBtcUxtoRows,
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
						K: "btc_seek_num",
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
