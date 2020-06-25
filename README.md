# go-dc-wallet

## 内容列表

- [背景](背景)
- [项目依赖](项目依赖)
- [使用说明](使用说明)
- [维护者](维护者)
- [使用许可](使用许可)

## 背景

很多数字货币相关的项目需要收提币的功能,这里提供了一个用于收提币服务的项目.目前支持的币种有

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

### 配置环境变量

### 初始化数据库

### 运行定时任务

### 运行API服务接口
   
## 维护者

[@moremorefun](https://github.com/moremorefun)
[那些年我们De过的Bug](https://www.jidangeng.com)

## 使用许可

[MIT](LICENSE) © moremorefun

这是一个交易所收提币功能的项目。

`.env-example`为系统的环境变量文件，程序的配置从系统变量读取，读取的内容从`app/env.go`中可以查看用途。该文件为示例文件，实际实际使用中，请拷贝为`.env`并根据自身实际情况进行修改。

`init`文件夹中为mysql数据库初始化文件，由于当前版本更新可能较频繁，所以目前只保留最新的数据格式和数据。

`cmd`文件夹下为每个小功能的执行入口，具体的小功能说明请参见[那些年我们De过的Bug](https://www.jidangeng.com)