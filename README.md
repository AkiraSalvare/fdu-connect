# FDU Connect

> 🚫 **免责声明**
>
> 本程序**按原样提供**，作者**不对程序的正确性或可靠性提供保证**，请使用者自行判断具体场景是否适合使用该程序，**使用该程序造成的问题或后果由使用者自行承担**！

---

中文 | [English](README_en.md)

**本程序基于 [ZJU Connect](https://github.com/Mythologyli/zju-connect) 修改，感谢原作者 [Mythologyli](https://github.com/Mythologyli) 以及原项目的贡献者。**

### 使用方法

+ 如果你是来自复旦大学的用户：

  1. 在 [Release](https://github.com/akirasalvare/fdu-connect/releases) 页面下载对应平台的最新版本。
      - Arch Linux 用户可直接安装 AUR 包 [fdu-connect-git](https://aur.archlinux.org/packages/fdu-connect-git): `[yay|paru] -S fdu-connect-git`

  2. 以 macOS 为例，解压出可执行文件 `fdu-connect`。

  3. macOS 需要先解除安全限制。命令行运行：`sudo xattr -rd com.apple.quarantine fdu-connect`。

  4. 命令行运行：`./fdu-connect -username <上网账户> -password <密码>`。

  5. 此时 `1080` 端口为 Socks5 代理，`1081` 端口为 HTTP 代理。如需更改默认端口，请参考参数说明。

+ 如果你是非复旦大学的用户：

  建议使用原版 [ZJU Connect](https://github.com/Mythologyli/zju-connect)

#### 作为服务运行

[链接](docs/service.md)

#### Docker 运行

[链接](docs/docker.md)

### ⚠️ 警告

1. 当使用其他开启了 TUN 模式的代理工具，同时配合 fdu-connect 作为下游代理时，请注意务必提供正确的分流规则，参考[此 issue](https://github.com/Mythologyli/zju-connect/issues/57)

### ⚠️ TUN 模式注意事项

1. 需要管理员权限运行

2. Windows 系统需要前往 [Wintun 官网](https://www.wintun.net)下载 `wintun.dll` 并放置于可执行文件同目录下

3. 为保证域名解析正确，建议配置 `dns-hijack` 劫持系统 DNS

### 参数说明

+ `server`: SSL VPN 服务端地址，默认为 `stuvpn.fudan.edu.cn`

+ `port`: SSL VPN 服务端端口，默认为 `443`

+ `username`: 网络账户。例如：学号

+ `password`: 网络账户密码

+ `totp-secret`: TOTP 密钥，可用于自动完成 TOTP 验证。如服务端无需 TOTP 验证或希望手动输入验证码，可不填

+ `cert-file`: p12 证书文件路径，如果服务器要求证书验证，需要配置此参数

+ `cert-password`: 证书密码

+ `disable-server-config`: 禁用服务端配置，一般不需要加此参数

+ `skip-domain-resource`: 不使用服务端下发的域名资源分流，一般不需要加此参数

+ `disable-fdu-config`: 禁用 FDU 相关配置，一般不需要加此参数

+ `disable-fdu-dns`: 禁用 FDU DNS 改用本地 DNS，一般不需要加此参数

+ `disable-multi-line`: 禁用自动根据延时选择线路。加此参数后，使用 `server` 参数指定的线路

+ `proxy-all`: 是否代理所有流量，一般不需要加此参数

+ `socks-bind`: SOCKS5 代理监听地址，默认为 `:1080`

+ `socks-user`: SOCKS5 代理用户名，不填则不需要认证

+ `socks-passwd`: SOCKS5 代理密码，不填则不需要认证

+ `http-bind`: HTTP 代理监听地址，默认为 `:1081`。为 `""` 时不启用 HTTP 代理

+ `shadowsocks-url`: Shadowsocks 服务端 URL。例如：`ss://aes-128-gcm:password@server:port`。格式[参考此处](https://github.com/shadowsocks/go-shadowsocks2)

+ `dial-direct-proxy`: 当URL未命中RVPN规则，切换到直连时使用代理，常用于与其他代理工具配合的场景，目前仅支持http代理。 例如：`http://127.0.0.1:7890"`，为 `""` 时不启用

+ `tun-mode`: TUN 模式（实验性）。请阅读后文中的 TUN 模式注意事项

+ `add-route`: 启用 TUN 模式时根据服务端下发配置添加路由

+ `dns-ttl`: DNS 缓存时间，默认为 `3600` 秒

+ `disable-keep-alive`: 禁用定时保活，一般不需要加此参数

+ `fdu-dns-server`: FDU DNS 服务器地址，默认为 `202.120.224.26`。设置为 auto 时使用从服务端获取的 DNS 服务器，如果未能获取则禁用 FDU DNS

+ `secondary-dns-server`: 当使用 FDU DNS 服务器无法解析时使用的备用 DNS 服务器，默认为 `223.5.5.5`。留空则使用系统默认 DNS，但在开启 `dns-hijack` 时必须设置

+ `dns-server-bind`: DNS 服务器监听地址，默认为空即禁用。例如，设置为 `127.0.0.1:53`，则可向 `127.0.0.1:53` 发起 DNS 请求

+ `dns-hijack`: 启用 TUN 模式时劫持 DNS 请求，建议在启用 TUN 模式时添加此参数

+ `debug-dump`: 是否开启调试，一般不需要加此参数

+ `tcp-port-forwarding`: TCP 端口转发，格式为 `本地地址-远程地址,本地地址-远程地址,...`，例如 `127.0.0.1:9898-10.10.98.98:80,0.0.0.0:9899-10.10.98.98:80`。多个转发用 `,` 分隔

+ `udp-port-forwarding`: UDP 端口转发，格式为 `本地地址-远程地址,本地地址-远程地址,...`，例如 `127.0.0.1:53-202.120.224.26:53`。多个转发用 `,` 分隔

+ `custom-dns`: 指定自定义DNS解析结果，格式为 `域名:IP,域名:IP,...`，例如 `zb.fudan.edu.cn:10.108.68.200,mirrors.fducslg.com:10.176.52.2`。多个解析用 `,` 分隔

+ `custom-proxy-domain`: 指定自定义域名使用RVPN代理，格式为 `域名,域名,...`，例如 `nature.com,science.org`。多个域名用 `,` 分隔

+ `twf-id`: twfID 登录，调试用途，一般不需要加此参数

+ `config`: 指定配置文件，内容参考 `config.toml.example`。启用配置文件时其他参数无效

### 感谢

+ [ZJU Connect](https://github.com/Mythologyli/zju-connect)

+ [EasierConnect](https://github.com/lyc8503/EasierConnect)

+ [socks2http](https://github.com/zenhack/socks2http)

+ [![image](docs/yxvm.png)](https://yxvm.com/)

  [NodeSupport](https://github.com/NodeSeekDev/NodeSupport) 赞助了本项目