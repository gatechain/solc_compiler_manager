# 提供 solc 下载管理，提供合约验证 RPC 服务

## 安装
```
make install
```

## 项目初始化
```
solc-compiler init --platform macosx-amd64/linux-amd64 -a
```

初始化需要指定项目运行平台，初始化可以指定是否下载运行平台所有的 `solc` 版本。

可以重复初始化，以覆盖当前配置，初始化过程中出现文件下载失败的情况，可以：
```
// 重试初始化，该过程会自动跳过已下载的编译器版本
solc-compiler init --platform macosx-amd64 -a

// 指定版本，使用可执行程序手动获取
solc-compiler fetch v0.8.0
```

## 管理编译器版本
```
// 下载
solc-compiler fetch v0.8.0

// 删除
solc-compiler delete v0.8.0
```

## 编译 solidity 文件
```
solc-compiler compile --scope bin,abi,hashes --name Counter v0.5.11 [path_to_file] --optimize --optimize-runs 200 --evm-version "default"
--scope         指定编译的输出内容
--name          指定编译合约
--optimize      指定编译是否进行优化
--optimize-runs 指定编译油画参数
--evm-version   指定evm版本
详细参数请参考帮助文档
```

## 启动 RPC 服务
```
solc-compiler rest-server
提供以下接口：

contract_ping: 心跳

curl --location --request POST 'http://127.0.0.1:1212' \
--header 'Content-Type: application/json' \
--data-raw '{
    "jsonrpc": "2.0",
    "method": "contract_ping",
    "params": [],
    "id": 1
}'

contract_verify: 返回 ABI、编译的binary hex string

curl --location --request POST 'http://127.0.0.1:1212' \
--header 'Content-Type: application/json' \
--data-raw '{
    "jsonrpc": "2.0",
    "method": "contract_verify",
    "params": [
        {
            "name": "Test",
            "compiler_version": "v0.5.11+commit.22be8592",
            "code":"pragma solidity ^0.5.11; contract Test {}",
	        "optimize": false,
            "optimize_runs": 0,
            "evm_version": "default"
        }
    ],
    "id": 1
}'

contract_listVersions: 返回当前平台所支持的编译器版本

curl --location --request GET 'http://127.0.0.1:1212' \
--header 'Content-Type: application/json' \
--data-raw '{
    "jsonrpc": "2.0",
    "method": "contract_listVersions",
    "params": [],
    "id": 1
}'
```