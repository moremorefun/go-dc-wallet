package web

import (
	"encoding/json"
	"go-dc-wallet/model"
	"go-dc-wallet/value"
	"go-dc-wallet/xenv"
	"io/ioutil"
	"strings"
	"time"

	"github.com/moremorefun/mcommon"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func productReq(c *gin.Context) {
	var req struct {
		AppName string `json:"app_name" binding:"required"`
		Nonce   string `json:"nonce" binding:"required" validate:"max=40"`
		Sign    string `json:"sign" binding:"required"`
	}
	err := c.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		mcommon.Log.Warnf("req args error: %#v", err)
		mcommon.GinFillBindError(c, err)
		c.Abort()
		return
	}
	// 获取产品信息
	productRow, err := model.SQLGetTProductColKV(
		c,
		xenv.DbCon,
		[]string{
			model.DBColTProductID,
			model.DBColTProductAppSk,
			model.DBColTProductWhitelistIP,
		},
		[]string{
			model.DBColShortTProductAppName,
		},
		[]interface{}{
			req.AppName,
		},
	)
	if err != nil {
		mcommon.Log.Warnf("err: [%T] %s", err.Error())
		mcommon.GinDoRespInternalErr(c)
		c.Abort()
		return
	}
	if productRow == nil {
		mcommon.Log.Warnf("no product of: %s", req.AppName)
		mcommon.GinDoRespErr(
			c,
			value.ErrorNoProduct,
			value.ErrorNoProductMsg,
			nil,
		)
		c.Abort()
		return
	}
	// 对比ip白名单
	if len(productRow.WhitelistIP) > 0 {
		if !strings.Contains(productRow.WhitelistIP, c.ClientIP()) {
			mcommon.Log.Warnf("no in ip list of: %s %s", req.AppName, c.ClientIP())
			mcommon.GinDoRespErr(
				c,
				value.ErrorIPLimit,
				value.ErrorIPLimitMsg,
				nil,
			)
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
			mcommon.GinDoRespInternalErr(c)
			c.Abort()
			return
		}
		c.Set(gin.BodyBytesKey, body)
	}
	oldObj := gin.H{}
	err = json.Unmarshal(body, &oldObj)
	if err != nil {
		mcommon.Log.Warnf("req body error")
		mcommon.GinDoRespInternalErr(c)
		c.Abort()
		return
	}
	checkObj := gin.H{}
	for k, v := range oldObj {
		if k != "sign" {
			checkObj[k] = v
		}
	}
	checkSign := mcommon.WechatGetSign(productRow.AppSk, checkObj)
	if checkSign == "" || checkSign != req.Sign {
		mcommon.Log.Warnf("sign error of: %s", req.AppName)
		mcommon.GinDoRespErr(
			c,
			value.ErrorSignWrong,
			value.ErrorSignWrongMsg,
			nil,
		)
		c.Abort()
		return
	}
	// 检测nonce
	count, err := model.SQLCreateTProductNonce(
		c,
		xenv.DbCon,
		&model.DBTProductNonce{
			C:          req.Nonce,
			CreateTime: time.Now().Unix(),
		},
		true,
	)
	if err != nil {
		mcommon.Log.Warnf("err: [%T] %s", err, err.Error())
		mcommon.GinDoRespInternalErr(c)
		c.Abort()
		return
	}
	if count <= 0 {
		mcommon.Log.Warnf("nonce repeated")
		mcommon.GinDoRespErr(
			c,
			value.ErrorNonceRepeat,
			value.ErrorNonceRepeatMsg,
			nil,
		)
		c.Abort()
		return
	}
	c.Set("product_id", productRow.ID)
}
