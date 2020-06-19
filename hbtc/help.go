package hbtc

import (
	"errors"
	"go-dc-wallet/hcommon"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
)

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
