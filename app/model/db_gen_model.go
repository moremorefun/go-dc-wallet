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

// const TAppLock
const (
	DBColTAppLockID         = "t_app_lock.id"
	DBColTAppLockK          = "t_app_lock.k"           // 上锁键值
	DBColTAppLockV          = "t_app_lock.v"           // 是否锁定
	DBColTAppLockCreateTime = "t_app_lock.create_time" // 上锁时间
)

// DBTAppLock t_app_lock 数据表
/*
   id,
   k,
   v,
   create_time
*/
type DBTAppLock struct {
	ID         int64  `db:"id" json:"id"`
	K          string `db:"k" json:"k"`                     // 上锁键值
	V          int64  `db:"v" json:"v"`                     // 是否锁定
	CreateTime int64  `db:"create_time" json:"create_time"` // 上锁时间
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

// const TProduct
const (
	DBColTProductID          = "t_product.id"
	DBColTProductAppName     = "t_product.app_name"     // 应用名
	DBColTProductAppSk       = "t_product.app_sk"       // 应用私钥
	DBColTProductCbURL       = "t_product.cb_url"       // 回调地址
	DBColTProductWhitelistIP = "t_product.whitelist_ip" // ip白名单
)

// DBTProduct t_product 数据表
/*
   id,
   app_name,
   app_sk,
   cb_url,
   whitelist_ip
*/
type DBTProduct struct {
	ID          int64  `db:"id" json:"id"`
	AppName     string `db:"app_name" json:"app_name"`         // 应用名
	AppSk       string `db:"app_sk" json:"app_sk"`             // 应用私钥
	CbURL       string `db:"cb_url" json:"cb_url"`             // 回调地址
	WhitelistIP string `db:"whitelist_ip" json:"whitelist_ip"` // ip白名单
}

// const TProductNonce
const (
	DBColTProductNonceID         = "t_product_nonce.id"
	DBColTProductNonceC          = "t_product_nonce.c"
	DBColTProductNonceCreateTime = "t_product_nonce.create_time"
)

// DBTProductNonce t_product_nonce 数据表
/*
   id,
   c,
   create_time
*/
type DBTProductNonce struct {
	ID         int64  `db:"id" json:"id"`
	C          string `db:"c" json:"c"`
	CreateTime int64  `db:"create_time" json:"create_time"`
}

// const TSend
const (
	DBColTSendID           = "t_send.id"
	DBColTSendRelatedType  = "t_send.related_type"  // 关联类型 1 零钱整理 2 提币
	DBColTSendRelatedID    = "t_send.related_id"    // 关联id
	DBColTSendTxID         = "t_send.tx_id"         // tx hash
	DBColTSendFromAddress  = "t_send.from_address"  // 打币地址
	DBColTSendToAddress    = "t_send.to_address"    // 收币地址
	DBColTSendBalance      = "t_send.balance"       // 打币金额 Wei
	DBColTSendBalanceReal  = "t_send.balance_real"  // 打币金额 Ether
	DBColTSendGas          = "t_send.gas"           // gas消耗
	DBColTSendGasPrice     = "t_send.gas_price"     // gasPrice
	DBColTSendNonce        = "t_send.nonce"         // nonce
	DBColTSendHex          = "t_send.hex"           // tx raw hex
	DBColTSendCreateTime   = "t_send.create_time"   // 创建时间
	DBColTSendHandleStatus = "t_send.handle_status" // 处理状态
	DBColTSendHandleMsg    = "t_send.handle_msg"    // 处理消息
	DBColTSendHandleTime   = "t_send.handle_time"   // 处理时间
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
	RelatedType  int64  `db:"related_type" json:"related_type"`   // 关联类型 1 零钱整理 2 提币
	RelatedID    int64  `db:"related_id" json:"related_id"`       // 关联id
	TxID         string `db:"tx_id" json:"tx_id"`                 // tx hash
	FromAddress  string `db:"from_address" json:"from_address"`   // 打币地址
	ToAddress    string `db:"to_address" json:"to_address"`       // 收币地址
	Balance      int64  `db:"balance" json:"balance"`             // 打币金额 Wei
	BalanceReal  string `db:"balance_real" json:"balance_real"`   // 打币金额 Ether
	Gas          int64  `db:"gas" json:"gas"`                     // gas消耗
	GasPrice     int64  `db:"gas_price" json:"gas_price"`         // gasPrice
	Nonce        int64  `db:"nonce" json:"nonce"`                 // nonce
	Hex          string `db:"hex" json:"hex"`                     // tx raw hex
	CreateTime   int64  `db:"create_time" json:"create_time"`     // 创建时间
	HandleStatus int64  `db:"handle_status" json:"handle_status"` // 处理状态
	HandleMsg    string `db:"handle_msg" json:"handle_msg"`       // 处理消息
	HandleTime   int64  `db:"handle_time" json:"handle_time"`     // 处理时间
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
	DBColTTxHandleMsg    = "t_tx.handle_msg"    // 处理消息
	DBColTTxHandleTime   = "t_tx.handle_time"   // 处理时间戳
	DBColTTxOrgStatus    = "t_tx.org_status"    // 零钱整理状态
	DBColTTxOrgMsg       = "t_tx.org_msg"       // 零钱整理消息
	DBColTTxOrgTime      = "t_tx.org_time"      // 零钱整理时间
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
	HandleMsg    string `db:"handle_msg" json:"handle_msg"`       // 处理消息
	HandleTime   int64  `db:"handle_time" json:"handle_time"`     // 处理时间戳
	OrgStatus    int64  `db:"org_status" json:"org_status"`       // 零钱整理状态
	OrgMsg       string `db:"org_msg" json:"org_msg"`             // 零钱整理消息
	OrgTime      int64  `db:"org_time" json:"org_time"`           // 零钱整理时间
}

// const TWithdraw
const (
	DBColTWithdrawID           = "t_withdraw.id"
	DBColTWithdrawProductID    = "t_withdraw.product_id"    // 产品id
	DBColTWithdrawOutSerial    = "t_withdraw.out_serial"    // 提币唯一标示
	DBColTWithdrawToAddress    = "t_withdraw.to_address"    // 提币地址
	DBColTWithdrawBalanceReal  = "t_withdraw.balance_real"  // 提币金额
	DBColTWithdrawTxHash       = "t_withdraw.tx_hash"       // 提币tx hash
	DBColTWithdrawCreateTime   = "t_withdraw.create_time"   // 创建时间
	DBColTWithdrawHandleStatus = "t_withdraw.handle_status" // 处理状态
	DBColTWithdrawHandleMsg    = "t_withdraw.handle_msg"    // 处理消息
	DBColTWithdrawHandleTime   = "t_withdraw.handle_time"   // 处理时间
)

// DBTWithdraw t_withdraw 数据表
/*
   id,
   product_id,
   out_serial,
   to_address,
   balance_real,
   tx_hash,
   create_time,
   handle_status,
   handle_msg,
   handle_time
*/
type DBTWithdraw struct {
	ID           int64  `db:"id" json:"id"`
	ProductID    int64  `db:"product_id" json:"product_id"`       // 产品id
	OutSerial    string `db:"out_serial" json:"out_serial"`       // 提币唯一标示
	ToAddress    string `db:"to_address" json:"to_address"`       // 提币地址
	BalanceReal  string `db:"balance_real" json:"balance_real"`   // 提币金额
	TxHash       string `db:"tx_hash" json:"tx_hash"`             // 提币tx hash
	CreateTime   int64  `db:"create_time" json:"create_time"`     // 创建时间
	HandleStatus int64  `db:"handle_status" json:"handle_status"` // 处理状态
	HandleMsg    string `db:"handle_msg" json:"handle_msg"`       // 处理消息
	HandleTime   int64  `db:"handle_time" json:"handle_time"`     // 处理时间
}
