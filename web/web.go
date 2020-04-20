package web

import (
	"encoding/json"
	"go-dc-wallet/app"
	"go-dc-wallet/app/model"
	"go-dc-wallet/hcommon"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

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
			model.DBColTProductAppSk,
			model.DBColTProductWhitelistIP,
		},
		req.AppName,
	)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error":   hcommon.ErrorInternal,
			"err_msg": hcommon.ErrorInternalMsg,
		})
		c.Abort()
		return
	}
	if productRow == nil {
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
		c.JSON(http.StatusOK, gin.H{
			"error":   hcommon.ErrorInternal,
			"err_msg": hcommon.ErrorInternalMsg,
		})
		c.Abort()
		return
	}
	if count <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"error":   hcommon.ErrorNonceRepeat,
			"err_msg": hcommon.ErrorNonceRepeatMsg,
		})
		c.Abort()
		return
	}
}

func postAddress(c *gin.Context) {
	var req struct {
		AppName string `json:"app_name" binding:"required"`
	}
	err := c.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		hcommon.Log.Warnf("req args error: %#v", err)
		hcommon.GinFillBindError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"error":   hcommon.ErrorSuccess,
		"err_msg": hcommon.ErrorSuccessMsg,
	})
}

func postWithdraw(c *gin.Context) {

}
