# kchain
基于tendermint的区块链


## 项目初始化

1. 请安装task项目管理工具

```
go get -u -v github.com/go-task/task/cmd/task
```

2. 初始化项目

```
task init
```

3. 安装依赖

```
task deps
```
其他依赖请把依赖添加到`scripts/deps.sh中`

4. Goland作为开发工具配置
请把当前目录`pwd`配置为`Project GOPATH`


## 编译

```
task build
```

