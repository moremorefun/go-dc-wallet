// 发送通知
package main

import (
	"go-dc-wallet/app"
)

func main() {
	app.EnvCreate()
	defer app.EnvDestroy()

	app.CheckDoNotify()
}
