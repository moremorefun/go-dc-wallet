package hcommon

import (
	"context"
	"strings"

	"github.com/regcostajr/go-web3"
	"github.com/regcostajr/go-web3/providers"
)

var client *web3.Web3

// EthInitClient 初始化接口对象
func EthInitClient(uri string) {
	isSecure := false
	if strings.HasPrefix(uri, "https://") {
		isSecure = true
		uri = strings.Replace(uri, "https://", "", 1)
	} else if strings.HasPrefix(uri, "http://") {
		uri = strings.Replace(uri, "http://", "", 1)
	}

	client = web3.NewWeb3(providers.NewHTTPProvider(uri, 60, isSecure))
}

// EthRpcBlockNumber 获取最新的block number
func EthRpcBlockNumber(ctx context.Context) (int64, error) {
	blockNum, err := client.Eth.GetBlockNumber()
	if nil != err {
		return 0, err
	}
	return blockNum.Int64(), nil
}
