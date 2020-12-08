## 两个任务

#### 1. 学习lightsocks

- 参考教程：[你也能写个 Shadowsocks](https://github.com/gwuhaolin/blog/issues/12)

### 2. 学习 `glider` 重写 `socks5` 协议

- 参考 [glider 代码](https://github.com/nadoo/glider/blob/e6e5c3d4b68bcb4d77c4403e8f55288bc1a5ef17/proxy/socks5/client.go#L102)



## 说明

客户端启动，监听端口（ `soocks` 协议），直接对转发的数据，先进行加密，再进行转发。

服务端启动，监听端口（ `soocks` 协议），按照 `socks5` 协议进行解析。解析后请求真正想访问的地址，加密后，转发。

## 运行

- 客户端：

```
go run cmd/socks-local/main.go
```

- 服务端：

```
go run cmd/socks-server/main.go
```

### 测试

浏览器代理使用 `SwitchyOmega`， 参考教程 [搭配 Chrome 使用](https://github.com/gwuhaolin/lightsocks/wiki/%E6%90%AD%E9%85%8D-Chrome-%E4%BD%BF%E7%94%A8).