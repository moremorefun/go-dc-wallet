package eth

import (
	"context"
	"crypto/ecdsa"
	"go-dc-wallet/app"
	"go-dc-wallet/app/model"
	"go-dc-wallet/hcommon"
	"go-dc-wallet/hcommon/eth"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ethereum/go-ethereum/crypto"
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
	rpcBlockNum, err := eth.RpcBlockNumber(context.Background())
	if err != nil {
		hcommon.Log.Warnf("RpcBlockNumber err: [%T] %s", err, err.Error())
		return
	}
	// 遍历获取需要查询的block信息
	for i := seekRow.V + 1; i < rpcBlockNum-confirmRow.V; i++ {
		rpcBlock, err := eth.RpcBlockByNum(context.Background(), i)
		if err != nil {
			hcommon.Log.Warnf("EthRpcBlockByNum err: [%T] %s", err, err.Error())
			return
		}
		hcommon.Log.Debugf("rpcBlock: %#v", rpcBlock)
	}
}
