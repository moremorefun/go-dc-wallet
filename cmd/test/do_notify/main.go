// 发送通知
package main

import (
	"go-dc-wallet/app"
	"go-dc-wallet/xenv"
)

func main() {
	xenv.EnvCreate()
	defer xenv.EnvDestroy()

	app.CheckDoNotify()
}
