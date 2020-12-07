## 两个任务

#### 1. 学习lightsocks

- 参考教程：[你也能写个 Shadowsocks](https://github.com/gwuhaolin/blog/issues/12)

### 2. 学习 `glider` 重写 `socks5` 协议

- 参考 [glider 代码](https://github.com/nadoo/glider/blob/e6e5c3d4b68bcb4d77c4403e8f55288bc1a5ef17/proxy/socks5/client.go#L102)



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