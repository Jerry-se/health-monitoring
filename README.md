# Health Monitoring

## Build

```shell
go build -ldflags "-X main.version=v0.1.5" -o .hm main.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=v0.1.5" -o .hm main.go
```

## Run

创建一个 JSON 配置文件，示例如下:

```json
{
  "Addr": "0.0.0.0:9521",
  "LogLevel": "info",
  "LogFile": "./test.log",
  "MongoURI": "mongodb://127.0.0.1:27017/",
  "MongoDB": "health_monitoring"
}
```

使用命令 `hm -config ./config.json` 运行即可。

程序会启动一个 WebSocket 服务，可以使用 `ws://localhost:9521/websocket` 连接。

## WebSocket

WebSocket 设置了心跳服务，即 client 发送 ping 消息，服务回复 pong 消息。
如果 30s 内没有任何 ping 消息，连接将被服务端断开。
请及时发送 ping 消息，既是一种心跳，又能保证长连接的稳定可靠。

WebSocket 消息采用 UTF-8 文本格式，主要使用 JSON 形式。具体示例请看 [测试用例](./ws/ws_test.go)

client 向 server 发送的请求消息主要由 Header 和 Body 两部分组成。

<table>
  <tr>
    <td></td>
    <td>字段</td>
    <td>描述</td>
    <td>类型</td>
    <td>备注</td>
  </tr>
  <tr>
    <td rowspan="6">Header</td>
    <td>version</td>
    <td>协议版本，暂时用 0</td>
    <td>uint32</td>
    <td></td>
  </tr>
  <tr>
    <td>timestamp</td>
    <td>时间戳</td>
    <td>int64</td>
    <td></td>
  </tr>
  <tr>
    <td>id</td>
    <td>消息序号，一对请求与应答的序号相同</td>
    <td>uint64</td>
    <td></td>
  </tr>
  <tr>
    <td>type</td>
    <td>消息体的类型，0 - 保留， 1 - online，2 - 机器信息</td>
    <td>uint32</td>
    <td></td>
  </tr>
  <tr>
    <td>pub_key</td>
    <td>公钥，验证消息安全完整，暂时不需要</td>
    <td>[]byte</td>
    <td></td>
  </tr>
  <tr>
    <td>sign</td>
    <td>签名，验证消息安全完整，暂时不需要</td>
    <td>[]byte</td>
    <td></td>
  </tr>
  <tr>
    <td>Body</td>
    <td>body</td>
    <td>消息体，真正的消息通过 JSON 编码，加密后的字节数组</td>
    <td>[]byte</td>
    <td></td>
  </tr>
</table>

消息体暂时有以下几种:
- 0 - 没有意义
- 1 - Online，表示 WebSocket 连接属于那个设备或者节点。
```json
{
  "node_id": "123456789"
}
```
- 2 - 设备信息，定时发送的模型和显卡使用信息。
```json
{
  "project": "DecentralGPT",
  "models": [
    {
      "model": "Codestral-22B-v0.1"
    }
  ],
  "gpu_name": "NVIDIA RTX A5000",
  "utilization_gpu": 30,
  "memory_total": 24564,
  "memory_used": 22128
}
```

server 向 client 返回的应答消息体格式结构相似，只比请求多了 Code 和 Message 两个字段。

<table>
  <tr>
    <td></td>
    <td>字段</td>
    <td>描述</td>
    <td>类型</td>
    <td>备注</td>
  </tr>
  <tr>
    <td rowspan="6">Header</td>
    <td>version</td>
    <td>协议版本，暂时用 0</td>
    <td>uint32</td>
    <td></td>
  </tr>
  <tr>
    <td>timestamp</td>
    <td>时间戳</td>
    <td>int64</td>
    <td></td>
  </tr>
  <tr>
    <td>id</td>
    <td>消息序号，一对请求与应答的序号相同</td>
    <td>uint64</td>
    <td></td>
  </tr>
  <tr>
    <td>type</td>
    <td>消息体的类型，与请求的类型相同</td>
    <td>uint32</td>
    <td></td>
  </tr>
  <tr>
    <td>pub_key</td>
    <td>公钥，验证消息安全完整，暂时不需要</td>
    <td>[]byte</td>
    <td></td>
  </tr>
  <tr>
    <td>sign</td>
    <td>签名，验证消息安全完整，暂时不需要</td>
    <td>[]byte</td>
    <td></td>
  </tr>
  <tr>
    <td>Code</td>
    <td>code</td>
    <td>错误码，0 表示正常</td>
    <td>uint32</td>
    <td></td>
  </tr>
  <tr>
    <td>Message</td>
    <td>message</td>
    <td>错误信息</td>
    <td>string</td>
    <td></td>
  </tr>
  <tr>
    <td>Body</td>
    <td>body</td>
    <td>消息体，真正的消息通过 JSON 编码，加密后的字节数组</td>
    <td>[]byte</td>
    <td></td>
  </tr>
</table>

