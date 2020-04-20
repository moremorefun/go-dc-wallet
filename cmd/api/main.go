// 对外服务接口
package main

import (
	"go-dc-wallet/app"
	"go-dc-wallet/hcommon"
	"go-dc-wallet/web"
	"time"

	"github.com/fvbock/endless"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
)

func main() {
	app.EnvCreate()
	defer app.EnvDestroy()
	// 初始化gin
	if !app.Cfg.IsDebug {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	if !app.Cfg.IsDebug {
		r.Use(gin.Logger(), gin.Recovery())

	} else {
		r.Use(ginzap.Ginzap(hcommon.ZapLog, time.StampMilli, true), gin.Recovery())
	}
	// 注册api
	web.Start(r)
	// 开始服务
	_ = endless.ListenAndServe("0.0.0.0:1000", r)
}
