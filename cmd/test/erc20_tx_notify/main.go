// 发送erc20冲币通知
package main

import (
	"go-dc-wallet/heth"
	"go-dc-wallet/xenv"
)

func main() {
	xenv.EnvCreate()
	defer xenv.EnvDestroy()

	heth.CheckErc20TxNotify()
}
