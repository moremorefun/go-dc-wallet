package web

import (
	"encoding/json"
	"go-dc-wallet/app"
	"go-dc-wallet/app/model"
	"go-dc-wallet/hbtc"
	"go-dc-wallet/hcommon"
	"go-dc-wallet/heth"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/btcsuite/btcutil"

	"github.com/shopspring/decimal"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func Start(r *gin.Engine) {
	r.POST("/api/address", productReq, postAddress)
	r.POST("/api/withdraw", productReq, postWithdraw)
}

func productReq(c *gin.Context) {
	var req struct {
		AppName string `json:"app_name" binding:"required"`
		Nonce   string `json:"nonce" binding:"required" validate:"max=40"`
		Sign    string `json:"sign" binding:"required"`
	}
	err := c.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		hcommon.Log.Warnf("req args error: %#v", err)
		hcommon.GinFillBindError(c, err)
		c.Abort()
		return
	}
	// 获取产品信息
	productRow, err := app.SQLGetTProductColByName(
		c,
		app.DbCon,
		[]string{
			model.DBColTProductID,
			model.DBColTProductAppSk,
			model.DBColTProductWhitelistIP,
		},
		req.AppName,
	)
	if err != nil {
		hcommon.Log.Warnf("err: [%T] %s", err.Error())
		c.JSON(http.StatusOK, gin.H{
			"error":   hcommon.ErrorInternal,
			"err_msg": hcommon.ErrorInternalMsg,
		})
		c.Abort()
		return
	}
	if productRow == nil {
		hcommon.Log.Warnf("no product of: %s", req.AppName)
		c.JSON(http.StatusOK, gin.H{
			"error":   hcommon.ErrorNoProduct,
			"err_msg": hcommon.ErrorNoProductMsg,
		})
		c.Abort()
		return
	}
	// 对比ip白名单
	if len(productRow.WhitelistIP) > 0 {
		if !strings.Contains(productRow.WhitelistIP, c.ClientIP()) {
			hcommon.Log.Warnf("no in ip list of: %s %s", req.AppName, c.ClientIP())
			c.JSON(http.StatusOK, gin.H{
				"error":   hcommon.ErrorIPLimit,
				"err_msg": hcommon.ErrorIPLimitMsg,
			})
			c.Abort()
			return
		}
	}
	// 验证签名
	var body []byte
	if cb, ok := c.Get(gin.BodyBytesKey); ok {
		if cbb, ok := cb.([]byte); ok {
			body = cbb
		}
	}
	if body == nil {
		body, err = ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error":   hcommon.ErrorInternal,
				"err_msg": hcommon.ErrorInternalMsg,
			})
			c.Abort()
			return
		}
		c.Set(gin.BodyBytesKey, body)
	}
	oldObj := gin.H{}
	err = json.Unmarshal(body, &oldObj)
	if err != nil {
		hcommon.Log.Warnf("req body error")
		c.JSON(http.StatusOK, gin.H{
			"error":   hcommon.ErrorInternal,
			"err_msg": hcommon.ErrorInternalMsg,
		})
		c.Abort()
		return
	}
	checkObj := gin.H{}
	for k, v := range oldObj {
		if k != "sign" {
			checkObj[k] = v
		}
	}
	checkSign := hcommon.GetSign(productRow.AppSk, checkObj)
	if checkSign == "" || checkSign != req.Sign {
		hcommon.Log.Warnf("sign error of: %s", req.AppName)
		c.JSON(http.StatusOK, gin.H{
			"error":   hcommon.ErrorSignWrong,
			"err_msg": hcommon.ErrorSignWrongMsg,
		})
		c.Abort()
		return
	}
	// 检测nonce
	count, err := model.SQLCreateIgnoreTProductNonce(
		c,
		app.DbCon,
		&model.DBTProductNonce{
			C:          req.Nonce,
			CreateTime: time.Now().Unix(),
		},
	)
	if err != nil {
		hcommon.Log.Warnf("err: [%T] %s", err, err.Error())
		c.JSON(http.StatusOK, gin.H{
			"error":   hcommon.ErrorInternal,
			"err_msg": hcommon.ErrorInternalMsg,
		})
		c.Abort()
		return
	}
	if count <= 0 {
		hcommon.Log.Warnf("nonce repeated")
		c.JSON(http.StatusOK, gin.H{
			"error":   hcommon.ErrorNonceRepeat,
			"err_msg": hcommon.ErrorNonceRepeatMsg,
		})
		c.Abort()
		return
	}
	c.Set("product_id", productRow.ID)
}

func postAddress(c *gin.Context) {
	var req struct {
		Symbol string `json:"symbol" binding:"required" validate:"oneof=eth btc eos"`
	}
	err := c.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		hcommon.Log.Warnf("req args error: %#v", err)
		hcommon.GinFillBindError(c, err)
		return
	}
	productID := c.GetInt64("product_id")
	if productID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"error":   hcommon.ErrorInternal,
			"err_msg": hcommon.ErrorInternalMsg,
		})
		return
	}
	// 开始事物
	isComment := false
	tx, err := app.DbCon.BeginTxx(c, nil)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		c.JSON(http.StatusOK, gin.H{
			"error":   hcommon.ErrorInternal,
			"err_msg": hcommon.ErrorInternalMsg,
		})
		return
	}
	defer func() {
		if !isComment {
			_ = tx.Rollback()
		}
	}()
	addressRow, err := app.SQLGetTAddressKeyColFreeForUpdate(
		c,
		tx,
		[]string{
			model.DBColTAddressKeyID,
			model.DBColTAddressKeyAddress,
		},
		req.Symbol,
	)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		c.JSON(http.StatusOK, gin.H{
			"error":   hcommon.ErrorInternal,
			"err_msg": hcommon.ErrorInternalMsg,
		})
		return
	}
	if addressRow == nil {
		c.JSON(http.StatusOK, gin.H{
			"error":   hcommon.ErrorNoFreeAddress,
			"err_msg": hcommon.ErrorNoFreeAddressMsg,
		})
		return
	}
	count, err := app.SQLUpdateTAddressKeyUseTag(
		c,
		tx,
		&model.DBTAddressKey{
			ID:     addressRow.ID,
			UseTag: productID,
		},
	)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		c.JSON(http.StatusOK, gin.H{
			"error":   hcommon.ErrorInternal,
			"err_msg": hcommon.ErrorInternalMsg,
		})
		return
	}
	if count <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"error":   hcommon.ErrorInternal,
			"err_msg": hcommon.ErrorInternalMsg,
		})
		return
	}
	// 获取冷钱包地址
	eosColdAddressValue, err := app.SQLGetTAppConfigStrValueByK(
		c,
		tx,
		"cold_wallet_address_eos",
	)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		c.JSON(http.StatusOK, gin.H{
			"error":   hcommon.ErrorInternal,
			"err_msg": hcommon.ErrorInternalMsg,
		})
		return
	}
	eosColdAddressValue = strings.TrimSpace(eosColdAddressValue)
	if eosColdAddressValue == "" {
		hcommon.Log.Errorf("eosColdAddressValue null")
		c.JSON(http.StatusOK, gin.H{
			"error":   hcommon.ErrorInternal,
			"err_msg": hcommon.ErrorInternalMsg,
		})
		return
	}
	// 提交事物
	err = tx.Commit()
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		c.JSON(http.StatusOK, gin.H{
			"error":   hcommon.ErrorInternal,
			"err_msg": hcommon.ErrorInternalMsg,
		})
		return
	}
	isComment = true
	c.JSON(http.StatusOK, gin.H{
		"error":       hcommon.ErrorSuccess,
		"err_msg":     hcommon.ErrorSuccessMsg,
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
	}
	err := c.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		hcommon.Log.Warnf("req args error: %#v", err)
		hcommon.GinFillBindError(c, err)
		return
	}
	req.Symbol = strings.ToLower(req.Symbol)

	productID := c.GetInt64("product_id")
	if productID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"error":   hcommon.ErrorInternal,
			"err_msg": hcommon.ErrorInternalMsg,
		})
		return
	}
	tokenDecimalsMap := make(map[string]int64)
	ethSymbols := []string{heth.CoinSymbol}
	tokenDecimalsMap[heth.CoinSymbol] = 18
	tokenRows, err := app.SQLSelectTAppConfigTokenColAll(
		c,
		app.DbCon,
		[]string{
			model.DBColTAppConfigTokenTokenSymbol,
			model.DBColTAppConfigTokenTokenDecimals,
		},
	)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		c.JSON(http.StatusOK, gin.H{
			"error":   hcommon.ErrorInternal,
			"err_msg": hcommon.ErrorInternalMsg,
		})
		return
	}
	for _, tokenRow := range tokenRows {
		tokenRow.TokenSymbol = strings.ToLower(tokenRow.TokenSymbol)

		ethSymbols = append(ethSymbols, tokenRow.TokenSymbol)
		tokenDecimalsMap[tokenRow.TokenSymbol] = tokenRow.TokenDecimals
	}
	btcSymbols := []string{hbtc.CoinSymbol}
	tokenDecimalsMap[hbtc.CoinSymbol] = 8
	tokenBtcRows, err := app.SQLSelectTAppConfigTokenBtcColAll(
		c,
		app.DbCon,
		[]string{
			model.DBColTAppConfigTokenBtcTokenSymbol,
		},
	)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		c.JSON(http.StatusOK, gin.H{
			"error":   hcommon.ErrorInternal,
			"err_msg": hcommon.ErrorInternalMsg,
		})
		return
	}
	for _, tokenRow := range tokenBtcRows {
		tokenRow.TokenSymbol = strings.ToLower(tokenRow.TokenSymbol)

		btcSymbols = append(btcSymbols, tokenRow.TokenSymbol)
		tokenDecimalsMap[tokenRow.TokenSymbol] = 8
	}
	tokenDecimals, ok := tokenDecimalsMap[req.Symbol]
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"error":   hcommon.ErrorSymbolNotSupport,
			"err_msg": hcommon.ErrorSymbolNotSupportMsg,
		})
		return
	}
	// 验证金额
	balanceObj, err := decimal.NewFromString(req.Balance)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		c.JSON(http.StatusOK, gin.H{
			"error":   hcommon.ErrorBalanceFormat,
			"err_msg": hcommon.ErrorBalanceFormatMsg,
		})
		return
	}
	if balanceObj.LessThanOrEqual(decimal.NewFromInt(0)) {
		c.JSON(http.StatusOK, gin.H{
			"error":   hcommon.ErrorBalanceFormat,
			"err_msg": hcommon.ErrorBalanceFormatMsg,
		})
		return
	}
	if balanceObj.Exponent() < -int32(tokenDecimals) {
		c.JSON(http.StatusOK, gin.H{
			"error":   hcommon.ErrorBalanceFormat,
			"err_msg": hcommon.ErrorBalanceFormatMsg,
		})
		return
	}
	if hcommon.IsStringInSlice(ethSymbols, req.Symbol) {
		// 验证地址
		req.Address = strings.ToLower(req.Address)
		re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
		if !re.MatchString(req.Address) {
			c.JSON(http.StatusOK, gin.H{
				"error":   hcommon.ErrorAddressWrong,
				"err_msg": hcommon.ErrorAddressWrongMsg,
			})
			return
		}

	} else if hcommon.IsStringInSlice(btcSymbols, req.Symbol) {
		// 验证地址
		_, err := btcutil.DecodeAddress(
			req.Address,
			hbtc.GetNetwork(app.Cfg.BtcNetworkType).Params,
		)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error":   hcommon.ErrorAddressWrong,
				"err_msg": hcommon.ErrorAddressWrongMsg,
			})
			return
		}
	} else {
		c.JSON(http.StatusOK, gin.H{
			"error":   hcommon.ErrorSymbolNotSupport,
			"err_msg": hcommon.ErrorSymbolNotSupportMsg,
		})
		return
	}
	now := time.Now().Unix()
	// 开始事物
	isComment := false
	tx, err := app.DbCon.BeginTxx(c, nil)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		c.JSON(http.StatusOK, gin.H{
			"error":   hcommon.ErrorInternal,
			"err_msg": hcommon.ErrorInternalMsg,
		})
		return
	}
	defer func() {
		if !isComment {
			_ = tx.Rollback()
		}
	}()
	_, err = model.SQLCreateIgnoreTWithdraw(
		c,
		tx,
		&model.DBTWithdraw{
			ProductID:    productID,
			OutSerial:    req.OutSerial,
			ToAddress:    req.Address,
			Symbol:       req.Symbol,
			BalanceReal:  req.Balance,
			TxHash:       "",
			CreateTime:   now,
			HandleStatus: 0,
			HandleMsg:    "",
			HandleTime:   now,
		},
	)
	// 提交事物
	err = tx.Commit()
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		c.JSON(http.StatusOK, gin.H{
			"error":   hcommon.ErrorInternal,
			"err_msg": hcommon.ErrorInternalMsg,
		})
		return
	}
	isComment = true
	c.JSON(http.StatusOK, gin.H{
		"error":   hcommon.ErrorSuccess,
		"err_msg": hcommon.ErrorSuccessMsg,
	})
}
