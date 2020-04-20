package web

import "github.com/gin-gonic/gin"

func Start(r *gin.Engine) {
	r.POST("/api/address", postAddress)
	r.POST("/api/withdraw", postWithdraw)
}

func postAddress(c *gin.Context) {

}

func postWithdraw(c *gin.Context) {

}
