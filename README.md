# nft-collect
![](https://img.shields.io/badge/license-MIT-green)
[![goreportcard for backend-go](https://goreportcard.com/badge/github.com/decert-me/nft-collect)](https://goreportcard.com/report/github.com/decert-me/nft-collect)
## 安装
```bash
git clone https://github.com/decert-me/nft-collect.git
```
## 编译
```bash
go build
```
## 配置
```bash
# 主程序配置
cp ./config/config.demo.yaml ./config/config.yaml
vi ./config/config.yaml
```

### 运行配置

配置项：

```
# system configuration
system:
  env: develop
  addr: 8888
  api-key: ""
```

env：运行环境，可选值为 develop、test、production

addr：运行端口

api-key：Zcloak 证书生成调用 API Key

### 数据库配置

配置项：
```yaml
# pgsql configuration
pgsql:
  path: "127.0.0.1"
  port: "5432"
  config: ""
  db-name: ""
  username: "postgres"
  password: "123456"
  auto-migrate: true
  prefix: ""
  slow-threshold: 200
  max-idle-conns: 10
  max-open-conns: 100
  log-mode: "info"
  log-zap: false
```

path：数据库地址

port：数据库端口

config：数据库配置

db-name：数据库名称

username：数据库用户名

password：数据库密码

auto-migrate：是否自动迁移数据库

prefix：数据库表前缀

slow-threshold：慢查询阈值，单位毫秒

max-idle-conns：最大空闲连接数

max-open-conns：最大连接数

log-mode：日志级别

log-zap：是否使用zap日志库

### 日志级别配置

配置项：
```yaml
# log configuration
log:
  level: info
  save: true
  format: console
  log-in-console: true
  prefix: '[backend-go]'
  director: log
  show-line: true
  encode-level: LowercaseColorLevelEncoder
  stacktrace-key: stacktrace
```

level：日志级别 debug、info、warn、error、dpanic、panic、fatal

save：是否保存日志

format：日志格式

log-in-console：是否在控制台输出日志

prefix：日志前缀

director：日志保存路径

show-line：是否显示行号

encode-level：日志编码级别

stacktrace-key：堆栈信息


### JWT 配置

配置保持 decert app 一致

配置项：

```yaml
# auth configuration
auth:
  signing-key: "Decert"
  expires-time: 86400
  issuer: "Decert"
```

signing-key：签名密钥

expires-time：过期时间，单位秒

issuer：签发人

### 文件上传配置

配置项：

```yaml
# local configuration
local:
  path: 'uploads/file'
  ipfs: 'uploads/ipfs'
```

path：本地文件保存路径

ipfs：IPFS文件保存路径

### NFT 配置

配置项：

```yaml
# nft configuration
nft:
  ens-rpc: "https://rpc.ankr.com/eth"
  api-key: ""
  cache-time: 15
  logo-path: "assets"
  api-config:
    - chain: "eth"
      chain-id: 1
      api-per-host: "restapi"
      symbol: "ETH"
```

ens-rpc：获取 ENS RPC 地址

api-key：API Key

cache-time：NFT 数据缓存时间

logo-path：NFT 项目图标保存路径

chain：链名称

chain-id：链 ID

api-per-host：NFTScan API 名称

symbol：链符号


### IPFS 配置

配置项：

```yaml
# ipfs configuration
ipfs:
  url: "https://dweb.link/ipfs"
```

url：IPFS 网关地址

## 运行
```bash
./nft-collect
```
