package app

// 交易状态
const (
	TxStatusInit   = 0
	TxStatusNotify = 1
)

// 零钱整理状态
const (
	TxOrgStatusInit       = 0
	TxOrgStatusHex        = 1
	TxOrgStatusSend       = 2
	TxOrgStatusConfirm    = 3
	TxOrgStatusFeeHex     = 4
	TxOrgStatusFeeSend    = 5
	TxOrgStatusFeeConfirm = 6
)

// 发送状态
const (
	SendStatusInit    = 0
	SendStatusSend    = 1
	SendStatusConfirm = 2
)

// 发送类型
const (
	SendRelationTypeTx         = 1
	SendRelationTypeWithdraw   = 2
	SendRelationTypeTxErc20    = 3
	SendRelationTypeTxErc20Fee = 4
	SendRelationTypeUXTOOrg    = 5
	SendRelationTypeOmniOrg    = 6
)

// 通知状态
const (
	NotifyStatusInit = 0
	NotifyStatusFail = 1
	NotifyStatusPass = 2
)

// 通知类型
const (
	NotifyTypeTx              = 1
	NotifyTypeWithdrawSend    = 2
	NotifyTypeWithdrawConfirm = 3
)

// 提币状态
const (
	WithdrawStatusInit    = 0
	WithdrawStatusHex     = 1
	WithdrawStatusSend    = 2
	WithdrawStatusConfirm = 3
)

// uxto 类型
const (
	UxtoTypeTx      = 1
	UxtoTypeHot     = 2
	UxtoTypeOmni    = 3
	UxtoTypeOmniHot = 4
)

// uxto 处理类型
const (
	UxtoHandleStatusInit    = 0
	UxtoHandleStatusUse     = 1
	UxtoHandleStatusConfirm = 2
)
