package eth

import (
	"context"
	"crypto/ecdsa"
	"go-dc-wallet/app"
	"go-dc-wallet/app/model"
	"go-dc-wallet/ethclient"
	"go-dc-wallet/hcommon"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/shopspring/decimal"
)

// CheckAddressCheck 检测是否有充足的备用地址
func CheckAddressCheck() {
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
	// 遍历获取需要查询的block信息
	for i := seekRow.V + 1; i < rpcBlockNum-confirmRow.V; i++ {
		rpcBlock, err := ethclient.RpcBlockByNum(context.Background(), i)
		if err != nil {
			hcommon.Log.Warnf("EthRpcBlockByNum err: [%T] %s", err, err.Error())
			return
		}
		// 将需要处理的数据存入字典和数组
		txMap := make(map[string][]*types.Transaction)
		var toAddresses []string
		for _, rpcTx := range rpcBlock.Transactions() {
			// 转账数额大于0
			if rpcTx.Value().Int64() > 0 && rpcTx.To() != nil {
				toAddress := strings.ToLower(rpcTx.To().Hex())
				txMap[toAddress] = append(txMap[toAddress], rpcTx)
				toAddresses = append(toAddresses, toAddress)
			}
		}
		// 从db中查询这些地址是否是冲币地址中的地址
		dbAddressRows, err := app.SQLSelectTAddressKeyColByAddress(
			context.Background(),
			app.DbCon,
			[]string{
				model.DBColTAddressKeyAddress,
			},
			toAddresses,
		)
		if err != nil {
			hcommon.Log.Warnf("dbAddressRows err: [%T] %s", err, err.Error())
			return
		}
		var dbTxRows []*model.DBTTx
		for _, dbAddressRow := range dbAddressRows {
			txes := txMap[dbAddressRow.Address]
			for _, tx := range txes {
				msg, err := tx.AsMessage(types.NewEIP155Signer(tx.ChainId()))
				if err != nil {
					hcommon.Log.Errorf("AsMessage err: [%T] %s", err, err.Error())
				}
				now := time.Now().Unix()
				balanceReal := decimal.NewFromInt(tx.Value().Int64()).Div(decimal.NewFromInt(1e18))
				dbTxRows = append(dbTxRows, &model.DBTTx{
					TxID:         tx.Hash().String(),
					FromAddress:  strings.ToLower(msg.From().Hex()),
					ToAddress:    strings.ToLower(tx.To().Hex()),
					Balance:      tx.Value().Int64(),
					BalanceReal:  balanceReal.String(),
					CreateTime:   now,
					HandleStatus: 0,
					HandleMsg:    "",
					HandleTime:   now,
				})
			}
		}
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
