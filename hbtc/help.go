package hbtc

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"go-dc-wallet/app"
	"go-dc-wallet/app/model"
	"go-dc-wallet/hcommon"
	"go-dc-wallet/xenv"

	"github.com/shopspring/decimal"

	"github.com/btcsuite/btcd/chaincfg/chainhash"

	"github.com/btcsuite/btcd/txscript"

	"github.com/btcsuite/btcd/wire"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
)

type StBtxTxIn struct {
	VinTxHash string
	VinTxN    int64
	VinScript string
	Balance   int64
	Wif       *btcutil.WIF
}

type StBtxTxOut struct {
	VoutAddress string
	Balance     int64
}

type Network struct {
	Params *chaincfg.Params
}

var network = map[string]Network{
	"btc":      {Params: &chaincfg.MainNetParams},
	"btc-test": {Params: &chaincfg.TestNet3Params},
}

func GetNetwork(coinType string) Network {
	n, ok := network[coinType]
	if !ok {
		hcommon.Log.Errorf("no network: %s get btc for replace", coinType)
		return network[CoinSymbol]
	}
	return n
}

func (network Network) GetNetworkParams() *chaincfg.Params {
	return network.Params
}

func (network Network) CreatePrivateKey() (*btcutil.WIF, error) {
	secret, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return nil, err
	}
	return btcutil.NewWIF(secret, network.GetNetworkParams(), true)
}

func (network Network) ImportWIF(wifStr string) (*btcutil.WIF, error) {
	wif, err := btcutil.DecodeWIF(wifStr)
	if err != nil {
		return nil, err
	}
	if !wif.IsForNet(network.GetNetworkParams()) {
		return nil, errors.New("The WIF string is not valid for the `" + network.Params.Name + "` network")
	}
	return wif, nil
}

func (network Network) GetAddress(wif *btcutil.WIF) (*btcutil.AddressPubKey, error) {
	return btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeCompressed(), network.GetNetworkParams())
}

func BtcAddTxOut(tx *wire.MsgTx, toAddress string, balance int64) error {
	addrTo, err := btcutil.DecodeAddress(
		toAddress,
		GetNetwork(xenv.Cfg.BtcNetworkType).Params,
	)
	if err != nil {
		return err
	}
	pkScriptf, err := txscript.PayToAddrScript(addrTo)
	if err != nil {
		return err
	}
	tx.AddTxOut(wire.NewTxOut(balance, pkScriptf))
	return nil
}

// BtcMakeTx 创建交易
func BtcMakeTx(vins []*StBtxTxIn, vouts []*StBtxTxOut, gasPrice int64, changeAddress string) (*wire.MsgTx, error) {
	inAmount := int64(0)
	outAmount := int64(0)
	tx := wire.NewMsgTx(wire.TxVersion)
	for _, vin := range vins {
		hash, err := chainhash.NewHashFromStr(vin.VinTxHash)
		if err != nil {
			return nil, err
		}
		outPoint := wire.NewOutPoint(hash, uint32(vin.VinTxN))
		txIn := wire.NewTxIn(outPoint, nil, nil)
		tx.AddTxIn(txIn)
		inAmount += vin.Balance
	}
	for _, vout := range vouts {
		err := BtcAddTxOut(tx, vout.VoutAddress, vout.Balance)
		if err != nil {
			return nil, err
		}
		outAmount += vout.Balance
	}
	// 计算手续费
	for i, vin := range vins {
		txinPkScript, err := hex.DecodeString(vin.VinScript)
		if err != nil {
			return nil, err
		}
		script, err := txscript.SignatureScript(
			tx,
			i,
			txinPkScript,
			txscript.SigHashAll,
			vin.Wif.PrivKey,
			true,
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return nil, err
		}
		tx.TxIn[i].SignatureScript = script
	}
	txSize := tx.SerializeSize()
	txFee := gasPrice * int64(txSize)
	change := inAmount - outAmount - txFee
	if change < 0 {
		// 数额不足
		return nil, errors.New("btc tx input amount not ok")
	}
	if change >= MinNondustOutput {
		// 创建找零
		err := BtcAddTxOut(tx, changeAddress, change)
		if err != nil {
			return nil, err
		}
	}
	if len(tx.TxOut) <= 0 {
		// 数额不足
		return nil, errors.New("btc tx input amount not ok")
	}
	if tx.SerializeSize() > MaxTxSize {
		// 长度过大
		return nil, errors.New("btc tx size too big")
	}
	// 重新签名
	for i, vin := range vins {
		txinPkScript, err := hex.DecodeString(vin.VinScript)
		if err != nil {
			return nil, err
		}
		script, err := txscript.SignatureScript(
			tx,
			i,
			txinPkScript,
			txscript.SigHashAll,
			vin.Wif.PrivKey,
			true,
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return nil, err
		}
		tx.TxIn[i].SignatureScript = script
		vm, err := txscript.NewEngine(
			txinPkScript,
			tx,
			i,
			txscript.StandardVerifyFlags,
			nil,
			nil,
			-1,
		)
		if err != nil {
			return nil, err
		}
		err = vm.Execute()
		if err != nil {
			return nil, err
		}
	}
	return tx, nil
}

// BtcTxSize 交易大小
func BtcTxSize(vins []*StBtxTxIn, vouts []*StBtxTxOut) (int64, error) {
	tx := wire.NewMsgTx(wire.TxVersion)
	for _, vin := range vins {
		hash, err := chainhash.NewHashFromStr(vin.VinTxHash)
		if err != nil {
			return 0, err
		}
		outPoint := wire.NewOutPoint(hash, uint32(vin.VinTxN))
		txIn := wire.NewTxIn(outPoint, nil, nil)
		tx.AddTxIn(txIn)
	}
	for _, vout := range vouts {
		err := BtcAddTxOut(tx, vout.VoutAddress, vout.Balance)
		if err != nil {
			return 0, err
		}
	}
	// 计算手续费
	for i, vin := range vins {
		txinPkScript, err := hex.DecodeString(vin.VinScript)
		if err != nil {
			return 0, err
		}
		script, err := txscript.SignatureScript(
			tx,
			i,
			txinPkScript,
			txscript.SigHashAll,
			vin.Wif.PrivKey,
			true,
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return 0, err
		}
		tx.TxIn[i].SignatureScript = script
	}
	txSize := int64(tx.SerializeSize())
	return txSize, nil
}

// BtcTxWithdrawSize 提币tx大小
func BtcTxWithdrawSize(vins []*model.DBTTxBtcUxto, vouts []*model.DBTWithdraw, keyMap map[string]*btcutil.WIF) (int64, error) {
	var argVins []*StBtxTxIn
	var argVouts []*StBtxTxOut
	firstAddress := ""
	for _, vin := range vins {
		wif, ok := keyMap[vin.VoutAddress]
		if !ok {
			return 0, errors.New("no key of wif")
		}
		balance, err := decimal.NewFromString(vin.VoutValue)
		if err != nil {
			return 0, err
		}
		argVins = append(argVins, &StBtxTxIn{
			VinTxHash: vin.TxID,
			VinTxN:    vin.VoutN,
			VinScript: vin.VoutScript,
			Balance:   balance.Mul(decimal.NewFromInt(1e8)).IntPart(),
			Wif:       wif,
		})
		if firstAddress == "" {
			firstAddress = vin.VoutAddress
		}
	}
	for _, vout := range vouts {
		balance, err := decimal.NewFromString(vout.BalanceReal)
		if err != nil {
			return 0, err
		}
		argVouts = append(argVouts, &StBtxTxOut{
			VoutAddress: vout.ToAddress,
			Balance:     balance.Mul(decimal.NewFromInt(1e8)).IntPart(),
		})
		if firstAddress == "" {
			firstAddress = vout.ToAddress
		}
	}
	argVouts = append(argVouts, &StBtxTxOut{
		VoutAddress: firstAddress,
		Balance:     0,
	})
	return BtcTxSize(argVins, argVouts)
}

// OmniTxSize 计算omni大小
func OmniTxSize(senderUxtoRow *model.DBTTxBtcUxto, toAddress string, tokenIndex int64, tokenBalance int64, keyMap map[string]*btcutil.WIF, inUxtoRows []*model.DBTTxBtcUxto, vouts []*StBtxTxOut) (int64, error) {
	tx := wire.NewMsgTx(wire.TxVersion)
	// 添加sender
	hash, err := chainhash.NewHashFromStr(senderUxtoRow.TxID)
	if err != nil {
		return 0, err
	}
	outPoint := wire.NewOutPoint(hash, uint32(senderUxtoRow.VoutN))
	txIn := wire.NewTxIn(outPoint, nil, nil)
	tx.AddTxIn(txIn)
	// 添加input
	for _, inUxtoRow := range inUxtoRows {
		hash, err := chainhash.NewHashFromStr(inUxtoRow.TxID)
		if err != nil {
			return 0, err
		}
		outPoint := wire.NewOutPoint(hash, uint32(inUxtoRow.VoutN))
		txIn := wire.NewTxIn(outPoint, nil, nil)
		tx.AddTxIn(txIn)
	}
	// 添加script
	sHex := fmt.Sprintf("%016x%016x", tokenIndex, tokenBalance)
	b, err := hex.DecodeString(omniHex + sHex)
	if err != nil {
		return 0, err
	}
	opreturnScript, err := txscript.NullDataScript(b)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return 0, err
	}
	tx.AddTxOut(wire.NewTxOut(0, opreturnScript))
	// 添加output
	for _, vout := range vouts {
		err := BtcAddTxOut(tx, vout.VoutAddress, vout.Balance)
		if err != nil {
			return 0, err
		}
	}
	// 添加 omni get
	err = BtcAddTxOut(tx, toAddress, MinNondustOutput)
	if err != nil {
		return 0, err
	}
	// 计算手续费
	for i := range tx.TxIn {
		scriptHex := ""
		address := ""
		if i == 0 {
			scriptHex = senderUxtoRow.VoutScript
			address = senderUxtoRow.VoutAddress
		} else {
			scriptHex = inUxtoRows[i-1].VoutScript
			address = inUxtoRows[i-1].VoutAddress
		}
		txinPkScript, err := hex.DecodeString(scriptHex)
		if err != nil {
			return 0, err
		}
		wif, ok := keyMap[address]
		if !ok {
			return 0, errors.New("no wif")
		}
		script, err := txscript.SignatureScript(
			tx,
			i,
			txinPkScript,
			txscript.SigHashAll,
			wif.PrivKey,
			true,
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return 0, err
		}
		tx.TxIn[i].SignatureScript = script
	}
	txSize := int64(tx.SerializeSize())
	return txSize, nil
}

// OmniTxMake 生成交易
func OmniTxMake(senderUxtoRow *model.DBTTxBtcUxto, toAddress string, changeAddress string, tokenIndex int64, tokenBalance int64, gasPrice int64, keyMap map[string]*btcutil.WIF, inUxtoRows []*model.DBTTxBtcUxto) (*wire.MsgTx, error) {
	inBalance := int64(0)
	outBalance := int64(0)

	// --- 生成基础交易 ---
	tx := wire.NewMsgTx(wire.TxVersion)
	// 添加sender
	hash, err := chainhash.NewHashFromStr(senderUxtoRow.TxID)
	if err != nil {
		return nil, err
	}
	outPoint := wire.NewOutPoint(hash, uint32(senderUxtoRow.VoutN))
	txIn := wire.NewTxIn(outPoint, nil, nil)
	tx.AddTxIn(txIn)
	balance, err := decimal.NewFromString(senderUxtoRow.VoutValue)
	if err != nil {
		return nil, err
	}
	inBalance += balance.Mul(decimal.NewFromInt(1e8)).IntPart()
	// 添加input
	for _, inUxtoRow := range inUxtoRows {
		hash, err := chainhash.NewHashFromStr(inUxtoRow.TxID)
		if err != nil {
			return nil, err
		}
		outPoint := wire.NewOutPoint(hash, uint32(inUxtoRow.VoutN))
		txIn := wire.NewTxIn(outPoint, nil, nil)
		tx.AddTxIn(txIn)
		balance, err := decimal.NewFromString(inUxtoRow.VoutValue)
		if err != nil {
			return nil, err
		}
		inBalance += balance.Mul(decimal.NewFromInt(1e8)).IntPart()
	}
	// 添加 out script
	sHex := fmt.Sprintf("%016x%016x", tokenIndex, tokenBalance)
	b, err := hex.DecodeString(omniHex + sHex)
	if err != nil {
		return nil, err
	}
	opreturnScript, err := txscript.NullDataScript(b)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return nil, err
	}
	tx.AddTxOut(wire.NewTxOut(0, opreturnScript))
	// 添加 out omni address
	err = BtcAddTxOut(tx, toAddress, MinNondustOutput)
	if err != nil {
		return nil, err
	}
	outBalance += MinNondustOutput
	// --- 计算拆分热钱包交易 ---
	tmpBalance := int64(50000)
	addOutPutCount := (inBalance - outBalance) / tmpBalance
	// 获取预估大小
	txSize, err := GetEstimateTxSize(int64(len(tx.TxIn)), 1+addOutPutCount+1, true)
	fee := txSize * gasPrice
	addOutPutCount = (inBalance - outBalance - fee) / tmpBalance
	// --- 添加拆分输出 ---
	for i := int64(0); i < addOutPutCount; i++ {
		err = BtcAddTxOut(tx, changeAddress, tmpBalance)
		if err != nil {
			return nil, err
		}
	}
	// --- 添加找零 ---
	err = BtcAddTxOut(tx, changeAddress, 0)
	if err != nil {
		return nil, err
	}
	// --- 重新计算找零 ---
	for i := range tx.TxIn {
		scriptHex := ""
		address := ""
		if i == 0 {
			scriptHex = senderUxtoRow.VoutScript
			address = senderUxtoRow.VoutAddress
		} else {
			scriptHex = inUxtoRows[i-1].VoutScript
			address = inUxtoRows[i-1].VoutAddress
		}
		txinPkScript, err := hex.DecodeString(scriptHex)
		if err != nil {
			return nil, err
		}
		wif, ok := keyMap[address]
		if !ok {
			return nil, errors.New("no wif")
		}
		script, err := txscript.SignatureScript(
			tx,
			i,
			txinPkScript,
			txscript.SigHashAll,
			wif.PrivKey,
			true,
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return nil, err
		}
		tx.TxIn[i].SignatureScript = script
	}
	txSize = int64(tx.SerializeSize())
	leaveInBalance := inBalance - txSize*gasPrice - MinNondustOutput - addOutPutCount*tmpBalance
	if leaveInBalance < 0 {
		return nil, errors.New("error input")
	}
	if leaveInBalance < MinNondustOutput {
		tx.TxOut = tx.TxOut[:len(tx.TxOut)-1]
	} else {
		tx.TxOut[len(tx.TxOut)-1].Value = leaveInBalance
	}
	// --- 重新签名 ---
	var newOut []*wire.TxOut
	for i, out := range tx.TxOut {
		if i != 1 {
			newOut = append(newOut, out)
		}
	}
	newOut = append(newOut, tx.TxOut[1])
	tx.TxOut = newOut
	for i := range tx.TxIn {
		scriptHex := ""
		address := ""
		if i == 0 {
			scriptHex = senderUxtoRow.VoutScript
			address = senderUxtoRow.VoutAddress
		} else {
			scriptHex = inUxtoRows[i-1].VoutScript
			address = inUxtoRows[i-1].VoutAddress
		}
		txinPkScript, err := hex.DecodeString(scriptHex)
		if err != nil {
			return nil, err
		}
		wif, ok := keyMap[address]
		if !ok {
			return nil, errors.New("no wif")
		}
		script, err := txscript.SignatureScript(
			tx,
			i,
			txinPkScript,
			txscript.SigHashAll,
			wif.PrivKey,
			true,
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return nil, err
		}
		tx.TxIn[i].SignatureScript = script
		vm, err := txscript.NewEngine(
			txinPkScript,
			tx,
			i,
			txscript.StandardVerifyFlags,
			nil,
			nil,
			-1,
		)
		if err != nil {
			return nil, err
		}
		err = vm.Execute()
		if err != nil {
			return nil, err
		}
	}
	return tx, nil
}

// GetEstimateTxSize 获取tx大小
func GetEstimateTxSize(fromAddressCount int64, toAddressCount int64, isOmniScript bool) (int64, error) {
	tmpTxID := "e326842c86612d9e3849825117839b40444e7e1066136afcc5e6b7757f9508e0"
	tmpAddress := "mnzBq3LMyq71maLMqzSKaPMSm2WuiJtZpQ"
	tmpWif := "cRdLxqbRbPFQB12XDKQNkD8fDsBCAqEbNJVn4Z6fT8qMFSih3AFm"
	wif, err := btcutil.DecodeWIF(tmpWif)
	if err != nil {
		return 0, err
	}
	addrTo, err := btcutil.DecodeAddress(
		tmpAddress,
		&chaincfg.TestNet3Params,
	)
	if err != nil {
		return 0, err
	}
	pkScriptf, err := txscript.PayToAddrScript(addrTo)
	if err != nil {
		return 0, err
	}
	tx := wire.NewMsgTx(wire.TxVersion)
	for i := int64(0); i < toAddressCount; i++ {
		tx.AddTxOut(wire.NewTxOut(0, pkScriptf))
	}
	if isOmniScript {
		// 添加script
		sHex := fmt.Sprintf("%016x%016x", 1, 0)
		b, err := hex.DecodeString(omniHex + sHex)
		if err != nil {
			return 0, err
		}
		opreturnScript, err := txscript.NullDataScript(b)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return 0, err
		}
		tx.AddTxOut(wire.NewTxOut(0, opreturnScript))
	}
	for i := int64(0); i < fromAddressCount; i++ {
		hash, err := chainhash.NewHashFromStr(tmpTxID)
		if err != nil {
			return 0, err
		}
		outPoint := wire.NewOutPoint(hash, 0)
		txIn := wire.NewTxIn(outPoint, nil, nil)
		tx.AddTxIn(txIn)
	}
	for i := range tx.TxIn {
		script, err := txscript.SignatureScript(
			tx,
			i,
			pkScriptf,
			txscript.SigHashAll,
			wif.PrivKey,
			true,
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return 0, err
		}
		tx.TxIn[i].SignatureScript = script
	}
	txSize := int64(tx.SerializeSize()) + fromAddressCount*2
	return txSize, nil
}

// RealStrToBalanceInt64 转换金额 real to balance
func RealStrToBalanceInt64(balanceRealStr string) (int64, error) {
	balanceReal, err := decimal.NewFromString(balanceRealStr)
	if err != nil {
		return 0, err
	}
	balance := balanceReal.Mul(decimal.NewFromInt(1e8))
	return balance.IntPart(), nil
}

// GetWifMapByAddresses 获取私钥
func GetWifMapByAddresses(ctx context.Context, db hcommon.DbExeAble, addresses []string) (map[string]*btcutil.WIF, error) {
	addressKeyMap, err := app.SQLGetAddressKeyMap(
		ctx,
		db,
		[]string{
			model.DBColTAddressKeyID,
			model.DBColTAddressKeyAddress,
			model.DBColTAddressKeyPwd,
		},
		addresses,
	)
	if err != nil {
		return nil, err
	}
	addressWifMap := make(map[string]*btcutil.WIF)
	for k, addressKey := range addressKeyMap {
		key := hcommon.AesDecrypt(addressKey.Pwd, xenv.Cfg.AESKey)
		if len(key) == 0 {
			return nil, err
		}
		wif, err := btcutil.DecodeWIF(key)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return nil, err
		}
		addressWifMap[k] = wif
	}
	return addressWifMap, nil
}
