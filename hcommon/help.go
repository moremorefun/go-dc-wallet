package hcommon

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/iancoleman/strcase"
	"gopkg.in/go-playground/validator.v8"
)

// IsStringInSlice 字符串是否在数组中
func IsStringInSlice(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

// GinFillBindError 检测gin输入绑定错误
func GinFillBindError(c *gin.Context, err error) {
	validatorError, ok := err.(validator.ValidationErrors)
	if ok {
		errMsgList := make([]string, 0, 16)
		for _, v := range validatorError {
			errMsgList = append(errMsgList, fmt.Sprintf("[%s] is %s", strcase.ToSnake(v.Field), v.ActualTag))
		}
		c.JSON(http.StatusOK, gin.H{"error": ErrorBind, "err_msg": strings.Join(errMsgList, ", ")})
		return
	}
	unmarshalError, ok := err.(*json.UnmarshalTypeError)
	if ok {
		c.JSON(http.StatusOK, gin.H{"error": ErrorBind, "err_msg": fmt.Sprintf("[%s] type error", unmarshalError.Field)})
		return
	}
	if err == io.EOF {
		c.JSON(http.StatusOK, gin.H{"error": ErrorBind, "err_msg": fmt.Sprintf("empty body")})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": ErrorInternal})
}

// GetSign 获取签名
func GetSign(appSecret string, paramsMap gin.H) string {
	var args []string
	var keys []string
	for k := range paramsMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := fmt.Sprintf("%s=%v", k, paramsMap[k])
		args = append(args, v)
	}
	baseString := strings.Join(args, "&")
	baseString += fmt.Sprintf("&key=%s", appSecret)
	data := []byte(baseString)
	r := md5.Sum(data)
	signedString := hex.EncodeToString(r[:])
	return strings.ToUpper(signedString)
}
