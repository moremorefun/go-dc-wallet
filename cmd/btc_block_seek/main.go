// 检测区块到账
package main

import (
	"go-dc-wallet/app"
	"go-dc-wallet/hbtc"
)

func main() {
	app.EnvCreate()
	defer app.EnvDestroy()

	hbtc.CheckBlockSeek()
}
