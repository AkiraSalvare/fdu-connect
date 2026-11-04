# FDU Connect

> 🚫 **Disclaimer**
>
> This program is provided **as is**, and the author **does not guarantee the correctness or reliability of the program**. Please judge whether the specific scenario is suitable for using this program. **The problems or consequences caused by using this program are borne by the user**!

---

[中文](README.md) | English

**This program is modified on [ZJU Connect](https://github.com/Mythologyli/zju-connect), thanks to the original author [Mythologyli](https://github.com/Mythologyli) and the contributors of the original project.**

### Usage

+ If you are from FDU:
  1. Download the latest version of the corresponding platform on the [Release](https://github.com/akirasalvare/fdu-connect/releases) page.
     - Arch Linux users can directly install the AUR package [fdu-connect-git](https://aur.archlinux.org/packages/fdu-connect-git): `[yay|paru] -S fdu-connect-git`

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

3. EasyConnect Protocol: The recommended configuration is: `-tun-mode -add-route -dns-hijack`

4. aTrust Protocol: The recommended configuration is: `-tun-mode -add-route -dns-hijack -fake-ip`. Note: When using the aTrust protocol, direct TCP traffic involving domain names via the TUN interface may fail if DNS Hijacking or Fake IP is not enabled

### Arguments

#### General Arguments

+ `protocol`: Login protocol, supports `easyconnect`/`atrust`, default is `easyconnect`

+ `server`: VPN server address, default is `stuvpn.fudan.edu.cn`

+ `port`: VPN server port, default is `443`

+ `username`: Network account. For example: student ID

+ `password`: Network account password. (Sensitive information is recommended to be provided via env `FDU_CONNECT_PASSWORD`)

+ `graph-code-file`: Graphic captcha file path, default is empty. In aTrust mode, if set, the program will save the captcha image to this file and ask user to input the JSON in terminal.

+ `totp-secret`: TOTP secret, can be used to automatically complete TOTP verification. If the server doesn't require TOTP or you want to manually enter the code, leave blank. (Sensitive information is recommended to be provided via env `FDU_CONNECT_TOTP_SECRET`)

+ `disable-remote-dns`: Disable remote DNS and use local DNS instead, generally no need to add this argument

+ `disable-server-config`: Disable server configuration, generally no need to add this argument

+ `socks-bind`: SOCKS5 proxy listening address, default is `:1080`

+ `socks-user`: SOCKS5 proxy username, leave blank if no authentication is required

+ `socks-passwd`: SOCKS5 proxy password, leave blank if no authentication is required

+ `http-bind`: HTTP proxy listening address, default is `:1081`. Set to `""` to disable HTTP proxy

+ `shadowsocks-url`: Shadowsocks server URL. For example: `ss://aes-128-gcm:password@server:port`. Format [refer to here](https://github.com/shadowsocks/go-shadowsocks2)

+ `dial-direct-proxy`: When a URL does not match rules and switches to direct connection, use a proxy. Typically used in conjunction with other proxy tools, currently supports only HTTP proxy. For example: `http://127.0.0.1:7890`. Set to `""` to disable

+ `tcp-tunnel-mode`: TCP tunnel mode, default is `false`. When enabled, only TCP traffic can be proxied through the TCP tunnel. Since only aTrust supports TCP tunneling, this mode is ineffective under EasyConnect. Enabling this will disable TUN mode

+ `tun-mode`: TUN mode (experimental). Please read the TUN mode precautions above

+ `add-route`: Add routes according to the configuration issued by the server when TUN mode is enabled

+ `dns-ttl`: DNS cache time, default is `3600` seconds

+ `disable-keep-alive`: Disable periodic keep-alive, generally no need to add this argument

+ `keep-alive-url`: Uses HTTP keep-alive, suitable for situations where the remote DNS is unavailable. Set the URL to visit, for example, `https://www.cnki.net/favicon.ico`. The default is empty, in which case the remote DNS is used for keep-alive.

+ `remote-dns-server`: Remote DNS server address, default is `202.120.224.26`. Set to `auto` to use the DNS server obtained from the server; disable remote DNS if it fails to obtain

+ `secondary-dns-server`: Standby DNS server used when the remote DNS server cannot resolve. The default `auto` uses the second server supplied by VPN policy, then falls back to `223.5.5.5`. Leave blank to use system default DNS, but it must be set when `dns-hijack` is enabled

+ `dns-server-bind`: DNS server listening address, default is empty (disabled). For example, set to `127.0.0.1:53`, then you can send DNS requests to `127.0.0.1:53`

+ `local-dns-server`: Local DNS server used to resolve the VPN server hostname, as IP or IP:port; when empty, the system DNS is used, routable DNS addresses are bound to the detected underlay interface, and local DNS stubs keep their loopback route

+ `dns-hijack`: Hijack DNS requests when TUN mode is enabled, it's recommended to add this argument when using TUN mode

+ `fake-ip`: Enable Fake IP mode. Works with dns-hijack. Don't enable it if you are using EasyConnect protocol

+ `debug-dump`: Whether to enable debugging, generally no need to add this argument

+ `debug-pcap-file`: Reconstruct a PCAP from data read and written on VPN underlay TCP connections; a full capture queue blocks network I/O, kernel handshakes and retransmissions are omitted, and TLS payloads remain encrypted

+ `debug-tls-log-file`: Export TLS session secrets in NSS key log format for decrypting `debug-pcap-file` traffic in Wireshark. The file contains session secrets and must only be used for debugging and stored securely

+ `bind-interface`: Manually bind VPN underlay connections to this network interface for either EasyConnect or aTrust. A non-empty value takes precedence over automatic detection.

+ `auto-detect-interface`: Automatically detect and bind the VPN underlay interface; defaults to `false`. Set it to `true` to enable automatic detection. If disabled and `bind-interface` is empty, underlay connections use system routing. **This feature may not work correctly while another VPN with Fake IP enabled is in use.**

+ `tcp-port-forwarding`: TCP port forwarding, format is `local address-remote address,local address-remote address,...`, for example `127.0.0.1:9898-10.10.98.98:80,0.0.0.0:9899-10.10.98.98:80`. Multiple forwardings are separated by `,`

+ `udp-port-forwarding`: UDP port forwarding, format is `local address-remote address,local address-remote address,...`, for example `127.0.0.1:53-202.120.224.26:53`. Multiple forwardings are separated by `,`

+ `custom-dns`: Specify custom DNS resolution results, format is `domain:IP,domain:IP,...`, for example `zb.fudan.edu.cn:10.108.68.200,mirrors.fducslg.com:10.176.52.2`. Multiple resolutions are separated by `,`

+ `proxy-all`: Whether to proxy all traffic, generally no need to add this argument

+ `config`: Specify the configuration file, referring to `config.toml.example`

Configuration sources are merged with the precedence “built-in defaults < TOML file < environment variables < command-line flags”. Only explicitly supplied command-line flags override other sources; for example, `-tun-mode=false` disables a value enabled by the file or environment.

Every configuration key can be set with an uppercase environment variable prefixed by `FDU_CONNECT_`. For example, `server_address` maps to `FDU_CONNECT_SERVER_ADDRESS`, and `tun_mode` maps to `FDU_CONNECT_TUN_MODE`. Collection values use the same names and syntax as command-line flags, for example `FDU_CONNECT_TCP_PORT_FORWARDING=127.0.0.1:9898-10.10.98.98:80`, `FDU_CONNECT_CUSTOM_DNS=zb.fudan.edu.cn:10.108.68.200`, and `FDU_CONNECT_CUSTOM_PROXY_DOMAIN=nature.com,science.org`.

#### EasyConnect Related Arguments

+ `cert-file`: p12 certificate file path, if the server requires certificate verification, this parameter needs to be configured

+ `cert-password`: Certificate password. (Sensitive information is recommended to be provided via env `FDU_CONNECT_CERT_PASSWORD`)

+ `skip-domain-resource`: Do not use the domain resource provided by the server for split tunneling, generally no need to add this argument

+ `disable-multi-line`: Disable automatic line selection based on latency. When added, use the line specified by the `server` parameter

+ `disable-fdu-config`: Disable FDU related configuration, generally no need to add this argument

+ `custom-proxy-domain`: Specify custom domains to use RVPN proxy, format is `domain,domain,...`, for example `nature.com,science.org`. Multiple domains are separated by `,`

+ `twf-id`: twfID login, for debugging purposes, generally no need to add this argument

#### aTrust Related Arguments

+ `auth-type`: aTrust login authentication type, supports `auth/psw` (password), `auth/cas` (CAS), `auth/smsCheckCode` (SMS verification code), default is auto detected.
+ `login-domain`: Login domain, default is `Radius`.
+ `client-data-file`: Client data file path, used to save login status to avoid repeated verification.
+ `cas-ticket`: CAS verification ticket, defaults to empty, which triggers interactive verification.
+ `phone`: Phone number used for SMS verification code login.
+ `update-best-nodes-interval`: Interval for updating the optimal line automatically, in seconds, default is `300`. Set to `0` to disable automatic optimal line selection.
+ `session-refresh-interval`: aTrust SID refresh interval in seconds, default `1800`. Set to `0` to disable periodic refresh. The TOML key is `session_refresh_interval`; the environment variable is `FDU_CONNECT_SESSION_REFRESH_INTERVAL`.
+ `auth-info`: Only get aTrust authentication information without logging in, generally no need to add this argument. Can be used to check supported authentication methods.
+ `trust-device`: Trust the current device (requires logged-in `-client-data-file`), does not start the tunnel.
+ `untrust-device`: Untrust the current device (requires logged-in `-client-data-file`), does not start the tunnel.
+ `sid`: aTrust SID, for debugging purposes, generally no need to add this argument.
+ `device-id`: aTrust device ID, for debugging purposes, generally no need to add this argument.
+ `sign-key`: aTrust signature key, for debugging purposes, generally no need to add this argument.
+ `resource-file`: aTrust resource file, for debugging purposes, generally no need to add this argument.

### Thanks

+ [ZJU Connect](https://github.com/Mythologyli/zju-connect)

+ [EasierConnect](https://github.com/lyc8503/EasierConnect)

+ [socks2http](https://github.com/zenhack/socks2http)

+ [![image](docs/yxvm.png)](https://yxvm.com/)

  [NodeSupport](https://github.com/NodeSeekDev/NodeSupport) sponsored this project
