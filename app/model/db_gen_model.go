package model

import "time"

// const TAddressKey
const (
	DBColTAddressKeyID      = "t_address_key.id"
	DBColTAddressKeyAddress = "t_address_key.address" // 地址
	DBColTAddressKeyPwd     = "t_address_key.pwd"     // 加密私钥
	DBColTAddressKeyUseTag  = "t_address_key.use_tag" // 占用标志 -1 作为热钱包占用-0 未占用->0 作为用户冲币地址占用
)

// DBTAddressKey t_address_key 数据表
/*
   id,
   address,
   pwd,
   use_tag
*/
type DBTAddressKey struct {
	ID      int64  `db:"id" json:"id"`
	Address string `db:"address" json:"address"` // 地址
	Pwd     string `db:"pwd" json:"pwd"`         // 加密私钥
	UseTag  int64  `db:"use_tag" json:"use_tag"` // 占用标志 -1 作为热钱包占用-0 未占用->0 作为用户冲币地址占用
}

// const TAppConfigInt
const (
	DBColTAppConfigIntID = "t_app_config_int.id"
	DBColTAppConfigIntK  = "t_app_config_int.k"
	DBColTAppConfigIntV  = "t_app_config_int.v"
)

// DBTAppConfigInt t_app_config_int 数据表
/*
   id,
   k,
   v
*/
type DBTAppConfigInt struct {
	ID int64  `db:"id" json:"id"`
	K  string `db:"k" json:"k"`
	V  int64  `db:"v" json:"v"`
}

// const TAppStatusInt
const (
	DBColTAppStatusIntID = "t_app_status_int.id"
	DBColTAppStatusIntK  = "t_app_status_int.k"
	DBColTAppStatusIntV  = "t_app_status_int.v"
)

// DBTAppStatusInt t_app_status_int 数据表
/*
   id,
   k,
   v
*/
type DBTAppStatusInt struct {
	ID int64  `db:"id" json:"id"`
	K  string `db:"k" json:"k"`
	V  int64  `db:"v" json:"v"`
}

// const TTx
const (
	DBColTTxID           = "t_tx.id"
	DBColTTxTxID         = "t_tx.tx_id"        // 交易id
	DBColTTxFromAddress  = "t_tx.from_address" // 来源地址
	DBColTTxToAddress    = "t_tx.to_address"   // 目标地址
	DBColTTxValue        = "t_tx.value"        // 到账金额
	DBColTTxCreateTime   = "t_tx.create_time"
	DBColTTxHandleStatus = "t_tx.handle_status" // 处理状态
	DBColTTxHandleMsg    = "t_tx.handle_msg"
	DBColTTxHandleTime   = "t_tx.handle_time"
)

// DBTTx t_tx 数据表
/*
   id,
   tx_id,
   from_address,
   to_address,
   value,
   create_time,
   handle_status,
   handle_msg,
   handle_time
*/
type DBTTx struct {
	ID           int64     `db:"id" json:"id"`
	TxID         string    `db:"tx_id" json:"tx_id"`               // 交易id
	FromAddress  string    `db:"from_address" json:"from_address"` // 来源地址
	ToAddress    string    `db:"to_address" json:"to_address"`     // 目标地址
	Value        string    `db:"value" json:"value"`               // 到账金额
	CreateTime   time.Time `db:"create_time" json:"create_time"`
	HandleStatus int64     `db:"handle_status" json:"handle_status"` // 处理状态
	HandleMsg    string    `db:"handle_msg" json:"handle_msg"`
	HandleTime   time.Time `db:"handle_time" json:"handle_time"`
}
