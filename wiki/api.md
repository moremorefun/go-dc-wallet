# go-dc-wallet 接口使用文档

## 目录

- [go-dc-wallet 接口使用文档](#go-dc-wallet-接口使用文档)
  - [目录](#目录)
  - [注意事项](#注意事项)
  - [签名规则](#签名规则)
  - [错误列表](#错误列表)
  - [接口列表](#接口列表)
    - [从地址池获取地址](#从地址池获取地址)
    - [申请提币](#申请提币)
  - [回调列表](#回调列表)
    - [充币到账通知](#充币到账通知)
    - [提币处理通知](#提币处理通知)

## 注意事项

1. app_name 和 key 在数据表`t_product`中配置,分别对应其中的字段为`app_name`和`app_sk`
2. API接口地址为api接口对外服务的地址,对应的代码入口文件为`cmd/api/main.go`
3. 接口nonce不可重复（可以使用uuid生成），重复将返回错误
4. 回调必须返回"Content-Type":"application/json"类型的数据，数据必须包含error字段，否则将以每两分钟的间隔重复发送通知，以避免通知遗漏
5. 由于需要做零钱整理，所以对不同币种需要做最低入账金额处理，在平台通知到应用的时候，请判断充币金额是否达到入账额度
6. 由于转账需要手续费，平台并不知道应用的手续费设置，请在发送提币时将提币金额减去手续费发送，平台将按照接口数额直接打币，不考虑手续费扣除

## 签名规则

签名生成的通用步骤如下：

1. 设所有发送或者接收到的数据为集合M，将集合M内非空参数值的参数按照参数名ASCII码从小到大排序（字典序），使用URL键值对的格式（即key1=value1&key2=value2…）拼接成字符串stringA。

    特别注意以下重要规则：
   * 参数名ASCII码从小到大排序（字典序）；
   * 参数名区分大小写；
   * 验证调用返回或主动通知签名时，传送的sign参数不参与签名，将生成的签名与该sign值作校验。
   * 接口可能增加字段，验证签名时必须支持增加的扩展字段
2. 在stringA最后拼接上key得到stringSignTemp字符串，并对stringSignTemp进行MD5运算，再将得到的字符串所有字符转换为大写，得到sign值signValue。
3. 举例：
    
    假设传送的参数如下：
    ```
    app_name: app_dc_client
    nonce: ibuaiVcKdpRxkhJA
    ```
    1. 对参数按照key=value的格式，并按照参数名ASCII字典序排序如下：
        `stringA="app_name=app_dc_client&nonce=ibuaiVcKdpRxkhJA"`
    2. 拼接API密钥：
        `stringSignTemp=stringA+"&key=192006250b4c09247ec02edce69f6a2d"` // 注：key为平台设置的密钥key
        
        `sign=MD5(stringSignTemp).toUpperCase()="D35F0447629711EE378F2FE1D26AB43C"` // 注：MD5签名方式



4. 接口签名校验工具
    [https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=20_1](https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=20_1)


## 错误列表

```golang
// ErrorSuccess 成功
ErrorSuccess = 0
ErrorSuccessMsg = "success"

// ErrorInternal 内部错误
ErrorInternal = -1
ErrorInternalMsg = "internal"

// ErrorBind 输入绑定错误
ErrorBind = -2
ErrorBindMsg = "input bind"

// ErrorNoProduct 没有该应用接口信息
ErrorNoProduct    = -3
ErrorNoProductMsg = "no product"

// ErrorIPLimit IP不符合要求
ErrorIPLimit    = -4
ErrorIPLimitMsg = "ip limit"

// ErrorSignWrong 签名错误
ErrorSignWrong    = -5
ErrorSignWrongMsg = "sign wrong"

// ErrorNonceRepeat nonce重复
ErrorNonceRepeat    = -6
ErrorNonceRepeatMsg = "nonce repeat"

// ErrorNoFreeAddress 没有剩余可用地址
ErrorNoFreeAddress    = -7
ErrorNoFreeAddressMsg = "no free address"

// ErrorAddressWrong 提币地址格式错误
ErrorAddressWrong    = -8
ErrorAddressWrongMsg = "address error"

// ErrorBalanceFormat 提币金额格式错误
ErrorBalanceFormat    = -9
ErrorBalanceFormatMsg = "balance format error"

// ErrorSymbolNotSupport 提币币种不支持
ErrorSymbolNotSupport    = -10
ErrorSymbolNotSupportMsg = "symbol not support"
```

## 接口列表

### 从地址池获取地址
```
/api/address

输入参数
POST "Content-Type":"application/json"
{
    // 币种 可选 [eth,btc,eos]
    "symbol": "eth",
	"app_name": "app_dc_client",
	"nonce":"ibuaiVcKdpRxkhJA",
	"sign":"XXXXXX"
}

输出参数
"Content-Type":"application/json"

成功返回
{
    "address": "0x48fbf3e686751cdd363225e3698daac4469e47d9",
    "eos_address": "",
    "error": 0,
    "error_msg": "success"
}
失败返回
{
    "error": -7,
    "error_msg": "no free address"
}
```

### 申请提币
```
/api/withdraw

输入参数
POST "Content-Type":"application/json"
{
    // 提币币种
    "symbol": "eth",
    // 商户订单号
    "out_serial": "7cfd51a2cc0d4e22aac842201eb695f2",
    // 提币地址
    "address": "0x4cd457c0a2ad63198c2da0ce1ba6a7823ffafed9",
    // 提币金额
    "balance": "0.01",
    // eos 提币memo
    "memo": "eos memo",
	"app_name": "app_dc_client",
	"nonce":"ibuaiVcKdpRxkhJB",
	"sign":"XXXXXX"
}

输出参数
"Content-Type":"application/json"

成功返回
{
    "error": 0,
    "error_msg": "success"
}
失败返回
{
    "error": -8,
    "error_msg": "address error"
}
```

## 回调列表

回调地址在数据表`t_product`中配置,对应其中的字段为`cb_url`

```
// 通知类型
const (
    // 冲币到账通知
	NotifyTypeTx              = 1
	// 提币已广播通知
    NotifyTypeWithdrawSend    = 2
	// 提币到账通知
    NotifyTypeWithdrawConfirm = 3
)
```

### 充币到账通知
```
输入参数
POST "Content-Type":"application/json"
{
    // 到账唯一标示，请确保同一tx_hash不会重复入账
    "tx_hash": "0x2be332373700ff87fe6ae2ec2777139ba6b655f49e8b9c0b354a30c52f71a097",
    // 请确保与自己的id是否相同
    "app_name": "app_dc_client",
    // 请务必对签名进行检测，避免攻击者伪造入账通知	
    "sign": "A070E36E9FB0C05DEFB49BA053068912",
    // 充币地址
    "address": "0x09370e3d54ebcb0ff8a399ab3975b74f74cab304",
    // 充币金额
    "balance": "100.100000000000000000",
    // 代币类型
    "symbol": "eth",	
    // 通知类型	NotifyTypeTx
    "notify_type":1
}

输出参数
POST "Content-Type":"application/json"
{
    // 0:  通知处理成功; 非0: 通知处理失败，但不需要再次发送通知
    "error": 0,
    // 如果回复中没有error字段，将重复发送通知
}
```

### 提币处理通知
```
输入参数
POST "Content-Type":"application/json"
{
    // 提币交易hash值
    "tx_hash": "0x9b9632a8509f38e080745cf7713619c62fa4df5e8f98886081bedfd90e209fb2",
    // 提币金额
    "balance": "1.1",
    // 请确保与自己的id是否相同
    "app_name": "app_dc_client",
    // 提币商户流水号，于提币请求中发送的字段对应
    "out_serial": "111666222",
    // 提币地址
    "address": "0xded99b580328671e77be756280d3b070bd371bae",
    // 请务必对签名进行检测，避免攻击者伪造通知
    "sign": "0D1EA3382D937DA292A1F771C0087A9F",
    // 代币类型，小写
    "symbol": "eth",
    // 通知类型 NotifyTypeWithdrawSend | NotifyTypeWithdrawConfirm
    "notify_type": 2,
}

输出参数
POST "Content-Type":"application/json"
{
    // 0:  通知处理成功; 非0: 通知处理失败，但不需要再次发送通知
    "error": 0,	
    // 如果回复中没有error字段，将重复发送通知
}
```


