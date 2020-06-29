package heos

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

const (
	CoinSymbol  = "eos"
	MiniAddress = 1000000000
)

// EosValueToDecimal 获取金额
func EosValueToDecimal(quantity string) (decimal.Decimal, error) {
	if quantity == "" {
		return decimal.NewFromInt(0).RoundBank(4), nil
	}
	quantitys := strings.Split(quantity, " ")
	if len(quantitys) != 2 {
		return decimal.NewFromInt(0), fmt.Errorf("error value: %s", quantity)
	}
	if quantitys[1] != "EOS" {
		return decimal.NewFromInt(0), fmt.Errorf("error value: %s", quantity)
	}
	v, err := decimal.NewFromString(quantitys[0])
	if err != nil {
		return decimal.NewFromInt(0), err
	}
	return v, nil
}

// EosValueToStr 获取金额
func EosValueToStr(quantity string) (string, error) {
	quantitys := strings.Split(quantity, " ")
	if len(quantitys) != 2 {
		return "0", fmt.Errorf("error value: %s")
	}
	if quantitys[1] != "EOS" {
		return "0", fmt.Errorf("error value: %s")
	}
	return quantitys[0], nil
}

// StrToEosDecimal 字符串转数额
func StrToEosDecimal(balanceReal string) (decimal.Decimal, error) {
	v, err := decimal.NewFromString(balanceReal)
	if err != nil {
		return decimal.NewFromInt(0), err
	}
	v = v.RoundBank(4)
	return v, nil
}
