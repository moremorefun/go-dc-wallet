// 发送erc20冲币通知
package main

import (
	"go-dc-wallet/app"
	"go-dc-wallet/heth"
)

func main() {
	app.EnvCreate()
	defer app.EnvDestroy()

	heth.CheckErc20TxNotify()
}
