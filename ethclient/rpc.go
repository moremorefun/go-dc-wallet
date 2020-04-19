package ethclient

import (
	"context"
	"go-dc-wallet/hcommon"
	"math/big"

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
