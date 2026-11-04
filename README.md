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

3. EasyConnect 协议下的推荐配置为 `-tun-mode -add-route -dns-hijack`

4. aTrust 协议下的推荐配置为 `-tun-mode -add-route -dns-hijack -fake-ip`。在使用 aTrust 协议时，如果不使用 DNS 劫持/Fake IP，直接通过 TUN 网卡的涉及域名的 TCP 流量可能会出错

### 参数说明

#### 通用参数

+ `protocol`: 登录协议，支持 `easyconnect`/`atrust`，默认为 `easyconnect`

+ `server`: VPN 服务端地址，默认为 `stuvpn.fudan.edu.cn`

+ `port`: VPN 服务端端口，默认为 `443`

+ `username`: 网络账户。例如：学号

+ `password`: 网络账户密码（敏感信息建议使用环境变量 `FDU_CONNECT_PASSWORD` 给定）

+ `graph-code-file`: 图形验证码文件路径。默认为空。在 aTrust 模式下，留空时使用浏览器完成验证码，设置路径则登录时会将图形验证码保存至该文件，由用户手动输入 JSON

+ `totp-secret`: TOTP 密钥，可用于自动完成 TOTP 验证。如服务端无需 TOTP 验证或希望手动输入验证码，可不填（敏感信息建议使用环境变量 `FDU_CONNECT_TOTP_SECRET` 给定）

+ `disable-remote-dns`: 禁用远端 DNS 改用本地 DNS，一般不需要加此参数

+ `disable-server-config`: 禁用服务端配置，一般不需要加此参数

+ `socks-bind`: SOCKS5 代理监听地址，默认为 `:1080`

+ `socks-user`: SOCKS5 代理用户名，不填则不需要认证

+ `socks-passwd`: SOCKS5 代理密码，不填则不需要认证

+ `http-bind`: HTTP 代理监听地址，默认为 `:1081`。为 `""` 时不启用 HTTP 代理

+ `shadowsocks-url`: Shadowsocks 服务端 URL。例如：`ss://aes-128-gcm:password@server:port`。格式[参考此处](https://github.com/shadowsocks/go-shadowsocks2)

+ `dial-direct-proxy`: 当 URL 未命中规则，切换到直连时使用代理，常用于与其他代理工具配合的场景，目前仅支持 http 代理。例如：`http://127.0.0.1:7890"`，为 `""` 时不启用

+ `tcp-tunnel-mode`: TCP 隧道模式，默认为 `false`。启用后仅可通过 TCP 隧道代理 TCP 流量。由于只有 aTrust 支持 TCP 隧道，此模式在 EasyConnect 下无效。启用后会禁用 TUN 模式

+ `tun-mode`: TUN 模式（实验性）。请阅读 TUN 模式注意事项

+ `add-route`: 启用 TUN 模式时根据服务端下发配置添加路由

+ `dns-ttl`: DNS 缓存时间，默认为 `3600` 秒

+ `disable-keep-alive`: 禁用定时保活，一般不需要加此参数

+ `keep-alive-url`: 使用 HTTP 保活，适用于远端 DNS 不可用的情况。填写要访问的 URL，例如 `https://www.cnki.net/favicon.ico` 。默认为空，此时使用远端 DNS 保活

+ `remote-dns-server`: 远端 DNS 服务器地址，默认为 `202.120.224.26`。设置为 `auto` 时使用从服务端获取的 DNS 服务器，如果未能获取则禁用远端 DNS

+ `secondary-dns-server`: 当远端 DNS 无法解析时使用的备用服务器。默认值 `auto` 优先采用 VPN 策略下发的第二 DNS，否则回退到 `223.5.5.5`。留空则使用系统默认 DNS，但在开启 `dns-hijack` 时必须设置

+ `dns-server-bind`: DNS 服务器监听地址，默认为空即禁用。例如，设置为 `127.0.0.1:53`，则可向 `127.0.0.1:53` 发起 DNS 请求

+ `local-dns-server`: 指定用于解析 VPN 服务器域名的本地 DNS，格式为 IP 或 IP:port；留空时使用系统 DNS，可路由的 DNS 地址在探测成功后绑定到底层网卡，本地 DNS stub 保持 loopback 路由

+ `dns-hijack`: 启用 TUN 模式时劫持 DNS 请求，建议在启用 TUN 模式时添加此参数

+ `fake-ip`: 启用 Fake IP 功能，与 dns-hijack 配合使用，建议在使用 aTrust 协议并启用 TUN 模式时添加此参数。此参数在 EasyConnect 协议下无效

+ `debug-dump`: 是否开启调试，一般不需要加此参数

+ `debug-pcap-file`: 根据 VPN 底层 TCP 连接实际收发的数据重建 PCAP 文件，仅用于调试；捕获队列满时会阻塞网络读写，不包含内核握手和重传，TLS 内容仍为密文

+ `debug-tls-log-file`: 将 TLS 会话密钥导出为 NSS key log 格式，可配合 `debug-pcap-file` 在 Wireshark 中解密 TLS 流量。该文件包含会话密钥，仅用于调试并应妥善保管

+ `bind-interface`: 手动指定 VPN 底层连接使用的网卡接口，支持 EasyConnect 和 aTrust。非空时优先使用该接口，不再自动探测

+ `auto-detect-interface`: 自动探测并绑定 VPN 底层网卡，默认为 `false`。设为 `true` 时启用自动探测；未启用且未指定 `bind-interface` 时，底层连接使用系统路由。**若同时使用其他启用了 Fake IP 的 VPN，此功能可能无法正常工作**

+ `tcp-port-forwarding`: TCP 端口转发，格式为 `本地地址-远程地址,本地地址-远程地址,...`，例如 `127.0.0.1:9898-10.10.98.98:80,0.0.0.0:9899-10.10.98.98:80`。多个转发用 `,` 分隔

+ `udp-port-forwarding`: UDP 端口转发，格式为 `本地地址-远程地址,本地地址-远程地址,...`，例如 `127.0.0.1:53-202.120.224.26:53`。多个转发用 `,` 分隔

+ `custom-dns`: 指定自定义 DNS 解析结果，格式为 `域名:IP,域名:IP,...`，例如 `zb.fudan.edu.cn:10.108.68.200,mirrors.fducslg.com:10.176.52.2`。多个解析用 `,` 分隔

+ `proxy-all`: 是否代理所有流量，一般不需要加此参数

+ `config`: 指定配置文件，内容参考 `config.toml.example`

配置按“内置默认值 < TOML 配置文件 < 环境变量 < 命令行参数”的优先级合并。只有显式提供的命令行参数会覆盖其他来源，例如 `-tun-mode=false` 可以关闭配置文件或环境变量中启用的选项。

配置项均可使用带 `FDU_CONNECT_` 前缀的大写环境变量设置，例如 `server_address` 对应 `FDU_CONNECT_SERVER_ADDRESS`，`tun_mode` 对应 `FDU_CONNECT_TUN_MODE`。集合项使用与命令行参数相同的语法和名称，例如 `FDU_CONNECT_TCP_PORT_FORWARDING=127.0.0.1:9898-10.10.98.98:80`、`FDU_CONNECT_CUSTOM_DNS=zb.fudan.edu.cn:10.108.68.200` 和 `FDU_CONNECT_CUSTOM_PROXY_DOMAIN=nature.com,science.org`。

#### EasyConnect 相关参数

+ `cert-file`: p12 证书文件路径，如果服务器要求证书验证，需要配置此参数

+ `cert-password`: 证书密码（敏感信息建议使用环境变量 `FDU_CONNECT_CERT_PASSWORD` 给定）

+ `skip-domain-resource`: 不使用服务端下发的域名资源分流，一般不需要加此参数

+ `disable-multi-line`: 禁用自动根据延时选择线路。加此参数后，使用 `server` 参数指定的线路

+ `disable-fdu-config`: 禁用 FDU 相关配置，一般不需要加此参数

+ `custom-proxy-domain`: 指定自定义域名使用 RVPN 代理，格式为 `域名,域名,...`，例如 `nature.com,science.org`。多个域名用 `,` 分隔

+ `twf-id`: twfID 登录，调试用途，一般不需要加此参数

#### aTrust 相关参数

+ `auth-type`: aTrust 登录验证类型，支持 `auth/psw`（密码验证）、`auth/cas`（CAS 验证）、`auth/smsCheckCode`（短信验证码验证），默认为自动检测

+ `login-domain`: 登录域，默认为 `Radius`

+ `client-data-file`: 客户端数据文件路径，可用于保存登录状态，避免重复验证

+ `cas-ticket`: CAS 验证票据，默认为空，此时进入交互式验证

+ `phone`: 短信验证码登录时使用的手机号

+ `update-best-nodes-interval`: 自动选择最优线路的更新间隔，单位为秒，默认为 `300` 秒。设置为 `0` 则禁用自动选择最优线路

+ `session-refresh-interval`: aTrust SID 刷新间隔，单位为秒，默认为 `1800`。设置为 `0` 则禁用定时刷新。

+ `auth-info`: 仅获取 aTrust 验证信息而不登录，一般不需要加此参数。可用于查看服务端支持的验证方式

+ `trust-device`: 设置当前设备为授信终端（需要已登录的 `-client-data-file`），不启用隧道

+ `untrust-device`: 从授信终端中移除当前设备（需要已登录的 `-client-data-file`），不启用隧道

+ `sid`: aTrust SID，调试用途，一般不需要加此参数

+ `device-id`: aTrust 设备 ID，调试用途，一般不需要加此参数

+ `sign-key`: aTrust 签名密钥，调试用途，一般不需要加此参数

+ `resource-file`: aTrust 资源文件，调试用途，一般不需要加此参数

### 感谢

+ [ZJU Connect](https://github.com/Mythologyli/zju-connect)

+ [EasierConnect](https://github.com/lyc8503/EasierConnect)

+ [socks2http](https://github.com/zenhack/socks2http)

+ [![image](docs/yxvm.png)](https://yxvm.com/)

  [NodeSupport](https://github.com/NodeSeekDev/NodeSupport) 赞助了本项目
