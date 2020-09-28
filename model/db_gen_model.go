package model

// TableNames 所有表名
var TableNames = []string{"t_address_key", "t_app_config_int", "t_app_config_str", "t_app_config_token", "t_app_config_token_btc", "t_app_lock", "t_app_status_int", "t_product", "t_product_nonce", "t_product_notify", "t_send", "t_send_btc", "t_send_eos", "t_tx", "t_tx_btc", "t_tx_btc_token", "t_tx_btc_uxto", "t_tx_eos", "t_tx_erc20", "t_withdraw"}

// 表名
const (
	DbTableTAddressKey        = "t_address_key"
	DbTableTAppConfigInt      = "t_app_config_int"
	DbTableTAppConfigStr      = "t_app_config_str"
	DbTableTAppConfigToken    = "t_app_config_token"
	DbTableTAppConfigTokenBtc = "t_app_config_token_btc"
	DbTableTAppLock           = "t_app_lock"
	DbTableTAppStatusInt      = "t_app_status_int"
	DbTableTProduct           = "t_product"
	DbTableTProductNonce      = "t_product_nonce"
	DbTableTProductNotify     = "t_product_notify"
	DbTableTSend              = "t_send"
	DbTableTSendBtc           = "t_send_btc"
	DbTableTSendEos           = "t_send_eos"
	DbTableTTx                = "t_tx"
	DbTableTTxBtc             = "t_tx_btc"
	DbTableTTxBtcToken        = "t_tx_btc_token"
	DbTableTTxBtcUxto         = "t_tx_btc_uxto"
	DbTableTTxEos             = "t_tx_eos"
	DbTableTTxErc20           = "t_tx_erc20"
	DbTableTWithdraw          = "t_withdraw"
)

// 字段名

// const TAddressKey full
const (
	DBColTAddressKeyID      = "t_address_key.id"
	DBColTAddressKeySymbol  = "t_address_key.symbol"  // 币种
	DBColTAddressKeyAddress = "t_address_key.address" // 地址
	DBColTAddressKeyPwd     = "t_address_key.pwd"     // 加密私钥
	DBColTAddressKeyUseTag  = "t_address_key.use_tag" // 占用标志 -1 作为热钱包占用-0 未占用->0 作为用户冲币地址占用
)

// const TAddressKey short
const (
	DBColShortTAddressKeyID      = "id"
	DBColShortTAddressKeySymbol  = "symbol"  // 币种
	DBColShortTAddressKeyAddress = "address" // 地址
	DBColShortTAddressKeyPwd     = "pwd"     // 加密私钥
	DBColShortTAddressKeyUseTag  = "use_tag" // 占用标志 -1 作为热钱包占用-0 未占用->0 作为用户冲币地址占用
)

// DBColTAddressKeyAll 所有字段
var DBColTAddressKeyAll = []string{
	"t_address_key.id",
	"t_address_key.symbol",
	"t_address_key.address",
	"t_address_key.pwd",
	"t_address_key.use_tag",
}

// 表结构
// DBTAddressKey t_address_key
/*
   id,
   symbol,
   address,
   pwd,
   use_tag
*/
type DBTAddressKey struct {
	ID      int64  `db:"id" json:"id"`
	Symbol  string `db:"symbol" json:"symbol"`   // 币种
	Address string `db:"address" json:"address"` // 地址
	Pwd     string `db:"pwd" json:"pwd"`         // 加密私钥
	UseTag  int64  `db:"use_tag" json:"use_tag"` // 占用标志 -1 作为热钱包占用-0 未占用->0 作为用户冲币地址占用
}

// const TAppConfigInt full
const (
	DBColTAppConfigIntID = "t_app_config_int.id"
	DBColTAppConfigIntK  = "t_app_config_int.k" // 配置键名
	DBColTAppConfigIntV  = "t_app_config_int.v" // 配置键值
)

// const TAppConfigInt short
const (
	DBColShortTAppConfigIntID = "id"
	DBColShortTAppConfigIntK  = "k" // 配置键名
	DBColShortTAppConfigIntV  = "v" // 配置键值
)

// DBColTAppConfigIntAll 所有字段
var DBColTAppConfigIntAll = []string{
	"t_app_config_int.id",
	"t_app_config_int.k",
	"t_app_config_int.v",
}

// 表结构
// DBTAppConfigInt t_app_config_int
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

// const TAppConfigStr full
const (
	DBColTAppConfigStrID = "t_app_config_str.id"
	DBColTAppConfigStrK  = "t_app_config_str.k" // 配置键名
	DBColTAppConfigStrV  = "t_app_config_str.v" // 配置键值
)

// const TAppConfigStr short
const (
	DBColShortTAppConfigStrID = "id"
	DBColShortTAppConfigStrK  = "k" // 配置键名
	DBColShortTAppConfigStrV  = "v" // 配置键值
)

// DBColTAppConfigStrAll 所有字段
var DBColTAppConfigStrAll = []string{
	"t_app_config_str.id",
	"t_app_config_str.k",
	"t_app_config_str.v",
}

// 表结构
// DBTAppConfigStr t_app_config_str
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

// const TAppConfigToken full
const (
	DBColTAppConfigTokenID            = "t_app_config_token.id"
	DBColTAppConfigTokenTokenAddress  = "t_app_config_token.token_address"
	DBColTAppConfigTokenTokenDecimals = "t_app_config_token.token_decimals"
	DBColTAppConfigTokenTokenSymbol   = "t_app_config_token.token_symbol"
	DBColTAppConfigTokenColdAddress   = "t_app_config_token.cold_address"
	DBColTAppConfigTokenHotAddress    = "t_app_config_token.hot_address"
	DBColTAppConfigTokenOrgMinBalance = "t_app_config_token.org_min_balance"
	DBColTAppConfigTokenCreateTime    = "t_app_config_token.create_time"
)

// const TAppConfigToken short
const (
	DBColShortTAppConfigTokenID            = "id"
	DBColShortTAppConfigTokenTokenAddress  = "token_address"
	DBColShortTAppConfigTokenTokenDecimals = "token_decimals"
	DBColShortTAppConfigTokenTokenSymbol   = "token_symbol"
	DBColShortTAppConfigTokenColdAddress   = "cold_address"
	DBColShortTAppConfigTokenHotAddress    = "hot_address"
	DBColShortTAppConfigTokenOrgMinBalance = "org_min_balance"
	DBColShortTAppConfigTokenCreateTime    = "create_time"
)

// DBColTAppConfigTokenAll 所有字段
var DBColTAppConfigTokenAll = []string{
	"t_app_config_token.id",
	"t_app_config_token.token_address",
	"t_app_config_token.token_decimals",
	"t_app_config_token.token_symbol",
	"t_app_config_token.cold_address",
	"t_app_config_token.hot_address",
	"t_app_config_token.org_min_balance",
	"t_app_config_token.create_time",
}

// 表结构
// DBTAppConfigToken t_app_config_token
/*
   id,
   token_address,
   token_decimals,
   token_symbol,
   cold_address,
   hot_address,
   org_min_balance,
   create_time
*/
type DBTAppConfigToken struct {
	ID            int64  `db:"id" json:"id"`
	TokenAddress  string `db:"token_address" json:"token_address"`
	TokenDecimals int64  `db:"token_decimals" json:"token_decimals"`
	TokenSymbol   string `db:"token_symbol" json:"token_symbol"`
	ColdAddress   string `db:"cold_address" json:"cold_address"`
	HotAddress    string `db:"hot_address" json:"hot_address"`
	OrgMinBalance string `db:"org_min_balance" json:"org_min_balance"`
	CreateTime    int64  `db:"create_time" json:"create_time"`
}

// const TAppConfigTokenBtc full
const (
	DBColTAppConfigTokenBtcID              = "t_app_config_token_btc.id"
	DBColTAppConfigTokenBtcTokenIndex      = "t_app_config_token_btc.token_index"
	DBColTAppConfigTokenBtcTokenSymbol     = "t_app_config_token_btc.token_symbol"
	DBColTAppConfigTokenBtcColdAddress     = "t_app_config_token_btc.cold_address"
	DBColTAppConfigTokenBtcHotAddress      = "t_app_config_token_btc.hot_address"
	DBColTAppConfigTokenBtcFeeAddress      = "t_app_config_token_btc.fee_address"
	DBColTAppConfigTokenBtcTxOrgMinBalance = "t_app_config_token_btc.tx_org_min_balance"
	DBColTAppConfigTokenBtcCreateAt        = "t_app_config_token_btc.create_at"
)

// const TAppConfigTokenBtc short
const (
	DBColShortTAppConfigTokenBtcID              = "id"
	DBColShortTAppConfigTokenBtcTokenIndex      = "token_index"
	DBColShortTAppConfigTokenBtcTokenSymbol     = "token_symbol"
	DBColShortTAppConfigTokenBtcColdAddress     = "cold_address"
	DBColShortTAppConfigTokenBtcHotAddress      = "hot_address"
	DBColShortTAppConfigTokenBtcFeeAddress      = "fee_address"
	DBColShortTAppConfigTokenBtcTxOrgMinBalance = "tx_org_min_balance"
	DBColShortTAppConfigTokenBtcCreateAt        = "create_at"
)

// DBColTAppConfigTokenBtcAll 所有字段
var DBColTAppConfigTokenBtcAll = []string{
	"t_app_config_token_btc.id",
	"t_app_config_token_btc.token_index",
	"t_app_config_token_btc.token_symbol",
	"t_app_config_token_btc.cold_address",
	"t_app_config_token_btc.hot_address",
	"t_app_config_token_btc.fee_address",
	"t_app_config_token_btc.tx_org_min_balance",
	"t_app_config_token_btc.create_at",
}

// 表结构
// DBTAppConfigTokenBtc t_app_config_token_btc
/*
   id,
   token_index,
   token_symbol,
   cold_address,
   hot_address,
   fee_address,
   tx_org_min_balance,
   create_at
*/
type DBTAppConfigTokenBtc struct {
	ID              int64  `db:"id" json:"id"`
	TokenIndex      int64  `db:"token_index" json:"token_index"`
	TokenSymbol     string `db:"token_symbol" json:"token_symbol"`
	ColdAddress     string `db:"cold_address" json:"cold_address"`
	HotAddress      string `db:"hot_address" json:"hot_address"`
	FeeAddress      string `db:"fee_address" json:"fee_address"`
	TxOrgMinBalance string `db:"tx_org_min_balance" json:"tx_org_min_balance"`
	CreateAt        int64  `db:"create_at" json:"create_at"`
}

// const TAppLock full
const (
	DBColTAppLockID         = "t_app_lock.id"
	DBColTAppLockK          = "t_app_lock.k"           // 上锁键值
	DBColTAppLockV          = "t_app_lock.v"           // 是否锁定
	DBColTAppLockCreateTime = "t_app_lock.create_time" // 上锁时间
)

// const TAppLock short
const (
	DBColShortTAppLockID         = "id"
	DBColShortTAppLockK          = "k"           // 上锁键值
	DBColShortTAppLockV          = "v"           // 是否锁定
	DBColShortTAppLockCreateTime = "create_time" // 上锁时间
)

// DBColTAppLockAll 所有字段
var DBColTAppLockAll = []string{
	"t_app_lock.id",
	"t_app_lock.k",
	"t_app_lock.v",
	"t_app_lock.create_time",
}

// 表结构
// DBTAppLock t_app_lock
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

// const TAppStatusInt full
const (
	DBColTAppStatusIntID = "t_app_status_int.id"
	DBColTAppStatusIntK  = "t_app_status_int.k" // 配置键名
	DBColTAppStatusIntV  = "t_app_status_int.v" // 配置键值
)

// const TAppStatusInt short
const (
	DBColShortTAppStatusIntID = "id"
	DBColShortTAppStatusIntK  = "k" // 配置键名
	DBColShortTAppStatusIntV  = "v" // 配置键值
)

// DBColTAppStatusIntAll 所有字段
var DBColTAppStatusIntAll = []string{
	"t_app_status_int.id",
	"t_app_status_int.k",
	"t_app_status_int.v",
}

// 表结构
// DBTAppStatusInt t_app_status_int
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

// const TProduct full
const (
	DBColTProductID          = "t_product.id"
	DBColTProductAppName     = "t_product.app_name"     // 应用名
	DBColTProductAppSk       = "t_product.app_sk"       // 应用私钥
	DBColTProductCbURL       = "t_product.cb_url"       // 回调地址
	DBColTProductWhitelistIP = "t_product.whitelist_ip" // ip白名单
)

// const TProduct short
const (
	DBColShortTProductID          = "id"
	DBColShortTProductAppName     = "app_name"     // 应用名
	DBColShortTProductAppSk       = "app_sk"       // 应用私钥
	DBColShortTProductCbURL       = "cb_url"       // 回调地址
	DBColShortTProductWhitelistIP = "whitelist_ip" // ip白名单
)

// DBColTProductAll 所有字段
var DBColTProductAll = []string{
	"t_product.id",
	"t_product.app_name",
	"t_product.app_sk",
	"t_product.cb_url",
	"t_product.whitelist_ip",
}

// 表结构
// DBTProduct t_product
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

// const TProductNonce full
const (
	DBColTProductNonceID         = "t_product_nonce.id"
	DBColTProductNonceC          = "t_product_nonce.c"
	DBColTProductNonceCreateTime = "t_product_nonce.create_time"
)

// const TProductNonce short
const (
	DBColShortTProductNonceID         = "id"
	DBColShortTProductNonceC          = "c"
	DBColShortTProductNonceCreateTime = "create_time"
)

// DBColTProductNonceAll 所有字段
var DBColTProductNonceAll = []string{
	"t_product_nonce.id",
	"t_product_nonce.c",
	"t_product_nonce.create_time",
}

// 表结构
// DBTProductNonce t_product_nonce
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

// const TProductNotify full
const (
	DBColTProductNotifyID           = "t_product_notify.id"
	DBColTProductNotifyNonce        = "t_product_notify.nonce"
	DBColTProductNotifyProductID    = "t_product_notify.product_id"
	DBColTProductNotifyItemType     = "t_product_notify.item_type"
	DBColTProductNotifyItemID       = "t_product_notify.item_id"
	DBColTProductNotifyNotifyType   = "t_product_notify.notify_type"
	DBColTProductNotifyTokenSymbol  = "t_product_notify.token_symbol"
	DBColTProductNotifyURL          = "t_product_notify.url"
	DBColTProductNotifyMsg          = "t_product_notify.msg"
	DBColTProductNotifyHandleStatus = "t_product_notify.handle_status"
	DBColTProductNotifyHandleMsg    = "t_product_notify.handle_msg"
	DBColTProductNotifyCreateTime   = "t_product_notify.create_time"
	DBColTProductNotifyUpdateTime   = "t_product_notify.update_time"
)

// const TProductNotify short
const (
	DBColShortTProductNotifyID           = "id"
	DBColShortTProductNotifyNonce        = "nonce"
	DBColShortTProductNotifyProductID    = "product_id"
	DBColShortTProductNotifyItemType     = "item_type"
	DBColShortTProductNotifyItemID       = "item_id"
	DBColShortTProductNotifyNotifyType   = "notify_type"
	DBColShortTProductNotifyTokenSymbol  = "token_symbol"
	DBColShortTProductNotifyURL          = "url"
	DBColShortTProductNotifyMsg          = "msg"
	DBColShortTProductNotifyHandleStatus = "handle_status"
	DBColShortTProductNotifyHandleMsg    = "handle_msg"
	DBColShortTProductNotifyCreateTime   = "create_time"
	DBColShortTProductNotifyUpdateTime   = "update_time"
)

// DBColTProductNotifyAll 所有字段
var DBColTProductNotifyAll = []string{
	"t_product_notify.id",
	"t_product_notify.nonce",
	"t_product_notify.product_id",
	"t_product_notify.item_type",
	"t_product_notify.item_id",
	"t_product_notify.notify_type",
	"t_product_notify.token_symbol",
	"t_product_notify.url",
	"t_product_notify.msg",
	"t_product_notify.handle_status",
	"t_product_notify.handle_msg",
	"t_product_notify.create_time",
	"t_product_notify.update_time",
}

// 表结构
// DBTProductNotify t_product_notify
/*
   id,
   nonce,
   product_id,
   item_type,
   item_id,
   notify_type,
   token_symbol,
   url,
   msg,
   handle_status,
   handle_msg,
   create_time,
   update_time
*/
type DBTProductNotify struct {
	ID           int64  `db:"id" json:"id"`
	Nonce        string `db:"nonce" json:"nonce"`
	ProductID    int64  `db:"product_id" json:"product_id"`
	ItemType     int64  `db:"item_type" json:"item_type"`
	ItemID       int64  `db:"item_id" json:"item_id"`
	NotifyType   int64  `db:"notify_type" json:"notify_type"`
	TokenSymbol  string `db:"token_symbol" json:"token_symbol"`
	URL          string `db:"url" json:"url"`
	Msg          string `db:"msg" json:"msg"`
	HandleStatus int64  `db:"handle_status" json:"handle_status"`
	HandleMsg    string `db:"handle_msg" json:"handle_msg"`
	CreateTime   int64  `db:"create_time" json:"create_time"`
	UpdateTime   int64  `db:"update_time" json:"update_time"`
}

// const TSend full
const (
	DBColTSendID           = "t_send.id"
	DBColTSendRelatedType  = "t_send.related_type" // 关联类型 1 零钱整理 2 提币
	DBColTSendRelatedID    = "t_send.related_id"   // 关联id
	DBColTSendTokenID      = "t_send.token_id"
	DBColTSendTxID         = "t_send.tx_id"         // tx hash
	DBColTSendFromAddress  = "t_send.from_address"  // 打币地址
	DBColTSendToAddress    = "t_send.to_address"    // 收币地址
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

// const TSend short
const (
	DBColShortTSendID           = "id"
	DBColShortTSendRelatedType  = "related_type" // 关联类型 1 零钱整理 2 提币
	DBColShortTSendRelatedID    = "related_id"   // 关联id
	DBColShortTSendTokenID      = "token_id"
	DBColShortTSendTxID         = "tx_id"         // tx hash
	DBColShortTSendFromAddress  = "from_address"  // 打币地址
	DBColShortTSendToAddress    = "to_address"    // 收币地址
	DBColShortTSendBalanceReal  = "balance_real"  // 打币金额 Ether
	DBColShortTSendGas          = "gas"           // gas消耗
	DBColShortTSendGasPrice     = "gas_price"     // gasPrice
	DBColShortTSendNonce        = "nonce"         // nonce
	DBColShortTSendHex          = "hex"           // tx raw hex
	DBColShortTSendCreateTime   = "create_time"   // 创建时间
	DBColShortTSendHandleStatus = "handle_status" // 处理状态
	DBColShortTSendHandleMsg    = "handle_msg"    // 处理消息
	DBColShortTSendHandleTime   = "handle_time"   // 处理时间
)

// DBColTSendAll 所有字段
var DBColTSendAll = []string{
	"t_send.id",
	"t_send.related_type",
	"t_send.related_id",
	"t_send.token_id",
	"t_send.tx_id",
	"t_send.from_address",
	"t_send.to_address",
	"t_send.balance_real",
	"t_send.gas",
	"t_send.gas_price",
	"t_send.nonce",
	"t_send.hex",
	"t_send.create_time",
	"t_send.handle_status",
	"t_send.handle_msg",
	"t_send.handle_time",
}

// 表结构
// DBTSend t_send
/*
   id,
   related_type,
   related_id,
   token_id,
   tx_id,
   from_address,
   to_address,
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
	RelatedType  int64  `db:"related_type" json:"related_type"` // 关联类型 1 零钱整理 2 提币
	RelatedID    int64  `db:"related_id" json:"related_id"`     // 关联id
	TokenID      int64  `db:"token_id" json:"token_id"`
	TxID         string `db:"tx_id" json:"tx_id"`                 // tx hash
	FromAddress  string `db:"from_address" json:"from_address"`   // 打币地址
	ToAddress    string `db:"to_address" json:"to_address"`       // 收币地址
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

// const TSendBtc full
const (
	DBColTSendBtcID           = "t_send_btc.id"
	DBColTSendBtcRelatedType  = "t_send_btc.related_type" // 关联类型 1 零钱整理 2 提币
	DBColTSendBtcRelatedID    = "t_send_btc.related_id"   // 关联id
	DBColTSendBtcTokenID      = "t_send_btc.token_id"
	DBColTSendBtcTxID         = "t_send_btc.tx_id"         // tx hash
	DBColTSendBtcFromAddress  = "t_send_btc.from_address"  // 打币地址
	DBColTSendBtcToAddress    = "t_send_btc.to_address"    // 收币地址
	DBColTSendBtcBalanceReal  = "t_send_btc.balance_real"  // 打币金额 Ether
	DBColTSendBtcGas          = "t_send_btc.gas"           // gas消耗
	DBColTSendBtcGasPrice     = "t_send_btc.gas_price"     // gasPrice
	DBColTSendBtcHex          = "t_send_btc.hex"           // tx raw hex
	DBColTSendBtcCreateTime   = "t_send_btc.create_time"   // 创建时间
	DBColTSendBtcHandleStatus = "t_send_btc.handle_status" // 处理状态
	DBColTSendBtcHandleMsg    = "t_send_btc.handle_msg"    // 处理消息
	DBColTSendBtcHandleTime   = "t_send_btc.handle_time"   // 处理时间
)

// const TSendBtc short
const (
	DBColShortTSendBtcID           = "id"
	DBColShortTSendBtcRelatedType  = "related_type" // 关联类型 1 零钱整理 2 提币
	DBColShortTSendBtcRelatedID    = "related_id"   // 关联id
	DBColShortTSendBtcTokenID      = "token_id"
	DBColShortTSendBtcTxID         = "tx_id"         // tx hash
	DBColShortTSendBtcFromAddress  = "from_address"  // 打币地址
	DBColShortTSendBtcToAddress    = "to_address"    // 收币地址
	DBColShortTSendBtcBalanceReal  = "balance_real"  // 打币金额 Ether
	DBColShortTSendBtcGas          = "gas"           // gas消耗
	DBColShortTSendBtcGasPrice     = "gas_price"     // gasPrice
	DBColShortTSendBtcHex          = "hex"           // tx raw hex
	DBColShortTSendBtcCreateTime   = "create_time"   // 创建时间
	DBColShortTSendBtcHandleStatus = "handle_status" // 处理状态
	DBColShortTSendBtcHandleMsg    = "handle_msg"    // 处理消息
	DBColShortTSendBtcHandleTime   = "handle_time"   // 处理时间
)

// DBColTSendBtcAll 所有字段
var DBColTSendBtcAll = []string{
	"t_send_btc.id",
	"t_send_btc.related_type",
	"t_send_btc.related_id",
	"t_send_btc.token_id",
	"t_send_btc.tx_id",
	"t_send_btc.from_address",
	"t_send_btc.to_address",
	"t_send_btc.balance_real",
	"t_send_btc.gas",
	"t_send_btc.gas_price",
	"t_send_btc.hex",
	"t_send_btc.create_time",
	"t_send_btc.handle_status",
	"t_send_btc.handle_msg",
	"t_send_btc.handle_time",
}

// 表结构
// DBTSendBtc t_send_btc
/*
   id,
   related_type,
   related_id,
   token_id,
   tx_id,
   from_address,
   to_address,
   balance_real,
   gas,
   gas_price,
   hex,
   create_time,
   handle_status,
   handle_msg,
   handle_time
*/
type DBTSendBtc struct {
	ID           int64  `db:"id" json:"id"`
	RelatedType  int64  `db:"related_type" json:"related_type"` // 关联类型 1 零钱整理 2 提币
	RelatedID    int64  `db:"related_id" json:"related_id"`     // 关联id
	TokenID      int64  `db:"token_id" json:"token_id"`
	TxID         string `db:"tx_id" json:"tx_id"`                 // tx hash
	FromAddress  string `db:"from_address" json:"from_address"`   // 打币地址
	ToAddress    string `db:"to_address" json:"to_address"`       // 收币地址
	BalanceReal  string `db:"balance_real" json:"balance_real"`   // 打币金额 Ether
	Gas          int64  `db:"gas" json:"gas"`                     // gas消耗
	GasPrice     int64  `db:"gas_price" json:"gas_price"`         // gasPrice
	Hex          string `db:"hex" json:"hex"`                     // tx raw hex
	CreateTime   int64  `db:"create_time" json:"create_time"`     // 创建时间
	HandleStatus int64  `db:"handle_status" json:"handle_status"` // 处理状态
	HandleMsg    string `db:"handle_msg" json:"handle_msg"`       // 处理消息
	HandleTime   int64  `db:"handle_time" json:"handle_time"`     // 处理时间
}

// const TSendEos full
const (
	DBColTSendEosID           = "t_send_eos.id"
	DBColTSendEosWithdrawID   = "t_send_eos.withdraw_id"   // 关联id
	DBColTSendEosTxHash       = "t_send_eos.tx_hash"       // tx hash
	DBColTSendEosLogIndex     = "t_send_eos.log_index"     // log_index
	DBColTSendEosFromAddress  = "t_send_eos.from_address"  // 打币地址
	DBColTSendEosToAddress    = "t_send_eos.to_address"    // 收币地址
	DBColTSendEosMemo         = "t_send_eos.memo"          // 收币地址
	DBColTSendEosBalanceReal  = "t_send_eos.balance_real"  // 打币金额 Ether
	DBColTSendEosHex          = "t_send_eos.hex"           // tx raw hex
	DBColTSendEosCreateTime   = "t_send_eos.create_time"   // 创建时间
	DBColTSendEosHandleStatus = "t_send_eos.handle_status" // 处理状态
	DBColTSendEosHandleMsg    = "t_send_eos.handle_msg"    // 处理消息
	DBColTSendEosHandleAt     = "t_send_eos.handle_at"     // 处理时间
)

// const TSendEos short
const (
	DBColShortTSendEosID           = "id"
	DBColShortTSendEosWithdrawID   = "withdraw_id"   // 关联id
	DBColShortTSendEosTxHash       = "tx_hash"       // tx hash
	DBColShortTSendEosLogIndex     = "log_index"     // log_index
	DBColShortTSendEosFromAddress  = "from_address"  // 打币地址
	DBColShortTSendEosToAddress    = "to_address"    // 收币地址
	DBColShortTSendEosMemo         = "memo"          // 收币地址
	DBColShortTSendEosBalanceReal  = "balance_real"  // 打币金额 Ether
	DBColShortTSendEosHex          = "hex"           // tx raw hex
	DBColShortTSendEosCreateTime   = "create_time"   // 创建时间
	DBColShortTSendEosHandleStatus = "handle_status" // 处理状态
	DBColShortTSendEosHandleMsg    = "handle_msg"    // 处理消息
	DBColShortTSendEosHandleAt     = "handle_at"     // 处理时间
)

// DBColTSendEosAll 所有字段
var DBColTSendEosAll = []string{
	"t_send_eos.id",
	"t_send_eos.withdraw_id",
	"t_send_eos.tx_hash",
	"t_send_eos.log_index",
	"t_send_eos.from_address",
	"t_send_eos.to_address",
	"t_send_eos.memo",
	"t_send_eos.balance_real",
	"t_send_eos.hex",
	"t_send_eos.create_time",
	"t_send_eos.handle_status",
	"t_send_eos.handle_msg",
	"t_send_eos.handle_at",
}

// 表结构
// DBTSendEos t_send_eos
/*
   id,
   withdraw_id,
   tx_hash,
   log_index,
   from_address,
   to_address,
   memo,
   balance_real,
   hex,
   create_time,
   handle_status,
   handle_msg,
   handle_at
*/
type DBTSendEos struct {
	ID           int64  `db:"id" json:"id"`
	WithdrawID   int64  `db:"withdraw_id" json:"withdraw_id"`     // 关联id
	TxHash       string `db:"tx_hash" json:"tx_hash"`             // tx hash
	LogIndex     int64  `db:"log_index" json:"log_index"`         // log_index
	FromAddress  string `db:"from_address" json:"from_address"`   // 打币地址
	ToAddress    string `db:"to_address" json:"to_address"`       // 收币地址
	Memo         string `db:"memo" json:"memo"`                   // 收币地址
	BalanceReal  string `db:"balance_real" json:"balance_real"`   // 打币金额 Ether
	Hex          string `db:"hex" json:"hex"`                     // tx raw hex
	CreateTime   int64  `db:"create_time" json:"create_time"`     // 创建时间
	HandleStatus int64  `db:"handle_status" json:"handle_status"` // 处理状态
	HandleMsg    string `db:"handle_msg" json:"handle_msg"`       // 处理消息
	HandleAt     int64  `db:"handle_at" json:"handle_at"`         // 处理时间
}

// const TTx full
const (
	DBColTTxID           = "t_tx.id"
	DBColTTxProductID    = "t_tx.product_id"
	DBColTTxTxID         = "t_tx.tx_id"         // 交易id
	DBColTTxFromAddress  = "t_tx.from_address"  // 来源地址
	DBColTTxToAddress    = "t_tx.to_address"    // 目标地址
	DBColTTxBalanceReal  = "t_tx.balance_real"  // 到账金额Ether
	DBColTTxCreateTime   = "t_tx.create_time"   // 创建时间戳
	DBColTTxHandleStatus = "t_tx.handle_status" // 处理状态
	DBColTTxHandleMsg    = "t_tx.handle_msg"    // 处理消息
	DBColTTxHandleTime   = "t_tx.handle_time"   // 处理时间戳
	DBColTTxOrgStatus    = "t_tx.org_status"    // 零钱整理状态
	DBColTTxOrgMsg       = "t_tx.org_msg"       // 零钱整理消息
	DBColTTxOrgTime      = "t_tx.org_time"      // 零钱整理时间
)

// const TTx short
const (
	DBColShortTTxID           = "id"
	DBColShortTTxProductID    = "product_id"
	DBColShortTTxTxID         = "tx_id"         // 交易id
	DBColShortTTxFromAddress  = "from_address"  // 来源地址
	DBColShortTTxToAddress    = "to_address"    // 目标地址
	DBColShortTTxBalanceReal  = "balance_real"  // 到账金额Ether
	DBColShortTTxCreateTime   = "create_time"   // 创建时间戳
	DBColShortTTxHandleStatus = "handle_status" // 处理状态
	DBColShortTTxHandleMsg    = "handle_msg"    // 处理消息
	DBColShortTTxHandleTime   = "handle_time"   // 处理时间戳
	DBColShortTTxOrgStatus    = "org_status"    // 零钱整理状态
	DBColShortTTxOrgMsg       = "org_msg"       // 零钱整理消息
	DBColShortTTxOrgTime      = "org_time"      // 零钱整理时间
)

// DBColTTxAll 所有字段
var DBColTTxAll = []string{
	"t_tx.id",
	"t_tx.product_id",
	"t_tx.tx_id",
	"t_tx.from_address",
	"t_tx.to_address",
	"t_tx.balance_real",
	"t_tx.create_time",
	"t_tx.handle_status",
	"t_tx.handle_msg",
	"t_tx.handle_time",
	"t_tx.org_status",
	"t_tx.org_msg",
	"t_tx.org_time",
}

// 表结构
// DBTTx t_tx
/*
   id,
   product_id,
   tx_id,
   from_address,
   to_address,
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
	ProductID    int64  `db:"product_id" json:"product_id"`
	TxID         string `db:"tx_id" json:"tx_id"`                 // 交易id
	FromAddress  string `db:"from_address" json:"from_address"`   // 来源地址
	ToAddress    string `db:"to_address" json:"to_address"`       // 目标地址
	BalanceReal  string `db:"balance_real" json:"balance_real"`   // 到账金额Ether
	CreateTime   int64  `db:"create_time" json:"create_time"`     // 创建时间戳
	HandleStatus int64  `db:"handle_status" json:"handle_status"` // 处理状态
	HandleMsg    string `db:"handle_msg" json:"handle_msg"`       // 处理消息
	HandleTime   int64  `db:"handle_time" json:"handle_time"`     // 处理时间戳
	OrgStatus    int64  `db:"org_status" json:"org_status"`       // 零钱整理状态
	OrgMsg       string `db:"org_msg" json:"org_msg"`             // 零钱整理消息
	OrgTime      int64  `db:"org_time" json:"org_time"`           // 零钱整理时间
}

// const TTxBtc full
const (
	DBColTTxBtcID           = "t_tx_btc.id"
	DBColTTxBtcProductID    = "t_tx_btc.product_id"
	DBColTTxBtcBlockHash    = "t_tx_btc.block_hash"
	DBColTTxBtcTxID         = "t_tx_btc.tx_id"
	DBColTTxBtcVoutN        = "t_tx_btc.vout_n"
	DBColTTxBtcVoutAddress  = "t_tx_btc.vout_address"
	DBColTTxBtcVoutValue    = "t_tx_btc.vout_value"
	DBColTTxBtcCreateTime   = "t_tx_btc.create_time"
	DBColTTxBtcHandleStatus = "t_tx_btc.handle_status"
	DBColTTxBtcHandleMsg    = "t_tx_btc.handle_msg"
	DBColTTxBtcHandleTime   = "t_tx_btc.handle_time"
)

// const TTxBtc short
const (
	DBColShortTTxBtcID           = "id"
	DBColShortTTxBtcProductID    = "product_id"
	DBColShortTTxBtcBlockHash    = "block_hash"
	DBColShortTTxBtcTxID         = "tx_id"
	DBColShortTTxBtcVoutN        = "vout_n"
	DBColShortTTxBtcVoutAddress  = "vout_address"
	DBColShortTTxBtcVoutValue    = "vout_value"
	DBColShortTTxBtcCreateTime   = "create_time"
	DBColShortTTxBtcHandleStatus = "handle_status"
	DBColShortTTxBtcHandleMsg    = "handle_msg"
	DBColShortTTxBtcHandleTime   = "handle_time"
)

// DBColTTxBtcAll 所有字段
var DBColTTxBtcAll = []string{
	"t_tx_btc.id",
	"t_tx_btc.product_id",
	"t_tx_btc.block_hash",
	"t_tx_btc.tx_id",
	"t_tx_btc.vout_n",
	"t_tx_btc.vout_address",
	"t_tx_btc.vout_value",
	"t_tx_btc.create_time",
	"t_tx_btc.handle_status",
	"t_tx_btc.handle_msg",
	"t_tx_btc.handle_time",
}

// 表结构
// DBTTxBtc t_tx_btc
/*
   id,
   product_id,
   block_hash,
   tx_id,
   vout_n,
   vout_address,
   vout_value,
   create_time,
   handle_status,
   handle_msg,
   handle_time
*/
type DBTTxBtc struct {
	ID           int64  `db:"id" json:"id"`
	ProductID    int64  `db:"product_id" json:"product_id"`
	BlockHash    string `db:"block_hash" json:"block_hash"`
	TxID         string `db:"tx_id" json:"tx_id"`
	VoutN        int64  `db:"vout_n" json:"vout_n"`
	VoutAddress  string `db:"vout_address" json:"vout_address"`
	VoutValue    string `db:"vout_value" json:"vout_value"`
	CreateTime   int64  `db:"create_time" json:"create_time"`
	HandleStatus int64  `db:"handle_status" json:"handle_status"`
	HandleMsg    string `db:"handle_msg" json:"handle_msg"`
	HandleTime   int64  `db:"handle_time" json:"handle_time"`
}

// const TTxBtcToken full
const (
	DBColTTxBtcTokenID           = "t_tx_btc_token.id"
	DBColTTxBtcTokenProductID    = "t_tx_btc_token.product_id"
	DBColTTxBtcTokenTokenIndex   = "t_tx_btc_token.token_index"
	DBColTTxBtcTokenTokenSymbol  = "t_tx_btc_token.token_symbol"
	DBColTTxBtcTokenBlockHash    = "t_tx_btc_token.block_hash"
	DBColTTxBtcTokenTxID         = "t_tx_btc_token.tx_id"
	DBColTTxBtcTokenFromAddress  = "t_tx_btc_token.from_address"
	DBColTTxBtcTokenToAddress    = "t_tx_btc_token.to_address"
	DBColTTxBtcTokenValue        = "t_tx_btc_token.value"
	DBColTTxBtcTokenBlocktime    = "t_tx_btc_token.blocktime"
	DBColTTxBtcTokenCreateAt     = "t_tx_btc_token.create_at"
	DBColTTxBtcTokenHandleStatus = "t_tx_btc_token.handle_status"
	DBColTTxBtcTokenHandleMsg    = "t_tx_btc_token.handle_msg"
	DBColTTxBtcTokenHandleAt     = "t_tx_btc_token.handle_at"
	DBColTTxBtcTokenOrgStatus    = "t_tx_btc_token.org_status"
	DBColTTxBtcTokenOrgMsg       = "t_tx_btc_token.org_msg"
	DBColTTxBtcTokenOrgAt        = "t_tx_btc_token.org_at"
)

// const TTxBtcToken short
const (
	DBColShortTTxBtcTokenID           = "id"
	DBColShortTTxBtcTokenProductID    = "product_id"
	DBColShortTTxBtcTokenTokenIndex   = "token_index"
	DBColShortTTxBtcTokenTokenSymbol  = "token_symbol"
	DBColShortTTxBtcTokenBlockHash    = "block_hash"
	DBColShortTTxBtcTokenTxID         = "tx_id"
	DBColShortTTxBtcTokenFromAddress  = "from_address"
	DBColShortTTxBtcTokenToAddress    = "to_address"
	DBColShortTTxBtcTokenValue        = "value"
	DBColShortTTxBtcTokenBlocktime    = "blocktime"
	DBColShortTTxBtcTokenCreateAt     = "create_at"
	DBColShortTTxBtcTokenHandleStatus = "handle_status"
	DBColShortTTxBtcTokenHandleMsg    = "handle_msg"
	DBColShortTTxBtcTokenHandleAt     = "handle_at"
	DBColShortTTxBtcTokenOrgStatus    = "org_status"
	DBColShortTTxBtcTokenOrgMsg       = "org_msg"
	DBColShortTTxBtcTokenOrgAt        = "org_at"
)

// DBColTTxBtcTokenAll 所有字段
var DBColTTxBtcTokenAll = []string{
	"t_tx_btc_token.id",
	"t_tx_btc_token.product_id",
	"t_tx_btc_token.token_index",
	"t_tx_btc_token.token_symbol",
	"t_tx_btc_token.block_hash",
	"t_tx_btc_token.tx_id",
	"t_tx_btc_token.from_address",
	"t_tx_btc_token.to_address",
	"t_tx_btc_token.value",
	"t_tx_btc_token.blocktime",
	"t_tx_btc_token.create_at",
	"t_tx_btc_token.handle_status",
	"t_tx_btc_token.handle_msg",
	"t_tx_btc_token.handle_at",
	"t_tx_btc_token.org_status",
	"t_tx_btc_token.org_msg",
	"t_tx_btc_token.org_at",
}

// 表结构
// DBTTxBtcToken t_tx_btc_token
/*
   id,
   product_id,
   token_index,
   token_symbol,
   block_hash,
   tx_id,
   from_address,
   to_address,
   value,
   blocktime,
   create_at,
   handle_status,
   handle_msg,
   handle_at,
   org_status,
   org_msg,
   org_at
*/
type DBTTxBtcToken struct {
	ID           int64  `db:"id" json:"id"`
	ProductID    int64  `db:"product_id" json:"product_id"`
	TokenIndex   int64  `db:"token_index" json:"token_index"`
	TokenSymbol  string `db:"token_symbol" json:"token_symbol"`
	BlockHash    string `db:"block_hash" json:"block_hash"`
	TxID         string `db:"tx_id" json:"tx_id"`
	FromAddress  string `db:"from_address" json:"from_address"`
	ToAddress    string `db:"to_address" json:"to_address"`
	Value        string `db:"value" json:"value"`
	Blocktime    int64  `db:"blocktime" json:"blocktime"`
	CreateAt     int64  `db:"create_at" json:"create_at"`
	HandleStatus int64  `db:"handle_status" json:"handle_status"`
	HandleMsg    string `db:"handle_msg" json:"handle_msg"`
	HandleAt     int64  `db:"handle_at" json:"handle_at"`
	OrgStatus    int64  `db:"org_status" json:"org_status"`
	OrgMsg       string `db:"org_msg" json:"org_msg"`
	OrgAt        int64  `db:"org_at" json:"org_at"`
}

// const TTxBtcUxto full
const (
	DBColTTxBtcUxtoID           = "t_tx_btc_uxto.id"
	DBColTTxBtcUxtoUxtoType     = "t_tx_btc_uxto.uxto_type"
	DBColTTxBtcUxtoBlockHash    = "t_tx_btc_uxto.block_hash"
	DBColTTxBtcUxtoTxID         = "t_tx_btc_uxto.tx_id"
	DBColTTxBtcUxtoVoutN        = "t_tx_btc_uxto.vout_n"
	DBColTTxBtcUxtoVoutAddress  = "t_tx_btc_uxto.vout_address"
	DBColTTxBtcUxtoVoutValue    = "t_tx_btc_uxto.vout_value"
	DBColTTxBtcUxtoVoutScript   = "t_tx_btc_uxto.vout_script"
	DBColTTxBtcUxtoCreateTime   = "t_tx_btc_uxto.create_time"
	DBColTTxBtcUxtoSpendTxID    = "t_tx_btc_uxto.spend_tx_id"
	DBColTTxBtcUxtoSpendN       = "t_tx_btc_uxto.spend_n"
	DBColTTxBtcUxtoHandleStatus = "t_tx_btc_uxto.handle_status"
	DBColTTxBtcUxtoHandleMsg    = "t_tx_btc_uxto.handle_msg"
	DBColTTxBtcUxtoHandleTime   = "t_tx_btc_uxto.handle_time"
)

// const TTxBtcUxto short
const (
	DBColShortTTxBtcUxtoID           = "id"
	DBColShortTTxBtcUxtoUxtoType     = "uxto_type"
	DBColShortTTxBtcUxtoBlockHash    = "block_hash"
	DBColShortTTxBtcUxtoTxID         = "tx_id"
	DBColShortTTxBtcUxtoVoutN        = "vout_n"
	DBColShortTTxBtcUxtoVoutAddress  = "vout_address"
	DBColShortTTxBtcUxtoVoutValue    = "vout_value"
	DBColShortTTxBtcUxtoVoutScript   = "vout_script"
	DBColShortTTxBtcUxtoCreateTime   = "create_time"
	DBColShortTTxBtcUxtoSpendTxID    = "spend_tx_id"
	DBColShortTTxBtcUxtoSpendN       = "spend_n"
	DBColShortTTxBtcUxtoHandleStatus = "handle_status"
	DBColShortTTxBtcUxtoHandleMsg    = "handle_msg"
	DBColShortTTxBtcUxtoHandleTime   = "handle_time"
)

// DBColTTxBtcUxtoAll 所有字段
var DBColTTxBtcUxtoAll = []string{
	"t_tx_btc_uxto.id",
	"t_tx_btc_uxto.uxto_type",
	"t_tx_btc_uxto.block_hash",
	"t_tx_btc_uxto.tx_id",
	"t_tx_btc_uxto.vout_n",
	"t_tx_btc_uxto.vout_address",
	"t_tx_btc_uxto.vout_value",
	"t_tx_btc_uxto.vout_script",
	"t_tx_btc_uxto.create_time",
	"t_tx_btc_uxto.spend_tx_id",
	"t_tx_btc_uxto.spend_n",
	"t_tx_btc_uxto.handle_status",
	"t_tx_btc_uxto.handle_msg",
	"t_tx_btc_uxto.handle_time",
}

// 表结构
// DBTTxBtcUxto t_tx_btc_uxto
/*
   id,
   uxto_type,
   block_hash,
   tx_id,
   vout_n,
   vout_address,
   vout_value,
   vout_script,
   create_time,
   spend_tx_id,
   spend_n,
   handle_status,
   handle_msg,
   handle_time
*/
type DBTTxBtcUxto struct {
	ID           int64  `db:"id" json:"id"`
	UxtoType     int64  `db:"uxto_type" json:"uxto_type"`
	BlockHash    string `db:"block_hash" json:"block_hash"`
	TxID         string `db:"tx_id" json:"tx_id"`
	VoutN        int64  `db:"vout_n" json:"vout_n"`
	VoutAddress  string `db:"vout_address" json:"vout_address"`
	VoutValue    string `db:"vout_value" json:"vout_value"`
	VoutScript   string `db:"vout_script" json:"vout_script"`
	CreateTime   int64  `db:"create_time" json:"create_time"`
	SpendTxID    string `db:"spend_tx_id" json:"spend_tx_id"`
	SpendN       int64  `db:"spend_n" json:"spend_n"`
	HandleStatus int64  `db:"handle_status" json:"handle_status"`
	HandleMsg    string `db:"handle_msg" json:"handle_msg"`
	HandleTime   int64  `db:"handle_time" json:"handle_time"`
}

// const TTxEos full
const (
	DBColTTxEosID           = "t_tx_eos.id"
	DBColTTxEosProductID    = "t_tx_eos.product_id"
	DBColTTxEosTxHash       = "t_tx_eos.tx_hash"
	DBColTTxEosLogIndex     = "t_tx_eos.log_index"
	DBColTTxEosFromAddress  = "t_tx_eos.from_address"
	DBColTTxEosToAddress    = "t_tx_eos.to_address"
	DBColTTxEosMemo         = "t_tx_eos.memo"
	DBColTTxEosBalanceReal  = "t_tx_eos.balance_real"
	DBColTTxEosCreateAt     = "t_tx_eos.create_at"
	DBColTTxEosHandleStatus = "t_tx_eos.handle_status"
	DBColTTxEosHandleMsg    = "t_tx_eos.handle_msg"
	DBColTTxEosHandleAt     = "t_tx_eos.handle_at"
)

// const TTxEos short
const (
	DBColShortTTxEosID           = "id"
	DBColShortTTxEosProductID    = "product_id"
	DBColShortTTxEosTxHash       = "tx_hash"
	DBColShortTTxEosLogIndex     = "log_index"
	DBColShortTTxEosFromAddress  = "from_address"
	DBColShortTTxEosToAddress    = "to_address"
	DBColShortTTxEosMemo         = "memo"
	DBColShortTTxEosBalanceReal  = "balance_real"
	DBColShortTTxEosCreateAt     = "create_at"
	DBColShortTTxEosHandleStatus = "handle_status"
	DBColShortTTxEosHandleMsg    = "handle_msg"
	DBColShortTTxEosHandleAt     = "handle_at"
)

// DBColTTxEosAll 所有字段
var DBColTTxEosAll = []string{
	"t_tx_eos.id",
	"t_tx_eos.product_id",
	"t_tx_eos.tx_hash",
	"t_tx_eos.log_index",
	"t_tx_eos.from_address",
	"t_tx_eos.to_address",
	"t_tx_eos.memo",
	"t_tx_eos.balance_real",
	"t_tx_eos.create_at",
	"t_tx_eos.handle_status",
	"t_tx_eos.handle_msg",
	"t_tx_eos.handle_at",
}

// 表结构
// DBTTxEos t_tx_eos
/*
   id,
   product_id,
   tx_hash,
   log_index,
   from_address,
   to_address,
   memo,
   balance_real,
   create_at,
   handle_status,
   handle_msg,
   handle_at
*/
type DBTTxEos struct {
	ID           int64  `db:"id" json:"id"`
	ProductID    int64  `db:"product_id" json:"product_id"`
	TxHash       string `db:"tx_hash" json:"tx_hash"`
	LogIndex     int64  `db:"log_index" json:"log_index"`
	FromAddress  string `db:"from_address" json:"from_address"`
	ToAddress    string `db:"to_address" json:"to_address"`
	Memo         string `db:"memo" json:"memo"`
	BalanceReal  string `db:"balance_real" json:"balance_real"`
	CreateAt     int64  `db:"create_at" json:"create_at"`
	HandleStatus int64  `db:"handle_status" json:"handle_status"`
	HandleMsg    string `db:"handle_msg" json:"handle_msg"`
	HandleAt     int64  `db:"handle_at" json:"handle_at"`
}

// const TTxErc20 full
const (
	DBColTTxErc20ID           = "t_tx_erc20.id"
	DBColTTxErc20TokenID      = "t_tx_erc20.token_id"
	DBColTTxErc20ProductID    = "t_tx_erc20.product_id"
	DBColTTxErc20TxID         = "t_tx_erc20.tx_id"         // 交易id
	DBColTTxErc20FromAddress  = "t_tx_erc20.from_address"  // 来源地址
	DBColTTxErc20ToAddress    = "t_tx_erc20.to_address"    // 目标地址
	DBColTTxErc20BalanceReal  = "t_tx_erc20.balance_real"  // 到账金额Ether
	DBColTTxErc20CreateTime   = "t_tx_erc20.create_time"   // 创建时间戳
	DBColTTxErc20HandleStatus = "t_tx_erc20.handle_status" // 处理状态
	DBColTTxErc20HandleMsg    = "t_tx_erc20.handle_msg"    // 处理消息
	DBColTTxErc20HandleTime   = "t_tx_erc20.handle_time"   // 处理时间戳
	DBColTTxErc20OrgStatus    = "t_tx_erc20.org_status"    // 零钱整理状态
	DBColTTxErc20OrgMsg       = "t_tx_erc20.org_msg"       // 零钱整理消息
	DBColTTxErc20OrgTime      = "t_tx_erc20.org_time"      // 零钱整理时间
)

// const TTxErc20 short
const (
	DBColShortTTxErc20ID           = "id"
	DBColShortTTxErc20TokenID      = "token_id"
	DBColShortTTxErc20ProductID    = "product_id"
	DBColShortTTxErc20TxID         = "tx_id"         // 交易id
	DBColShortTTxErc20FromAddress  = "from_address"  // 来源地址
	DBColShortTTxErc20ToAddress    = "to_address"    // 目标地址
	DBColShortTTxErc20BalanceReal  = "balance_real"  // 到账金额Ether
	DBColShortTTxErc20CreateTime   = "create_time"   // 创建时间戳
	DBColShortTTxErc20HandleStatus = "handle_status" // 处理状态
	DBColShortTTxErc20HandleMsg    = "handle_msg"    // 处理消息
	DBColShortTTxErc20HandleTime   = "handle_time"   // 处理时间戳
	DBColShortTTxErc20OrgStatus    = "org_status"    // 零钱整理状态
	DBColShortTTxErc20OrgMsg       = "org_msg"       // 零钱整理消息
	DBColShortTTxErc20OrgTime      = "org_time"      // 零钱整理时间
)

// DBColTTxErc20All 所有字段
var DBColTTxErc20All = []string{
	"t_tx_erc20.id",
	"t_tx_erc20.token_id",
	"t_tx_erc20.product_id",
	"t_tx_erc20.tx_id",
	"t_tx_erc20.from_address",
	"t_tx_erc20.to_address",
	"t_tx_erc20.balance_real",
	"t_tx_erc20.create_time",
	"t_tx_erc20.handle_status",
	"t_tx_erc20.handle_msg",
	"t_tx_erc20.handle_time",
	"t_tx_erc20.org_status",
	"t_tx_erc20.org_msg",
	"t_tx_erc20.org_time",
}

// 表结构
// DBTTxErc20 t_tx_erc20
/*
   id,
   token_id,
   product_id,
   tx_id,
   from_address,
   to_address,
   balance_real,
   create_time,
   handle_status,
   handle_msg,
   handle_time,
   org_status,
   org_msg,
   org_time
*/
type DBTTxErc20 struct {
	ID           int64  `db:"id" json:"id"`
	TokenID      int64  `db:"token_id" json:"token_id"`
	ProductID    int64  `db:"product_id" json:"product_id"`
	TxID         string `db:"tx_id" json:"tx_id"`                 // 交易id
	FromAddress  string `db:"from_address" json:"from_address"`   // 来源地址
	ToAddress    string `db:"to_address" json:"to_address"`       // 目标地址
	BalanceReal  string `db:"balance_real" json:"balance_real"`   // 到账金额Ether
	CreateTime   int64  `db:"create_time" json:"create_time"`     // 创建时间戳
	HandleStatus int64  `db:"handle_status" json:"handle_status"` // 处理状态
	HandleMsg    string `db:"handle_msg" json:"handle_msg"`       // 处理消息
	HandleTime   int64  `db:"handle_time" json:"handle_time"`     // 处理时间戳
	OrgStatus    int64  `db:"org_status" json:"org_status"`       // 零钱整理状态
	OrgMsg       string `db:"org_msg" json:"org_msg"`             // 零钱整理消息
	OrgTime      int64  `db:"org_time" json:"org_time"`           // 零钱整理时间
}

// const TWithdraw full
const (
	DBColTWithdrawID           = "t_withdraw.id"
	DBColTWithdrawProductID    = "t_withdraw.product_id" // 产品id
	DBColTWithdrawOutSerial    = "t_withdraw.out_serial" // 提币唯一标示
	DBColTWithdrawToAddress    = "t_withdraw.to_address" // 提币地址
	DBColTWithdrawMemo         = "t_withdraw.memo"
	DBColTWithdrawSymbol       = "t_withdraw.symbol"
	DBColTWithdrawBalanceReal  = "t_withdraw.balance_real"  // 提币金额
	DBColTWithdrawTxHash       = "t_withdraw.tx_hash"       // 提币tx hash
	DBColTWithdrawCreateTime   = "t_withdraw.create_time"   // 创建时间
	DBColTWithdrawHandleStatus = "t_withdraw.handle_status" // 处理状态
	DBColTWithdrawHandleMsg    = "t_withdraw.handle_msg"    // 处理消息
	DBColTWithdrawHandleTime   = "t_withdraw.handle_time"   // 处理时间
)

// const TWithdraw short
const (
	DBColShortTWithdrawID           = "id"
	DBColShortTWithdrawProductID    = "product_id" // 产品id
	DBColShortTWithdrawOutSerial    = "out_serial" // 提币唯一标示
	DBColShortTWithdrawToAddress    = "to_address" // 提币地址
	DBColShortTWithdrawMemo         = "memo"
	DBColShortTWithdrawSymbol       = "symbol"
	DBColShortTWithdrawBalanceReal  = "balance_real"  // 提币金额
	DBColShortTWithdrawTxHash       = "tx_hash"       // 提币tx hash
	DBColShortTWithdrawCreateTime   = "create_time"   // 创建时间
	DBColShortTWithdrawHandleStatus = "handle_status" // 处理状态
	DBColShortTWithdrawHandleMsg    = "handle_msg"    // 处理消息
	DBColShortTWithdrawHandleTime   = "handle_time"   // 处理时间
)

// DBColTWithdrawAll 所有字段
var DBColTWithdrawAll = []string{
	"t_withdraw.id",
	"t_withdraw.product_id",
	"t_withdraw.out_serial",
	"t_withdraw.to_address",
	"t_withdraw.memo",
	"t_withdraw.symbol",
	"t_withdraw.balance_real",
	"t_withdraw.tx_hash",
	"t_withdraw.create_time",
	"t_withdraw.handle_status",
	"t_withdraw.handle_msg",
	"t_withdraw.handle_time",
}

// 表结构
// DBTWithdraw t_withdraw
/*
   id,
   product_id,
   out_serial,
   to_address,
   memo,
   symbol,
   balance_real,
   tx_hash,
   create_time,
   handle_status,
   handle_msg,
   handle_time
*/
type DBTWithdraw struct {
	ID           int64  `db:"id" json:"id"`
	ProductID    int64  `db:"product_id" json:"product_id"` // 产品id
	OutSerial    string `db:"out_serial" json:"out_serial"` // 提币唯一标示
	ToAddress    string `db:"to_address" json:"to_address"` // 提币地址
	Memo         string `db:"memo" json:"memo"`
	Symbol       string `db:"symbol" json:"symbol"`
	BalanceReal  string `db:"balance_real" json:"balance_real"`   // 提币金额
	TxHash       string `db:"tx_hash" json:"tx_hash"`             // 提币tx hash
	CreateTime   int64  `db:"create_time" json:"create_time"`     // 创建时间
	HandleStatus int64  `db:"handle_status" json:"handle_status"` // 处理状态
	HandleMsg    string `db:"handle_msg" json:"handle_msg"`       // 处理消息
	HandleTime   int64  `db:"handle_time" json:"handle_time"`     // 处理时间
}
