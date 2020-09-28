// 检测区块到账
package main

import (
	"go-dc-wallet/hbtc"
	"go-dc-wallet/xenv"
)

func main() {
	xenv.EnvCreate()
	defer xenv.EnvDestroy()

	hbtc.CheckBlockSeek()
}
