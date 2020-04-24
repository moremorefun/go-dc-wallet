package eth

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
	"regexp"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/gin-gonic/gin"

	"github.com/ethereum/go-ethereum/rlp"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/shopspring/decimal"
)

// CheckAddressCheck 检测是否有充足的备用地址
func CheckAddressCheck() {
	lockKey := "EthCheckAddressCheck"
	ok, err := app.GetLock(
		context.Background(),
		app.DbCon,
		lockKey,
	)
	if err != nil {
		hcommon.Log.Warnf("GetLock err: [%T] %s", err, err.Error())
		return
	}
	if !ok {
		return
	}
	defer func() {
		err := app.ReleaseLock(
			context.Background(),
			app.DbCon,
			lockKey,
		)
		if err != nil {
			hcommon.Log.Warnf("ReleaseLock err: [%T] %s", err, err.Error())
			return
		}
	}()

	// 获取配置 允许的最小剩余地址数
	minFreeRow, err := app.SQLGetTAppConfigIntByK(
		context.Background(),
		app.DbCon,
		"min_free_address",
	)
	if err != nil {
		hcommon.Log.Warnf("SQLGetTAppConfigInt err: [%T] %s", err, err.Error())
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
	)
	if err != nil {
		hcommon.Log.Warnf("SQLGetTAddressKeyFreeCount err: [%T] %s", err, err.Error())
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
				hcommon.Log.Warnf("GenerateKey err: [%T] %s", err, err.Error())
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
				hcommon.Log.Warnf("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
				return
			}
			// 地址全部储存为小写方便处理
			address := strings.ToLower(crypto.PubkeyToAddress(*publicKeyECDSA).Hex())
			// 存入待添加队列
			rows = append(rows, &model.DBTAddressKey{
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
			hcommon.Log.Warnf("SQLCreateIgnoreManyTAddressKey err: [%T] %s", err, err.Error())
			return
		}
	}
}

// CheckBlockSeek 检测到账
func CheckBlockSeek() {
	lockKey := "CheckBlockSeek"
	ok, err := app.GetLock(
		context.Background(),
		app.DbCon,
		lockKey,
	)
	if err != nil {
		hcommon.Log.Warnf("GetLock err: [%T] %s", err, err.Error())
		return
	}
	if !ok {
		return
	}
	defer func() {
		err := app.ReleaseLock(
			context.Background(),
			app.DbCon,
			lockKey,
		)
		if err != nil {
			hcommon.Log.Warnf("ReleaseLock err: [%T] %s", err, err.Error())
			return
		}
	}()

	// 获取配置 延迟确认数
	confirmRow, err := app.SQLGetTAppConfigIntByK(
		context.Background(),
		app.DbCon,
		"block_confirm_num",
	)
	if err != nil {
		hcommon.Log.Warnf("SQLGetTAppConfigInt err: [%T] %s", err, err.Error())
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
		hcommon.Log.Warnf("SQLGetTAppStatusIntByK err: [%T] %s", err, err.Error())
		return
	}
	if seekRow == nil {
		hcommon.Log.Errorf("no config int of seek_num")
		return
	}
	// rpc 获取当前最新区块数
	rpcBlockNum, err := ethclient.RpcBlockNumber(context.Background())
	if err != nil {
		hcommon.Log.Warnf("RpcBlockNumber err: [%T] %s", err, err.Error())
		return
	}
	startI := seekRow.V + 1
	endI := rpcBlockNum - confirmRow.V
	if startI < endI {
		// 手续费钱包列表
		feeRow, err := app.SQLGetTAppConfigStrByK(
			context.Background(),
			app.DbCon,
			"fee_wallet_address_list",
		)
		if err != nil {
			hcommon.Log.Warnf("SQLGetTAppConfigInt err: [%T] %s", err, err.Error())
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
			rpcBlock, err := ethclient.RpcBlockByNum(context.Background(), i)
			if err != nil {
				hcommon.Log.Warnf("EthRpcBlockByNum err: [%T] %s", err, err.Error())
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
					if hcommon.IsStringInSlice(feeAddresses, strings.ToLower(msg.From().Hex())) {
						// 如果打币地址在手续费热钱包地址则不处理
						continue
					}
					toAddress := strings.ToLower(rpcTx.To().Hex())
					toAddressTxMap[toAddress] = append(toAddressTxMap[toAddress], rpcTx)
					toAddresses = append(toAddresses, toAddress)
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
				hcommon.Log.Warnf("dbAddressRows err: [%T] %s", err, err.Error())
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
				// 获取地址对应的交易列表
				txes := toAddressTxMap[dbAddressRow.Address]
				for _, tx := range txes {
					msg, err := tx.AsMessage(types.NewEIP155Signer(tx.ChainId()))
					if err != nil {
						hcommon.Log.Errorf("AsMessage err: [%T] %s", err, err.Error())
						return
					}
					toAddress := strings.ToLower(tx.To().Hex())
					balanceReal := decimal.NewFromInt(tx.Value().Int64()).Div(decimal.NewFromInt(1e18))
					dbTxRows = append(dbTxRows, &model.DBTTx{
						ProductID:    addressProductMap[toAddress],
						TxID:         tx.Hash().String(),
						FromAddress:  strings.ToLower(msg.From().Hex()),
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
				hcommon.Log.Warnf("SQLCreateIgnoreManyTTx err: [%T] %s", err, err.Error())
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
				hcommon.Log.Warnf("SQLUpdateTAppStatusIntByK err: [%T] %s", err, err.Error())
				return
			}
		}
	}
}

// CheckAddressOrg 零钱整理到冷钱包
func CheckAddressOrg() {
	lockKey := "EthCheckAddressOrg"
	ok, err := app.GetLock(
		context.Background(),
		app.DbCon,
		lockKey,
	)
	if err != nil {
		hcommon.Log.Warnf("GetLock err: [%T] %s", err, err.Error())
		return
	}
	if !ok {
		return
	}
	defer func() {
		err := app.ReleaseLock(
			context.Background(),
			app.DbCon,
			lockKey,
		)
		if err != nil {
			hcommon.Log.Warnf("ReleaseLock err: [%T] %s", err, err.Error())
			return
		}
	}()

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
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	if !re.MatchString(coldRow.V) {
		hcommon.Log.Errorf("config int cold_wallet_address err: %s", coldRow.V)
		return
	}
	coldAddress := common.HexToAddress(coldRow.V)
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
		hcommon.Log.Warnf("SQLSelectTTxColByOrg err: [%T] %s", err, err.Error())
		return
	}
	if len(txRows) <= 0 {
		// 没有要处理的信息
		return
	}
	// 将待整理地址按地址做归并处理
	type OrgInfo struct {
		RowIDs  []int64
		Balance int64
	}
	addressMap := make(map[string]*OrgInfo)
	// 获取gap price
	gasRow, err := app.SQLGetTAppStatusIntByK(
		context.Background(),
		app.DbCon,
		"to_cold_gas_price",
	)
	if err != nil {
		hcommon.Log.Warnf("SQLGetTAppStatusIntByK err: [%T] %s", err, err.Error())
		return
	}
	if gasRow == nil {
		hcommon.Log.Errorf("no config int of to_cold_gas_price")
		return
	}
	gasPrice := gasRow.V
	gasLimit := int64(21000)
	feeValue := gasLimit * gasPrice
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
			hcommon.Log.Warnf("SQLGetTAddressKeyColByAddress err: [%T] %s", err, err.Error())
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
			hcommon.Log.Warnf("HexToECDSA err: [%T] %s", err, err.Error())
			continue
		}
		// 获取nonce值
		nonce, err := GetNonce(app.DbCon, address)
		if err != nil {
			hcommon.Log.Warnf("GetNonce err: [%T] %s", err, err.Error())
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
		chainID, err := ethclient.RpcNetworkID(context.Background())
		if err != nil {
			hcommon.Log.Warnf("RpcNetworkID err: [%T] %s", err, err.Error())
			return
		}
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
				sendRows = append(sendRows, &model.DBTSend{
					RelatedType:  1,
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
		_, err = model.SQLCreateManyTSend(
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
}

// CheckRawTxSend 发送交易
func CheckRawTxSend() {
	lockKey := "EthCheckRawTxSend"
	ok, err := app.GetLock(
		context.Background(),
		app.DbCon,
		lockKey,
	)
	if err != nil {
		hcommon.Log.Warnf("GetLock err: [%T] %s", err, err.Error())
		return
	}
	if !ok {
		return
	}
	defer func() {
		err := app.ReleaseLock(
			context.Background(),
			app.DbCon,
			lockKey,
		)
		if err != nil {
			hcommon.Log.Warnf("ReleaseLock err: [%T] %s", err, err.Error())
			return
		}
	}()
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
	// 提币
	var withdrawIDs []int64
	withdrawMap := make(map[int64]*model.DBTWithdraw)
	for _, sendRow := range sendRows {
		if sendRow.RelatedType == app.SendRelationTypeWithdraw {
			// 提币
			if !hcommon.IsIntInSlice(withdrawIDs, sendRow.RelatedID) {
				withdrawIDs = append(withdrawIDs, sendRow.RelatedID)
			}
		}
	}
	withdrawRows, err := model.SQLSelectTWithdrawCol(
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
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
	var productIDs []int64
	for _, withdrawRow := range withdrawRows {
		withdrawMap[withdrawRow.ID] = withdrawRow

		if !hcommon.IsIntInSlice(productIDs, withdrawRow.ProductID) {
			productIDs = append(productIDs, withdrawRow.ProductID)
		}
	}
	productRows, err := model.SQLSelectTProductCol(
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
	productMap := make(map[int64]*model.DBTProduct)
	for _, productRow := range productRows {
		productMap[productRow.ID] = productRow
	}
	// 执行发送
	var txHashes []string
	var notifyRows []*model.DBTProductNotify
	now := time.Now().Unix()
	for _, sendRow := range sendRows {
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
		if sendRow.RelatedType == app.SendRelationTypeWithdraw {
			// 提币
			withdrawRow := withdrawMap[sendRow.RelatedID]
			productRow := productMap[withdrawRow.ProductID]
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
				URL:          productRow.CbURL,
				Msg:          string(req),
				HandleStatus: app.NotifyStatusInit,
				HandleMsg:    "",
				CreateTime:   now,
				UpdateTime:   now,
			})
		}
		txHashes = append(txHashes, strings.ToLower(tx.Hash().Hex()))
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
	_, err = app.SQLUpdateTWithdrawStatusByTxIDs(
		context.Background(),
		app.DbCon,
		txHashes,
		model.DBTWithdraw{
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
	_, err = app.SQLUpdateTSendStatusByTxIDs(
		context.Background(),
		app.DbCon,
		txHashes,
		model.DBTSend{
			HandleStatus: app.TxOrgStatusSend,
			HandleMsg:    "send",
			HandleTime:   now,
		},
	)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
}

// CheckRawTxConfirm 确认tx是否打包完成
func CheckRawTxConfirm() {
	lockKey := "EthCheckRawTxConfirm"
	ok, err := app.GetLock(
		context.Background(),
		app.DbCon,
		lockKey,
	)
	if err != nil {
		hcommon.Log.Warnf("GetLock err: [%T] %s", err, err.Error())
		return
	}
	if !ok {
		return
	}
	defer func() {
		err := app.ReleaseLock(
			context.Background(),
			app.DbCon,
			lockKey,
		)
		if err != nil {
			hcommon.Log.Warnf("ReleaseLock err: [%T] %s", err, err.Error())
			return
		}
	}()

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
	withdrawMap := make(map[int64]*model.DBTWithdraw)
	for _, sendRow := range sendRows {
		if sendRow.RelatedType == app.SendRelationTypeWithdraw {
			// 提币
			if !hcommon.IsIntInSlice(withdrawIDs, sendRow.RelatedID) {
				withdrawIDs = append(withdrawIDs, sendRow.RelatedID)
			}
		}
	}
	withdrawRows, err := model.SQLSelectTWithdrawCol(
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
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
	var productIDs []int64
	for _, withdrawRow := range withdrawRows {
		withdrawMap[withdrawRow.ID] = withdrawRow

		if !hcommon.IsIntInSlice(productIDs, withdrawRow.ProductID) {
			productIDs = append(productIDs, withdrawRow.ProductID)
		}
	}
	productRows, err := model.SQLSelectTProductCol(
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
	productMap := make(map[int64]*model.DBTProduct)
	for _, productRow := range productRows {
		productMap[productRow.ID] = productRow
	}

	now := time.Now().Unix()
	var notifyRows []*model.DBTProductNotify
	var sendIDs []int64
	var withdrawUpdateIDs []int64
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
			withdrawRow := withdrawMap[sendRow.RelatedID]
			productRow := productMap[withdrawRow.ProductID]
			nonce := hcommon.GetUUIDStr()
			reqObj := gin.H{
				"tx_hash":     sendRow.TxID,
				"balance":     withdrawRow.BalanceReal,
				"app_name":    productRow.AppName,
				"out_serial":  withdrawRow.OutSerial,
				"address":     withdrawRow.ToAddress,
				"symbol":      "eth",
				"notify_type": app.NotifyTypeWithdrawConfirm,
			}
			reqObj["sign"] = hcommon.GetSign(productRow.AppSk, reqObj)
			req, err := json.Marshal(reqObj)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				continue
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
			withdrawUpdateIDs = append(withdrawUpdateIDs, sendRow.RelatedID)
		}
		// 完成
		sendIDs = append(sendIDs, sendRow.ID)
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
		withdrawUpdateIDs,
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
}

// CheckWithdraw 检测提现
func CheckWithdraw() {
	lockKey := "EthCheckWithdraw"
	ok, err := app.GetLock(
		context.Background(),
		app.DbCon,
		lockKey,
	)
	if err != nil {
		hcommon.Log.Warnf("GetLock err: [%T] %s", err, err.Error())
		return
	}
	if !ok {
		return
	}
	defer func() {
		err := app.ReleaseLock(
			context.Background(),
			app.DbCon,
			lockKey,
		)
		if err != nil {
			hcommon.Log.Warnf("ReleaseLock err: [%T] %s", err, err.Error())
			return
		}
	}()

	// 获取热钱包地址
	hotRow, err := app.SQLGetTAppConfigStrByK(
		context.Background(),
		app.DbCon,
		"hot_wallet_address",
	)
	if err != nil {
		hcommon.Log.Warnf("SQLGetTAppConfigInt err: [%T] %s", err, err.Error())
		return
	}
	if hotRow == nil {
		hcommon.Log.Errorf("no config int of hot_wallet_address")
		return
	}
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	if !re.MatchString(hotRow.V) {
		hcommon.Log.Errorf("config int hot_wallet_address err: %s", hotRow.V)
		return
	}
	hotAddress := common.HexToAddress(hotRow.V)
	hcommon.Log.Debugf("hotAddress: %s", hotAddress)
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
		hcommon.Log.Warnf("SQLGetTAddressKeyColByAddress err: [%T] %s", err, err.Error())
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
		hcommon.Log.Warnf("HexToECDSA err: [%T] %s", err, err.Error())
		return
	}
	hcommon.Log.Debugf("privateKey: %v", privateKey)
	withdrawRows, err := app.SQLSelectTWithdrawColByStatus(
		context.Background(),
		app.DbCon,
		[]string{
			model.DBColTWithdrawID,
		},
		app.WithdrawStatusInit,
	)
	if err != nil {
		hcommon.Log.Warnf("SQLSelectTWithdrawColByStatus err: [%T] %s", err, err.Error())
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
		hcommon.Log.Warnf("RpcBalanceAt err: [%T] %s", err, err.Error())
		return
	}
	pendingBalance, err := app.SQLGetTSendPendingBalance(
		context.Background(),
		app.DbCon,
		hotRow.V,
	)
	if err != nil {
		hcommon.Log.Warnf("SQLGetTSendPendingBalance err: [%T] %s", err, err.Error())
		return
	}
	hotAddressBalance -= pendingBalance
	hcommon.Log.Debugf("hotAddressBalance: %d", hotAddressBalance)
	// 获取gap price
	gasRow, err := app.SQLGetTAppStatusIntByK(
		context.Background(),
		app.DbCon,
		"to_user_gas_price",
	)
	if err != nil {
		hcommon.Log.Warnf("SQLGetTAppStatusIntByK err: [%T] %s", err, err.Error())
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
		hcommon.Log.Warnf("RpcNetworkID err: [%T] %s", err, err.Error())
		return
	}

	for _, withdrawRow := range withdrawRows {
		err = handleWithdraw(withdrawRow.ID, chainID, hotRow.V, privateKey, &hotAddressBalance, gasLimit, gasPrice, feeValue)
		if err != nil {
			hcommon.Log.Warnf("RpcBalanceAt err: [%T] %s", err, err.Error())
			return
		}
	}
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
	hcommon.Log.Debugf("balance: %d", balance)
	*hotAddressBalance -= balance + feeValue
	if *hotAddressBalance < 0 {
		hcommon.Log.Warnf("hot balance limit")
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
	tx := types.NewTransaction(
		uint64(nonce),
		common.HexToAddress(withdrawRow.ToAddress),
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

func CheckTxNotify() {
	lockKey := "CheckTxNotify"
	ok, err := app.GetLock(
		context.Background(),
		app.DbCon,
		lockKey,
	)
	if err != nil {
		hcommon.Log.Warnf("GetLock err: [%T] %s", err, err.Error())
		return
	}
	if !ok {
		return
	}
	defer func() {
		err := app.ReleaseLock(
			context.Background(),
			app.DbCon,
			lockKey,
		)
		if err != nil {
			hcommon.Log.Warnf("ReleaseLock err: [%T] %s", err, err.Error())
			return
		}
	}()

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
	productRows, err := model.SQLSelectTProductCol(
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
	productMap := make(map[int64]*model.DBTProduct)
	for _, productRow := range productRows {
		productMap[productRow.ID] = productRow
	}
	var notifyTxIDs []int64
	var notifyRows []*model.DBTProductNotify
	now := time.Now().Unix()
	for _, txRow := range txRows {
		productRow, ok := productMap[txRow.ProductID]
		if !ok {
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
			return
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
}

func CheckErc20BlockSeek() {
	lockKey := "CheckErc20BlockSeek"
	ok, err := app.GetLock(
		context.Background(),
		app.DbCon,
		lockKey,
	)
	if err != nil {
		hcommon.Log.Warnf("GetLock err: [%T] %s", err, err.Error())
		return
	}
	if !ok {
		return
	}
	defer func() {
		err := app.ReleaseLock(
			context.Background(),
			app.DbCon,
			lockKey,
		)
		if err != nil {
			hcommon.Log.Warnf("ReleaseLock err: [%T] %s", err, err.Error())
			return
		}
	}()

	// 获取配置 延迟确认数
	confirmRow, err := app.SQLGetTAppConfigIntByK(
		context.Background(),
		app.DbCon,
		"block_confirm_num",
	)
	if err != nil {
		hcommon.Log.Warnf("err: [%T] %s", err, err.Error())
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
		hcommon.Log.Warnf("err: [%T] %s", err, err.Error())
		return
	}
	if seekRow == nil {
		hcommon.Log.Errorf("no config int of erc20_seek_num")
		return
	}
	// rpc 获取当前最新区块数
	rpcBlockNum, err := ethclient.RpcBlockNumber(context.Background())
	if err != nil {
		hcommon.Log.Warnf("err: [%T] %s", err, err.Error())
		return
	}
	startI := seekRow.V + 1
	endI := rpcBlockNum - confirmRow.V
	if startI < endI {
		// 读取abi
		type LogTransfer struct {
			From   string
			To     string
			Tokens *big.Int
		}
		contractAbi, err := abi.JSON(strings.NewReader(ethclient.EthABI))
		if err != nil {
			hcommon.Log.Warnf("err: [%T] %s", err, err.Error())
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
			hcommon.Log.Warnf("err: [%T] %s", err, err.Error())
			return
		}
		for _, contractRow := range configTokenRows {
			configTokenRowAddresses = append(configTokenRowAddresses, contractRow.TokenAddress)
			configTokenRowMap[contractRow.TokenAddress] = contractRow
		}
		// 遍历获取需要查询的block信息
		for i := startI; i < endI; i++ {
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
					hcommon.Log.Warnf("SQLSelectTAppConfigTokenColAll err: [%T] %s", err, err.Error())
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
					toAddress := strings.ToLower(common.HexToAddress(log.Topics[2].Hex()).Hex())
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
					hcommon.Log.Warnf("err: [%T] %s", err, err.Error())
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
							hcommon.Log.Warnf("err: [%T] %s", err, err.Error())
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
							hcommon.Log.Warnf("err: [%T] %s", err, err.Error())
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
							hcommon.Log.Warnf("err: [%T] %s", err, err.Error())
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
					hcommon.Log.Warnf("err: [%T] %s", err, err.Error())
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
				hcommon.Log.Warnf("SQLUpdateTAppStatusIntByK err: [%T] %s", err, err.Error())
				return
			}
		}
	}
}

func CheckErc20TxNotify() {
	lockKey := "CheckErc20TxNotify"
	ok, err := app.GetLock(
		context.Background(),
		app.DbCon,
		lockKey,
	)
	if err != nil {
		hcommon.Log.Warnf("GetLock err: [%T] %s", err, err.Error())
		return
	}
	if !ok {
		return
	}
	defer func() {
		err := app.ReleaseLock(
			context.Background(),
			app.DbCon,
			lockKey,
		)
		if err != nil {
			hcommon.Log.Warnf("ReleaseLock err: [%T] %s", err, err.Error())
			return
		}
	}()

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
	productRows, err := model.SQLSelectTProductCol(
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
	productMap := make(map[int64]*model.DBTProduct)
	for _, productRow := range productRows {
		productMap[productRow.ID] = productRow
	}
	tokenRows, err := model.SQLSelectTAppConfigTokenCol(
		context.Background(),
		app.DbCon,
		[]string{
			model.DBColTAppConfigTokenID,
			model.DBColTAppConfigTokenTokenSymbol,
		},
		tokenIDs,
	)
	tokenMap := make(map[int64]*model.DBTAppConfigToken)
	for _, tokenRow := range tokenRows {
		tokenMap[tokenRow.ID] = tokenRow
	}
	var notifyTxIDs []int64
	var notifyRows []*model.DBTProductNotify
	now := time.Now().Unix()
	for _, txRow := range txRows {
		productRow, ok := productMap[txRow.ProductID]
		if !ok {
			continue
		}
		tokenRow, ok := tokenMap[txRow.TokenID]
		if !ok {
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
}

func CheckErc20TxOrg() {
	lockKey := "CheckErc20TxOrg"
	ok, err := app.GetLock(
		context.Background(),
		app.DbCon,
		lockKey,
	)
	if err != nil {
		hcommon.Log.Warnf("GetLock err: [%T] %s", err, err.Error())
		return
	}
	if !ok {
		return
	}
	defer func() {
		err := app.ReleaseLock(
			context.Background(),
			app.DbCon,
			lockKey,
		)
		if err != nil {
			hcommon.Log.Warnf("ReleaseLock err: [%T] %s", err, err.Error())
			return
		}
	}()

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
	var productIDs []int64
	var toAddresses []string
	tokenMap := make(map[int64]*model.DBTAppConfigToken)
	productMap := make(map[int64]*model.DBTProduct)
	addressMap := make(map[string]*model.DBTAddressKey)

	for _, txRow := range txRows {
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

		// 转换为map
		txMap[txRow.ID] = txRow
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
		orgInfo.TxIDs = append(orgInfo.TxIDs, txRow.TokenID)
		orgInfo.TokenBalance += txRow.Balance

		// 待查询id
		if hcommon.IsIntInSlice(tokenIDs, txRow.TokenID) {
			tokenIDs = append(tokenIDs, txRow.TokenID)
		}
		if hcommon.IsIntInSlice(productIDs, txRow.ProductID) {
			productIDs = append(productIDs, txRow.ProductID)
		}
		if hcommon.IsStringInSlice(toAddresses, txRow.ToAddress) {
			toAddresses = append(toAddresses, txRow.ToAddress)
		}
	}
	tokenRows, err := model.SQLSelectTAppConfigTokenCol(
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
	for _, tokenRow := range tokenRows {
		tokenMap[tokenRow.ID] = tokenRow
	}
	productRows, err := model.SQLSelectTProductCol(
		context.Background(),
		app.DbCon,
		[]string{
			model.DBColTProductID,
		},
		productIDs,
	)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
	for _, productRow := range productRows {
		productMap[productRow.ID] = productRow
	}
	addressRows, err := app.SQLSelectTAddressKeyColByAddress(
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
	for _, addressRow := range addressRows {
		addressMap[addressRow.Address] = addressRow
	}

	// 计算转账token所需的手续费
	erc20GasRow, err := app.SQLGetTAppConfigIntByK(
		context.Background(),
		app.DbCon,
		"erc20_gas_use",
	)
	if err != nil {
		hcommon.Log.Warnf("SQLGetTAppConfigInt err: [%T] %s", err, err.Error())
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
		hcommon.Log.Warnf("SQLGetTAppStatusIntByK err: [%T] %s", err, err.Error())
		return
	}
	if gasPriceRow == nil {
		hcommon.Log.Errorf("no config int of to_cold_gas_price")
		return
	}
	chainID, err := ethclient.RpcNetworkID(context.Background())
	if err != nil {
		hcommon.Log.Warnf("RpcNetworkID err: [%T] %s", err, err.Error())
		return
	}
	erc20Fee := erc20GasRow.V * gasPriceRow.V
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
		}
		// 处理token转账
		addressRow, ok := addressMap[toAddress]
		if !ok {
			hcommon.Log.Errorf("addressMap no: %s", toAddress)
			return
		}
		key := hcommon.AesDecrypt(addressRow.Pwd, app.Cfg.AESKey)
		if len(key) == 0 {
			hcommon.Log.Errorf("error key of: %s", toAddress)
			return
		}
		if strings.HasPrefix(key, "0x") {
			key = key[2:]
		}
		privateKey, err := crypto.HexToECDSA(key)
		if err != nil {
			hcommon.Log.Warnf("HexToECDSA err: [%T] %s", err, err.Error())
			return
		}
		// 获取nonce值
		nonce, err := GetNonce(app.DbCon, toAddress)
		if err != nil {
			hcommon.Log.Warnf("GetNonce err: [%T] %s", err, err.Error())
			return
		}
		rpcTx, err := ethclient.RpcGenTokenTransfer(
			context.Background(),
			tokenRow.TokenAddress,
			&bind.TransactOpts{
				Nonce:    big.NewInt(nonce),
				GasPrice: big.NewInt(gasPriceRow.V),
				GasLimit: uint64(erc20GasRow.V),
			},
			tokenRow.ColdAddress,
			orgInfo.TokenBalance,
		)
		if err != nil {
			hcommon.Log.Warnf("err: [%T] %s", err, err.Error())
			return
		}
		signedTx, err := types.SignTx(rpcTx, types.NewEIP155Signer(big.NewInt(chainID)), privateKey)
		if err != nil {
			hcommon.Log.Warnf("RpcNetworkID err: [%T] %s", err, err.Error())
			return
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
		var sendTxIDs []int64
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
			sendTxIDs = append(sendTxIDs, txID)
		}
		_, err = model.SQLCreateIgnoreManyTSend(
			context.Background(),
			app.DbCon,
			sendRows,
		)
		if err != nil {
			hcommon.Log.Warnf("err: [%T] %s", err, err.Error())
			return
		}
		_, err = app.SQLUpdateTTxErc20OrgStatusByAddresses(
			context.Background(),
			app.DbCon,
			sendTxIDs,
			model.DBTTxErc20{
				OrgStatus: app.TxOrgStatusHex,
				OrgMsg:    "hex",
				OrgTime:   now,
			},
		)
		if err != nil {
			hcommon.Log.Warnf("err: [%T] %s", err, err.Error())
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
			hcommon.Log.Warnf("SQLGetTAppConfigInt err: [%T] %s", err, err.Error())
			return
		}
		if feeRow == nil {
			hcommon.Log.Errorf("no config int of fee_wallet_address")
			return
		}
		re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
		if !re.MatchString(feeRow.V) {
			hcommon.Log.Errorf("config int fee_wallet_address err: %s", feeRow.V)
			return
		}
		hotAddress := common.HexToAddress(feeRow.V)
		hcommon.Log.Debugf("hotAddress: %s", hotAddress)
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
			hcommon.Log.Warnf("SQLGetTAddressKeyColByAddress err: [%T] %s", err, err.Error())
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
			hcommon.Log.Warnf("HexToECDSA err: [%T] %s", err, err.Error())
			return
		}
		hcommon.Log.Debugf("privateKey: %v", privateKey)
		feeAddressBalance, err := ethclient.RpcBalanceAt(
			context.Background(),
			feeRow.V,
		)
		if err != nil {
			hcommon.Log.Warnf("RpcBalanceAt err: [%T] %s", err, err.Error())
			return
		}
		pendingBalance, err := app.SQLGetTSendPendingBalance(
			context.Background(),
			app.DbCon,
			feeRow.V,
		)
		if err != nil {
			hcommon.Log.Warnf("SQLGetTSendPendingBalance err: [%T] %s", err, err.Error())
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
			hcommon.Log.Warnf("SQLGetTAppStatusIntByK err: [%T] %s", err, err.Error())
			return
		}
		if gasRow == nil {
			hcommon.Log.Errorf("no config int of to_cold_gas_price")
			return
		}
		gasPrice := gasRow.V
		gasLimit := int64(21000)
		ethFee := gasLimit * gasPrice
		chainID, err := ethclient.RpcNetworkID(context.Background())
		if err != nil {
			hcommon.Log.Warnf("RpcNetworkID err: [%T] %s", err, err.Error())
			return
		}
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
				hcommon.Log.Warnf("err: [%T] %s", err, err.Error())
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
				hcommon.Log.Warnf("err: [%T] %s", err, err.Error())
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
			var sendTxIDs []int64
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
				sendTxIDs = append(sendTxIDs, txID)
			}
			_, err = model.SQLCreateIgnoreManyTSend(
				context.Background(),
				app.DbCon,
				sendRows,
			)
			if err != nil {
				hcommon.Log.Warnf("err: [%T] %s", err, err.Error())
				return
			}
			_, err = app.SQLUpdateTTxErc20OrgStatusByAddresses(
				context.Background(),
				app.DbCon,
				sendTxIDs,
				model.DBTTxErc20{
					OrgStatus: app.TxOrgStatusFeeHex,
					OrgMsg:    "fee hex",
					OrgTime:   now,
				},
			)
			if err != nil {
				hcommon.Log.Warnf("err: [%T] %s", err, err.Error())
				return
			}
		}
	}
}
