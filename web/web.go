package web

import (
	"fmt"
	"go-dc-wallet/app"
	"go-dc-wallet/eosclient"
	"go-dc-wallet/hbtc"
	"go-dc-wallet/heos"
	"go-dc-wallet/heth"
	"go-dc-wallet/model"
	"go-dc-wallet/value"
	"go-dc-wallet/xenv"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/moremorefun/mcommon"

	"github.com/btcsuite/btcutil"

	"github.com/shopspring/decimal"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func Start(r *gin.Engine) {
	r.POST("/api/address", productReq, postAddress)
	r.POST("/api/withdraw", productReq, postWithdraw)
}

func postAddress(c *gin.Context) {
	var req struct {
		Symbol string `json:"symbol" binding:"required" validate:"oneof=eth btc eos"`
	}
	err := c.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		mcommon.Log.Warnf("req args error: %#v", err)
		mcommon.GinFillBindError(c, err)
		return
	}
	productID := c.GetInt64("product_id")
	if productID == 0 {
		mcommon.GinDoRespInternalErr(c)
		return
	}
	var addressRow *model.DBTAddressKey
	var eosColdAddressValue string
	// 开始事物
	isUseGinErr := true
	err = mcommon.DbTransaction(c, xenv.DbCon, func(tx mcommon.DbExeAble) error {
		// 获取可用地址
		var err error
		addressRow, err = app.SQLGetTAddressKeyColFreeForUpdate(
			c,
			tx,
			[]string{
				model.DBColTAddressKeyID,
				model.DBColTAddressKeyAddress,
			},
			req.Symbol,
		)
		if err != nil {
			return err
		}
		if addressRow == nil {
			// 没有可用地址了
			mcommon.GinDoRespErr(
				c,
				value.ErrorNoFreeAddress,
				value.ErrorNoFreeAddressMsg,
				nil,
			)
			isUseGinErr = false
			return fmt.Errorf("no free address")
		}
		// 更新获取到的地址的使用状态
		count, err := mcommon.DbUpdateKV(
			c,
			tx,
			model.DbTableTAddressKey,
			mcommon.H{
				model.DBColShortTAddressKeyUseTag: productID,
			},
			[]string{
				model.DBColShortTAddressKeyID,
			},
			[]interface{}{
				addressRow.ID,
			},
		)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return err
		}
		if count <= 0 {
			return fmt.Errorf("update address use tag error")
		}
		if req.Symbol == "eos" {
			// 获取冷钱包地址
			eosColdAddressValue, err = app.SQLGetTAppConfigStrValueByK(
				c,
				tx,
				"cold_wallet_address_eos",
			)
			if err != nil {
				mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return err
			}
			eosColdAddressValue = strings.TrimSpace(eosColdAddressValue)
			if eosColdAddressValue == "" {
				mcommon.Log.Errorf("eosColdAddressValue null")
				return fmt.Errorf("eosColdAddressValue null")
			}
		}
		return nil
	})
	if err != nil {
		mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		if isUseGinErr {
			mcommon.GinDoRespInternalErr(c)
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"error":       mcommon.ErrorSuccess,
		"err_msg":     mcommon.ErrorSuccessMsg,
		"address":     addressRow.Address,
		"eos_address": eosColdAddressValue,
	})
}

func postWithdraw(c *gin.Context) {
	var req struct {
		Symbol    string `json:"symbol" binding:"required"`
		OutSerial string `json:"out_serial" binding:"required" validate:"max=40"`
		Address   string `json:"address" binding:"required"`
		Balance   string `json:"balance" binding:"required"`
		Memo      string `json:"memo" binding:"omitempty"`
	}
	err := c.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		mcommon.Log.Warnf("req args error: %#v", err)
		mcommon.GinFillBindError(c, err)
		return
	}
	// 将币种小写
	req.Symbol = strings.ToLower(req.Symbol)
	// 获取产品id
	productID := c.GetInt64("product_id")
	if productID == 0 {
		mcommon.GinDoRespInternalErr(c)
		return
	}
	// eth 信息
	tokenDecimalsMap := make(map[string]int64)
	ethSymbols := []string{heth.CoinSymbol}
	tokenDecimalsMap[heth.CoinSymbol] = 18
	// 获取所有eth代币币种
	tokenRows, err := model.SQLSelectTAppConfigTokenColKV(
		c,
		xenv.DbCon,
		[]string{
			model.DBColTAppConfigTokenTokenSymbol,
			model.DBColTAppConfigTokenTokenDecimals,
		},
		nil,
		nil,
		nil,
		nil,
	)
	if err != nil {
		mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		mcommon.GinDoRespInternalErr(c)
		return
	}
	for _, tokenRow := range tokenRows {
		tokenRow.TokenSymbol = strings.ToLower(tokenRow.TokenSymbol)

		ethSymbols = append(ethSymbols, tokenRow.TokenSymbol)
		tokenDecimalsMap[tokenRow.TokenSymbol] = tokenRow.TokenDecimals
	}
	// btc 信息
	btcSymbols := []string{hbtc.CoinSymbol}
	tokenDecimalsMap[hbtc.CoinSymbol] = 8
	tokenBtcRows, err := model.SQLSelectTAppConfigTokenBtcColKV(
		c,
		xenv.DbCon,
		[]string{
			model.DBColTAppConfigTokenBtcTokenSymbol,
		},
		nil,
		nil,
		nil,
		nil,
	)
	if err != nil {
		mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		mcommon.GinDoRespInternalErr(c)
		return
	}
	for _, tokenRow := range tokenBtcRows {
		tokenRow.TokenSymbol = strings.ToLower(tokenRow.TokenSymbol)

		btcSymbols = append(btcSymbols, tokenRow.TokenSymbol)
		tokenDecimalsMap[tokenRow.TokenSymbol] = 8
	}
	// eos 信息
	tokenDecimalsMap[heos.CoinSymbol] = 4
	// 验证金额
	tokenDecimals, ok := tokenDecimalsMap[req.Symbol]
	if !ok {
		mcommon.GinDoRespErr(
			c,
			value.ErrorSymbolNotSupport,
			value.ErrorSymbolNotSupportMsg,
			nil,
		)
		return
	}
	balanceObj, err := decimal.NewFromString(req.Balance)
	if err != nil {
		mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		mcommon.GinDoRespErr(
			c,
			value.ErrorBalanceFormat,
			value.ErrorBalanceFormatMsg,
			nil,
		)
		return
	}
	if balanceObj.LessThanOrEqual(decimal.NewFromInt(0)) {
		mcommon.GinDoRespErr(
			c,
			value.ErrorBalanceFormat,
			value.ErrorBalanceFormatMsg,
			nil,
		)
		return
	}
	if balanceObj.Exponent() < -int32(tokenDecimals) {
		mcommon.GinDoRespErr(
			c,
			value.ErrorBalanceFormat,
			value.ErrorBalanceFormatMsg,
			nil,
		)
		return
	}
	if mcommon.IsStringInSlice(ethSymbols, req.Symbol) {
		// 验证地址
		req.Address = strings.ToLower(req.Address)
		re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
		if !re.MatchString(req.Address) {
			mcommon.GinDoRespErr(
				c,
				value.ErrorAddressWrong,
				value.ErrorAddressWrongMsg,
				nil,
			)
			return
		}

	} else if mcommon.IsStringInSlice(btcSymbols, req.Symbol) {
		// 验证地址
		_, err := btcutil.DecodeAddress(
			req.Address,
			hbtc.GetNetwork(xenv.Cfg.BtcNetworkType).Params,
		)
		if err != nil {
			mcommon.GinDoRespErr(
				c,
				value.ErrorAddressWrong,
				value.ErrorAddressWrongMsg,
				nil,
			)
			return
		}
	} else if req.Symbol == heos.CoinSymbol {
		// eos
		// 验证地址
		_, err := eosclient.RpcChainGetAccount(
			req.Address,
		)
		if err != nil {
			mcommon.GinDoRespErr(
				c,
				value.ErrorAddressWrong,
				value.ErrorAddressWrongMsg,
				nil,
			)
			return
		}
	} else {
		mcommon.GinDoRespErr(
			c,
			value.ErrorSymbolNotSupport,
			value.ErrorSymbolNotSupportMsg,
			nil,
		)
		return
	}
	now := time.Now().Unix()
	_, err = model.SQLCreateTWithdraw(
		c,
		xenv.DbCon,
		&model.DBTWithdraw{
			ProductID:    productID,
			OutSerial:    req.OutSerial,
			ToAddress:    req.Address,
			Memo:         req.Memo,
			Symbol:       req.Symbol,
			BalanceReal:  req.Balance,
			TxHash:       "",
			CreateTime:   now,
			HandleStatus: 0,
			HandleMsg:    "",
			HandleTime:   now,
		},
		true,
	)
	if err != nil {
		mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		mcommon.GinDoRespInternalErr(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"error":   mcommon.ErrorSuccess,
		"err_msg": mcommon.ErrorSuccessMsg,
	})
}
