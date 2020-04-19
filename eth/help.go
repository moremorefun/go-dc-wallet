package eth

import (
	"context"
	"go-dc-wallet/app"
	"go-dc-wallet/ethclient"
)

func GetNonce(address string) (int64, error) {
	// 通过rpc获取
	rpcNonce, err := ethclient.RpcNonceAt(
		context.Background(),
		address,
	)
	if nil != err {
		return 0, err
	}
	// 获取db nonce
	dbNonce, err := app.SQLGetTSendMaxNonce(
		context.Background(),
		app.DbCon,
		address,
	)
	if nil != err {
		return 0, err
	}
	if dbNonce > rpcNonce {
		rpcNonce = dbNonce
	}
	return rpcNonce, nil
}
