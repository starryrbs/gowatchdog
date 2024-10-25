# watchdog
watchdog 是一个轻量级的服务监控工具，使用 Go 语言编写。它旨在测试各种服务、中间件的容灾能力，包括 Zookeeper、MySQL、Redis 、Kafka等。

## 特性

- 定时连接和发送心跳，检测 Zookeeper 的状态。
- 可扩展支持 MySQL、Redis 和 Kafka 的健康检查。
- 可配置的参数，方便用户自定义监控频率和连接信息。

## 安装

### 前提条件

- Go 1.21.5 及以上版本


### 配置文件示例

当前目录下新建config.toml
```yaml
redis:
  addr: "127.0.0.1:6379"     # Redis 服务器地址
  password: "123456"               # Redis 密码（如果没有则为空）
  db: 0                      # Redis 数据库编号
  timeout: 5s                # 超时时间，单位为秒
zookeeper:
  hosts:
    - "127.0.0.1:21811"      # Zookeeper 主机地址和端口
    - "127.0.0.1:21812"      # 如果有多个节点，可以依次列出
    - "127.0.0.1:21813"      # 如果有多个节点，可以依次列出
  timeout: 5s       # 会话超时时间，单位为秒
interval: 100ms
```

### 运行

```shell
go run cmd/main.go
```

可以模拟关闭zk集群中的一个节点来测试容灾场景下zk集群的高可用

## 贡献
欢迎贡献！请遵循以下步骤：

1. Fork 本项目
2. 创建你的特性分支 (git checkout -b feature/YourFeature)
3. 提交更改 (git commit -m 'Add some feature')
4. 推送到分支 (git push origin feature/YourFeature)
5. 创建 Pull Request

## 许可证
此项目遵循 MIT 许可证。有关详细信息，请查看 LICENSE 文件。

