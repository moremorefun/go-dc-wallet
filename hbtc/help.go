package hbtc

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"go-dc-wallet/app"
	"go-dc-wallet/model"
	"go-dc-wallet/omniclient"
	"go-dc-wallet/xenv"
	"math"

	"github.com/moremorefun/mcommon"
	"github.com/shopspring/decimal"

	"github.com/btcsuite/btcd/chaincfg/chainhash"

	"github.com/btcsuite/btcd/txscript"

	"github.com/btcsuite/btcd/wire"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
)

// StBtxTxIn 输入信息
type StBtxTxIn struct {
	VinTxHash string
	VinTxN    int64
	VinScript string
	Balance   int64
	Wif       *btcutil.WIF
}

// StBtxTxOut 输出信息
type StBtxTxOut struct {
	VoutAddress string
	Balance     int64
}

// Network 类型
type Network struct {
	Params *chaincfg.Params
}

var network = map[string]Network{
	"btc":      {Params: &chaincfg.MainNetParams},
	"btc-test": {Params: &chaincfg.TestNet3Params},
}

// GetNetwork 获取对象
func GetNetwork(coinType string) Network {
	n, ok := network[coinType]
	if !ok {
		mcommon.Log.Errorf("no network: %s get btc for replace", coinType)
		return network[CoinSymbol]
	}
	return n
}

// GetNetworkParams 获取网络参数
func (network Network) GetNetworkParams() *chaincfg.Params {
	return network.Params
}

// CreatePrivateKey 创建私钥
func (network Network) CreatePrivateKey() (*btcutil.WIF, error) {
	secret, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return nil, err
	}
	return btcutil.NewWIF(secret, network.GetNetworkParams(), true)
}

// ImportWIF 导入私钥
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

// GetAddress 获取地址
func (network Network) GetAddress(wif *btcutil.WIF) (*btcutil.AddressPubKey, error) {
	return btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeCompressed(), network.GetNetworkParams())
}

// GetAddressSegwitNested 获取隔离见证地址
func (network Network) GetAddressSegwitNested(wif *btcutil.WIF) (*btcutil.AddressScriptHash, error) {
	witnessProg := btcutil.Hash160(wif.PrivKey.PubKey().SerializeCompressed())
	addressWitnessPubKeyHash, err := btcutil.NewAddressWitnessPubKeyHash(witnessProg, network.GetNetworkParams())
	if err != nil {
		return nil, err
	}
	serializedScript, err := txscript.PayToAddrScript(addressWitnessPubKeyHash)
	if err != nil {
		return nil, err
	}
	addressScriptHash, err := btcutil.NewAddressScriptHash(serializedScript, network.GetNetworkParams())
	if err != nil {
		return nil, err
	}
	return addressScriptHash, nil
}

// BtcAddTxOut 添加一个输出
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
func BtcMakeTx(chainParams *chaincfg.Params, vins []*StBtxTxIn, vouts []*StBtxTxOut, gasPrice int64, changeAddress string) (*wire.MsgTx, error) {
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
	// 添加预找零信息
	err := BtcAddTxOut(tx, changeAddress, 0)
	if err != nil {
		return nil, err
	}
	// 计算手续费
	err = SigVins(chainParams, tx, vins)
	if err != nil {
		return nil, err
	}
	txSize := GetTxVsize(tx)
	txFee := gasPrice * txSize
	change := inAmount - outAmount - txFee
	if change < 0 {
		// 数额不足
		return nil, errors.New("btc tx input amount not ok")
	}
	if change >= MinNondustOutput {
		// 设置预找零数额
		tx.TxOut[len(tx.TxOut)-1].Value = change
	} else {
		// 删除预找零
		tx.TxOut = tx.TxOut[:len(tx.TxOut)-1]
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
	err = SigVins(chainParams, tx, vins)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// BtcTxSize 交易大小
func BtcTxSize(chainParams *chaincfg.Params, vins []*StBtxTxIn, vouts []*StBtxTxOut) (int64, error) {
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
	err := SigVins(chainParams, tx, vins)
	if err != nil {
		return 0, err
	}
	txSize := GetTxVsize(tx)
	return txSize, nil
}

// BtcTxWithdrawSize 提币tx大小
func BtcTxWithdrawSize(chainParams *chaincfg.Params, vins []*model.DBTTxBtcUxto, vouts []*model.DBTWithdraw, keyMap map[string]*btcutil.WIF) (int64, error) {
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
	return BtcTxSize(chainParams, argVins, argVouts)
}

// OmniTxMake 生成交易
func OmniTxMake(chainParams *chaincfg.Params, senderUxtoRow *model.DBTTxBtcUxto, toAddress string, changeAddress string, tokenIndex int64, tokenBalance int64, gasPrice int64, keyMap map[string]*btcutil.WIF, inUxtoRows []*model.DBTTxBtcUxto) (*wire.MsgTx, error) {
	inBalance := int64(0)
	// 输入数据
	var vins []*StBtxTxIn

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
	// 设置输入
	balance, err := decimal.NewFromString(senderUxtoRow.VoutValue)
	if err != nil {
		return nil, err
	}
	inBalance += balance.Mul(decimal.NewFromInt(1e8)).IntPart()
	wif, ok := keyMap[senderUxtoRow.VoutAddress]
	if !ok {
		return nil, errors.New("no wif")
	}
	vins = append(
		vins,
		&StBtxTxIn{
			VinTxHash: senderUxtoRow.TxID,
			VinTxN:    senderUxtoRow.VoutN,
			VinScript: senderUxtoRow.VoutScript,
			Balance:   balance.Mul(decimal.NewFromInt(1e8)).IntPart(),
			Wif:       wif,
		},
	)
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
		// 加入输入列表
		wif, ok := keyMap[inUxtoRow.VoutAddress]
		if !ok {
			return nil, errors.New("no wif")
		}
		vins = append(
			vins,
			&StBtxTxIn{
				VinTxHash: inUxtoRow.TxID,
				VinTxN:    inUxtoRow.VoutN,
				VinScript: inUxtoRow.VoutScript,
				Balance:   balance.Mul(decimal.NewFromInt(1e8)).IntPart(),
				Wif:       wif,
			},
		)
	}
	// 添加 out script
	sHex := fmt.Sprintf("%016x%016x", tokenIndex, tokenBalance)
	b, err := hex.DecodeString(omniHex + sHex)
	if err != nil {
		return nil, err
	}
	opreturnScript, err := txscript.NullDataScript(b)
	if err != nil {
		mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return nil, err
	}
	tx.AddTxOut(wire.NewTxOut(0, opreturnScript))
	// --- 计算拆分热钱包交易 ---
	tmpBalance := int64(200000)
	addOutPutCount := (inBalance - MinNondustOutput) / tmpBalance
	// 获取预估大小
	txSize, err := GetEstimateTxSize(int64(len(tx.TxIn)), addOutPutCount+1+1, true)
	fee := txSize * gasPrice
	addOutPutCount = (inBalance - MinNondustOutput - fee) / tmpBalance
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
	// 添加 out omni address
	err = BtcAddTxOut(tx, toAddress, MinNondustOutput)
	if err != nil {
		return nil, err
	}
	// --- 重新计算找零 ---
	err = SigVins(
		chainParams,
		tx,
		vins,
	)
	if err != nil {
		return nil, err
	}
	txSize = GetTxVsize(tx)
	leaveInBalance := inBalance - addOutPutCount*tmpBalance - MinNondustOutput - txSize*gasPrice
	if leaveInBalance < 0 {
		return nil, errors.New("error input")
	}
	changeIndex := len(tx.TxOut) - 2
	if leaveInBalance < MinNondustOutput {
		// 删除找零 out len-2
		tx.TxOut = append(tx.TxOut[:changeIndex], tx.TxOut[changeIndex+1:]...)
	} else {
		// 重置找零金额
		tx.TxOut[changeIndex].Value = leaveInBalance
	}
	// --- 重新签名 ---
	err = SigVins(
		chainParams,
		tx,
		vins,
	)
	if err != nil {
		return nil, err
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
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
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
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
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
func GetWifMapByAddresses(ctx context.Context, db mcommon.DbExeAble, addresses []string) (map[string]*btcutil.WIF, error) {
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
		key, err := mcommon.AesDecrypt(addressKey.Pwd, xenv.Cfg.AESKey)
		if err != nil {
			return nil, err
		}
		if len(key) == 0 {
			return nil, err
		}
		wif, err := btcutil.DecodeWIF(key)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return nil, err
		}
		addressWifMap[k] = wif
	}
	return addressWifMap, nil
}

// SigVins 对vin进行签名
func SigVins(chainParams *chaincfg.Params, tx *wire.MsgTx, vins []*StBtxTxIn) error {
	txSigHash := txscript.NewTxSigHashes(tx)
	for i, vin := range vins {
		// 重置sig
		tx.TxIn[i].SignatureScript = nil

		var setSignatureScript []byte
		// 解析vin的script字符串
		txInPkScript, err := hex.DecodeString(vin.VinScript)
		if err != nil {
			return err
		}
		// 获取vin的script的类型
		scriptClass, _, _, err := txscript.ExtractPkScriptAddrs(txInPkScript, chainParams)
		if err != nil {
			return err
		}
		// 设置的vin的签名字段
		switch scriptClass {
		case txscript.PubKeyHashTy:
			// vin为转账到地址
			script, err := txscript.SignatureScript(
				tx,
				i,
				txInPkScript,
				txscript.SigHashAll,
				vin.Wif.PrivKey,
				true,
			)
			if err != nil {
				return err
			}
			tx.TxIn[i].SignatureScript = script
		case txscript.ScriptHashTy:
			// vin为转账到script
			witnessProg := btcutil.Hash160(vin.Wif.PrivKey.PubKey().SerializeCompressed())
			addressWitnessPubKeyHash, err := btcutil.NewAddressWitnessPubKeyHash(witnessProg, chainParams)
			if err != nil {
				return err
			}
			// 对签名使用
			inAddressPKSHForSig, err := txscript.NewScriptBuilder().AddOp(txscript.OP_0).AddData(addressWitnessPubKeyHash.ScriptAddress()).Script()
			if err != nil {
				return err
			}
			// 对设置input使用
			inAddressPKSHForPK, err := txscript.NewScriptBuilder().AddOp(txscript.OP_DATA_22).AddOp(txscript.OP_0).AddData(addressWitnessPubKeyHash.ScriptAddress()).Script()
			if err != nil {
				return err
			}
			// 计算见证
			w, err := txscript.WitnessSignature(
				tx,
				txSigHash,
				i,
				vin.Balance,
				inAddressPKSHForSig,
				txscript.SigHashAll,
				vin.Wif.PrivKey,
				true,
			)
			if err != nil {
				return err
			}
			// 设置见证
			tx.TxIn[i].Witness = w
			// 设置数据, 需要验证完交易后再赋值
			txInPkScript = inAddressPKSHForSig
			setSignatureScript = inAddressPKSHForPK
		case txscript.WitnessV0PubKeyHashTy:
			// 计算见证
			w, err := txscript.WitnessSignature(
				tx,
				txSigHash,
				i,
				vin.Balance,
				txInPkScript,
				txscript.SigHashAll,
				vin.Wif.PrivKey,
				true,
			)
			if err != nil {
				return err
			}
			// 设置见证
			tx.TxIn[i].Witness = w
		default:
			return fmt.Errorf("error script type: %s", scriptClass.String())
		}
		vm, err := txscript.NewEngine(
			txInPkScript,
			tx,
			i,
			txscript.StandardVerifyFlags,
			nil,
			txSigHash,
			vin.Balance,
		)
		if err != nil {
			return err
		}
		err = vm.Execute()
		if err != nil {
			return err
		}
		if len(setSignatureScript) > 0 {
			tx.TxIn[i].SignatureScript = setSignatureScript
		}
	}
	return nil
}

// GetTxVsize 获取tx vsize
func GetTxVsize(tx *wire.MsgTx) int64 {
	s := math.Ceil(float64(1.0*3*tx.SerializeSizeStripped()+tx.SerializeSize()) / 4)
	return int64(s)
}

// GetAddressesOfVin 获取vin的地址
func GetAddressesOfVin(chainParams *chaincfg.Params, vin omniclient.StTxResultVin) ([]string, error) {
	var adds []btcutil.Address
	if len(vin.ScriptSig.Hex) > 0 {
		sigScript, err := hex.DecodeString(vin.ScriptSig.Hex)
		if err != nil {
			return nil, err
		}
		pk, err := txscript.ComputePkScript(sigScript, nil)
		if err != nil {
			return nil, err
		}
		_, adds, _, err = txscript.ExtractPkScriptAddrs(pk.Script(), chainParams)
		if err != nil {
			return nil, err
		}

	} else if len(vin.TxinWitness) > 0 {
		if len(vin.TxinWitness) != 2 {
			return nil, fmt.Errorf("error witness len")
		}
		w1, err := hex.DecodeString(vin.TxinWitness[0])
		if err != nil {
			return nil, err
		}
		w2, err := hex.DecodeString(vin.TxinWitness[1])
		if err != nil {
			return nil, err
		}
		tw := wire.TxWitness{
			w1,
			w2,
		}
		sig, err := txscript.ComputePkScript(nil, tw)
		if err != nil {
			return nil, err
		}
		_, adds, _, err = txscript.ExtractPkScriptAddrs(sig.Script(), chainParams)
		if err != nil {
			return nil, err
		}
	} else if len(vin.Coinbase) > 0 {

	} else {
		return nil, fmt.Errorf("error vin info")
	}
	var strs []string
	for _, address := range adds {
		strs = append(strs, address.String())
	}
	return strs, nil
}

// GetAddressesOfVinMsg 获取vin的地址
func GetAddressesOfVinMsg(chainParams *chaincfg.Params, txIn *wire.TxIn) ([]string, error) {
	var adds []btcutil.Address
	if len(txIn.SignatureScript) > 0 {
		pk, err := txscript.ComputePkScript(txIn.SignatureScript, nil)
		if err != nil {
			return nil, err
		}
		_, adds, _, err = txscript.ExtractPkScriptAddrs(pk.Script(), chainParams)
		if err != nil {
			return nil, err
		}

	} else if len(txIn.Witness[0]) > 0 {
		sig, err := txscript.ComputePkScript(nil, txIn.Witness)
		if err != nil {
			return nil, err
		}
		_, adds, _, err = txscript.ExtractPkScriptAddrs(sig.Script(), chainParams)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("error vin info")
	}
	var strs []string
	for _, address := range adds {
		strs = append(strs, address.String())
	}
	return strs, nil
}
