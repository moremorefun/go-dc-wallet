package heth

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"go-dc-wallet/app"
	"go-dc-wallet/ethclient"
	"go-dc-wallet/model"
	"go-dc-wallet/xenv"
	"math/big"
	"regexp"
	"strings"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/moremorefun/mcommon"

	"github.com/shopspring/decimal"

	"github.com/ethereum/go-ethereum/common"
)

const (
	// EthToWei 数据单位
	EthToWei = 1e18
	// CoinSymbol 单位标志
	CoinSymbol = "eth"
)

// ethToWeiDecimal 转换单位
var ethToWeiDecimal decimal.Decimal

func init() {
	ethToWeiDecimal = decimal.NewFromInt(EthToWei)
}

// GetNonce 获取nonce值
func GetNonce(tx mcommon.DbExeAble, address string) (int64, error) {
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
func IsValidAddress(iaddress string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(iaddress)
}

// AddressBytesToStr 地址转化为字符串
func AddressBytesToStr(addressBytes common.Address) string {
	return strings.ToLower(addressBytes.Hex())
}

// StrToAddressBytes 字符串转化为地址
func StrToAddressBytes(str string) (common.Address, error) {
	if !IsValidAddress(str) {
		return common.Address{}, fmt.Errorf("str not address: %s", str)
	}
	return common.HexToAddress(str), nil
}

// CommonHashToAddrss 字节转换为地址
func CommonHashToAddrss(bytes common.Hash) common.Address {
	var b common.Address
	b.SetBytes(bytes[:])
	return b
}

// CommonHashToAddrssString 字节转换为地址
func CommonHashToAddrssString(bytes common.Hash) string {
	return CommonHashToAddrss(bytes).Hex()
}

// CommonHashToAddrssStringLower 字节转换为地址
func CommonHashToAddrssStringLower(bytes common.Hash) string {
	return strings.ToLower(CommonHashToAddrss(bytes).Hex())
}

// EthStrToWeiBigInit 转换金额 eth to wei
func EthStrToWeiBigInit(balanceRealStr string) (*big.Int, error) {
	balanceReal, err := decimal.NewFromString(balanceRealStr)
	if err != nil {
		return nil, err
	}
	balanceStr := balanceReal.Mul(ethToWeiDecimal).StringFixed(0)
	b := new(big.Int)
	_, ok := b.SetString(balanceStr, 10)
	if !ok {
		return nil, errors.New("error str to bigint")
	}
	return b, nil
}

// WeiBigIntToEthStr 转换金额 wei to eth
func WeiBigIntToEthStr(wei *big.Int) (string, error) {
	balance, err := decimal.NewFromString(wei.String())
	if err != nil {
		return "0", err
	}
	balanceStr := balance.Div(ethToWeiDecimal).StringFixed(18)
	return balanceStr, nil
}

// TokenEthStrToWeiBigInit 转换金额 eth to wei
func TokenEthStrToWeiBigInit(balanceRealStr string, tokenDecimals int64) (*big.Int, error) {
	balanceReal, err := decimal.NewFromString(balanceRealStr)
	if err != nil {
		return nil, err
	}
	balanceStr := balanceReal.Mul(decimal.NewFromInt(10).Pow(decimal.NewFromInt(tokenDecimals))).StringFixed(0)
	b := new(big.Int)
	_, ok := b.SetString(balanceStr, 10)
	if !ok {
		return nil, errors.New("error str to bigint")
	}
	return b, nil
}

// TokenWeiBigIntToEthStr 转换金额 wei to eth
func TokenWeiBigIntToEthStr(wei *big.Int, tokenDecimals int64) (string, error) {
	balance, err := decimal.NewFromString(wei.String())
	if err != nil {
		return "0", err
	}
	balanceStr := balance.Div(decimal.NewFromInt(10).Pow(decimal.NewFromInt(tokenDecimals))).StringFixed(int32(tokenDecimals))
	return balanceStr, nil
}

// GetPKMapOfAddresses 获取地址私钥
func GetPKMapOfAddresses(ctx context.Context, db mcommon.DbExeAble, addresses []string) (map[string]*ecdsa.PrivateKey, error) {
	addressPKMap := make(map[string]*ecdsa.PrivateKey)
	addressKeyMap, err := app.SQLGetAddressKeyMap(
		ctx,
		db,
		[]string{
			model.DBColTAddressKeyID,
			model.DBColTAddressKeyAddress,
			model.DBColTAddressKeyPwd,
		},
		addresses,
	)
	if err != nil {
		mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return nil, err
	}
	for k, v := range addressKeyMap {
		key, err := mcommon.AesDecrypt(v.Pwd, xenv.Cfg.AESKey)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			continue
		}
		if len(key) == 0 {
			mcommon.Log.Errorf("error key of: %s", k)
			continue
		}
		key = strings.TrimPrefix(key, "0x")
		privateKey, err := crypto.HexToECDSA(key)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			continue
		}
		addressPKMap[k] = privateKey
	}
	return addressPKMap, nil
}

// GetPkOfAddress 获取地址私钥
func GetPkOfAddress(ctx context.Context, db mcommon.DbExeAble, address string) (*ecdsa.PrivateKey, error) {
	// 获取私钥
	keyRow, err := app.SQLGetTAddressKeyColByAddress(
		ctx,
		db,
		[]string{
			model.DBColTAddressKeyPwd,
		},
		address,
	)
	if err != nil {
		mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return nil, err
	}
	if keyRow == nil {
		mcommon.Log.Errorf("no key of: %s", address)
		return nil, fmt.Errorf("no key of: %s", address)
	}
	key, err := mcommon.AesDecrypt(keyRow.Pwd, xenv.Cfg.AESKey)
	if err != nil {
		mcommon.Log.Errorf("HexToECDSA err: [%T] %s", err, err.Error())
		return nil, err
	}
	if len(key) == 0 {
		mcommon.Log.Errorf("error key of: %s", address)
		return nil, fmt.Errorf("no key of: %s", address)
	}
	key = strings.TrimPrefix(key, "0x")
	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		mcommon.Log.Errorf("HexToECDSA err: [%T] %s", err, err.Error())
		return nil, err
	}
	return privateKey, nil
}

func NewSignTransaction(chainID int64, nonce uint64, to common.Address, amount *big.Int, gasLimit uint64, gasFeeCap *big.Int, gasTipCap *big.Int, data []byte, prv *ecdsa.PrivateKey) (*types.Transaction, error) {
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   big.NewInt(chainID),
		Nonce:     nonce,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       gasLimit,
		To:        &to,
		Value:     amount,
		Data:      data,
	})
	signedTx, err := types.SignTx(tx, types.NewLondonSigner(big.NewInt(chainID)), prv)
	if err != nil {
		mcommon.Log.Warnf("err: [%T] %s", err, err.Error())
		return nil, err
	}
	return signedTx, nil
}
