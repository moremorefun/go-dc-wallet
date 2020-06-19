// 检测eth gas price
package main

import (
	"go-dc-wallet/app"
	"go-dc-wallet/heth"
)

func main() {
	app.EnvCreate()
	defer app.EnvDestroy()

	heth.CheckGasPrice()
}
