package omniclient

import (
	"encoding/json"
	"fmt"
	"go-dc-wallet/hcommon"

	"github.com/parnurzeal/gorequest"
)

var client *gorequest.SuperAgent
var rpcURI string

type StRpcRespError struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

func (e *StRpcRespError) Error() string {
	return fmt.Sprintf("%d %s", e.Code, e.Message)
}

type StRpcReq struct {
	Jsonrpc string        `json:"jsonrpc"`
	ID      string        `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type StRpcResp struct {
	ID    string          `json:"id"`
	Error *StRpcRespError `json:"error"`
}

type StTxResult struct {
	Txid     string `json:"txid"`
	Hash     string `json:"hash"`
	Version  int64  `json:"version"`
	Size     int64  `json:"size"`
	Vsize    int64  `json:"vsize"`
	Weight   int64  `json:"weight"`
	Locktime int64  `json:"locktime"`
	Vin      []struct {
		Coinbase  string `json:"coinbase"`
		Txid      string `json:"txid"`
		Vout      int64  `json:"vout"`
		ScriptSig struct {
			Asm string `json:"asm"`
			Hex string `json:"hex"`
		} `json:"scriptSig"`
		Sequence int64 `json:"sequence"`
	} `json:"vin"`
	Vout []struct {
		Value        float64 `json:"value"`
		N            int64   `json:"n"`
		ScriptPubKey struct {
			Asm       string   `json:"asm"`
			Hex       string   `json:"hex"`
			ReqSigs   int64    `json:"reqSigs"`
			Type      string   `json:"type"`
			Addresses []string `json:"addresses"`
		} `json:"scriptPubKey,omitempty"`
	} `json:"vout"`
	Hex           string `json:"hex"`
	Blockhash     string `json:"blockhash"`
	Confirmations int64  `json:"confirmations"`
	Time          int64  `json:"time"`
	Blocktime     int64  `json:"blocktime"`
}

type StBlockResult struct {
	Hash              string        `json:"hash"`
	Confirmations     int64         `json:"confirmations"`
	Strippedsize      int64         `json:"strippedsize"`
	Size              int64         `json:"size"`
	Weight            int64         `json:"weight"`
	Height            int64         `json:"height"`
	Version           int64         `json:"version"`
	VersionHex        string        `json:"versionHex"`
	Merkleroot        string        `json:"merkleroot"`
	Tx                []*StTxResult `json:"tx"`
	Time              int64         `json:"time"`
	Mediantime        int64         `json:"mediantime"`
	Nonce             int64         `json:"nonce"`
	Bits              string        `json:"bits"`
	Difficulty        float64       `json:"difficulty"`
	Chainwork         string        `json:"chainwork"`
	NTx               int64         `json:"nTx"`
	Previousblockhash string        `json:"previousblockhash"`
	Nextblockhash     string        `json:"nextblockhash"`
}

// InitClient 初始化客户端
func InitClient(omniRPCHost, omniRPCUser, omniRPCPwd string) {
	rpcURI = omniRPCHost
	client = gorequest.New().SetBasicAuth(omniRPCUser, omniRPCPwd)
}

func doReq(method string, arqs []interface{}, resp interface{}) error {
	_, body, errs := client.Post(rpcURI).Send(StRpcReq{
		Jsonrpc: "1.0",
		ID:      hcommon.GetUUIDStr(),
		Method:  method,
		Params:  arqs,
	}).EndBytes()
	if errs != nil {
		return errs[0]
	}
	err := json.Unmarshal(body, resp)
	if err != nil {
		return err
	}
	return nil
}

// RpcGetBlockCount 获取block number
func RpcGetBlockCount() (int64, error) {
	resp := struct {
		StRpcResp
		Result int64 `json:"result"`
	}{}
	err := doReq(
		"getblockcount",
		nil,
		&resp,
	)
	if err != nil {
		return 0, err
	}
	if resp.Error != nil {
		return 0, resp.Error
	}
	return resp.Result, nil
}

// RpcGetBlockHash 获取block hash
func RpcGetBlockHash(blockHeight int64) (string, error) {
	resp := struct {
		StRpcResp
		Result string `json:"result"`
	}{}
	err := doReq(
		"getblockhash",
		[]interface{}{blockHeight},
		&resp,
	)
	if err != nil {
		return "", err
	}
	if resp.Error != nil {
		return "", resp.Error
	}
	return resp.Result, nil
}

// RpcGetBlockVerbose 获取block 内容
func RpcGetBlockVerbose(blockHash string) (*StBlockResult, error) {
	resp := struct {
		StRpcResp
		Result *StBlockResult `json:"result"`
	}{}
	err := doReq(
		"getblock",
		[]interface{}{blockHash, 2},
		&resp,
	)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}
	return resp.Result, nil
}

// RpcGetRawTransactionVerbose 获取tx
func RpcGetRawTransactionVerbose(txHash string) (*StTxResult, error) {
	resp := struct {
		StRpcResp
		Result *StTxResult `json:"result"`
	}{}
	err := doReq(
		"getrawtransaction",
		[]interface{}{txHash, 1},
		&resp,
	)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}
	return resp.Result, nil
}
