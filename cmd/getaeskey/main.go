package main

import (
	"flag"
	"fmt"
	"go-dc-wallet/app"
	"go-dc-wallet/hcommon"
	"strings"
)

func main() {
	// 读取运行参数
	var sourceKey = flag.String("k", "", "原始私钥")
	var h = flag.Bool("h", false, "help message")
	flag.Parse()
	if *h {
		flag.Usage()
		return
	}
	*sourceKey = strings.TrimSpace(*sourceKey)
	if *sourceKey == "" {
		flag.Usage()
		return
	}
	app.EnvCreate()
	defer app.EnvDestroy()

	// 加密密钥
	privateKeyStrEn := hcommon.AesEncrypt(*sourceKey, app.Cfg.AESKey)
	fmt.Printf("%s\n", privateKeyStrEn)
}
