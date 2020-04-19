package model

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
	DBColTAppConfigIntK  = "t_app_config_int.k" // 配置键名
	DBColTAppConfigIntV  = "t_app_config_int.v" // 配置键值
)

// DBTAppConfigInt t_app_config_int 数据表
/*
   id,
   k,
   v
*/
type DBTAppConfigInt struct {
	ID int64  `db:"id" json:"id"`
	K  string `db:"k" json:"k"` // 配置键名
	V  int64  `db:"v" json:"v"` // 配置键值
}

// const TAppConfigStr
const (
	DBColTAppConfigStrID = "t_app_config_str.id"
	DBColTAppConfigStrK  = "t_app_config_str.k" // 配置键名
	DBColTAppConfigStrV  = "t_app_config_str.v" // 配置键值
)

// DBTAppConfigStr t_app_config_str 数据表
/*
   id,
   k,
   v
*/
type DBTAppConfigStr struct {
	ID int64  `db:"id" json:"id"`
	K  string `db:"k" json:"k"` // 配置键名
	V  string `db:"v" json:"v"` // 配置键值
}

// const TAppStatusInt
const (
	DBColTAppStatusIntID = "t_app_status_int.id"
	DBColTAppStatusIntK  = "t_app_status_int.k" // 配置键名
	DBColTAppStatusIntV  = "t_app_status_int.v" // 配置键值
)

// DBTAppStatusInt t_app_status_int 数据表
/*
   id,
   k,
   v
*/
type DBTAppStatusInt struct {
	ID int64  `db:"id" json:"id"`
	K  string `db:"k" json:"k"` // 配置键名
	V  int64  `db:"v" json:"v"` // 配置键值
}

// const TSend
const (
	DBColTSendID           = "t_send.id"
	DBColTSendRelatedType  = "t_send.related_type"
	DBColTSendRelatedID    = "t_send.related_id"
	DBColTSendTxID         = "t_send.tx_id"
	DBColTSendFromAddress  = "t_send.from_address"
	DBColTSendToAddress    = "t_send.to_address"
	DBColTSendBalance      = "t_send.balance"
	DBColTSendBalanceReal  = "t_send.balance_real"
	DBColTSendGas          = "t_send.gas"
	DBColTSendGasPrice     = "t_send.gas_price"
	DBColTSendNonce        = "t_send.nonce"
	DBColTSendHex          = "t_send.hex"
	DBColTSendCreateTime   = "t_send.create_time"
	DBColTSendHandleStatus = "t_send.handle_status"
	DBColTSendHandleMsg    = "t_send.handle_msg"
	DBColTSendHandleTime   = "t_send.handle_time"
)

// DBTSend t_send 数据表
/*
   id,
   related_type,
   related_id,
   tx_id,
   from_address,
   to_address,
   balance,
   balance_real,
   gas,
   gas_price,
   nonce,
   hex,
   create_time,
   handle_status,
   handle_msg,
   handle_time
*/
type DBTSend struct {
	ID           int64  `db:"id" json:"id"`
	RelatedType  int64  `db:"related_type" json:"related_type"`
	RelatedID    int64  `db:"related_id" json:"related_id"`
	TxID         string `db:"tx_id" json:"tx_id"`
	FromAddress  string `db:"from_address" json:"from_address"`
	ToAddress    string `db:"to_address" json:"to_address"`
	Balance      int64  `db:"balance" json:"balance"`
	BalanceReal  string `db:"balance_real" json:"balance_real"`
	Gas          int64  `db:"gas" json:"gas"`
	GasPrice     int64  `db:"gas_price" json:"gas_price"`
	Nonce        int64  `db:"nonce" json:"nonce"`
	Hex          string `db:"hex" json:"hex"`
	CreateTime   int64  `db:"create_time" json:"create_time"`
	HandleStatus int64  `db:"handle_status" json:"handle_status"`
	HandleMsg    string `db:"handle_msg" json:"handle_msg"`
	HandleTime   int64  `db:"handle_time" json:"handle_time"`
}

// const TTx
const (
	DBColTTxID           = "t_tx.id"
	DBColTTxTxID         = "t_tx.tx_id"         // 交易id
	DBColTTxFromAddress  = "t_tx.from_address"  // 来源地址
	DBColTTxToAddress    = "t_tx.to_address"    // 目标地址
	DBColTTxBalance      = "t_tx.balance"       // 到账金额Wei
	DBColTTxBalanceReal  = "t_tx.balance_real"  // 到账金额Ether
	DBColTTxCreateTime   = "t_tx.create_time"   // 创建时间戳
	DBColTTxHandleStatus = "t_tx.handle_status" // 处理状态
	DBColTTxHandleMsg    = "t_tx.handle_msg"
	DBColTTxHandleTime   = "t_tx.handle_time" // 处理时间戳
	DBColTTxOrgStatus    = "t_tx.org_status"
	DBColTTxOrgMsg       = "t_tx.org_msg"
	DBColTTxOrgTime      = "t_tx.org_time"
)

// DBTTx t_tx 数据表
/*
   id,
   tx_id,
   from_address,
   to_address,
   balance,
   balance_real,
   create_time,
   handle_status,
   handle_msg,
   handle_time,
   org_status,
   org_msg,
   org_time
*/
type DBTTx struct {
	ID           int64  `db:"id" json:"id"`
	TxID         string `db:"tx_id" json:"tx_id"`                 // 交易id
	FromAddress  string `db:"from_address" json:"from_address"`   // 来源地址
	ToAddress    string `db:"to_address" json:"to_address"`       // 目标地址
	Balance      int64  `db:"balance" json:"balance"`             // 到账金额Wei
	BalanceReal  string `db:"balance_real" json:"balance_real"`   // 到账金额Ether
	CreateTime   int64  `db:"create_time" json:"create_time"`     // 创建时间戳
	HandleStatus int64  `db:"handle_status" json:"handle_status"` // 处理状态
	HandleMsg    string `db:"handle_msg" json:"handle_msg"`
	HandleTime   int64  `db:"handle_time" json:"handle_time"` // 处理时间戳
	OrgStatus    int64  `db:"org_status" json:"org_status"`
	OrgMsg       string `db:"org_msg" json:"org_msg"`
	OrgTime      int64  `db:"org_time" json:"org_time"`
}
