// 检测入账广播
package main

import (
	"go-dc-wallet/app"
	"go-dc-wallet/eth"
)

func main() {
	app.EnvCreate()
	defer app.EnvDestroy()

	eth.CheckTxNotify()
}
