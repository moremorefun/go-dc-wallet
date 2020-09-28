package main

import (
	"go-dc-wallet/hbtc"
	"go-dc-wallet/xenv"
)

func main() {
	xenv.EnvCreate()
	defer xenv.EnvDestroy()
	hbtc.OmniCheckTxNotify()
}
