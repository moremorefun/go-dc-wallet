package eosclient

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/parnurzeal/gorequest"
)

var client *gorequest.SuperAgent
var rpcURI string

type StRpcRespError struct {
	Code     int64  `json:"code"`
	Message  string `json:"message"`
	ErrorInv struct {
		Code    int64  `json:"code"`
		Name    string `json:"name"`
		What    string `json:"what"`
		Details []struct {
			Message    string `json:"message"`
			File       string `json:"file"`
			LineNumber int64  `json:"line_number"`
			Method     string `json:"method"`
		} `json:"details"`
	} `json:"error"`
}

func (e *StRpcRespError) Error() string {
	if len(e.ErrorInv.Details) > 0 {
		return fmt.Sprintf("%d[%d] %s-%s-%s", e.Code, e.ErrorInv.Code, e.Message, e.ErrorInv.What, e.ErrorInv.Details[0].Message)
	}
	return fmt.Sprintf("%d[%d] %s-%s", e.Code, e.ErrorInv.Code, e.Message, e.ErrorInv.What)
}

type StChainGetInfo struct {
	ServerVersion            string `json:"server_version"`
	ChainID                  string `json:"chain_id"`
	HeadBlockNum             int64  `json:"head_block_num"`
	LastIrreversibleBlockNum int64  `json:"last_irreversible_block_num"`
	LastIrreversibleBlockID  string `json:"last_irreversible_block_id"`
	HeadBlockID              string `json:"head_block_id"`
	HeadBlockTime            string `json:"head_block_time"`
	HeadBlockProducer        string `json:"head_block_producer"`
	VirtualBlockCPULimit     int64  `json:"virtual_block_cpu_limit"`
	VirtualBlockNetLimit     int64  `json:"virtual_block_net_limit"`
	BlockCPULimit            int64  `json:"block_cpu_limit"`
	BlockNetLimit            int64  `json:"block_net_limit"`
	ServerVersionString      string `json:"server_version_string"`
	ForkDbHeadBlockNum       int64  `json:"fork_db_head_block_num"`
	ForkDbHeadBlockID        string `json:"fork_db_head_block_id"`
	ServerFullVersionString  string `json:"server_full_version_string"`
}
type StAccount struct {
	AccountName       string `json:"account_name"`
	HeadBlockNum      int64  `json:"head_block_num"`
	HeadBlockTime     string `json:"head_block_time"`
	Privileged        bool   `json:"privileged"`
	LastCodeUpdate    string `json:"last_code_update"`
	Created           string `json:"created"`
	CoreLiquidBalance string `json:"core_liquid_balance"`
	RAMQuota          int64  `json:"ram_quota"`
	NetWeight         int64  `json:"net_weight"`
	CPUWeight         int64  `json:"cpu_weight"`
	NetLimit          struct {
		Used      int64 `json:"used"`
		Available int64 `json:"available"`
		Max       int64 `json:"max"`
	} `json:"net_limit"`
	CPULimit struct {
		Used      int64 `json:"used"`
		Available int64 `json:"available"`
		Max       int64 `json:"max"`
	} `json:"cpu_limit"`
	RAMUsage       int64 `json:"ram_usage"`
	TotalResources struct {
		Owner     string `json:"owner"`
		NetWeight string `json:"net_weight"`
		CPUWeight string `json:"cpu_weight"`
		RAMBytes  int64  `json:"ram_bytes"`
	} `json:"total_resources"`
	SelfDelegatedBandwidth struct {
		From      string `json:"from"`
		To        string `json:"to"`
		NetWeight string `json:"net_weight"`
		CPUWeight string `json:"cpu_weight"`
	} `json:"self_delegated_bandwidth"`
	RefundRequest struct {
		Owner       string `json:"owner"`
		RequestTime string `json:"request_time"`
		NetAmount   string `json:"net_amount"`
		CPUAmount   string `json:"cpu_amount"`
	} `json:"refund_request"`
}

type StAction struct {
	Account       string `json:"account"`
	Name          string `json:"name"`
	Authorization []struct {
		Actor      string `json:"actor"`
		Permission string `json:"permission"`
	} `json:"authorization"`
	Data json.RawMessage `json:"data"`
}

type StActionData struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Quantity string `json:"quantity"`
	Memo     string `json:"memo"`
}

type StTransactionTrx struct {
	ID                    string        `json:"id"`
	Signatures            []string      `json:"signatures"`
	Compression           string        `json:"compression"`
	PackedContextFreeData string        `json:"packed_context_free_data"`
	ContextFreeData       []interface{} `json:"context_free_data"`
	Transaction           struct {
		Expiration         string        `json:"expiration"`
		RefBlockNum        int           `json:"ref_block_num"`
		RefBlockPrefix     int64         `json:"ref_block_prefix"`
		MaxNetUsageWords   int           `json:"max_net_usage_words"`
		MaxCPUUsageMs      int           `json:"max_cpu_usage_ms"`
		DelaySec           int           `json:"delay_sec"`
		ContextFreeActions []interface{} `json:"context_free_actions"`
		Actions            []StAction    `json:"actions"`
	} `json:"transaction"`
}

type StTransaction struct {
	Status        string          `json:"status"`
	CPUUsageUs    int             `json:"cpu_usage_us"`
	NetUsageWords int             `json:"net_usage_words"`
	Trx           json.RawMessage `json:"trx"`
}

type StBlock struct {
	Timestamp         string          `json:"timestamp"`
	Producer          string          `json:"producer"`
	Confirmed         int             `json:"confirmed"`
	Previous          string          `json:"previous"`
	TransactionMroot  string          `json:"transaction_mroot"`
	ActionMroot       string          `json:"action_mroot"`
	ScheduleVersion   int             `json:"schedule_version"`
	NewProducers      interface{}     `json:"new_producers"`
	ProducerSignature string          `json:"producer_signature"`
	Transactions      []StTransaction `json:"transactions"`
	ID                string          `json:"id"`
	BlockNum          int             `json:"block_num"`
	RefBlockPrefix    int             `json:"ref_block_prefix"`
}

type StPushTransaction struct {
	TransactionID string `json:"transaction_id"`
	Processed     struct {
		ID        string `json:"id"`
		BlockNum  int64  `json:"block_num"`
		BlockTime string `json:"block_time"`
		Receipt   struct {
			Status        string `json:"status"`
			CPUUsageUs    int64  `json:"cpu_usage_us"`
			NetUsageWords int64  `json:"net_usage_words"`
		} `json:"receipt"`
		Elapsed         int64       `json:"elapsed"`
		NetUsage        int64       `json:"net_usage"`
		Scheduled       bool        `json:"scheduled"`
		AccountRAMDelta interface{} `json:"account_ram_delta"`
		Except          interface{} `json:"except"`
		ErrorCode       interface{} `json:"error_code"`
	} `json:"processed"`
}

type StPushTransactionArg struct {
	Signatures            []string `json:"signatures"`
	Compression           string   `json:"compression"`
	PackedContextFreeData string   `json:"packed_context_free_data"`
	PackedTrx             string   `json:"packed_trx"`
}

type StGetTransaction struct {
	BlockTime             time.Time `json:"block_time"`
	BlockNum              int64     `json:"block_num"`
	HeadBlockNum          int64     `json:"head_block_num"`
	ID                    string    `json:"id"`
	LastIrreversibleBlock int64     `json:"last_irreversible_block"`
	Irreversible          bool      `json:"irreversible"`
	Traces                []struct {
		AccountRAMDeltas                       []interface{} `json:"account_ram_deltas"`
		ActionOrdinal                          int64         `json:"action_ordinal"`
		BlockNum                               int64         `json:"block_num"`
		BlockTime                              time.Time     `json:"block_time"`
		ClosestUnnotifiedAncestorActionOrdinal int64         `json:"closest_unnotified_ancestor_action_ordinal"`
		ContextFree                            bool          `json:"context_free"`
		CreatorActionOrdinal                   int64         `json:"creator_action_ordinal"`
		Elapsed                                int64         `json:"elapsed"`
		Except                                 string        `json:"except"`
		ProducerBlockID                        string        `json:"producer_block_id"`
		Receipt                                struct {
			AbiSequence    int64  `json:"abi_sequence"`
			CodeSequence   int64  `json:"code_sequence"`
			GlobalSequence int64  `json:"global_sequence"`
			Receiver       string `json:"receiver"`
		} `json:"receipt"`
		Receiver string `json:"receiver"`
		TrxID    string `json:"trx_id"`
	} `json:"traces"`
	Trx struct {
		Receipt struct {
			CPUUsageUs    int64  `json:"cpu_usage_us"`
			NetUsageWords int64  `json:"net_usage_words"`
			Status        string `json:"status"`
		} `json:"receipt"`
		Trx struct {
			Actions []struct {
				Account       string `json:"account"`
				Authorization []struct {
					Actor      string `json:"actor"`
					Permission string `json:"permission"`
				} `json:"authorization"`
				HexData string `json:"hex_data"`
				Name    string `json:"name"`
			} `json:"actions"`
			Expiration  time.Time `json:"expiration"`
			RefBlockNum int64     `json:"ref_block_num"`
		} `json:"trx"`
	} `json:"trx"`
}

// InitClient 初始化客户端
func InitClient(uri string) {
	rpcURI = uri
	client = gorequest.New().Timeout(time.Minute * 5)
}

func doReq(funURI string, arqs interface{}, resp interface{}) error {
	_, body, errs := client.Post(rpcURI + funURI).Send(arqs).EndBytes()
	if errs != nil {
		return errs[0]
	}
	err := json.Unmarshal(body, resp)
	if err != nil {
		return err
	}
	return nil
}

// RpcChainGetInfo 获取链信息
func RpcChainGetInfo() (*StChainGetInfo, error) {
	resp := struct {
		StRpcRespError
		StChainGetInfo
	}{}
	err := doReq(
		"/v1/chain/get_info",
		nil,
		&resp,
	)
	if err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("%#v", resp)
	}
	return &resp.StChainGetInfo, nil
}

// RpcChainGetAccount 获取账户信息
func RpcChainGetAccount(account string) (*StAccount, error) {
	resp := struct {
		StRpcRespError
		StAccount
	}{}
	err := doReq(
		"/v1/chain/get_account",
		gin.H{
			"account_name": account,
		},
		&resp,
	)
	if err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("%#v", resp)
	}
	return &resp.StAccount, nil
}

// RpcChainGetBlock 获取链信息
func RpcChainGetBlock(blockNum int64) (*StBlock, error) {
	resp := struct {
		StRpcRespError
		StBlock
	}{}
	err := doReq(
		"/v1/chain/get_block",
		gin.H{
			"block_num_or_id": blockNum,
		},
		&resp,
	)
	if err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, &resp
	}
	return &resp.StBlock, nil
}

// RpcChainPushTransaction 推送交易
func RpcChainPushTransaction(arg StPushTransactionArg) (*StPushTransaction, error) {
	resp := struct {
		StRpcRespError
		StPushTransaction
	}{}
	err := doReq(
		"/v1/chain/push_transaction",
		arg,
		&resp,
	)
	if err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, &(resp.StRpcRespError)
	}
	return &resp.StPushTransaction, nil
}

// RpcChainPushTransaction 推送交易
func RpcHistoryGetTransaction(id string) (*StPushTransaction, error) {
	resp := struct {
		StRpcRespError
		StPushTransaction
	}{}
	err := doReq(
		"/v1/history/get_transaction",
		gin.H{
			"id": id,
		},
		&resp,
	)
	if err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, &(resp.StRpcRespError)
	}
	return &resp.StPushTransaction, nil
}
