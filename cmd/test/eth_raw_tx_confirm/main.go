// 检测到账情况
package main

import (
	"go-dc-wallet/heth"
	"go-dc-wallet/xenv"
)

func main() {
	xenv.EnvCreate()
	defer xenv.EnvDestroy()

	heth.CheckRawTxConfirm()
}
