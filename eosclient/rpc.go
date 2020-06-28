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
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Error   struct {
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

// InitClient 初始化客户端
func InitClient(uri string) {
	rpcURI = uri
	client = gorequest.New().Timeout(time.Minute * 5)
}

func doReq(funURI string, arqs gin.H, resp interface{}) error {
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
		return nil, fmt.Errorf("%#v", resp)
	}
	return &resp.StBlock, nil
}
