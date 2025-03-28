# FDU Connect

> 🚫 **Disclaimer**
>
> This program is provided **as is**, and the author **does not guarantee the correctness or reliability of the program**. Please judge whether the specific scenario is suitable for using this program. **The problems or consequences caused by using this program are borne by the user**!

---

[中文](README.md) | English

**This program is modified on [ZJU Connect](https://github.com/Mythologyli/zju-connect), thanks to the original author [Mythologyli](https://github.com/Mythologyli) and the contributors of the original project.**

**QQ group: 1037726410**, welcome to join the discussion.

### Usage

+ If you are from FDU:
  1. Download the latest version of the corresponding platform on the [Release](https://github.com/akirasalvare/fdu-connect/releases) page.

  2. Take macOS as an example, unzip the executable file `fdu-connect`.

  3. macOS needs to remove security restrictions first. Run: `sudo xattr -rd com.apple.quarantine fdu-connect`.

  4. Run: `./fdu-connect -username <username> -password <password>`.

  5. At this time, port `1080` is the Socks5 proxy, and port `1081` is the HTTP proxy. If you need to change the default port, please refer to [Arguments](#Arguments).

+ If you are not from FDU:

  We suppose you to try the original [ZJU Connect](https://github.com/Mythologyli/zju-connect)

#### Run as a service

[Link](docs/service_en.md)

#### Run in Docker

[Link](docs/docker_en.md)

### ⚠️ Warning

1. When using other proxy tools with TUN mode enabled and fdu-connect as a downstream proxy, please be sure to provide the correct network diversion rules, refer to [this issue](https://github.com/Mythologyli/zju-connect/issues/57)

### ⚠️ TUN mode precautions

1. Need to run with administrator privileges

2. Windows system needs to go to [Wintun official website](https://www.wintun.net) to download `wintun.dll` and place it in the same directory as the executable file

3. To ensure that domains are resolved correctly, it's recommended to configure `dns-hijack` to hijack the system DNS

### Arguments

+ `server`: SSL VPN server address, default is `stuvpn.fudan.edu.cn`

+ `port`: SSL VPN server port, default is `443`

+ `username`: Network account. For example: student ID

+ `password`: Network account password

+ `totp-secret`: TOTP secret. If the server doesn't need TOTP verification, or you want to manually enter the verification code, no need to add this argument

+ `cert-file`: p12 certificate file path, if the server requires certificate verification, this parameter needs to be configured

+ `cert-password`: Certificate password

+ `disable-server-config`: Disable server configuration, generally no need to add this argument

+ `skip-domain-resource`: Do not use the domain resource provided by the server to decide whether to use RVPN, generally no need to add this argument

+ `disable-fdu-config`: Disable FDU related configuration, generally no need to add this argument

+ `disable-fdu-dns`: Disable FDU DNS and use local DNS, generally no need to add this argument

+ `disable-multi-line`: Disable automatic line selection based on latency. After adding this argument, use the line specified by the `server` parameter

+ `proxy-all`: Whether to proxy all traffic, generally no need to add this argument

+ `socks-bind`: SOCKS5 proxy listening address, default is `:1080`

+ `socks-user`: SOCKS5 proxy username, leave blank if no authentication is required

+ `socks-passwd`: SOCKS5 proxy password, leave blank if no authentication is required

+ `http-bind`: HTTP proxy listening address, default is `:1081`. Set to `""` to disable HTTP proxy

+ `shadowsocks-url`: Shadowsocks server URL. For example: `ss://aes-128-gcm:password@server:port`. Format [refer to here](https://github.com/shadowsocks/go-shadowsocks2)

+ `dial-direct-proxy`: When a URL does not match RVPN rules and switches to direct connection, it uses a proxy, typically in scenarios where it works in conjunction with other proxy tools. Currently, only HTTP proxies are supported. For example: `http://127.0.0.1:7890`, setting it to empty string (`""`) will disable its use.

+ `tun-mode`: TUN mode (experimental). Please read the TUN mode precautions below

+ `add-route`: Add route according to the configuration issued by the server when TUN mode is enabled

+ `dns-ttl`: DNS cache time, default is `3600` seconds

+ `disable-keep-alive`: Disable periodic keep-alive, generally no need to add this argument

+ `fdu-dns-server`: FDU DNS server address, default is `202.120.224.26`. Set to `auto` to use the DNS server obtained from the server, and disable FDU DNS if it fails to obtain

+ `secondary-dns-server`: Standby DNS server used when FDU DNS server cannot be used, default is `223.5.5.5`. Leave blank to use the system default DNS, but must be set when `dns-hijack` is enabled

+ `dns-server-bind`: DNS server listening address, default is empty to disable. For example, set to `127.0.0.1:53`, then you can send DNS requests to `127.0.0.1:53`

+ `dns-hijack`: Hijack DNS requests when TUN mode is enabled, it's recommended to add this argument when TUN mode is enabled

+ `debug-dump`: Whether to enable debugging, generally no need to add this argument

+ `tcp-port-forwarding`: TCP port forwarding, format is `local address-remote address,local address-remote address,...`, for example `127.0.0.1:9898-10.10.98.98:80,0.0.0.0:9899-10.10.98.98:80`. Multiple forwarding is separated by `,`

+ `udp-port-forwarding`: UDP port forwarding, format is `local address-remote address,local address-remote address,...`, for example `127.0.0.1:53-202.120.224.26:53`. Multiple forwarding is separated by `,`

+ `custom-dns`: Specify custom DNS resolution results, format is `domain name:IP,domain name:IP,...`, for example `zb.fudan.edu.cn:10.108.68.200,mirrors.fducslg.com:10.176.52.2`. Multiple resolutions are separated by `,`

+ `custom-proxy-domain`: Specify custom domains which use RVPN proxy, format is `domain,domain,...`, for example `nature.com,science.org`. Multiple resolutions are separated by `,`

+ `twf-id`: twfID login, for debugging purposes, generally no need to add this argument

+ `config`: Specify the configuration file, the content refers to `config.toml.example`. Other parameters are ignored when the configuration file is enabled

### Thanks

+ [ZJU Connect](https://github.com/Mythologyli/zju-connect)

+ [EasierConnect](https://github.com/lyc8503/EasierConnect)

+ [socks2http](https://github.com/zenhack/socks2http)
