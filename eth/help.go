package eth

import (
	"context"
	"fmt"
	"go-dc-wallet/app"
	"go-dc-wallet/ethclient"
	"go-dc-wallet/hcommon"
	"regexp"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

// GetNonce 获取nonce值
func GetNonce(tx hcommon.DbExeAble, address string) (int64, error) {
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
		tx,
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

// IsValidAddress validate hex address
func IsValidAddress(iaddress interface{}) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	switch v := iaddress.(type) {
	case string:
		return re.MatchString(v)
	case common.Address:
		return re.MatchString(v.Hex())
	default:
		return false
	}
}

// AddressBytesToStr 地址转化为字符串
func AddressBytesToStr(addressBytes common.Address) string {
	return strings.ToLower(addressBytes.Hex())
}

// StrToAddressBytes 字符串转化为地址
func StrToAddressBytes(str string) (common.Address, error) {
	if !IsValidAddress(str) {
		return common.HexToAddress("0x0"), fmt.Errorf("str not address: %s", str)
	}
	return common.HexToAddress(str), nil
}
