# xcaddy kill的信号插件

监听caddy被kill的信号,写入到文件,使用场景: 当docker内收到kill信号之后不再接受服务,配合容器健康检查一并使用

- 安装插件

```shell
go install github.com/caddyserver/xcaddy/cmd/xcaddy@latest
xcaddy build --with github.com/jcleng/caddy-pl@v0.1.0
```

- caddyfile配置

```caddy
{
  order pl last
}
:9199 {
    pl {
        shutdown_file "D:\work\go_test2\2.txt"
    }
}

```

- 开发

```shell
go mod tidy

xcaddy run --config .\Caddyfile
# 使用ctrl-c停止xcaddy服务,查看写入文件,如果文件存在会返回500,所以是配置docker进行使用
```
