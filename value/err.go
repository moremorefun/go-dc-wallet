package value

const (
	// ErrorSuccess 成功
	ErrorSuccess = 200
	// ErrorSuccessMsg 成功
	ErrorSuccessMsg = "success"

	// ErrorInternal 内部错误
	ErrorInternal = -1
	// ErrorInternalMsg 内部错误
	ErrorInternalMsg = "internal"

	// ErrorBind 输入绑定错误
	ErrorBind = -2
	// ErrorBindMsg 输入绑定错误
	ErrorBindMsg = "input bind"

	ErrorNoProduct    = -3
	ErrorNoProductMsg = "no product"

	ErrorIPLimit    = -4
	ErrorIPLimitMsg = "ip limit"

	ErrorSignWrong    = -5
	ErrorSignWrongMsg = "sign wrong"

	ErrorNonceRepeat    = -6
	ErrorNonceRepeatMsg = "nonce repeat"

	ErrorNoFreeAddress    = -7
	ErrorNoFreeAddressMsg = "no free address"

	ErrorAddressWrong    = -8
	ErrorAddressWrongMsg = "address error"

	ErrorBalanceFormat    = -9
	ErrorBalanceFormatMsg = "balance format error"

	ErrorSymbolNotSupport    = -10
	ErrorSymbolNotSupportMsg = "symbol not support"
)
