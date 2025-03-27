## Run in Docker

```shell
docker run -d --name fdu-connect -v $PWD/config.toml:/home/nonroot/config.toml -p 1080:1080 -p 1081:1081 --restart unless-stopped akirasalvare/fdu-connect
```

You can also use Docker Compose. Create `docker-compose.yml` file with the following content:

```yaml
version: '3'

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

Additionally, you can also use [configs top-level elements](https://docs.docker.com/compose/compose-file/08-configs/) to directly write the configuration files of fdu-connect into docker-compose.yml, as shown below:

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

And run the following command in the same directory:

```shell
docker compose up -d
```