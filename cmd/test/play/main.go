package main

import (
	"bytes"
	"encoding/hex"
	"go-dc-wallet/hbtc"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"

	"github.com/btcsuite/btcd/txscript"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/moremorefun/mcommon"
)

func main() {
	// mz2vpprjaUUNDfjrbQdJcBzdxcu8RwwsQ9 cRNRzDgCRNsaSPK1rmRVLPUKVdqqQcBmKG6BMtd3ZLXc2AoVMebM
	// 08799fa6294190a28c4f7390987f69fae3cc99e784eda95ac1a3d0f4efc317b8:0 0.00030000 76a914cb1d842e00a716935a1d24b79fc22fa1fa9070d488ac

	// 2N4g57dpDTVAUrZUFv43y768oX8e14L6964 cV7e525yTXCdasTuJhPfmRY3E9S9WtnnkpsmPW3DrJiyKrS2zD2R
	// 9cc6aaa6c58cf1f36a1ac97a04e8349edc21f5146924e9774230993cb7b12f05:0 0.00012000 a9147d5c602ab21e160c40da35841496012bc1909fe787

	// tb1q9l3cn0phxg0hdfsr25spr72xvfamfsgyhpm32e cRk6nhaWswnRbW8XJdmW3CFCQjspcwaHaEfLWMe5d2myh6h8ndx8

	wif1, err := btcutil.DecodeWIF("cRNRzDgCRNsaSPK1rmRVLPUKVdqqQcBmKG6BMtd3ZLXc2AoVMebM")
	if err != nil {
		mcommon.Log.Fatalf("err: [%T] %s", err, err.Error())
	}
	wif2, err := btcutil.DecodeWIF("cV7e525yTXCdasTuJhPfmRY3E9S9WtnnkpsmPW3DrJiyKrS2zD2R")
	if err != nil {
		mcommon.Log.Fatalf("err: [%T] %s", err, err.Error())
	}
	wif3, err := btcutil.DecodeWIF("cRk6nhaWswnRbW8XJdmW3CFCQjspcwaHaEfLWMe5d2myh6h8ndx8")
	if err != nil {
		mcommon.Log.Fatalf("err: [%T] %s", err, err.Error())
	}

	wifMap := map[string]*btcutil.WIF{
		"mz2vpprjaUUNDfjrbQdJcBzdxcu8RwwsQ9":         wif1,
		"2N4g57dpDTVAUrZUFv43y768oX8e14L6964":        wif2,
		"tb1q9l3cn0phxg0hdfsr25spr72xvfamfsgyhpm32e": wif3,
	}
	_ = wifMap

	vins := []*hbtc.StBtxTxIn{
		{
			VinTxHash: "ff94c4764964c6613627e2579f4b8f22d1c248f797c4a9ef1552168e69cff964",
			VinTxN:    1,
			VinScript: "00142fe389bc37321f76a603552011f946627bb4c104",
			Balance:   1024663,
			Wif:       wif3,
		},
	}
	// 创建tx
	// input
	inAmount := int64(0)
	tx := wire.NewMsgTx(wire.TxVersion)
	for _, vin := range vins {
		hash, err := chainhash.NewHashFromStr(vin.VinTxHash)
		if err != nil {
			mcommon.Log.Fatalf("err: [%T] %s", err, err.Error())
		}
		outPoint := wire.NewOutPoint(hash, uint32(vin.VinTxN))
		txIn := wire.NewTxIn(outPoint, nil, nil)
		tx.AddTxIn(txIn)
		inAmount += vin.Balance
	}
	// output
	addrTo, err := btcutil.DecodeAddress(
		"2N4g57dpDTVAUrZUFv43y768oX8e14L6964",
		&chaincfg.TestNet3Params,
	)
	if err != nil {
		mcommon.Log.Fatalf("err: [%T] %s", err, err.Error())
	}
	pkScriptf, err := txscript.PayToAddrScript(addrTo)
	if err != nil {
		mcommon.Log.Fatalf("err: [%T] %s", err, err.Error())
	}
	tx.AddTxOut(wire.NewTxOut(inAmount-500, pkScriptf))
	// 签名
	err = hbtc.SigVins(&chaincfg.TestNet3Params, tx, vins)
	if err != nil {
		mcommon.Log.Fatalf("err: [%T] %s", err, err.Error())
	}
	// 序列化
	txSize := tx.SerializeSize()
	b := new(bytes.Buffer)
	b.Grow(txSize)
	err = tx.Serialize(b)
	if err != nil {
		mcommon.Log.Fatalf("err: [%T] %s", err, err.Error())
	}
	mcommon.Log.Debugf("raw tx: %s", hex.EncodeToString(b.Bytes()))
}
