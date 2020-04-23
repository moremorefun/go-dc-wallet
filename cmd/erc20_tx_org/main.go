// erc20零钱整理
package main

import (
	"go-dc-wallet/app"
	"go-dc-wallet/eth"
)

func main() {
	app.EnvCreate()
	defer app.EnvDestroy()

	eth.CheckErc20TxOrg()
}
