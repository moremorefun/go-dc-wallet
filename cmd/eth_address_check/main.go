// 检测eth剩余可用地址是否满足需求，
// 如果不足则创建地址
package main

import (
	"go-dc-wallet/app"
	"go-dc-wallet/heth"
)

func main() {
	app.EnvCreate()
	defer app.EnvDestroy()

	heth.CheckAddressFree()
}
