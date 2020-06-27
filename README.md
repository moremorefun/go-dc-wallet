# go-dc-wallet

[![build](https://github.com/moremorefun/go-dc-wallet/workflows/build/badge.svg)](https://github.com/moremorefun/go-dc-wallet/actions?query=workflow%3Abuild)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/moremorefun/go-dc-wallet/blob/master/LICENSE)
[![blog](https://img.shields.io/badge/blog-@moremorefun-green.svg)](https://www.jidangeng.com)


## 目录

- [背景](#背景)
- [项目依赖](#项目依赖)
- [使用说明](#使用说明)
  - [拉取代码并获取依赖库](#拉取代码并获取依赖库)
  - [配置环境变量](#配置环境变量)
  - [初始化数据库](#初始化数据库)
  - [初始化基础数据](#初始化基础数据)
  - [手动添加自身需要设置的数据](#手动添加自身需要设置的数据)
  - [运行定时任务](#运行定时任务)
- [运行API服务接口](#运行API服务接口)
- [维护者](#维护者)
- [使用许可](#使用许可)

## 背景

很多加密货币相关的项目需要收提币的功能,这里提供了一个用于收提币服务的项目.目前支持的币种有:

- Ethereum(以太坊)
- Erc20(以太坊代币)
- Bitcoin(比特币)
- OmniLayer(比特币代币)

## 项目依赖

- 项目使用`Golang`编写
- 数据库使用`MySQL`
- `Ethereum`的RPC服务
- `OmniLayer`的RPC服务

## 使用说明

### 拉取代码并获取依赖库

```
# 拉取代吗
git clone https://github.com/moremorefun/go-dc-wallet.git
# 切换到项目目录
cd go-dc-wallet
# 获取依赖
go mod download
```

### 配置环境变量
```
cp .env-example .env
```
根据自身情况编辑 `.env` 文件
```
### 是否是debug,用来设置日志等级
IS-DEBUG=true

### mysql配置
MYSQL=root:123456@tcp(127.0.0.1:3306)/dc-wallet-prod?parseTime=true&loc=Local&charset=utf8mb4&collation=utf8mb4_unicode_ci&tx_isolation=%27READ-COMMITTED%27&sql_mode=%27STRICT_TRANS_TABLES%2cNO_ZERO_IN_DATE%2cERROR_FOR_DIVISION_BY_ZERO%2cNO_AUTO_CREATE_USER%2cNO_ENGINE_SUBSTITUTION%27
# 是否在日志中输出sql请求语句
MYSQL-IS-SHOW-SQL=false

### 私钥加密的key
AES-KEY=123

### eth rpc 接口
ETH_RPC=https://mainnet.infura.io/v3/0b359d2406a6492fb53883d46921d775

### btc rpc 接口
# btc接口类型可选值为 btc 和 btc-test
BTC-NETWORK-TYPE=btc
OMNI_RPC_HOST=http://127.0.0.1:18332
OMNI_RPC_USER=omni
OMNI_RPC_PWD=omni
```

### 初始化数据库

#### 创建配置的数据库,例如上面的`.env`中设置的数据库名为`dc-wallet-prod`.

#### 运行数据库表差异SQL生成工具,生成数据库表结构SQL文件

```
go run cmd/db/main.go
```

确认生成的SQL无问题后在MySQL中执行,生成相应的表格.

#### 初始化基础数据:

```
go run cmd/dbinit/main.go
```

#### 手动添加自身需要设置的数据
```
# eth冷钱包地址
t_app_config_str.cold_wallet_address
# btc冷钱包地址
t_app_config_str.cold_wallet_address_btc

# erc20 token 冷钱包地址
t_app_config_token[].cold_address

# omni token 冷钱包地址
t_app_config_token_btc[].cold_address

# 用于提供api服务的相关数据
t_product
```

### 运行定时任务

```
go run cmd/crontab/main.go
```

### 运行API服务接口

```
go run cmd/api/main.go
```
   
## 维护者

[@moremorefun](https://github.com/moremorefun)
[那些年我们De过的Bug](https://www.jidangeng.com)

## 使用许可

[MIT](LICENSE) © moremorefun

这是一个交易所收提币功能的项目。

`.env-example`为系统的环境变量文件，程序的配置从系统变量读取，读取的内容从`app/env.go`中可以查看用途。该文件为示例文件，实际实际使用中，请拷贝为`.env`并根据自身实际情况进行修改。

`init`文件夹中为mysql数据库初始化文件，由于当前版本更新可能较频繁，所以目前只保留最新的数据格式和数据。

`cmd`文件夹下为每个小功能的执行入口，具体的小功能说明请参见[那些年我们De过的Bug](https://www.jidangeng.com)