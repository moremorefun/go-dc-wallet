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
	keySignTx, err := kb.Sign(signTx, chainID, keys[0])
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return err
	}
	// 打包tx
	packedTx, err := keySignTx.Pack(eos.CompressionNone)
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
			Hex:          hex.EncodeToString(packedTxBs),
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
