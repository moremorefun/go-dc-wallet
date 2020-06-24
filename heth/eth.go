package heth

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"go-dc-wallet/app"
	"go-dc-wallet/app/model"
	"go-dc-wallet/ethclient"
	"go-dc-wallet/hcommon"
	"math"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/parnurzeal/gorequest"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/gin-gonic/gin"

	"github.com/ethereum/go-ethereum/rlp"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/shopspring/decimal"
)

// CheckAddressFree 检测是否有充足的备用地址
func CheckAddressFree() {
	lockKey := "EthCheckAddressFree"
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
			"eth",
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
				privateKey, err := crypto.GenerateKey()
				if err != nil {
					hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
				privateKeyBytes := crypto.FromECDSA(privateKey)
				privateKeyStr := hexutil.Encode(privateKeyBytes)
				// 加密密钥
				privateKeyStrEn := hcommon.AesEncrypt(privateKeyStr, app.Cfg.AESKey)
				// 获取地址
				publicKey := privateKey.Public()
				publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
				if !ok {
					hcommon.Log.Errorf("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
					return
				}
				// 地址全部储存为小写方便处理
				address := AddressBytesToStr(crypto.PubkeyToAddress(*publicKeyECDSA))
				// 存入待添加队列
				rows = append(rows, &model.DBTAddressKey{
					Symbol:  "eth",
					Address: address,
					Pwd:     privateKeyStrEn,
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
	lockKey := "EthCheckBlockSeek"
	app.LockWrap(lockKey, func() {
		// 获取配置 延迟确认数
		confirmRow, err := app.SQLGetTAppConfigIntByK(
			context.Background(),
			app.DbCon,
			"block_confirm_num",
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		if confirmRow == nil {
			hcommon.Log.Errorf("no config int of min_free_address")
			return
		}
		// 获取状态 当前处理完成的最新的block number
		seekRow, err := app.SQLGetTAppStatusIntByK(
			context.Background(),
			app.DbCon,
			"seek_num",
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		if seekRow == nil {
			hcommon.Log.Errorf("no config int of seek_num")
			return
		}
		// rpc 获取当前最新区块数
		rpcBlockNum, err := ethclient.RpcBlockNumber(context.Background())
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		startI := seekRow.V + 1
		endI := rpcBlockNum - confirmRow.V + 1
		if startI < endI {
			// 手续费钱包列表
			feeRow, err := app.SQLGetTAppConfigStrByK(
				context.Background(),
				app.DbCon,
				"fee_wallet_address_list",
			)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			if feeRow == nil {
				hcommon.Log.Errorf("no config int of fee_wallet_address_list")
				return
			}
			addresses := strings.Split(feeRow.V, ",")
			var feeAddresses []string
			for _, address := range addresses {
				if address == "" {
					continue
				}
				feeAddresses = append(feeAddresses, address)
			}
			// 遍历获取需要查询的block信息
			for i := startI; i < endI; i++ {
				// rpc获取block信息
				//hcommon.Log.Debugf("eth check block: %d", i)
				rpcBlock, err := ethclient.RpcBlockByNum(context.Background(), i)
				if err != nil {
					hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
				// 接收地址列表
				var toAddresses []string
				// map[接收地址] => []交易信息
				toAddressTxMap := make(map[string][]*types.Transaction)
				// 遍历block中的tx
				for _, rpcTx := range rpcBlock.Transactions() {
					// 转账数额大于0 and 不是创建合约交易
					if rpcTx.Value().Int64() > 0 && rpcTx.To() != nil {
						msg, err := rpcTx.AsMessage(types.NewEIP155Signer(rpcTx.ChainId()))
						if err != nil {
							hcommon.Log.Errorf("AsMessage err: [%T] %s", err, err.Error())
							return
						}
						if hcommon.IsStringInSlice(feeAddresses, AddressBytesToStr(msg.From())) {
							// 如果打币地址在手续费热钱包地址则不处理
							continue
						}
						toAddress := AddressBytesToStr(*(rpcTx.To()))
						toAddressTxMap[toAddress] = append(toAddressTxMap[toAddress], rpcTx)
						if !hcommon.IsStringInSlice(toAddresses, toAddress) {
							toAddresses = append(toAddresses, toAddress)
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
				// 待插入数据
				var dbTxRows []*model.DBTTx
				// map[接收地址] => 产品id
				addressProductMap := make(map[string]int64)
				for _, dbAddressRow := range dbAddressRows {
					addressProductMap[dbAddressRow.Address] = dbAddressRow.UseTag
				}
				// 时间
				now := time.Now().Unix()
				// 遍历数据库中有交易的地址
				for _, dbAddressRow := range dbAddressRows {
					if dbAddressRow.UseTag < 0 {
						continue
					}
					// 获取地址对应的交易列表
					txes := toAddressTxMap[dbAddressRow.Address]
					for _, tx := range txes {
						msg, err := tx.AsMessage(types.NewEIP155Signer(tx.ChainId()))
						if err != nil {
							hcommon.Log.Errorf("AsMessage err: [%T] %s", err, err.Error())
							return
						}
						fromAddress := AddressBytesToStr(msg.From())
						toAddress := AddressBytesToStr(*(tx.To()))
						balanceReal := decimal.NewFromInt(tx.Value().Int64()).Div(decimal.NewFromInt(1e18))
						dbTxRows = append(dbTxRows, &model.DBTTx{
							ProductID:    addressProductMap[toAddress],
							TxID:         tx.Hash().String(),
							FromAddress:  fromAddress,
							ToAddress:    toAddress,
							Balance:      tx.Value().Int64(),
							BalanceReal:  balanceReal.String(),
							CreateTime:   now,
							HandleStatus: app.TxStatusInit,
							HandleMsg:    "",
							HandleTime:   now,
							OrgStatus:    app.TxOrgStatusInit,
							OrgMsg:       "",
							OrgTime:      now,
						})
					}
				}
				// 插入交易数据
				_, err = model.SQLCreateIgnoreManyTTx(
					context.Background(),
					app.DbCon,
					dbTxRows,
				)
				if err != nil {
					hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
				// 更新检查到的最新区块数
				_, err = app.SQLUpdateTAppStatusIntByK(
					context.Background(),
					app.DbCon,
					&model.DBTAppStatusInt{
						K: "seek_num",
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

// CheckAddressOrg 零钱整理到冷钱包
func CheckAddressOrg() {
	lockKey := "EthCheckAddressOrg"
	app.LockWrap(lockKey, func() {
		// 获取冷钱包地址
		coldRow, err := app.SQLGetTAppConfigStrByK(
			context.Background(),
			app.DbCon,
			"cold_wallet_address",
		)
		if err != nil {
			hcommon.Log.Warnf("SQLGetTAppConfigInt err: [%T] %s", err, err.Error())
			return
		}
		if coldRow == nil {
			hcommon.Log.Errorf("no config int of cold_wallet_address")
			return
		}
		coldAddress, err := StrToAddressBytes(coldRow.V)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		// 获取待整理的地址列表
		txRows, err := app.SQLSelectTTxColByOrg(
			context.Background(),
			app.DbCon,
			[]string{
				model.DBColTTxID,
				model.DBColTTxToAddress,
				model.DBColTTxBalance,
			},
			app.TxOrgStatusInit,
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		if len(txRows) <= 0 {
			// 没有要处理的信息
			return
		}
		// 将待整理地址按地址做归并处理
		type OrgInfo struct {
			RowIDs  []int64 // db t_tx.id
			Balance int64   // 金额
		}
		// addressMap map[地址] = []整理信息
		addressMap := make(map[string]*OrgInfo)
		// 获取gap price
		gasRow, err := app.SQLGetTAppStatusIntByK(
			context.Background(),
			app.DbCon,
			"to_cold_gas_price",
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		if gasRow == nil {
			hcommon.Log.Errorf("no config int of to_cold_gas_price")
			return
		}
		gasPrice := gasRow.V
		gasLimit := int64(21000)
		feeValue := gasLimit * gasPrice
		// addresses 需要整理的地址列表
		var addresses []string
		for _, txRow := range txRows {
			info := addressMap[txRow.ToAddress]
			if info == nil {
				info = &OrgInfo{
					RowIDs:  []int64{},
					Balance: 0,
				}
				addressMap[txRow.ToAddress] = info
			}
			info.RowIDs = append(info.RowIDs, txRow.ID)
			info.Balance += txRow.Balance
			if !hcommon.IsStringInSlice(addresses, txRow.ToAddress) {
				addresses = append(addresses, txRow.ToAddress)
			}
		}
		chainID, err := ethclient.RpcNetworkID(context.Background())
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		now := time.Now().Unix()
		for address, info := range addressMap {
			// 获取私钥
			keyRow, err := app.SQLGetTAddressKeyColByAddress(
				context.Background(),
				app.DbCon,
				[]string{
					model.DBColTAddressKeyPwd,
				},
				address,
			)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			if keyRow == nil {
				hcommon.Log.Errorf("no key of: %s", address)
				continue
			}
			key := hcommon.AesDecrypt(keyRow.Pwd, app.Cfg.AESKey)
			if len(key) == 0 {
				hcommon.Log.Errorf("error key of: %s", address)
				continue
			}
			if strings.HasPrefix(key, "0x") {
				key = key[2:]
			}
			privateKey, err := crypto.HexToECDSA(key)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				continue
			}
			// 获取nonce值
			nonce, err := GetNonce(app.DbCon, address)
			if err != nil {
				hcommon.Log.Errorf("GetNonce err: [%T] %s", err, err.Error())
				return
			}
			// 发送数量
			sendBalance := info.Balance - feeValue
			if sendBalance <= 0 {
				continue
			}
			sendBalanceReal := decimal.NewFromInt(sendBalance).Div(decimal.NewFromInt(1e18))
			// 生成tx
			var data []byte
			tx := types.NewTransaction(
				uint64(nonce),
				coldAddress,
				big.NewInt(sendBalance),
				uint64(gasLimit),
				big.NewInt(gasPrice),
				data,
			)
			// 签名
			signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(chainID)), privateKey)
			if err != nil {
				hcommon.Log.Warnf("RpcNetworkID err: [%T] %s", err, err.Error())
				return
			}
			ts := types.Transactions{signedTx}
			rawTxBytes := ts.GetRlp(0)
			rawTxHex := hex.EncodeToString(rawTxBytes)
			txHash := strings.ToLower(signedTx.Hash().Hex())
			// 创建存入数据
			var sendRows []*model.DBTSend
			for rowIndex, rowID := range info.RowIDs {
				if rowIndex == 0 {
					// 只有第一条数据需要发送，其余数据为占位数据
					sendRows = append(sendRows, &model.DBTSend{
						RelatedType:  app.SendRelationTypeTx,
						RelatedID:    rowID,
						TxID:         txHash,
						FromAddress:  address,
						ToAddress:    coldRow.V,
						Balance:      sendBalance,
						BalanceReal:  sendBalanceReal.String(),
						Gas:          gasLimit,
						GasPrice:     gasPrice,
						Nonce:        nonce,
						Hex:          rawTxHex,
						CreateTime:   now,
						HandleStatus: app.SendStatusInit,
						HandleMsg:    "",
						HandleTime:   now,
					})
				} else {
					// 占位数据
					sendRows = append(sendRows, &model.DBTSend{
						RelatedType:  app.SendRelationTypeTx,
						RelatedID:    rowID,
						TxID:         txHash,
						FromAddress:  address,
						ToAddress:    coldRow.V,
						Balance:      0,
						BalanceReal:  "",
						Gas:          0,
						GasPrice:     0,
						Nonce:        -1,
						Hex:          "",
						CreateTime:   now,
						HandleStatus: app.SendStatusInit,
						HandleMsg:    "",
						HandleTime:   now,
					})
				}
			}
			// 插入数据
			_, err = model.SQLCreateIgnoreManyTSend(
				context.Background(),
				app.DbCon,
				sendRows,
			)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			// 更改状态
			_, err = app.SQLUpdateTTxOrgStatusByIDs(
				context.Background(),
				app.DbCon,
				info.RowIDs,
				model.DBTTx{
					OrgStatus: app.TxOrgStatusHex,
					OrgMsg:    "gen raw tx",
					OrgTime:   now,
				},
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
	lockKey := "EthCheckRawTxSend"
	app.LockWrap(lockKey, func() {
		// 获取待发送的数据
		sendRows, err := app.SQLSelectTSendColByStatus(
			context.Background(),
			app.DbCon,
			[]string{
				model.DBColTSendID,
				model.DBColTSendTxID,
				model.DBColTSendHex,
				model.DBColTSendRelatedType,
				model.DBColTSendRelatedID,
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
		var txIDs []int64
		var erc20TxIDs []int64
		var erc20TxFeeIDs []int64
		withdrawIDs = []int64{}
		// 通知数据
		var notifyRows []*model.DBTProductNotify
		now := time.Now().Unix()
		for _, sendRow := range sendRows {
			// 发送数据中需要排除占位数据
			if sendRow.Hex != "" {
				rawTxBytes, err := hex.DecodeString(sendRow.Hex)
				if err != nil {
					hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					continue
				}
				tx := new(types.Transaction)
				err = rlp.DecodeBytes(rawTxBytes, &tx)
				if err != nil {
					hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					continue
				}
				err = ethclient.RpcSendTransaction(
					context.Background(),
					tx,
				)
				if err != nil {
					if !strings.Contains(err.Error(), "known transaction") {
						hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
						continue
					}
				}
			}
			// 将发送成功和占位数据计入数组
			if !hcommon.IsIntInSlice(sendIDs, sendRow.ID) {
				sendIDs = append(sendIDs, sendRow.ID)
			}
			switch sendRow.RelatedType {
			case app.SendRelationTypeTx:
				if !hcommon.IsIntInSlice(txIDs, sendRow.RelatedID) {
					txIDs = append(txIDs, sendRow.RelatedID)
				}
			case app.SendRelationTypeWithdraw:
				if !hcommon.IsIntInSlice(withdrawIDs, sendRow.RelatedID) {
					withdrawIDs = append(withdrawIDs, sendRow.RelatedID)
				}
			case app.SendRelationTypeTxErc20:
				if !hcommon.IsIntInSlice(erc20TxIDs, sendRow.RelatedID) {
					erc20TxIDs = append(erc20TxIDs, sendRow.RelatedID)
				}
			case app.SendRelationTypeTxErc20Fee:
				if !hcommon.IsIntInSlice(erc20TxFeeIDs, sendRow.RelatedID) {
					erc20TxFeeIDs = append(erc20TxFeeIDs, sendRow.RelatedID)
				}
			}
			// 如果是提币，创建通知信息
			if sendRow.RelatedType == app.SendRelationTypeWithdraw {
				withdrawRow, ok := withdrawMap[sendRow.RelatedID]
				if !ok {
					hcommon.Log.Errorf("withdrawMap no: %d", sendRow.RelatedID)
					continue
				}
				productRow, ok := productMap[withdrawRow.ProductID]
				if !ok {
					hcommon.Log.Errorf("productMap no: %d", withdrawRow.ProductID)
					continue
				}
				nonce := hcommon.GetUUIDStr()
				reqObj := gin.H{
					"tx_hash":     sendRow.TxID,
					"balance":     withdrawRow.BalanceReal,
					"app_name":    productRow.AppName,
					"out_serial":  withdrawRow.OutSerial,
					"address":     withdrawRow.ToAddress,
					"symbol":      "eth",
					"notify_type": app.NotifyTypeWithdrawSend,
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
					NotifyType:   app.NotifyTypeWithdrawSend,
					TokenSymbol:  "eth",
					URL:          productRow.CbURL,
					Msg:          string(req),
					HandleStatus: app.NotifyStatusInit,
					HandleMsg:    "",
					CreateTime:   now,
					UpdateTime:   now,
				})
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
		// 更新eth零钱整理状态
		_, err = app.SQLUpdateTTxOrgStatusByIDs(
			context.Background(),
			app.DbCon,
			txIDs,
			model.DBTTx{
				OrgStatus: app.TxOrgStatusSend,
				OrgMsg:    "send",
				OrgTime:   now,
			},
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		// 更新erc20零钱整理状态
		_, err = app.SQLUpdateTTxOrgStatusByIDs(
			context.Background(),
			app.DbCon,
			erc20TxIDs,
			model.DBTTx{
				OrgStatus: app.TxOrgStatusSend,
				OrgMsg:    "send",
				OrgTime:   now,
			},
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		// 更新erc20手续费状态
		_, err = app.SQLUpdateTTxOrgStatusByIDs(
			context.Background(),
			app.DbCon,
			erc20TxFeeIDs,
			model.DBTTx{
				OrgStatus: app.TxOrgStatusFeeSend,
				OrgMsg:    "send",
				OrgTime:   now,
			},
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		// 更新发送状态
		_, err = app.SQLUpdateTSendStatusByIDs(
			context.Background(),
			app.DbCon,
			sendIDs,
			model.DBTSend{
				HandleStatus: app.SendStatusSend,
				HandleMsg:    "send",
				HandleTime:   now,
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
	lockKey := "EthCheckRawTxConfirm"
	app.LockWrap(lockKey, func() {
		sendRows, err := app.SQLSelectTSendColByStatus(
			context.Background(),
			app.DbCon,
			[]string{
				model.DBColTSendID,
				model.DBColTSendRelatedType,
				model.DBColTSendRelatedID,
				model.DBColTSendID,
				model.DBColTSendTxID,
			},
			app.SendStatusSend,
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
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
		var txIDs []int64
		var erc20TxIDs []int64
		var erc20TxFeeIDs []int64
		withdrawIDs = []int64{}
		for _, sendRow := range sendRows {
			rpcTx, err := ethclient.RpcTransactionByHash(
				context.Background(),
				sendRow.TxID,
			)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				continue
			}
			if rpcTx == nil {
				continue
			}
			if sendRow.RelatedType == app.SendRelationTypeWithdraw {
				// 提币
				withdrawRow, ok := withdrawMap[sendRow.RelatedID]
				if !ok {
					hcommon.Log.Errorf("no withdrawMap: %d", sendRow.RelatedID)
					return
				}
				productRow, ok := productMap[withdrawRow.ProductID]
				if !ok {
					hcommon.Log.Errorf("no productMap: %d", withdrawRow.ProductID)
					return
				}
				nonce := hcommon.GetUUIDStr()
				reqObj := gin.H{
					"tx_hash":     sendRow.TxID,
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

			}
			// 将发送成功和占位数据计入数组
			if !hcommon.IsIntInSlice(sendIDs, sendRow.ID) {
				sendIDs = append(sendIDs, sendRow.ID)
			}
			switch sendRow.RelatedType {
			case app.SendRelationTypeTx:
				if !hcommon.IsIntInSlice(txIDs, sendRow.RelatedID) {
					txIDs = append(txIDs, sendRow.RelatedID)
				}
			case app.SendRelationTypeWithdraw:
				if !hcommon.IsIntInSlice(withdrawIDs, sendRow.RelatedID) {
					withdrawIDs = append(withdrawIDs, sendRow.RelatedID)
				}
			case app.SendRelationTypeTxErc20:
				if !hcommon.IsIntInSlice(erc20TxIDs, sendRow.RelatedID) {
					erc20TxIDs = append(erc20TxIDs, sendRow.RelatedID)
				}
			case app.SendRelationTypeTxErc20Fee:
				if !hcommon.IsIntInSlice(erc20TxFeeIDs, sendRow.RelatedID) {
					erc20TxFeeIDs = append(erc20TxFeeIDs, sendRow.RelatedID)
				}
			}
		}
		// 通知信息
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
		// 更新eth零钱整理状态
		_, err = app.SQLUpdateTTxOrgStatusByIDs(
			context.Background(),
			app.DbCon,
			txIDs,
			model.DBTTx{
				OrgStatus: app.TxOrgStatusConfirm,
				OrgMsg:    "confirm",
				OrgTime:   now,
			},
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		// 更新erc20零钱整理状态
		_, err = app.SQLUpdateTTxOrgStatusByIDs(
			context.Background(),
			app.DbCon,
			erc20TxIDs,
			model.DBTTx{
				OrgStatus: app.TxOrgStatusConfirm,
				OrgMsg:    "confirm",
				OrgTime:   now,
			},
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		// 更新erc20零钱整理eth手续费状态
		_, err = app.SQLUpdateTTxErc20OrgStatusByIDs(
			context.Background(),
			app.DbCon,
			erc20TxFeeIDs,
			model.DBTTxErc20{
				OrgStatus: app.TxOrgStatusFeeConfirm,
				OrgMsg:    "eth fee confirmed",
				OrgTime:   now,
			},
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		// 更新发送状态
		_, err = app.SQLUpdateTSendStatusByIDs(
			context.Background(),
			app.DbCon,
			sendIDs,
			model.DBTSend{
				HandleStatus: app.SendStatusConfirm,
				HandleMsg:    "confirmed",
				HandleTime:   now,
			},
		)
	})
}

// CheckWithdraw 检测提现
func CheckWithdraw() {
	lockKey := "EthCheckWithdraw"
	app.LockWrap(lockKey, func() {
		// 获取热钱包地址
		hotRow, err := app.SQLGetTAppConfigStrByK(
			context.Background(),
			app.DbCon,
			"hot_wallet_address",
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		if hotRow == nil {
			hcommon.Log.Errorf("no config int of hot_wallet_address")
			return
		}
		_, err = StrToAddressBytes(hotRow.V)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		// 获取私钥
		keyRow, err := app.SQLGetTAddressKeyColByAddress(
			context.Background(),
			app.DbCon,
			[]string{
				model.DBColTAddressKeyPwd,
			},
			hotRow.V,
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		if keyRow == nil {
			hcommon.Log.Errorf("no key of: %s", hotRow.V)
			return
		}
		key := hcommon.AesDecrypt(keyRow.Pwd, app.Cfg.AESKey)
		if len(key) == 0 {
			hcommon.Log.Errorf("error key of: %s", hotRow.V)
			return
		}
		if strings.HasPrefix(key, "0x") {
			key = key[2:]
		}
		privateKey, err := crypto.HexToECDSA(key)
		if err != nil {
			hcommon.Log.Errorf("HexToECDSA err: [%T] %s", err, err.Error())
			return
		}
		withdrawRows, err := app.SQLSelectTWithdrawColByStatus(
			context.Background(),
			app.DbCon,
			[]string{
				model.DBColTWithdrawID,
			},
			app.WithdrawStatusInit,
			[]string{"eth"},
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		if len(withdrawRows) == 0 {
			return
		}
		hotAddressBalance, err := ethclient.RpcBalanceAt(
			context.Background(),
			hotRow.V,
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		pendingBalance, err := app.SQLGetTSendPendingBalance(
			context.Background(),
			app.DbCon,
			hotRow.V,
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		hotAddressBalance -= pendingBalance
		// 获取gap price
		gasRow, err := app.SQLGetTAppStatusIntByK(
			context.Background(),
			app.DbCon,
			"to_user_gas_price",
		)
		if err != nil {
			hcommon.Log.Warnf("err: [%T] %s", err, err.Error())
			return
		}
		if gasRow == nil {
			hcommon.Log.Errorf("no config int of to_user_gas_price")
			return
		}
		gasPrice := gasRow.V
		gasLimit := int64(21000)
		feeValue := gasLimit * gasPrice
		chainID, err := ethclient.RpcNetworkID(context.Background())
		if err != nil {
			hcommon.Log.Warnf("err: [%T] %s", err, err.Error())
			return
		}
		for _, withdrawRow := range withdrawRows {
			err = handleWithdraw(withdrawRow.ID, chainID, hotRow.V, privateKey, &hotAddressBalance, gasLimit, gasPrice, feeValue)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				continue
			}
		}
	})
}

func handleWithdraw(withdrawID int64, chainID int64, hotAddress string, privateKey *ecdsa.PrivateKey, hotAddressBalance *int64, gasLimit, gasPrice, feeValue int64) error {
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
	balanceObj, err := decimal.NewFromString(withdrawRow.BalanceReal)
	if err != nil {
		return err
	}
	balance := balanceObj.Mul(decimal.NewFromInt(1e18)).IntPart()
	*hotAddressBalance -= balance + feeValue
	if *hotAddressBalance < 0 {
		hcommon.Log.Errorf("hot balance limit")
		return nil
	}
	// nonce
	nonce, err := GetNonce(
		dbTx,
		hotAddress,
	)
	if err != nil {
		return err
	}
	// 创建交易
	var data []byte
	toAddress, err := StrToAddressBytes(withdrawRow.ToAddress)
	if err != nil {
		return err
	}
	tx := types.NewTransaction(
		uint64(nonce),
		toAddress,
		big.NewInt(balance),
		uint64(gasLimit),
		big.NewInt(gasPrice),
		data,
	)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(chainID)), privateKey)
	if err != nil {
		return err
	}
	ts := types.Transactions{signedTx}
	rawTxBytes := ts.GetRlp(0)
	rawTxHex := hex.EncodeToString(rawTxBytes)
	txHash := strings.ToLower(signedTx.Hash().Hex())
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
	_, err = model.SQLCreateTSend(
		context.Background(),
		dbTx,
		&model.DBTSend{
			RelatedType:  app.SendRelationTypeWithdraw,
			RelatedID:    withdrawID,
			TxID:         txHash,
			FromAddress:  hotAddress,
			ToAddress:    withdrawRow.ToAddress,
			Balance:      balance,
			BalanceReal:  withdrawRow.BalanceReal,
			Gas:          gasLimit,
			GasPrice:     gasPrice,
			Nonce:        nonce,
			Hex:          rawTxHex,
			HandleStatus: app.SendStatusInit,
			HandleMsg:    "init",
			HandleTime:   now,
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

// CheckTxNotify 创建eth冲币通知
func CheckTxNotify() {
	lockKey := "EthCheckTxNotify"
	app.LockWrap(lockKey, func() {
		txRows, err := app.SQLSelectTTxColByStatus(
			context.Background(),
			app.DbCon,
			[]string{
				model.DBColTTxID,
				model.DBColTTxProductID,
				model.DBColTTxTxID,
				model.DBColTTxToAddress,
				model.DBColTTxBalanceReal,
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
				"tx_hash":     txRow.TxID,
				"app_name":    productRow.AppName,
				"address":     txRow.ToAddress,
				"balance":     txRow.BalanceReal,
				"symbol":      "eth",
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
				TokenSymbol:  "eth",
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
		_, err = app.SQLUpdateTTxStatusByIDs(
			context.Background(),
			app.DbCon,
			notifyTxIDs,
			model.DBTTx{
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

// CheckErc20BlockSeek 检测erc20到账
func CheckErc20BlockSeek() {
	lockKey := "Erc20CheckBlockSeek"
	app.LockWrap(lockKey, func() {
		// 获取配置 延迟确认数
		confirmRow, err := app.SQLGetTAppConfigIntByK(
			context.Background(),
			app.DbCon,
			"block_confirm_num",
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		if confirmRow == nil {
			hcommon.Log.Errorf("no config int of block_confirm_num")
			return
		}
		// 获取状态 当前处理完成的最新的block number
		seekRow, err := app.SQLGetTAppStatusIntByK(
			context.Background(),
			app.DbCon,
			"erc20_seek_num",
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		if seekRow == nil {
			hcommon.Log.Errorf("no config int of erc20_seek_num")
			return
		}
		// rpc 获取当前最新区块数
		rpcBlockNum, err := ethclient.RpcBlockNumber(context.Background())
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		startI := seekRow.V + 1
		endI := rpcBlockNum - confirmRow.V + 1
		if startI < endI {
			// 读取abi
			type LogTransfer struct {
				From   string
				To     string
				Tokens *big.Int
			}
			contractAbi, err := abi.JSON(strings.NewReader(ethclient.EthABI))
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			// 获取所有token
			var configTokenRowAddresses []string
			configTokenRowMap := make(map[string]*model.DBTAppConfigToken)
			configTokenRows, err := app.SQLSelectTAppConfigTokenColAll(
				context.Background(),
				app.DbCon,
				[]string{
					model.DBColTAppConfigTokenID,
					model.DBColTAppConfigTokenTokenAddress,
					model.DBColTAppConfigTokenTokenDecimals,
					model.DBColTAppConfigTokenTokenSymbol,
				},
			)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			for _, contractRow := range configTokenRows {
				configTokenRowAddresses = append(configTokenRowAddresses, contractRow.TokenAddress)
				configTokenRowMap[contractRow.TokenAddress] = contractRow
			}
			// 遍历获取需要查询的block信息
			for i := startI; i < endI; i++ {
				//hcommon.Log.Debugf("erc20 check block: %d", i)
				if len(configTokenRowAddresses) > 0 {
					// rpc获取block信息
					logs, err := ethclient.RpcFilterLogs(
						context.Background(),
						i,
						i,
						configTokenRowAddresses,
						contractAbi.Events["Transfer"],
					)
					if err != nil {
						hcommon.Log.Warnf("err: [%T] %s", err, err.Error())
						return
					}
					// 接收地址列表
					var toAddresses []string
					// map[接收地址] => []交易信息
					toAddressLogMap := make(map[string][]types.Log)
					for _, log := range logs {
						if log.Removed {
							continue
						}
						toAddress := AddressBytesToStr(common.HexToAddress(log.Topics[2].Hex()))
						if !hcommon.IsStringInSlice(toAddresses, toAddress) {
							toAddresses = append(toAddresses, toAddress)
						}
						toAddressLogMap[toAddress] = append(toAddressLogMap[toAddress], log)
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
					// map[接收地址] => 产品id
					addressProductMap := make(map[string]int64)
					for _, dbAddressRow := range dbAddressRows {
						addressProductMap[dbAddressRow.Address] = dbAddressRow.UseTag
					}
					// 时间
					now := time.Now().Unix()
					// 待添加数组
					var txErc20Rows []*model.DBTTxErc20
					// 遍历数据库中有交易的地址
					for _, dbAddressRow := range dbAddressRows {
						if dbAddressRow.UseTag < 0 {
							continue
						}
						// 获取地址对应的交易列表
						logs, ok := toAddressLogMap[dbAddressRow.Address]
						if !ok {
							hcommon.Log.Errorf("toAddressLogMap no: %s", dbAddressRow.Address)
							return
						}
						for _, log := range logs {
							var transferEvent LogTransfer
							err := contractAbi.Unpack(&transferEvent, "Transfer", log.Data)
							if err != nil {
								hcommon.Log.Warnf("err: [%T] %s", err, err.Error())
								return
							}
							transferEvent.From = strings.ToLower(common.HexToAddress(log.Topics[1].Hex()).Hex())
							transferEvent.To = strings.ToLower(common.HexToAddress(log.Topics[2].Hex()).Hex())
							contractAddress := strings.ToLower(log.Address.Hex())
							configTokenRow, ok := configTokenRowMap[contractAddress]
							if !ok {
								hcommon.Log.Errorf("no configTokenRowMap of: %s", contractAddress)
								return
							}
							rpcTxReceipt, err := ethclient.RpcTransactionReceipt(
								context.Background(),
								log.TxHash.Hex(),
							)
							if err != nil {
								hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
								return
							}
							if rpcTxReceipt.Status <= 0 {
								continue
							}
							rpcTx, err := ethclient.RpcTransactionByHash(
								context.Background(),
								log.TxHash.Hex(),
							)
							if err != nil {
								hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
								return
							}
							if strings.ToLower(rpcTx.To().Hex()) != contractAddress {
								// 合约地址和tx的to地址不匹配
								continue
							}
							// 检测input
							input, err := contractAbi.Pack(
								"transfer",
								common.HexToAddress(log.Topics[2].Hex()),
								transferEvent.Tokens,
							)
							if err != nil {
								hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
								return
							}
							if hexutil.Encode(input) != hexutil.Encode(rpcTx.Data()) {
								// input 不匹配
								continue
							}
							balanceReal := decimal.NewFromInt(transferEvent.Tokens.Int64()).Div(decimal.NewFromInt(int64(math.Pow10(int(configTokenRow.TokenDecimals))))).String()
							// 放入待插入数组
							txErc20Rows = append(txErc20Rows, &model.DBTTxErc20{
								TokenID:      configTokenRow.ID,
								ProductID:    addressProductMap[transferEvent.To],
								TxID:         log.TxHash.Hex(),
								FromAddress:  transferEvent.From,
								ToAddress:    transferEvent.To,
								Balance:      transferEvent.Tokens.Int64(),
								BalanceReal:  balanceReal,
								CreateTime:   now,
								HandleStatus: app.TxStatusInit,
								HandleMsg:    "",
								HandleTime:   now,
								OrgStatus:    app.TxOrgStatusInit,
								OrgMsg:       "",
								OrgTime:      now,
							})
						}
					}
					_, err = model.SQLCreateIgnoreManyTTxErc20(
						context.Background(),
						app.DbCon,
						txErc20Rows,
					)
					if err != nil {
						hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
						return
					}
				}
				// 更新检查到的最新区块数
				_, err = app.SQLUpdateTAppStatusIntByK(
					context.Background(),
					app.DbCon,
					&model.DBTAppStatusInt{
						K: "erc20_seek_num",
						V: i,
					},
				)
				if err != nil {
					hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
			}
		}
	})
}

// CheckErc20TxNotify 创建erc20冲币通知
func CheckErc20TxNotify() {
	lockKey := "Erc20CheckTxNotify"
	app.LockWrap(lockKey, func() {
		txRows, err := app.SQLSelectTTxErc20ColByStatus(
			context.Background(),
			app.DbCon,
			[]string{
				model.DBColTTxErc20ID,
				model.DBColTTxErc20TokenID,
				model.DBColTTxErc20ProductID,
				model.DBColTTxErc20TxID,
				model.DBColTTxErc20ToAddress,
				model.DBColTTxErc20BalanceReal,
			},
			app.TxStatusInit,
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		var productIDs []int64
		var tokenIDs []int64
		for _, txRow := range txRows {
			if !hcommon.IsIntInSlice(productIDs, txRow.ProductID) {
				productIDs = append(productIDs, txRow.ProductID)
			}
			if !hcommon.IsIntInSlice(tokenIDs, txRow.TokenID) {
				tokenIDs = append(tokenIDs, txRow.TokenID)
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
		tokenMap, err := app.SQLGetAppConfigTokenMap(
			context.Background(),
			app.DbCon,
			[]string{
				model.DBColTAppConfigTokenID,
				model.DBColTAppConfigTokenTokenSymbol,
			},
			tokenIDs,
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
				hcommon.Log.Warnf("productMap no: %d", txRow.ProductID)
				notifyTxIDs = append(notifyTxIDs, txRow.ID)
				continue
			}
			tokenRow, ok := tokenMap[txRow.TokenID]
			if !ok {
				hcommon.Log.Errorf("tokenMap no: %d", txRow.TokenID)
				continue
			}
			nonce := hcommon.GetUUIDStr()
			reqObj := gin.H{
				"tx_hash":     txRow.TxID,
				"app_name":    productRow.AppName,
				"address":     txRow.ToAddress,
				"balance":     txRow.BalanceReal,
				"symbol":      tokenRow.TokenSymbol,
				"notify_type": app.NotifyTypeTx,
			}
			reqObj["sign"] = hcommon.GetSign(productRow.AppSk, reqObj)
			req, err := json.Marshal(reqObj)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			notifyRows = append(notifyRows, &model.DBTProductNotify{
				Nonce:        nonce,
				ProductID:    txRow.ProductID,
				ItemType:     app.SendRelationTypeTx,
				ItemID:       txRow.ID,
				NotifyType:   app.NotifyTypeTx,
				TokenSymbol:  tokenRow.TokenSymbol,
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
		_, err = app.SQLUpdateTTxErc20StatusByIDs(
			context.Background(),
			app.DbCon,
			notifyTxIDs,
			model.DBTTxErc20{
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

// CheckErc20TxOrg erc20零钱整理
func CheckErc20TxOrg() {
	lockKey := "Erc20CheckTxOrg"
	app.LockWrap(lockKey, func() {
		txRows, err := app.SQLSelectTTxErc20ColByOrg(
			context.Background(),
			app.DbCon,
			[]string{
				model.DBColTTxErc20ID,
				model.DBColTTxErc20TokenID,
				model.DBColTTxErc20ProductID,
				model.DBColTTxErc20ToAddress,
				model.DBColTTxErc20Balance,
			},
			[]int64{app.TxOrgStatusInit, app.TxOrgStatusFeeConfirm},
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		if len(txRows) <= 0 {
			return
		}
		addressEthBalanceMap := make(map[string]int64)
		txMap := make(map[int64]*model.DBTTxErc20)
		type StOrgInfo struct {
			TxIDs        []int64
			ToAddress    string
			TokenID      int64
			TokenBalance int64
		}
		orgMap := make(map[string]*StOrgInfo)

		var tokenIDs []int64
		var toAddresses []string

		for _, txRow := range txRows {
			// 转换为map
			txMap[txRow.ID] = txRow
			// 读取eth余额
			_, ok := addressEthBalanceMap[txRow.ToAddress]
			if !ok {
				balance, err := ethclient.RpcBalanceAt(
					context.Background(),
					txRow.ToAddress,
				)
				if err != nil {
					hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
				addressEthBalanceMap[txRow.ToAddress] = balance
			}
			// 整理信息
			orgKey := fmt.Sprintf("%s-%d", txRow.ToAddress, txRow.TokenID)
			orgInfo, ok := orgMap[orgKey]
			if !ok {
				orgInfo = &StOrgInfo{
					TokenID:   txRow.TokenID,
					ToAddress: txRow.ToAddress,
				}
				orgMap[orgKey] = orgInfo
			}
			orgInfo.TxIDs = append(orgInfo.TxIDs, txRow.ID)
			orgInfo.TokenBalance += txRow.Balance

			// 待查询id
			if !hcommon.IsIntInSlice(tokenIDs, txRow.TokenID) {
				tokenIDs = append(tokenIDs, txRow.TokenID)
			}
			if !hcommon.IsStringInSlice(toAddresses, txRow.ToAddress) {
				toAddresses = append(toAddresses, txRow.ToAddress)
			}
		}
		tokenMap, err := app.SQLGetAppConfigTokenMap(
			context.Background(),
			app.DbCon,
			[]string{
				model.DBColTAppConfigTokenID,
				model.DBColTAppConfigTokenTokenAddress,
				model.DBColTAppConfigTokenTokenDecimals,
				model.DBColTAppConfigTokenTokenSymbol,
				model.DBColTAppConfigTokenColdAddress,
				model.DBColTAppConfigTokenOrgMinBalance,
			},
			tokenIDs,
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		addressMap, err := app.SQLGetAddressKeyMap(
			context.Background(),
			app.DbCon,
			[]string{
				model.DBColTAddressKeyAddress,
				model.DBColTAddressKeyPwd,
			},
			toAddresses,
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}

		// 计算转账token所需的手续费
		erc20GasRow, err := app.SQLGetTAppConfigIntByK(
			context.Background(),
			app.DbCon,
			"erc20_gas_use",
		)
		if err != nil {
			hcommon.Log.Warnf("err: [%T] %s", err, err.Error())
			return
		}
		if erc20GasRow == nil {
			hcommon.Log.Errorf("no config int of erc20_gas_use")
			return
		}
		gasPriceRow, err := app.SQLGetTAppStatusIntByK(
			context.Background(),
			app.DbCon,
			"to_cold_gas_price",
		)
		if err != nil {
			hcommon.Log.Warnf("err: [%T] %s", err, err.Error())
			return
		}
		if gasPriceRow == nil {
			hcommon.Log.Errorf("no config int of to_cold_gas_price")
			return
		}
		chainID, err := ethclient.RpcNetworkID(context.Background())
		if err != nil {
			hcommon.Log.Warnf("err: [%T] %s", err, err.Error())
			return
		}
		erc20Fee := erc20GasRow.V * gasPriceRow.V
		// 需要手续费的整理信息
		needEthFeeMap := make(map[string]*StOrgInfo)
		for k, orgInfo := range orgMap {
			toAddress := orgInfo.ToAddress
			addressEthBalanceMap[toAddress] -= erc20Fee
			if addressEthBalanceMap[toAddress] < 0 {
				// eth手续费不足
				// 处理添加手续费
				needEthFeeMap[k] = orgInfo
				continue
			}
			tokenRow, ok := tokenMap[orgInfo.TokenID]
			if !ok {
				hcommon.Log.Errorf("no tokenMap: %d", orgInfo.TokenID)
				continue
			}
			orgMinBalanceObj, err := decimal.NewFromString(tokenRow.OrgMinBalance)
			if err != nil {
				hcommon.Log.Warnf("err: [%T] %s", err, err.Error())
				continue
			}
			orgMinBalance := orgMinBalanceObj.Mul(decimal.NewFromInt(int64(math.Pow10(int(tokenRow.TokenDecimals))))).IntPart()
			if orgInfo.TokenBalance < orgMinBalance {
				hcommon.Log.Errorf("token balance < org min balance")
				continue
			}
			// 处理token转账
			addressRow, ok := addressMap[toAddress]
			if !ok {
				hcommon.Log.Errorf("addressMap no: %s", toAddress)
				continue
			}
			key := hcommon.AesDecrypt(addressRow.Pwd, app.Cfg.AESKey)
			if len(key) == 0 {
				hcommon.Log.Errorf("error key of: %s", toAddress)
				continue
			}
			if strings.HasPrefix(key, "0x") {
				key = key[2:]
			}
			privateKey, err := crypto.HexToECDSA(key)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				continue
			}
			// 获取nonce值
			nonce, err := GetNonce(app.DbCon, toAddress)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				continue
			}
			// 生成交易
			contractAbi, err := abi.JSON(strings.NewReader(ethclient.EthABI))
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			input, err := contractAbi.Pack(
				"transfer",
				common.HexToAddress(tokenRow.ColdAddress),
				big.NewInt(orgInfo.TokenBalance),
			)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			rpcTx := types.NewTransaction(
				uint64(nonce),
				common.HexToAddress(tokenRow.TokenAddress),
				big.NewInt(0),
				uint64(erc20GasRow.V),
				big.NewInt(gasPriceRow.V),
				input,
			)
			signedTx, err := types.SignTx(rpcTx, types.NewEIP155Signer(big.NewInt(chainID)), privateKey)
			if err != nil {
				hcommon.Log.Warnf("err: [%T] %s", err, err.Error())
				continue
			}
			ts := types.Transactions{signedTx}
			rawTxBytes := ts.GetRlp(0)
			rawTxHex := hex.EncodeToString(rawTxBytes)
			txHash := strings.ToLower(signedTx.Hash().Hex())
			// 创建存入数据
			now := time.Now().Unix()
			balanceReal := decimal.NewFromInt(orgInfo.TokenBalance).Div(decimal.NewFromInt(int64(math.Pow10(int(tokenRow.TokenDecimals))))).String()
			// 待插入数据
			var sendRows []*model.DBTSend
			for rowIndex, txID := range orgInfo.TxIDs {
				if rowIndex == 0 {
					sendRows = append(sendRows, &model.DBTSend{
						RelatedType:  app.SendRelationTypeTxErc20,
						RelatedID:    txID,
						TokenID:      orgInfo.TokenID,
						TxID:         txHash,
						FromAddress:  toAddress,
						ToAddress:    tokenRow.ColdAddress,
						Balance:      orgInfo.TokenBalance,
						BalanceReal:  balanceReal,
						Gas:          erc20GasRow.V,
						GasPrice:     gasPriceRow.V,
						Nonce:        nonce,
						Hex:          rawTxHex,
						CreateTime:   now,
						HandleStatus: app.SendStatusInit,
						HandleMsg:    "",
						HandleTime:   now,
					})
				} else {
					sendRows = append(sendRows, &model.DBTSend{
						RelatedType:  app.SendRelationTypeTxErc20,
						RelatedID:    txID,
						TokenID:      orgInfo.TokenID,
						TxID:         txHash,
						FromAddress:  toAddress,
						ToAddress:    tokenRow.ColdAddress,
						Balance:      0,
						BalanceReal:  "",
						Gas:          0,
						GasPrice:     0,
						Nonce:        -1,
						Hex:          "",
						CreateTime:   now,
						HandleStatus: app.SendStatusInit,
						HandleMsg:    "",
						HandleTime:   now,
					})
				}
			}
			// 插入发送队列
			_, err = model.SQLCreateIgnoreManyTSend(
				context.Background(),
				app.DbCon,
				sendRows,
			)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			// 更新整理状态
			_, err = app.SQLUpdateTTxErc20OrgStatusByIDs(
				context.Background(),
				app.DbCon,
				orgInfo.TxIDs,
				model.DBTTxErc20{
					OrgStatus: app.TxOrgStatusHex,
					OrgMsg:    "hex",
					OrgTime:   now,
				},
			)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
		}
		// 生成eth转账
		if len(needEthFeeMap) > 0 {
			// 获取热钱包地址
			feeRow, err := app.SQLGetTAppConfigStrByK(
				context.Background(),
				app.DbCon,
				"fee_wallet_address",
			)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			if feeRow == nil {
				hcommon.Log.Errorf("no config int of fee_wallet_address")
				return
			}
			_, err = StrToAddressBytes(feeRow.V)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			// 获取私钥
			keyRow, err := app.SQLGetTAddressKeyColByAddress(
				context.Background(),
				app.DbCon,
				[]string{
					model.DBColTAddressKeyPwd,
				},
				feeRow.V,
			)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			if keyRow == nil {
				hcommon.Log.Errorf("no key of: %s", feeRow.V)
				return
			}
			key := hcommon.AesDecrypt(keyRow.Pwd, app.Cfg.AESKey)
			if len(key) == 0 {
				hcommon.Log.Errorf("error key of: %s", feeRow.V)
				return
			}
			if strings.HasPrefix(key, "0x") {
				key = key[2:]
			}
			privateKey, err := crypto.HexToECDSA(key)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			feeAddressBalance, err := ethclient.RpcBalanceAt(
				context.Background(),
				feeRow.V,
			)
			if err != nil {
				hcommon.Log.Errorf("RpcBalanceAt err: [%T] %s", err, err.Error())
				return
			}
			pendingBalance, err := app.SQLGetTSendPendingBalance(
				context.Background(),
				app.DbCon,
				feeRow.V,
			)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			feeAddressBalance -= pendingBalance
			// 获取gap price
			gasRow, err := app.SQLGetTAppStatusIntByK(
				context.Background(),
				app.DbCon,
				"to_cold_gas_price",
			)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			if gasRow == nil {
				hcommon.Log.Errorf("no config int of to_cold_gas_price")
				return
			}
			gasPrice := gasRow.V
			gasLimit := int64(21000)
			ethFee := gasLimit * gasPrice
			tokenFee := erc20GasRow.V * gasPrice
			for _, orgInfo := range needEthFeeMap {
				feeAddressBalance -= ethFee + tokenFee
				if feeAddressBalance < 0 {
					hcommon.Log.Errorf("eth fee balance limit")
					return
				}
				// nonce
				nonce, err := GetNonce(
					app.DbCon,
					feeRow.V,
				)
				if err != nil {
					hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
				// 创建交易
				var data []byte
				tx := types.NewTransaction(
					uint64(nonce),
					common.HexToAddress(orgInfo.ToAddress),
					big.NewInt(tokenFee),
					uint64(gasLimit),
					big.NewInt(gasPrice),
					data,
				)
				signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(chainID)), privateKey)
				if err != nil {
					hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
				ts := types.Transactions{signedTx}
				rawTxBytes := ts.GetRlp(0)
				rawTxHex := hex.EncodeToString(rawTxBytes)
				txHash := strings.ToLower(signedTx.Hash().Hex())
				now := time.Now().Unix()
				balanceReal := decimal.NewFromInt(tokenFee).Div(decimal.NewFromInt(int64(math.Pow10(18)))).String()
				// 待插入数据
				var sendRows []*model.DBTSend
				for rowIndex, txID := range orgInfo.TxIDs {
					if rowIndex == 0 {
						sendRows = append(sendRows, &model.DBTSend{
							RelatedType:  app.SendRelationTypeTxErc20Fee,
							RelatedID:    txID,
							TokenID:      0,
							TxID:         txHash,
							FromAddress:  feeRow.V,
							ToAddress:    orgInfo.ToAddress,
							Balance:      feeAddressBalance,
							BalanceReal:  balanceReal,
							Gas:          gasLimit,
							GasPrice:     gasPriceRow.V,
							Nonce:        nonce,
							Hex:          rawTxHex,
							CreateTime:   now,
							HandleStatus: app.SendStatusInit,
							HandleMsg:    "",
							HandleTime:   now,
						})
					} else {
						sendRows = append(sendRows, &model.DBTSend{
							RelatedType:  app.SendRelationTypeTxErc20Fee,
							RelatedID:    txID,
							TokenID:      0,
							TxID:         txHash,
							FromAddress:  feeRow.V,
							ToAddress:    orgInfo.ToAddress,
							Balance:      0,
							BalanceReal:  "",
							Gas:          0,
							GasPrice:     0,
							Nonce:        -1,
							Hex:          "",
							CreateTime:   now,
							HandleStatus: app.SendStatusInit,
							HandleMsg:    "",
							HandleTime:   now,
						})
					}
				}
				_, err = model.SQLCreateIgnoreManyTSend(
					context.Background(),
					app.DbCon,
					sendRows,
				)
				if err != nil {
					hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
				_, err = app.SQLUpdateTTxErc20OrgStatusByIDs(
					context.Background(),
					app.DbCon,
					orgInfo.TxIDs,
					model.DBTTxErc20{
						OrgStatus: app.TxOrgStatusFeeHex,
						OrgMsg:    "fee hex",
						OrgTime:   now,
					},
				)
				if err != nil {
					hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
			}
		}
	})
}

// CheckErc20Withdraw erc20提币
func CheckErc20Withdraw() {
	lockKey := "Erc20CheckWithdraw"
	app.LockWrap(lockKey, func() {
		var tokenSymbols []string
		tokenMap := make(map[string]*model.DBTAppConfigToken)
		addressKeyMap := make(map[string]*ecdsa.PrivateKey)
		addressEthBalanceMap := make(map[string]int64)
		addressTokenBalanceMap := make(map[string]int64)
		tokenRows, err := app.SQLSelectTAppConfigTokenColAll(
			context.Background(),
			app.DbCon,
			[]string{
				model.DBColTAppConfigTokenID,
				model.DBColTAppConfigTokenTokenAddress,
				model.DBColTAppConfigTokenTokenDecimals,
				model.DBColTAppConfigTokenTokenSymbol,
				model.DBColTAppConfigTokenHotAddress,
			},
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		for _, tokenRow := range tokenRows {
			tokenMap[tokenRow.TokenSymbol] = tokenRow
			if !hcommon.IsStringInSlice(tokenSymbols, tokenRow.TokenSymbol) {
				tokenSymbols = append(tokenSymbols, tokenRow.TokenSymbol)
			}
			// 获取私钥
			_, err = StrToAddressBytes(tokenRow.HotAddress)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			hotAddress := tokenRow.HotAddress
			_, ok := addressKeyMap[hotAddress]
			if !ok {
				// 获取私钥
				keyRow, err := app.SQLGetTAddressKeyColByAddress(
					context.Background(),
					app.DbCon,
					[]string{
						model.DBColTAddressKeyPwd,
					},
					hotAddress,
				)
				if err != nil {
					hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
				if keyRow == nil {
					hcommon.Log.Errorf("no key of: %s", hotAddress)
					return
				}
				key := hcommon.AesDecrypt(keyRow.Pwd, app.Cfg.AESKey)
				if len(key) == 0 {
					hcommon.Log.Errorf("error key of: %s", hotAddress)
					return
				}
				if strings.HasPrefix(key, "0x") {
					key = key[2:]
				}
				privateKey, err := crypto.HexToECDSA(key)
				if err != nil {
					hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
				addressKeyMap[hotAddress] = privateKey
			}
			_, ok = addressEthBalanceMap[hotAddress]
			if !ok {
				hotAddressBalance, err := ethclient.RpcBalanceAt(
					context.Background(),
					hotAddress,
				)
				if err != nil {
					hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
				pendingBalance, err := app.SQLGetTSendPendingBalance(
					context.Background(),
					app.DbCon,
					hotAddress,
				)
				if err != nil {
					hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
				hotAddressBalance -= pendingBalance
				addressEthBalanceMap[hotAddress] = hotAddressBalance
			}
			tokenBalanceKey := fmt.Sprintf("%s-%s", tokenRow.HotAddress, tokenRow.TokenSymbol)
			_, ok = addressTokenBalanceMap[tokenBalanceKey]
			if !ok {
				tokenBalance, err := ethclient.RpcTokenBalance(
					context.Background(),
					tokenRow.TokenAddress,
					tokenRow.HotAddress,
				)
				if err != nil {
					hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
				addressTokenBalanceMap[tokenBalanceKey] = tokenBalance
			}
		}
		withdrawRows, err := app.SQLSelectTWithdrawColByStatus(
			context.Background(),
			app.DbCon,
			[]string{
				model.DBColTWithdrawID,
				model.DBColTWithdrawSymbol,
			},
			app.WithdrawStatusInit,
			tokenSymbols,
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		if len(withdrawRows) == 0 {
			return
		}
		// 获取gap price
		gasRow, err := app.SQLGetTAppStatusIntByK(
			context.Background(),
			app.DbCon,
			"to_user_gas_price",
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		if gasRow == nil {
			hcommon.Log.Errorf("no config int of to_user_gas_price")
			return
		}
		gasPrice := gasRow.V
		erc20GasRow, err := app.SQLGetTAppConfigIntByK(
			context.Background(),
			app.DbCon,
			"erc20_gas_use",
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		if erc20GasRow == nil {
			hcommon.Log.Errorf("no config int of erc20_gas_use")
			return
		}
		gasLimit := erc20GasRow.V
		// eth fee
		feeValue := gasLimit * gasPrice
		chainID, err := ethclient.RpcNetworkID(context.Background())
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		for _, withdrawRow := range withdrawRows {
			err = handleErc20Withdraw(withdrawRow.ID, chainID, &tokenMap, &addressKeyMap, &addressEthBalanceMap, &addressTokenBalanceMap, gasLimit, gasPrice, feeValue)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				continue
			}
		}
	})
}

func handleErc20Withdraw(withdrawID int64, chainID int64, tokenMap *map[string]*model.DBTAppConfigToken, addressKeyMap *map[string]*ecdsa.PrivateKey, addressEthBalanceMap *map[string]int64, addressTokenBalanceMap *map[string]int64, gasLimit, gasPrice, feeValue int64) error {
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
			model.DBColTWithdrawSymbol,
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
	tokenRow, ok := (*tokenMap)[withdrawRow.Symbol]
	if !ok {
		hcommon.Log.Errorf("no tokenMap: %s", withdrawRow.Symbol)
		return nil
	}
	hotAddress := tokenRow.HotAddress
	key, ok := (*addressKeyMap)[hotAddress]
	if !ok {
		hcommon.Log.Errorf("no addressKeyMap: %s", hotAddress)
		return nil
	}
	(*addressEthBalanceMap)[hotAddress] -= feeValue
	if (*addressEthBalanceMap)[hotAddress] < 0 {
		hcommon.Log.Errorf("%s eth limit", hotAddress)
		return nil
	}
	tokenBalanceKey := fmt.Sprintf("%s-%s", tokenRow.HotAddress, tokenRow.TokenSymbol)
	tokenBalanceObj, err := decimal.NewFromString(withdrawRow.BalanceReal)
	if err != nil {
		return err
	}
	tokenBalance := tokenBalanceObj.Mul(decimal.NewFromInt(int64(math.Pow10(int(tokenRow.TokenDecimals))))).IntPart()
	(*addressTokenBalanceMap)[tokenBalanceKey] -= tokenBalance
	if (*addressTokenBalanceMap)[tokenBalanceKey] < 0 {
		hcommon.Log.Errorf("%s token limit", tokenBalanceKey)
		return nil
	}
	// 获取nonce值
	nonce, err := GetNonce(dbTx, hotAddress)
	if err != nil {
		return err
	}
	// 生成交易
	contractAbi, err := abi.JSON(strings.NewReader(ethclient.EthABI))
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return err
	}
	input, err := contractAbi.Pack(
		"transfer",
		common.HexToAddress(withdrawRow.ToAddress),
		big.NewInt(tokenBalance),
	)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return err
	}
	rpcTx := types.NewTransaction(
		uint64(nonce),
		common.HexToAddress(tokenRow.TokenAddress),
		big.NewInt(0),
		uint64(gasLimit),
		big.NewInt(gasPrice),
		input,
	)
	signedTx, err := types.SignTx(rpcTx, types.NewEIP155Signer(big.NewInt(chainID)), key)
	if err != nil {
		return err
	}
	ts := types.Transactions{signedTx}
	rawTxBytes := ts.GetRlp(0)
	rawTxHex := hex.EncodeToString(rawTxBytes)
	txHash := strings.ToLower(signedTx.Hash().Hex())
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
	_, err = model.SQLCreateTSend(
		context.Background(),
		dbTx,
		&model.DBTSend{
			RelatedType:  app.SendRelationTypeWithdraw,
			RelatedID:    withdrawID,
			TxID:         txHash,
			FromAddress:  hotAddress,
			ToAddress:    withdrawRow.ToAddress,
			Balance:      tokenBalance,
			BalanceReal:  withdrawRow.BalanceReal,
			Gas:          gasLimit,
			GasPrice:     gasPrice,
			Nonce:        nonce,
			Hex:          rawTxHex,
			HandleStatus: app.SendStatusInit,
			HandleMsg:    "init",
			HandleTime:   now,
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

// CheckGasPrice 检测gas price
func CheckGasPrice() {
	lockKey := "EthCheckGasPrice"
	app.LockWrap(lockKey, func() {
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
		toUserGasPrice := resp.Fast * int64(math.Pow10(8))
		toColdGasPrice := resp.Average * int64(math.Pow10(8))
		_, err = app.SQLUpdateTAppStatusIntByK(
			context.Background(),
			app.DbCon,
			&model.DBTAppStatusInt{
				K: "to_user_gas_price",
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
				K: "to_cold_gas_price",
				V: toColdGasPrice,
			},
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
	})
}
