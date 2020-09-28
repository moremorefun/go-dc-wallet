package main

import (
	"go-dc-wallet/heos"
	"go-dc-wallet/xenv"
)

func main() {
	xenv.EnvCreate()
	defer xenv.EnvDestroy()

	heos.CheckAddressFree()
}
