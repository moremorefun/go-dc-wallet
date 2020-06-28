package main

import (
	"go-dc-wallet/app"
	"go-dc-wallet/heos"
)

func main() {
	app.EnvCreate()
	defer app.EnvDestroy()

	heos.CheckAddressFree()
}
