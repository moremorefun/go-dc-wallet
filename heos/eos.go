package heos

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"go-dc-wallet/app"
	"go-dc-wallet/app/model"
	"go-dc-wallet/eosclient"
	"go-dc-wallet/hcommon"
	"time"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/token"
	"github.com/shopspring/decimal"

	"github.com/gin-gonic/gin"
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
								hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
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
								rpcActionData.Quantity, err = EosValueToStr(rpcActionData.Quantity)
								if err != nil {
									hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
									return
								}
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
				_, err = app.SQLUpdateTAppStatusIntByKGreater(
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

// CheckTxNotify 创建冲币通知
func CheckTxNotify() {
	lockKey := "EosCheckTxNotify"
	app.LockWrap(lockKey, func() {
		txRows, err := app.SQLSelectTTxEosColByStatus(
			context.Background(),
			app.DbCon,
			[]string{
				model.DBColTTxEosID,
				model.DBColTTxEosProductID,
				model.DBColTTxEosTxHash,
				model.DBColTTxEosLogIndex,
				model.DBColTTxEosFromAddress,
				model.DBColTTxEosToAddress,
				model.DBColTTxEosMemo,
				model.DBColTTxEosBalanceReal,
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
				"tx_hash":     fmt.Sprintf("%s_%d", txRow.TxHash, txRow.LogIndex),
				"app_name":    productRow.AppName,
				"address":     txRow.ToAddress,
				"memo":        txRow.Memo,
				"balance":     txRow.BalanceReal,
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
				TokenSymbol:  CoinSymbol,
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
		_, err = app.SQLUpdateTTxEosStatusByIDs(
			context.Background(),
			app.DbCon,
			notifyTxIDs,
			model.DBTTxEos{
				HandleStatus: app.TxStatusNotify,
				HandleMsg:    "notify",
				HandleAt:     now,
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
	lockKey := "EosCheckWithdraw"
	app.LockWrap(lockKey, func() {
		// 获取需要处理的提币数据
		withdrawRows, err := app.SQLSelectTWithdrawColByStatus(
			context.Background(),
			app.DbCon,
			[]string{
				model.DBColTWithdrawID,
			},
			app.WithdrawStatusInit,
			[]string{CoinSymbol},
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		if len(withdrawRows) == 0 {
			// 没有要处理的提币
			return
		}
		// 获取热钱包地址
		hotAddressValue, err := app.SQLGetTAppConfigStrValueByK(
			context.Background(),
			app.DbCon,
			"hot_wallet_address_eos",
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		// 获取热钱包私钥
		hotKeyValue, err := app.SQLGetTAppConfigStrValueByK(
			context.Background(),
			app.DbCon,
			"hot_wallet_key_eos",
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		key := hcommon.AesDecrypt(hotKeyValue, app.Cfg.AESKey)
		if len(key) == 0 {
			hcommon.Log.Errorf("error key of eos")
			return
		}
		// 获取热钱包余额
		rpcAccount, err := eosclient.RpcChainGetAccount(
			hotAddressValue,
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		rpcHotBalance, err := EosValueToDecimal(rpcAccount.CoreLiquidBalance)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		pendingBalanceRealStr, err := app.SQLGetTSendEosPendingBalanceReal(
			context.Background(),
			app.DbCon,
			hotAddressValue,
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		pendingBalanceReal, err := StrToEosDecimal(pendingBalanceRealStr)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		rpcHotBalance = rpcHotBalance.Sub(pendingBalanceReal)
		// 获取链信息
		rpcChainInfo, err := eosclient.RpcChainGetInfo()
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		for _, withdrawRow := range withdrawRows {
			err = handleWithdraw(rpcChainInfo, withdrawRow.ID, hotAddressValue, key, &rpcHotBalance)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				continue
			}
		}
	})
}

func handleWithdraw(rpcChainInfo *eosclient.StChainGetInfo, withdrawID int64, hotAddressValue string, hotKey string, hotBalance *decimal.Decimal) error {
	isComment := false
	dbTx, err := app.DbCon.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer func() {
		if !isComment {
			_ = dbTx.Rollback()
		}
	}()
	// 处理业务
	withdrawRow, err := app.SQLGetTWithdrawColForUpdate(
		context.Background(),
		dbTx,
		[]string{
			model.DBColTWithdrawID,
			model.DBColTWithdrawBalanceReal,
			model.DBColTWithdrawToAddress,
			model.DBColTWithdrawMemo,
		},
		withdrawID,
		app.WithdrawStatusInit,
	)
	if err != nil {
		return err
	}
	if withdrawRow == nil {
		return nil
	}
	withdrawBalance, err := StrToEosDecimal(withdrawRow.BalanceReal)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return nil
	}
	*hotBalance = (*hotBalance).Sub(withdrawBalance)
	if (*hotBalance).Cmp(decimal.NewFromInt(0)) < 0 {
		// 金额不够
		hcommon.Log.Errorf("eos hot balance limit")
		*hotBalance = (*hotBalance).Add(withdrawBalance)
		return nil
	}
	eosAesset, err := eos.NewEOSAssetFromString(withdrawRow.BalanceReal)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return err
	}
	action := token.NewTransfer(
		eos.AccountName(hotAddressValue),
		eos.AccountName(withdrawRow.ToAddress),
		eosAesset,
		withdrawRow.Memo,
	)
	action.Account = "eosio.token"
	actions := []*eos.Action{action}
	// 设置tx属性
	chainID, err := hex.DecodeString(rpcChainInfo.ChainID)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return err
	}
	headBlockID, err := hex.DecodeString(rpcChainInfo.HeadBlockID)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return err
	}
	opts := &eos.TxOptions{
		ChainID:     chainID,
		HeadBlockID: headBlockID,
	}
	// 创建tx
	tx := eos.NewTransaction(actions, opts)
	tx.SetExpiration(time.Second * 3000)
	tx.ContextFreeActions = append(
		tx.ContextFreeActions,
		&eos.Action{
			Account:    "eosio.null",
			Name:       "nonce",
			ActionData: eos.NewActionDataFromHexData([]byte(fmt.Sprintf("%d", withdrawRow.ID))),
		},
	)
	// 生成待签名tx
	signTx := eos.NewSignedTransaction(tx)
	// 创建密钥对
	kb := eos.NewKeyBag()
	err = kb.Add(hotKey)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return err
	}
	keys, err := kb.AvailableKeys()
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return err
	}
	_, err = kb.Sign(signTx, chainID, keys[0])
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return err
	}
	// 打包tx
	packedTx, err := signTx.Pack(eos.CompressionNone)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return err
	}
	packedTxBs, err := json.Marshal(packedTx)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return err
	}
	packedHash, err := packedTx.ID()
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return err
	}
	txHash := packedHash.String()
	now := time.Now().Unix()
	_, err = app.SQLUpdateTWithdrawGenTx(
		context.Background(),
		dbTx,
		&model.DBTWithdraw{
			ID:           withdrawID,
			TxHash:       txHash,
			HandleStatus: app.WithdrawStatusHex,
			HandleMsg:    "gen tx hex",
			HandleTime:   now,
		},
	)
	if err != nil {
		return err
	}
	_, err = model.SQLCreateTSendEos(
		context.Background(),
		dbTx,
		&model.DBTSendEos{
			WithdrawID:   withdrawID,
			TxHash:       txHash,
			FromAddress:  hotAddressValue,
			ToAddress:    withdrawRow.ToAddress,
			Memo:         withdrawRow.Memo,
			BalanceReal:  withdrawRow.BalanceReal,
			Hex:          string(packedTxBs),
			CreateTime:   now,
			HandleStatus: app.SendStatusInit,
			HandleMsg:    "init",
			HandleAt:     now,
		},
	)
	if err != nil {
		return err
	}
	// 处理完成
	err = dbTx.Commit()
	if err != nil {
		return err
	}
	isComment = true
	return nil
}

// CheckRawTxSend 发送交易
func CheckRawTxSend() {
	lockKey := "EosCheckRawTxSend"
	app.LockWrap(lockKey, func() {
		// 获取待发送的数据
		sendRows, err := app.SQLSelectTSendEosColByStatus(
			context.Background(),
			app.DbCon,
			[]string{
				model.DBColTSendEosID,
				model.DBColTSendEosTxHash,
				model.DBColTSendEosHex,
				model.DBColTSendEosWithdrawID,
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
			if !hcommon.IsIntInSlice(withdrawIDs, sendRow.WithdrawID) {
				withdrawIDs = append(withdrawIDs, sendRow.WithdrawID)
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
		// 执行发送
		var sendIDs []int64
		withdrawIDs = []int64{}
		// 通知数据
		var notifyRows []*model.DBTProductNotify
		now := time.Now().Unix()
		onSendOk := func(sendRow *model.DBTSendEos) error {
			// 将发送成功和占位数据计入数组
			if !hcommon.IsIntInSlice(sendIDs, sendRow.ID) {
				sendIDs = append(sendIDs, sendRow.ID)
			}
			// 如果是提币，创建通知信息
			withdrawRow, ok := withdrawMap[sendRow.WithdrawID]
			if !ok {
				hcommon.Log.Errorf("withdrawMap no: %d", sendRow.WithdrawID)
				return nil
			}
			productRow, ok := productMap[withdrawRow.ProductID]
			if !ok {
				hcommon.Log.Errorf("productMap no: %d", withdrawRow.ProductID)
				return nil
			}
			nonce := hcommon.GetUUIDStr()
			reqObj := gin.H{
				"tx_hash":     sendRow.TxHash,
				"balance":     withdrawRow.BalanceReal,
				"app_name":    productRow.AppName,
				"out_serial":  withdrawRow.OutSerial,
				"address":     withdrawRow.ToAddress,
				"symbol":      withdrawRow.Symbol,
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
				TokenSymbol:  withdrawRow.Symbol,
				URL:          productRow.CbURL,
				Msg:          string(req),
				HandleStatus: app.NotifyStatusInit,
				HandleMsg:    "",
				CreateTime:   now,
				UpdateTime:   now,
			})
			withdrawIDs = append(withdrawIDs, withdrawRow.ID)
			return nil
		}
		for _, sendRow := range sendRows {
			// 判定是否已经发送过
			isSend := false
			_, err := eosclient.RpcHistoryGetTransaction(sendRow.TxHash)
			if err != nil {
				rpcErr, ok := err.(*eosclient.StRpcRespError)
				if !ok {
					hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
				if rpcErr.ErrorInv.Code == 3040011 {
					// tx_not_found
					// 还没有发送
				} else {
					hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
			} else {
				// 已经发送
				isSend = true
				err = onSendOk(sendRow)
				if err != nil {
					hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				}
			}
			// 发送数据中需要排除占位数据
			if !isSend && sendRow.Hex != "" {
				var args eosclient.StPushTransactionArg
				err := json.Unmarshal([]byte(sendRow.Hex), &args)
				if err != nil {
					hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					continue
				}
				_, err = eosclient.RpcChainPushTransaction(
					args,
				)
				if err != nil {
					rpcErr, ok := err.(*eosclient.StRpcRespError)
					if !ok {
						hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
						continue
					}
					switch rpcErr.ErrorInv.Code {
					case 3080001:
						// account using more than allotted RAM usage
						// rom 不足
						hcommon.Log.Errorf("eos hot rom limit")
						return
					case 3080002:
						// Transaction exceeded the current network usage limit imposed on the transaction
						// net 不足
						hcommon.Log.Errorf("eos hot net limit")
						return
					case 3080004:
						// Transaction exceeded the current CPU usage limit imposed on the transaction
						// cpu 不足
						hcommon.Log.Errorf("eos hot cpu limit")
						return
					case 3040008:
						// Duplicate transaction
						// 已经发送
					default:
						hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
						continue
					}
				}
				err = onSendOk(sendRow)
				if err != nil {
					hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
			}
		}
		// 插入通知
		_, err = model.SQLCreateIgnoreManyTProductNotify(
			context.Background(),
			app.DbCon,
			notifyRows,
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
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
		// 更新发送状态
		_, err = app.SQLUpdateTSendEosStatusByIDs(
			context.Background(),
			app.DbCon,
			sendIDs,
			model.DBTSendEos{
				HandleStatus: app.SendStatusSend,
				HandleMsg:    "send",
				HandleAt:     now,
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
	lockKey := "EosCheckRawTxConfirm"
	app.LockWrap(lockKey, func() {
		// 获取待发送的数据
		sendRows, err := app.SQLSelectTSendEosColByStatus(
			context.Background(),
			app.DbCon,
			[]string{
				model.DBColTSendEosID,
				model.DBColTSendEosTxHash,
				model.DBColTSendEosHex,
				model.DBColTSendEosWithdrawID,
			},
			app.SendStatusSend,
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		var withdrawIDs []int64
		for _, sendRow := range sendRows {
			// 提币
			if !hcommon.IsIntInSlice(withdrawIDs, sendRow.WithdrawID) {
				withdrawIDs = append(withdrawIDs, sendRow.WithdrawID)
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

		now := time.Now().Unix()
		var notifyRows []*model.DBTProductNotify
		var sendIDs []int64
		withdrawIDs = []int64{}
		for _, sendRow := range sendRows {
			_, err := eosclient.RpcHistoryGetTransaction(
				sendRow.TxHash,
			)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				continue
			}
			// 提币
			withdrawRow, ok := withdrawMap[sendRow.WithdrawID]
			if !ok {
				hcommon.Log.Errorf("no withdrawMap: %d", sendRow.WithdrawID)
				return
			}
			productRow, ok := productMap[withdrawRow.ProductID]
			if !ok {
				hcommon.Log.Errorf("no productMap: %d", withdrawRow.ProductID)
				return
			}
			nonce := hcommon.GetUUIDStr()
			reqObj := gin.H{
				"tx_hash":     sendRow.WithdrawID,
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
				return
			}
			notifyRows = append(notifyRows, &model.DBTProductNotify{
				Nonce:        nonce,
				ProductID:    withdrawRow.ProductID,
				ItemType:     app.SendRelationTypeWithdraw,
				ItemID:       withdrawRow.ID,
				NotifyType:   app.NotifyTypeWithdrawConfirm,
				TokenSymbol:  withdrawRow.Symbol,
				URL:          productRow.CbURL,
				Msg:          string(req),
				HandleStatus: app.NotifyStatusInit,
				HandleMsg:    "",
				CreateTime:   now,
				UpdateTime:   now,
			})
			// 将发送成功和占位数据计入数组
			if !hcommon.IsIntInSlice(sendIDs, sendRow.ID) {
				sendIDs = append(sendIDs, sendRow.ID)
			}
			if !hcommon.IsIntInSlice(withdrawIDs, sendRow.WithdrawID) {
				withdrawIDs = append(withdrawIDs, sendRow.WithdrawID)
			}
		}
		// 添加通知信息
		_, err = model.SQLCreateIgnoreManyTProductNotify(
			context.Background(),
			app.DbCon,
			notifyRows,
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		// 更新提币状态
		_, err = app.SQLUpdateTWithdrawStatusByIDs(
			context.Background(),
			app.DbCon,
			withdrawIDs,
			&model.DBTWithdraw{
				HandleStatus: app.WithdrawStatusConfirm,
				HandleMsg:    "confirmed",
				HandleTime:   now,
			},
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		// 更新发送状态
		_, err = app.SQLUpdateTSendEosStatusByIDs(
			context.Background(),
			app.DbCon,
			sendIDs,
			model.DBTSendEos{
				HandleStatus: app.SendStatusConfirm,
				HandleMsg:    "confirmed",
				HandleAt:     now,
			},
		)
	})
}
