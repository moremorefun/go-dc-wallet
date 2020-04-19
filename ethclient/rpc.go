package ethclient

import (
	"context"
	"go-dc-wallet/hcommon"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/core/types"
)

var client *Client

// InitClient 初始化接口对象
func InitClient(uri string) {
	var err error
	client, err = Dial(uri)
	if err != nil {
		hcommon.Log.Fatalf("eth client dial error: [%T] %s", err, err.Error())
	}
}

// RpcBlockNumber 获取最新的block number
func RpcBlockNumber(ctx context.Context) (int64, error) {
	blockNum, err := client.GetBlockNumber(ctx)
	if nil != err {
		return 0, err
	}
	return int64(blockNum), nil
}

// RpcBlockByNum 获取block信息
func RpcBlockByNum(ctx context.Context, blockNum int64) (*types.Block, error) {
	resp, err := client.BlockByNumber(ctx, big.NewInt(blockNum))
	if nil != err {
		return nil, err
	}
	return resp, nil
}

// RpcNonceAt 获取nonce
func RpcNonceAt(ctx context.Context, address string) (int64, error) {
	count, err := client.NonceAt(
		ctx,
		common.HexToAddress(address),
		nil,
	)
	if nil != err {
		return 0, err
	}
	return int64(count), nil
}

// RpcNetworkID 获取block信息
func RpcNetworkID(ctx context.Context) (int64, error) {
	resp, err := client.NetworkID(ctx)
	if nil != err {
		return 0, err
	}
	return resp.Int64(), nil
}

// RpcSendTransaction 发送交易
func RpcSendTransaction(ctx context.Context, tx *types.Transaction) error {
	err := client.SendTransaction(
		ctx,
		tx,
	)
	if nil != err {
		return err
	}
	return nil
}

// RpcTransactionByHash 确认交易是否打包完成
func RpcTransactionByHash(ctx context.Context, txHashStr string) (*types.Transaction, error) {
	txHash := common.HexToHash(txHashStr)
	tx, isPending, err := client.TransactionByHash(ctx, txHash)
	if err != nil {
		return nil, err
	}
	if isPending {
		return nil, nil
	}
	return tx, nil
}
