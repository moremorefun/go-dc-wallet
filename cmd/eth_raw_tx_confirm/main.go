// 检测到账情况
package main

import (
	"go-dc-wallet/app"
	"go-dc-wallet/heth"
)

func main() {
	app.EnvCreate()
	defer app.EnvDestroy()

	heth.CheckRawTxConfirm()
}
