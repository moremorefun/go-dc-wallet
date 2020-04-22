// 发送erc20冲币通知
package main

import (
	"go-dc-wallet/app"
	"go-dc-wallet/eth"
)

func main() {
	app.EnvCreate()
	defer app.EnvDestroy()

	eth.CheckErc20TxNotify()
}
