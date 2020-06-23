package hbtc

import (
	"encoding/hex"
	"errors"
	"go-dc-wallet/app"
	"go-dc-wallet/app/model"
	"go-dc-wallet/hcommon"

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
	params *chaincfg.Params
}

var network = map[string]Network{
	"btc":      {params: &chaincfg.MainNetParams},
	"btc-test": {params: &chaincfg.TestNet3Params},
}

func GetNetwork(coinType string) Network {
	n, ok := network[coinType]
	if !ok {
		hcommon.Log.Errorf("no network: %s get btc for replace", coinType)
		return network["btc"]
	}
	return n
}

func (network Network) GetNetworkParams() *chaincfg.Params {
	return network.params
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
		return nil, errors.New("The WIF string is not valid for the `" + network.params.Name + "` network")
	}
	return wif, nil
}

func (network Network) GetAddress(wif *btcutil.WIF) (*btcutil.AddressPubKey, error) {
	return btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeCompressed(), network.GetNetworkParams())
}

func BtcAddTxOut(tx *wire.MsgTx, toAddress string, balance int64) error {
	addrTo, err := btcutil.DecodeAddress(
		toAddress,
		GetNetwork(app.Cfg.BtcNetworkType).params,
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
	}
	argVouts = append(argVouts, &StBtxTxOut{
		VoutAddress: firstAddress,
		Balance:     0,
	})
	return BtcTxSize(argVins, argVouts)
}
