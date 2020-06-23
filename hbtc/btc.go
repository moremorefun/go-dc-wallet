package hbtc

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"go-dc-wallet/app"
	"go-dc-wallet/app/model"
	"go-dc-wallet/hcommon"
	"go-dc-wallet/omniclient"
	"net/http"
	"strings"
	"time"

	"github.com/parnurzeal/gorequest"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"github.com/gin-gonic/gin"

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
		endI := rpcBlockNum - confirmRow.V + 1
		hcommon.Log.Debugf("btc block seek %d->%d", startI, endI)
		if startI < endI {
			// 获取所有token
			var tokenHotAddresses []string
			tokenRows, err := app.SQLSelectTAppConfigTokenBtcColAll(
				context.Background(),
				app.DbCon,
				[]string{
					model.DBColTAppConfigTokenBtcID,
					model.DBColTAppConfigTokenBtcHotAddress,
				},
			)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			for _, tokenRow := range tokenRows {
				if !hcommon.IsStringInSlice(tokenHotAddresses, tokenRow.HotAddress) {
					tokenHotAddresses = append(tokenHotAddresses, tokenRow.HotAddress)
				}
			}
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
					RpcTx    *omniclient.StTxResult
					Index    int64
					IsOmniTx bool
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
					omniScript := "6a146f6d6e69"
					isOmniTx := false
					for _, vout := range rpcTx.Vout {
						if strings.HasPrefix(vout.ScriptPubKey.Hex, omniScript) {
							isOmniTx = true
						}
					}
					for _, vout := range rpcTx.Vout {
						if len(vout.ScriptPubKey.Addresses) == 1 {
							toAddress := vout.ScriptPubKey.Addresses[0]
							if !hcommon.IsStringInSlice(toAddresses, toAddress) {
								toAddresses = append(toAddresses, toAddress)
							}
							toAddressTxMap[toAddress] = append(toAddressTxMap[toAddress], &StTxWithIndex{
								RpcTx:    rpcTx,
								Index:    vout.N,
								IsOmniTx: isOmniTx,
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
						omniVinAddress := ""
						if rpcTxWithIndex.IsOmniTx {
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
								if len(rpcVinTx.Vout[vin.Vout].ScriptPubKey.Addresses) > 0 {
									omniVinAddress = strings.Join(rpcVinTx.Vout[vin.Vout].ScriptPubKey.Addresses, ",")
									break
								}
							}
						}
						value := decimal.NewFromFloat(checkVout.Value).String()
						if dbAddressRow.UseTag > 0 &&
							!rpcTxWithIndex.IsOmniTx {
							// 记录数据 只记录已经获取，并且输入没有输出的记录
							txBtcRows = append(
								txBtcRows,
								&model.DBTTxBtc{
									ProductID:    dbAddressRow.UseTag,
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
						if rpcTxWithIndex.IsOmniTx {
							omniOutAddress := ""
							for i := len(rpcTx.Vout) - 1; i >= 0; i-- {
								vout := rpcTx.Vout[i]
								if len(vout.ScriptPubKey.Addresses) > 0 {
									toAddress := strings.Join(vout.ScriptPubKey.Addresses, ",")
									if toAddress != omniVinAddress {
										omniOutAddress = toAddress
										break
									}
								}
							}
							if omniOutAddress == voutAddress {
								uxtoType = app.UxtoTypeOmni
							}
						}
						if hcommon.IsStringInSlice(tokenHotAddresses, voutAddress) {
							uxtoType = app.UxtoTypeOmniHot
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
		allUxtoRows, err := app.SQLSelectTTxBtcUxtoColToOrg(
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
			app.UxtoTypeTx,
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		if len(allUxtoRows) <= 0 {
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
		for _, uxtoRow := range allUxtoRows {
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
		// 按5000个in拆分
		step := 5000
		for i := 0; i < len(allUxtoRows); i += step {
			endI := i + step
			if len(allUxtoRows) < endI {
				endI = len(allUxtoRows)
			}
			uxtoRows := allUxtoRows[i:endI]
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
				balance := int64(0)
				balanceReal := "0"
				gas := int64(0)
				gasPrice := int64(0)
				if i == 0 {
					sendHex = hex.EncodeToString(b.Bytes())
					balance = tx.TxOut[0].Value
					balanceReal = (decimal.NewFromInt(tx.TxOut[0].Value).Div(decimal.NewFromInt(1e8))).String()
					gas = int64(txSize)
					gasPrice = feeRow.V
				}
				sendRows = append(sendRows, &model.DBTSendBtc{
					RelatedType:  app.SendRelationTypeUXTOOrg,
					RelatedID:    uxtoRow.ID,
					TokenID:      0,
					TxID:         tx.TxHash().String(),
					FromAddress:  uxtoRow.VoutAddress,
					ToAddress:    coldRow.V,
					Balance:      balance,
					BalanceReal:  balanceReal,
					Gas:          gas,
					GasPrice:     gasPrice,
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
				model.DBColTSendBtcRelatedType,
				model.DBColTSendBtcRelatedID,
			},
			app.SendStatusInit,
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		// 首先单独处理提币，提取提币通知要使用的数据
		var withdrawIDs []int64
		for _, sendRow := range sendRows {
			switch sendRow.RelatedType {
			case app.SendRelationTypeWithdraw:
				if !hcommon.IsIntInSlice(withdrawIDs, sendRow.RelatedID) {
					withdrawIDs = append(withdrawIDs, sendRow.RelatedID)
				}
			}
		}
		withdrawMap, err := app.SQLGetWithdrawMap(
			context.Background(),
			app.DbCon,
			[]string{
				model.DBColTWithdrawID,
				model.DBColTWithdrawProductID,
				model.DBColTWithdrawOutSerial,
				model.DBColTWithdrawToAddress,
				model.DBColTWithdrawBalanceReal,
				model.DBColTWithdrawTxHash,
			},
			withdrawIDs,
		)
		// 产品
		var productIDs []int64
		for _, withdrawRow := range withdrawMap {
			if !hcommon.IsIntInSlice(productIDs, withdrawRow.ProductID) {
				productIDs = append(productIDs, withdrawRow.ProductID)
			}
		}
		productMap, err := app.SQLGetProductMap(
			context.Background(),
			app.DbCon,
			[]string{
				model.DBColTProductID,
				model.DBColTProductAppName,
				model.DBColTProductCbURL,
				model.DBColTProductAppSk,
			},
			productIDs,
		)
		// 发送
		// 通知数据
		var notifyRows []*model.DBTProductNotify
		withdrawIDs = []int64{}
		now := time.Now().Unix()
		addNotifyRow := func(sendRow *model.DBTSendBtc) error {
			// 如果是提币，创建通知信息
			if sendRow.RelatedType == app.SendRelationTypeWithdraw {
				withdrawRow, ok := withdrawMap[sendRow.RelatedID]
				if !ok {
					hcommon.Log.Errorf("withdrawMap no: %d", sendRow.RelatedID)
					return nil
				}
				productRow, ok := productMap[withdrawRow.ProductID]
				if !ok {
					hcommon.Log.Errorf("productMap no: %d", withdrawRow.ProductID)
					return nil
				}
				nonce := hcommon.GetUUIDStr()
				reqObj := gin.H{
					"tx_hash":     withdrawRow.TxHash,
					"balance":     withdrawRow.BalanceReal,
					"app_name":    productRow.AppName,
					"out_serial":  withdrawRow.OutSerial,
					"address":     withdrawRow.ToAddress,
					"symbol":      CoinSymbol,
					"notify_type": app.NotifyTypeWithdrawSend,
				}
				reqObj["sign"] = hcommon.GetSign(productRow.AppSk, reqObj)
				req, err := json.Marshal(reqObj)
				if err != nil {
					hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return err
				}
				notifyRows = append(notifyRows, &model.DBTProductNotify{
					Nonce:        nonce,
					ProductID:    withdrawRow.ProductID,
					ItemType:     app.SendRelationTypeWithdraw,
					ItemID:       withdrawRow.ID,
					NotifyType:   app.NotifyTypeWithdrawSend,
					URL:          productRow.CbURL,
					Msg:          string(req),
					HandleStatus: app.NotifyStatusInit,
					HandleMsg:    "",
					CreateTime:   now,
					UpdateTime:   now,
				})
				withdrawIDs = append(withdrawIDs, withdrawRow.ID)
			}
			return nil
		}

		var sendIDs []int64
		var sendTxHashes []string
		for _, sendRow := range sendRows {
			if sendRow.Hex == "" {
				continue
			}
			_, err := omniclient.RpcSendRawTransaction(sendRow.Hex)
			if err != nil && !strings.Contains(err.Error(), "already in block chain") {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				continue
			}
			sendIDs = append(sendIDs, sendRow.ID)
			sendTxHashes = append(sendTxHashes, sendRow.TxID)
			err = addNotifyRow(sendRow)
			if err != nil {
				return
			}
		}
		for _, sendRow := range sendRows {
			if hcommon.IsStringInSlice(sendTxHashes, sendRow.TxID) {
				if !hcommon.IsIntInSlice(sendIDs, sendRow.ID) {
					sendIDs = append(sendIDs, sendRow.ID)
					err = addNotifyRow(sendRow)
					if err != nil {
						return
					}
				}
			}
		}
		// 更新提币状态
		_, err = app.SQLUpdateTWithdrawStatusByIDs(
			context.Background(),
			app.DbCon,
			withdrawIDs,
			&model.DBTWithdraw{
				HandleStatus: app.WithdrawStatusSend,
				HandleMsg:    "send",
				HandleTime:   now,
			},
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		// 添加发送通知
		_, err = model.SQLCreateIgnoreManyTProductNotify(
			context.Background(),
			app.DbCon,
			notifyRows,
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
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

// CheckRawTxConfirm 确认tx是否打包完成
func CheckRawTxConfirm() {
	lockKey := "BtcCheckRawTxConfirm"
	app.LockWrap(lockKey, func() {
		sendRows, err := app.SQLSelectTSendBtcColByStatus(
			context.Background(),
			app.DbCon,
			[]string{
				model.DBColTSendBtcID,
				model.DBColTSendBtcTxID,
				model.DBColTSendBtcHex,
				model.DBColTSendBtcRelatedType,
				model.DBColTSendBtcRelatedID,
			},
			app.SendStatusSend,
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		// 获取提币信息
		var withdrawIDs []int64
		for _, sendRow := range sendRows {
			if sendRow.RelatedType == app.SendRelationTypeWithdraw {
				// 提币
				if !hcommon.IsIntInSlice(withdrawIDs, sendRow.RelatedID) {
					withdrawIDs = append(withdrawIDs, sendRow.RelatedID)
				}
			}
		}
		withdrawMap, err := app.SQLGetWithdrawMap(
			context.Background(),
			app.DbCon,
			[]string{
				model.DBColTWithdrawID,
				model.DBColTWithdrawProductID,
				model.DBColTWithdrawOutSerial,
				model.DBColTWithdrawToAddress,
				model.DBColTWithdrawBalanceReal,
				model.DBColTWithdrawSymbol,
				model.DBColTWithdrawTxHash,
			},
			withdrawIDs,
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		var productIDs []int64
		for _, withdrawRow := range withdrawMap {
			if !hcommon.IsIntInSlice(productIDs, withdrawRow.ProductID) {
				productIDs = append(productIDs, withdrawRow.ProductID)
			}
		}
		productMap, err := app.SQLGetProductMap(
			context.Background(),
			app.DbCon,
			[]string{
				model.DBColTProductID,
				model.DBColTProductAppName,
				model.DBColTProductCbURL,
				model.DBColTProductAppSk,
			},
			productIDs,
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}

		var notifyRows []*model.DBTProductNotify
		withdrawIDs = []int64{}
		now := time.Now().Unix()
		addWithdrawNotify := func(sendRow *model.DBTSendBtc) error {
			if sendRow.RelatedType == app.SendRelationTypeWithdraw {
				// 提币
				withdrawRow, ok := withdrawMap[sendRow.RelatedID]
				if !ok {
					hcommon.Log.Errorf("no withdrawMap: %d", sendRow.RelatedID)
					return nil
				}
				productRow, ok := productMap[withdrawRow.ProductID]
				if !ok {
					hcommon.Log.Errorf("no productMap: %d", withdrawRow.ProductID)
					return nil
				}
				nonce := hcommon.GetUUIDStr()
				reqObj := gin.H{
					"tx_hash":     withdrawRow.TxHash,
					"balance":     withdrawRow.BalanceReal,
					"app_name":    productRow.AppName,
					"out_serial":  withdrawRow.OutSerial,
					"address":     withdrawRow.ToAddress,
					"symbol":      withdrawRow.Symbol,
					"notify_type": app.NotifyTypeWithdrawConfirm,
				}
				reqObj["sign"] = hcommon.GetSign(productRow.AppSk, reqObj)
				req, err := json.Marshal(reqObj)
				if err != nil {
					hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return err
				}
				notifyRows = append(notifyRows, &model.DBTProductNotify{
					Nonce:        nonce,
					ProductID:    withdrawRow.ProductID,
					ItemType:     app.SendRelationTypeWithdraw,
					ItemID:       withdrawRow.ID,
					NotifyType:   app.NotifyTypeWithdrawConfirm,
					URL:          productRow.CbURL,
					Msg:          string(req),
					HandleStatus: app.NotifyStatusInit,
					HandleMsg:    "",
					CreateTime:   now,
					UpdateTime:   now,
				})
				withdrawIDs = append(withdrawIDs, withdrawRow.ID)
			}
			return nil
		}

		var sendIDs []int64
		var confirmHashes []string
		for _, sendRow := range sendRows {
			if !hcommon.IsStringInSlice(confirmHashes, sendRow.TxID) {
				rpcTx, err := omniclient.RpcGetRawTransactionVerbose(sendRow.TxID)
				if err != nil {
					hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					continue
				}
				if rpcTx.Confirmations <= 0 {
					continue
				}
				confirmHashes = append(confirmHashes, sendRow.TxID)
			}
			// 已经确认
			sendIDs = append(sendIDs, sendRow.ID)
			err = addWithdrawNotify(sendRow)
			if err != nil {
				return
			}
		}
		// 更新提币状态
		_, err = app.SQLUpdateTWithdrawStatusByIDs(
			context.Background(),
			app.DbCon,
			withdrawIDs,
			&model.DBTWithdraw{
				HandleStatus: app.WithdrawStatusConfirm,
				HandleMsg:    "confirm",
				HandleTime:   now,
			},
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		_, err = model.SQLCreateIgnoreManyTProductNotify(
			context.Background(),
			app.DbCon,
			notifyRows,
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		_, err = app.SQLUpdateTSendBtcByIDs(
			context.Background(),
			app.DbCon,
			sendIDs,
			&model.DBTSendBtc{
				HandleStatus: app.SendStatusConfirm,
				HandleTime:   now,
				HandleMsg:    "confirm",
			},
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
	})
}

// CheckWithdraw 检测提现
func CheckWithdraw() {
	lockKey := "BtcCheckWithdraw"
	app.LockWrap(lockKey, func() {
		// 获取提币信息
		withdrawRows, err := app.SQLSelectTWithdrawColByStatus(
			context.Background(),
			app.DbCon,
			[]string{
				model.DBColTWithdrawID,
				model.DBColTWithdrawProductID,
				model.DBColTWithdrawOutSerial,
				model.DBColTWithdrawToAddress,
				model.DBColTWithdrawBalanceReal,
			},
			app.WithdrawStatusInit,
			[]string{"btc"},
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		if len(withdrawRows) == 0 {
			return
		}
		// 获取热钱包uxto
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
			app.UxtoTypeHot,
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
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
		// 获取手续费配置
		feeRow, err := app.SQLGetTAppStatusIntByK(
			context.Background(),
			app.DbCon,
			"to_user_gas_price_btc",
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		if feeRow == nil {
			hcommon.Log.Errorf("no config int of to_user_gas_price_btc")
			return
		}
		// 生成交易
		// 输入金额
		inBalance := int64(0)
		// 输出金额
		outBalance := int64(0)
		// 初算手续费
		inputCount := int64(0)
		outputCount := int64(1)
		// 使用到的uxto的索引
		uxtoUseIndex := 0
		// 输入信息
		var inUxtoRows []*model.DBTTxBtcUxto
		// 输出信息
		var outWithdrawRows []*model.DBTWithdraw
		// 计算交易
		for _, withdrawRow := range withdrawRows {
			var tmpInUxtoRows []*model.DBTTxBtcUxto
			// 添加输出
			withdrawBalance, err := decimal.NewFromString(withdrawRow.BalanceReal)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			outBalance += withdrawBalance.Mul(decimal.NewFromInt(1e8)).IntPart()
			// 处理金额
			outputCount++
			for outBalance+(inputCount*180+outputCount*34+10+40)*feeRow.V > inBalance {
				if uxtoUseIndex >= len(uxtoRows) {
					break
				}
				if inputCount*180+outputCount*34+10+40 >= 1000000 {
					break
				}
				uxtoRow := uxtoRows[uxtoUseIndex]
				uxtoBalance, err := decimal.NewFromString(uxtoRow.VoutValue)
				if err != nil {
					hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
				inBalance += uxtoBalance.Mul(decimal.NewFromInt(1e8)).IntPart()
				uxtoUseIndex++
				inputCount++
				tmpInUxtoRows = append(tmpInUxtoRows, uxtoRow)
			}
			if outBalance+(inputCount*180+outputCount*34+10+40)*feeRow.V <= inBalance {
				// 处理这些数据
				outWithdrawRows = append(outWithdrawRows, withdrawRow)
				if len(tmpInUxtoRows) > 0 {
					inUxtoRows = append(inUxtoRows, tmpInUxtoRows...)
				}
			} else {
				break
			}
		}
		hcommon.Log.Debugf("inUxtoRows: %#v, outWithdrawRows: %#v", inUxtoRows, outWithdrawRows)
		// 创建交易
		realInBalance := int64(0)
		realOutBalance := int64(0)
		tx := wire.NewMsgTx(wire.TxVersion)
		// 提币输出
		for _, outWithdrawRow := range outWithdrawRows {
			addrTo, err := btcutil.DecodeAddress(
				outWithdrawRow.ToAddress,
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
			outBalance, err := decimal.NewFromString(outWithdrawRow.BalanceReal)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			baf := outBalance.Mul(decimal.NewFromInt(1e8)).IntPart()
			tx.AddTxOut(wire.NewTxOut(baf, pkScriptf))
			realOutBalance += baf
		}
		// 找零输出
		if len(inUxtoRows) <= 0 || inUxtoRows == nil {
			hcommon.Log.Errorf("error input arr")
			return
		}
		addrTo, err := btcutil.DecodeAddress(
			inUxtoRows[0].VoutAddress,
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
		tx.AddTxOut(wire.NewTxOut(0, pkScriptf))
		// 创建输入
		var inputAddress []string
		for _, uxtoRow := range inUxtoRows {
			hash, _ := chainhash.NewHashFromStr(uxtoRow.TxID)
			outPoint := wire.NewOutPoint(hash, uint32(uxtoRow.VoutN))
			txIn := wire.NewTxIn(outPoint, nil, nil)
			tx.AddTxIn(txIn)
			if !hcommon.IsStringInSlice(inputAddress, uxtoRow.VoutAddress) {
				inputAddress = append(inputAddress, uxtoRow.VoutAddress)
			}
		}
		// 签名,用于计算手续费
		for i, uxtoRow := range inUxtoRows {
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
			inBalance, err := decimal.NewFromString(uxtoRow.VoutValue)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			baf := inBalance.Mul(decimal.NewFromInt(1e8)).IntPart()
			realInBalance += baf
		}
		// 计算手续费
		txSize := tx.SerializeSize()
		txFee := feeRow.V * int64(txSize)
		lastOut := realInBalance - realOutBalance - txFee
		if lastOut < 0 {
			hcommon.Log.Errorf("error out put")
			return
		}
		if lastOut < 5460 {
			var newOut []*wire.TxOut
			for i, out := range tx.TxOut {
				if i != len(tx.TxOut)-2 {
					newOut = append(newOut, out)
				}
			}
			tx.TxOut = newOut
		} else {
			tx.TxOut[len(tx.TxOut)-1].Value = lastOut
		}
		// 重新计算签名
		for i, uxtoRow := range inUxtoRows {
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
		now := time.Now().Unix()
		var sendRows []*model.DBTSendBtc
		var updateUxtoRows []*model.DBTTxBtcUxto
		var updateWithdrawRows []*model.DBTWithdraw
		for i, outWithdrawRow := range outWithdrawRows {
			balance, err := decimal.NewFromString(outWithdrawRow.BalanceReal)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			gas := int64(0)
			gasPrice := int64(0)
			sendHex := ""
			if i == 0 {
				sendHex = hex.EncodeToString(b.Bytes())
				gas = int64(tx.SerializeSize())
				gasPrice = feeRow.V
			}
			inputAddressStr := strings.Join(inputAddress, ",")
			sendRows = append(sendRows, &model.DBTSendBtc{
				RelatedType:  app.SendRelationTypeWithdraw,
				RelatedID:    outWithdrawRow.ID,
				TokenID:      0,
				TxID:         tx.TxHash().String(),
				FromAddress:  inputAddressStr,
				ToAddress:    outWithdrawRow.ToAddress,
				Balance:      balance.Mul(decimal.NewFromInt(1e8)).IntPart(),
				BalanceReal:  outWithdrawRow.BalanceReal,
				Gas:          gas,
				GasPrice:     gasPrice,
				Hex:          sendHex,
				CreateTime:   now,
				HandleStatus: 0,
				HandleMsg:    "",
				HandleTime:   now,
			})
			updateWithdrawRows = append(updateWithdrawRows, &model.DBTWithdraw{
				ID:           outWithdrawRow.ID,
				TxHash:       fmt.Sprintf("%s_%d", tx.TxHash().String(), i),
				HandleStatus: app.WithdrawStatusHex,
				HandleMsg:    "hex",
				HandleTime:   now,
			})
		}
		for i, uxtoRow := range inUxtoRows {
			updateUxtoRows = append(updateUxtoRows, &model.DBTTxBtcUxto{
				ID:           uxtoRow.ID,
				TxID:         uxtoRow.TxID,
				VoutN:        uxtoRow.VoutN,
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
		// 更新withdraw
		_, err = app.SQLCreateManyTWithdrawUpdate(
			context.Background(),
			app.DbCon,
			updateWithdrawRows,
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
	})
}

// CheckTxNotify 创建btc冲币通知
func CheckTxNotify() {
	lockKey := "BtcCheckTxNotify"
	app.LockWrap(lockKey, func() {
		txRows, err := app.SQLSelectTTxBtcColByStatus(
			context.Background(),
			app.DbCon,
			[]string{
				model.DBColTTxBtcID,
				model.DBColTTxBtcProductID,
				model.DBColTTxBtcTxID,
				model.DBColTTxBtcVoutAddress,
				model.DBColTTxBtcVoutN,
				model.DBColTTxBtcVoutValue,
			},
			app.TxStatusInit,
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		var productIDs []int64
		for _, txRow := range txRows {
			if !hcommon.IsIntInSlice(productIDs, txRow.ProductID) {
				productIDs = append(productIDs, txRow.ProductID)
			}
		}
		productMap, err := app.SQLGetProductMap(
			context.Background(),
			app.DbCon,
			[]string{
				model.DBColTProductID,
				model.DBColTProductAppName,
				model.DBColTProductCbURL,
				model.DBColTProductAppSk,
			},
			productIDs,
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		var notifyTxIDs []int64
		var notifyRows []*model.DBTProductNotify
		now := time.Now().Unix()
		for _, txRow := range txRows {
			productRow, ok := productMap[txRow.ProductID]
			if !ok {
				hcommon.Log.Warnf("no productMap: %d", txRow.ProductID)
				notifyTxIDs = append(notifyTxIDs, txRow.ID)
				continue
			}
			nonce := hcommon.GetUUIDStr()
			reqObj := gin.H{
				"tx_hash":     fmt.Sprintf("%s_%d", txRow.TxID, txRow.VoutN),
				"app_name":    productRow.AppName,
				"address":     txRow.VoutAddress,
				"balance":     txRow.VoutValue,
				"symbol":      CoinSymbol,
				"notify_type": app.NotifyTypeTx,
			}
			reqObj["sign"] = hcommon.GetSign(productRow.AppSk, reqObj)
			req, err := json.Marshal(reqObj)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				continue
			}
			notifyRows = append(notifyRows, &model.DBTProductNotify{
				Nonce:        nonce,
				ProductID:    txRow.ProductID,
				ItemType:     app.SendRelationTypeTx,
				ItemID:       txRow.ID,
				NotifyType:   app.NotifyTypeTx,
				URL:          productRow.CbURL,
				Msg:          string(req),
				HandleStatus: app.NotifyStatusInit,
				HandleMsg:    "",
				CreateTime:   now,
				UpdateTime:   now,
			})
			notifyTxIDs = append(notifyTxIDs, txRow.ID)
		}
		_, err = model.SQLCreateIgnoreManyTProductNotify(
			context.Background(),
			app.DbCon,
			notifyRows,
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		_, err = app.SQLUpdateTTxBtcStatusByIDs(
			context.Background(),
			app.DbCon,
			notifyTxIDs,
			model.DBTTxBtc{
				HandleStatus: app.TxStatusNotify,
				HandleMsg:    "notify",
				HandleTime:   now,
			},
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
	})

}

// CheckGasPrice 检测gas price
func CheckGasPrice() {
	lockKey := "BtcCheckGasPrice"
	app.LockWrap(lockKey, func() {
		type StRespGasPrice struct {
			FastestFee  int64 `json:"fastestFee"`
			HalfHourFee int64 `json:"halfHourFee"`
			HourFee     int64 `json:"hourFee"`
		}
		gresp, body, errs := gorequest.New().
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
		var resp StRespGasPrice
		err := json.Unmarshal([]byte(body), &resp)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		toUserGasPrice := resp.FastestFee
		toColdGasPrice := resp.HalfHourFee
		_, err = app.SQLUpdateTAppStatusIntByK(
			context.Background(),
			app.DbCon,
			&model.DBTAppStatusInt{
				K: "to_user_gas_price_btc",
				V: toUserGasPrice,
			},
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		_, err = app.SQLUpdateTAppStatusIntByK(
			context.Background(),
			app.DbCon,
			&model.DBTAppStatusInt{
				K: "to_cold_gas_price_btc",
				V: toColdGasPrice,
			},
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
	})
}

// OmniCheckBlockSeek 检测到账
func OmniCheckBlockSeek() {
	lockKey := "OmniCheckBlockSeek"
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
			"omni_seek_num",
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		if seekRow == nil {
			hcommon.Log.Errorf("no config int of omni_seek_num")
			return
		}
		rpcBlockNum, err := omniclient.RpcGetBlockCount()
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		startI := seekRow.V + 1
		endI := rpcBlockNum - confirmRow.V + 1
		hcommon.Log.Debugf("omni block seek %d->%d", startI, endI)
		if startI < endI {
			// 获取所有token
			var tokenIndexes []int64
			tokenMap := make(map[int64]*model.DBTAppConfigTokenBtc)
			tokenRows, err := app.SQLSelectTAppConfigTokenBtcColAll(
				context.Background(),
				app.DbCon,
				[]string{
					model.DBColTAppConfigTokenBtcID,
					model.DBColTAppConfigTokenBtcTokenIndex,
					model.DBColTAppConfigTokenBtcTokenSymbol,
				},
			)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			for _, tokenRow := range tokenRows {
				tokenIndexes = append(tokenIndexes, tokenRow.TokenIndex)
				tokenMap[tokenRow.TokenIndex] = tokenRow
			}
			// 遍历获取需要查询的block信息
			for i := startI; i < endI; i++ {
				rpcTransactionHashes, err := omniclient.RpcOmniListBlockTransactions(i)
				if err != nil {
					hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
				var toAddresses []string
				toAddressMap := make(map[string][]*omniclient.StOmniTx)

				for _, rpcTransactionHash := range rpcTransactionHashes {
					rpcTx, err := omniclient.RpcOmniGetTransaction(rpcTransactionHash)
					if err != nil {
						hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
						return
					}
					// type_int 0 Simple Send
					if rpcTx.TypeInt == 0 && rpcTx.Valid && rpcTx.Confirmations > 0 {
						// 验证成功
						if hcommon.IsIntInSlice(tokenIndexes, rpcTx.Propertyid) {
							// 是关注的代币类型
							if !hcommon.IsStringInSlice(toAddresses, rpcTx.Referenceaddress) {
								toAddresses = append(toAddresses, rpcTx.Referenceaddress)
							}
							toAddressMap[rpcTx.Referenceaddress] = append(
								toAddressMap[rpcTx.Referenceaddress],
								rpcTx,
							)
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
					toAddresses,
				)
				if err != nil {
					hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
				now := time.Now().Unix()
				var txTokenRows []*model.DBTTxBtcToken
				for _, dbAddressRow := range dbAddressRows {
					if dbAddressRow.UseTag >= 0 {
						toAddress := dbAddressRow.Address
						rpcTxes := toAddressMap[toAddress]
						for _, rpcTx := range rpcTxes {
							tokenRow, ok := tokenMap[rpcTx.Propertyid]
							if !ok {
								hcommon.Log.Errorf("no btc token: %d", rpcTx.Propertyid)
								return
							}
							txTokenRows = append(txTokenRows, &model.DBTTxBtcToken{
								TokenIndex:   rpcTx.Propertyid,
								TokenSymbol:  tokenRow.TokenSymbol,
								BlockHash:    rpcTx.Blockhash,
								TxID:         rpcTx.Txid,
								FromAddress:  rpcTx.Sendingaddress,
								ToAddress:    rpcTx.Referenceaddress,
								Value:        rpcTx.Amount,
								Blocktime:    rpcTx.Blocktime,
								CreateAt:     now,
								HandleStatus: 0,
								HandleMsg:    "",
								HandleAt:     0,
								OrgStatus:    0,
								OrgMsg:       "",
								OrgAt:        0,
							})
						}
					}

				}
				_, err = model.SQLCreateIgnoreManyTTxBtcToken(
					context.Background(),
					app.DbCon,
					txTokenRows,
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
