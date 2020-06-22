package hbtc

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"go-dc-wallet/app"
	"go-dc-wallet/app/model"
	"go-dc-wallet/hcommon"
	"go-dc-wallet/omniclient"
	"strings"
	"time"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"

	"github.com/btcsuite/btcd/chaincfg/chainhash"

	"github.com/btcsuite/btcd/wire"

	"github.com/shopspring/decimal"
)

const (
	CoinSymbol = "btc"
)

// CheckAddressFree 检测剩余地址数
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

// CheckBlockSeek 检测到账
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
				// 目标地址
				var toAddresses []string
				type StTxWithIndex struct {
					RpcTx *omniclient.StTxResult
					Index int64
				}
				toAddressTxMap := make(map[string][]*StTxWithIndex)
				// 来源hash
				var fromTxHashes []string
				type StVinWithIndex struct {
					TxHash string
					VoutN  int64

					SpendTxHash string
					SpendN      int64
				}
				vinMap := make(map[string]*StVinWithIndex)
				// 所有tx
				for _, rpcTx := range rpcBlock.Tx {
					for i, vin := range rpcTx.Vin {
						fromTxHash := vin.Txid
						if !hcommon.IsStringInSlice(fromTxHashes, fromTxHash) {
							fromTxHashes = append(fromTxHashes, fromTxHash)
						}
						key := fmt.Sprintf("%s-%d", vin.Txid, vin.Vout)
						vinMap[key] = &StVinWithIndex{
							TxHash:      vin.Txid,
							VoutN:       vin.Vout,
							SpendTxHash: rpcTx.Txid,
							SpendN:      int64(i),
						}
					}
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
						value := decimal.NewFromFloat(checkVout.Value).String()
						if !isVoutAddressInVin && dbAddressRow.UseTag > 0 {
							// 记录数据 只记录已经获取，并且输入没有输出的记录
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
						}
						uxtoType := int64(app.UxtoTypeTx)
						if dbAddressRow.UseTag < 0 {
							uxtoType = app.UxtoTypeHot
						}
						txBtcUxtoRows = append(
							txBtcUxtoRows,
							&model.DBTTxBtcUxto{
								UxtoType:     uxtoType,
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

				// 从uxto中查询txhash
				var updateUxtoRows []*model.DBTTxBtcUxto
				uxtoRows, err := app.SQLSelectTTxBtcUxtoColByTxIDs(
					context.Background(),
					app.DbCon,
					[]string{
						model.DBColTTxBtcUxtoID,
						model.DBColTTxBtcUxtoTxID,
						model.DBColTTxBtcUxtoVoutN,
					},
					fromTxHashes,
				)
				if err != nil {
					hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
				for _, uxtoRow := range uxtoRows {
					key := fmt.Sprintf("%s-%d", uxtoRow.TxID, uxtoRow.VoutN)
					rpcVin, ok := vinMap[key]
					if ok {
						updateUxtoRows = append(updateUxtoRows, &model.DBTTxBtcUxto{
							ID:           uxtoRow.ID,
							TxID:         uxtoRow.TxID,
							VoutN:        uxtoRow.VoutN,
							SpendTxID:    rpcVin.SpendTxHash,
							SpendN:       rpcVin.SpendN,
							HandleStatus: app.UxtoHandleStatusConfirm,
							HandleMsg:    "confirm",
							HandleTime:   now,
						})
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
				// 更新uxto状态
				_, err = app.SQLCreateManyTTxBtcUxtoUpdate(
					context.Background(),
					app.DbCon,
					updateUxtoRows,
				)
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

// CheckTxOrg 检测零钱整理
func CheckTxOrg() {
	lockKey := "BtcCheckTxOrg"
	app.LockWrap(lockKey, func() {
		uxtoRows, err := app.SQLSelectTTxBtcUxtoColToOrg(
			context.Background(),
			app.DbCon,
			[]string{
				model.DBColTTxBtcUxtoID,
				model.DBColTTxBtcUxtoTxID,
				model.DBColTTxBtcUxtoVoutN,
				model.DBColTTxBtcUxtoVoutAddress,
				model.DBColTTxBtcUxtoVoutValue,
				model.DBColTTxBtcUxtoVoutScript,
			},
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		if len(uxtoRows) <= 0 {
			return
		}
		// 获取冷包地址
		coldRow, err := app.SQLGetTAppConfigStrByK(
			context.Background(),
			app.DbCon,
			"cold_wallet_address_btc",
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		if coldRow == nil {
			hcommon.Log.Errorf("no config int of cold_wallet_address_btc")
			return
		}
		// 获取手续费配置
		feeRow, err := app.SQLGetTAppStatusIntByK(
			context.Background(),
			app.DbCon,
			"to_cold_gas_price_btc",
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		if feeRow == nil {
			hcommon.Log.Errorf("no config int of to_cold_gas_price_btc")
			return
		}
		// 获取私钥
		var addresses []string
		for _, uxtoRow := range uxtoRows {
			addresses = append(addresses, uxtoRow.VoutAddress)
		}
		addressKeyMap, err := app.SQLGetAddressKeyMap(
			context.Background(),
			app.DbCon,
			[]string{
				model.DBColTAddressKeyID,
				model.DBColTAddressKeyAddress,
				model.DBColTAddressKeyPwd,
			},
			addresses,
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		// 生成交易
		tx := wire.NewMsgTx(wire.TxVersion)
		voutAmount := decimal.NewFromInt(0)
		for _, uxtoRow := range uxtoRows {
			hash, _ := chainhash.NewHashFromStr(uxtoRow.TxID)
			outPoint := wire.NewOutPoint(hash, uint32(uxtoRow.VoutN))
			txIn := wire.NewTxIn(outPoint, nil, nil)
			tx.AddTxIn(txIn)

			amount, err := decimal.NewFromString(uxtoRow.VoutValue)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			voutAmount = voutAmount.Add(amount)
		}
		addrTo, err := btcutil.DecodeAddress(
			coldRow.V,
			GetNetwork(app.Cfg.BtcNetworkType).params,
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		pkScriptf, err := txscript.PayToAddrScript(addrTo)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		baf := voutAmount.Mul(decimal.NewFromInt(1e8)).IntPart()
		tx.AddTxOut(wire.NewTxOut(baf, pkScriptf))
		// 签名,用于计算手续费
		for i, uxtoRow := range uxtoRows {
			addressKey, ok := addressKeyMap[uxtoRow.VoutAddress]
			if !ok {
				hcommon.Log.Errorf("no address key: %s", uxtoRow.VoutAddress)
				return
			}
			key := hcommon.AesDecrypt(addressKey.Pwd, app.Cfg.AESKey)
			if len(key) == 0 {
				hcommon.Log.Errorf("error key of: %s", uxtoRow.VoutAddress)
				return
			}
			wif, err := btcutil.DecodeWIF(key)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			txinPkScript, err := hex.DecodeString(uxtoRow.VoutScript)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			script, err := txscript.SignatureScript(
				tx,
				i,
				txinPkScript,
				txscript.SigHashAll,
				wif.PrivKey,
				true,
			)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			tx.TxIn[i].SignatureScript = script
			vm, err := txscript.NewEngine(
				txinPkScript,
				tx,
				i,
				txscript.StandardVerifyFlags,
				nil,
				nil,
				-1,
			)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			err = vm.Execute()
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
		}
		// 计算手续费
		txSize := tx.SerializeSize()
		txFee := feeRow.V * int64(txSize)
		tx.TxOut[0].Value -= txFee
		// 重新计算签名
		for i, uxtoRow := range uxtoRows {
			addressKey, ok := addressKeyMap[uxtoRow.VoutAddress]
			if !ok {
				hcommon.Log.Errorf("no address key: %s", uxtoRow.VoutAddress)
				return
			}
			key := hcommon.AesDecrypt(addressKey.Pwd, app.Cfg.AESKey)
			if len(key) == 0 {
				hcommon.Log.Errorf("error key of: %s", uxtoRow.VoutAddress)
				return
			}
			wif, err := btcutil.DecodeWIF(key)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			txinPkScript, err := hex.DecodeString(uxtoRow.VoutScript)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			script, err := txscript.SignatureScript(
				tx,
				i,
				txinPkScript,
				txscript.SigHashAll,
				wif.PrivKey,
				true,
			)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			tx.TxIn[i].SignatureScript = script
			vm, err := txscript.NewEngine(
				txinPkScript,
				tx,
				i,
				txscript.StandardVerifyFlags,
				nil,
				nil,
				-1,
			)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			err = vm.Execute()
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
		}
		b := new(bytes.Buffer)
		b.Grow(tx.SerializeSize())
		err = tx.Serialize(b)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		hcommon.Log.Debugf("raw tx: %s", hex.EncodeToString(b.Bytes()))
		// 准备插入数据
		now := time.Now().Unix()
		var sendRows []*model.DBTSendBtc
		var updateUxtoRows []*model.DBTTxBtcUxto
		for i, uxtoRow := range uxtoRows {
			sendHex := ""
			if i == 0 {
				sendHex = hex.EncodeToString(b.Bytes())
			}
			sendRows = append(sendRows, &model.DBTSendBtc{
				RelatedType:  app.SendRelationTypeUXTOOrg,
				RelatedID:    uxtoRow.ID,
				TokenID:      0,
				TxID:         tx.TxHash().String(),
				FromAddress:  uxtoRow.VoutAddress,
				ToAddress:    coldRow.V,
				Balance:      tx.TxOut[0].Value,
				BalanceReal:  (decimal.NewFromInt(tx.TxOut[0].Value).Div(decimal.NewFromInt(1e8))).String(),
				Gas:          int64(txSize),
				GasPrice:     feeRow.V,
				Hex:          sendHex,
				CreateTime:   now,
				HandleStatus: 0,
				HandleMsg:    "",
				HandleTime:   now,
			})
			updateUxtoRows = append(updateUxtoRows, &model.DBTTxBtcUxto{
				ID:           uxtoRow.ID,
				SpendTxID:    tx.TxHash().String(),
				SpendN:       int64(i),
				HandleStatus: app.UxtoHandleStatusUse,
				HandleMsg:    "use",
				HandleTime:   now,
			})
		}
		// 插入数据
		_, err = model.SQLCreateIgnoreManyTSendBtc(
			context.Background(),
			app.DbCon,
			sendRows,
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		// 更新uxto状态
		_, err = app.SQLCreateManyTTxBtcUxtoUpdate(
			context.Background(),
			app.DbCon,
			updateUxtoRows,
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
	})
}

// CheckRawTxSend 发送交易
func CheckRawTxSend() {
	lockKey := "BtcCheckRawTxSend"
	app.LockWrap(lockKey, func() {
		sendRows, err := app.SQLSelectTSendBtcColByStatus(
			context.Background(),
			app.DbCon,
			[]string{
				model.DBColTSendBtcID,
				model.DBColTSendBtcTxID,
				model.DBColTSendBtcHex,
			},
			app.SendStatusInit,
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		var sendIDs []int64
		var sendTxHashes []string
		for _, sendRow := range sendRows {
			if sendRow.Hex == "" {
				continue
			}
			_, err := omniclient.RpcSendRawTransaction(sendRow.Hex)
			if err != nil && strings.Contains(err.Error(), "already in block chain") {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				continue
			}
			sendIDs = append(sendIDs, sendRow.ID)
			sendTxHashes = append(sendTxHashes, sendRow.TxID)
		}
		for _, sendRow := range sendRows {
			if hcommon.IsStringInSlice(sendTxHashes, sendRow.TxID) {
				if !hcommon.IsIntInSlice(sendIDs, sendRow.ID) {
					sendIDs = append(sendIDs, sendRow.ID)
				}
			}
		}
		now := time.Now().Unix()
		_, err = app.SQLUpdateTSendBtcByIDs(
			context.Background(),
			app.DbCon,
			sendIDs,
			&model.DBTSendBtc{
				HandleStatus: app.SendStatusSend,
				HandleTime:   now,
				HandleMsg:    "send",
			},
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
	})
}
