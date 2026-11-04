## Docker 运行

```zsh
$ docker run -d --name fdu-connect -v $PWD/config.toml:/home/nonroot/config.toml -p 1080:1080 -p 1081:1081 --restart unless-stopped akirasalvare/fdu-connect
```

也可以使用 Docker Compose。创建 `docker-compose.yml` 文件，内容如下：

```yaml
services:
   fdu-connect:
      image: akirasalvare/fdu-connect
      container_name: fdu-connect
      restart: unless-stopped
      ports:
         - 1080:1080
         - 1081:1081
      volumes:
         - ./config.toml:/home/nonroot/config.toml
```

另外，你还可以使用 [configs top-level elements](https://docs.docker.com/compose/compose-file/08-configs/) 将 fdu-connect 的配置文件直接写入 docker-compose.yml，如下：

```yaml
services:
   fdu-connect:
      container_name: fdu-connect
      image: akirasalvare/fdu-connect
      restart: unless-stopped
      ports: [1080:1080, 1081:1081]
      configs: [{ source: fdu-connect-config, target: /home/nonroot/config.toml }]

configs:
   fdu-connect-config:
      content: |
         username = ""
         password = ""
         # other configs ...
```

并在同目录下运行

```zsh
docker compose up -d
```

### aTrust 协议登录和数据持久化

aTrust 协议客户端需要在容器内持久化登录凭据（`client_data.json`），可参考以下示例配置，将凭据文件保存在 Docker 卷中：

```
services:
   fdu-connect:
      container_name: fdu-connect
      image: akirasalvare/fdu-connect
      restart: unless-stopped
      ports:
         - 1080:1080
         - 1081:1081
      configs:
         - source: fdu-connect-config
           target: /home/nonroot/config.toml
      volumes:
         - fdu-connect-data:/home/nonroot/data
configs:
   fdu-connect-config:
      content: |
         username = ""
         password = ""
         client_data_file = "/home/nonroot/data/client_data.json"
         # other configs ...
volumes:
   fdu-connect-data:
```

使用 aTrust 协议初次登录时需要人工验证，应当使用 `network_mode: host` 和交互模式临时运行：

```zsh
docker compose run --rm -it fdu-connect
```

以下以需要图形验证码和短信验证码的服务端为例：

**1. 图形验证码**

程序会启动一个临时 HTTP 服务器显示验证码，端口号由系统随机分配。请从日志中查看实际地址：

```
2026/03/25 09:10:46 Captcha server started at http://127.0.0.1:40855
2026/03/25 09:10:46 Failed to open browser: exec: "xdg-open": executable file not found in $PATH. Please visit: http://127.0.0.1:40855
```

访问日志中显示的地址完成图形验证码验证后，程序会提示输入短信验证码：

**2. 短信验证码**

```
2026/03/25 09:10:57 SMS message sent successfully: 验证码已发送到您的手机：***, 请查收！
fdu-connect-1  | 2026/03/25 09:11:50 Please enter the SMS verification code:
```

出现上述提示后，请输入收到的短信验证码并回车。登录成功后会显示如下日志：

```
2026/03/25 09:12:12 Client data saved to /home/nonroot/data/client_data.json
2026/03/25 09:12:12 VPN client started
```

完成登录后按 `Ctrl+C` 退出容器，之后就可以使用 `docker compose up -d` 正常启动了，登录状态会被持久化在 Docker 卷中，无需重复验证。
