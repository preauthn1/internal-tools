# terpro_local_dcp

原始路径：`terpro_local_dcp`　|　文件：**32**　|　主语言：Shell 6, YAML 3

## 目录构成

```
docker/apis/  — 5
docker/API_External/  — 5
docker/api-internal/  — 5
docker/API_Frontend/  — 5
docker/API_Agent/  — 5
docker/WEB_AgentBackend/  — 4
phpactor/WEB_AgentBackend/  — 1
phpactor/API_Agent/  — 1
```

## 关键文件

- CI/构建：`docker-compose.yml`、`docker/apis/Dockerfile`、`docker/WEB_AgentBackend/Dockerfile`

## 主要文件

- `docker/API_Frontend/fpm.conf`
- `docker/API_Agent/fpm.conf`
- `docker/API_External/fpm.conf`
- `docker/api-internal/fpm.conf`
- `docker/WEB_AgentBackend/fpm.conf`
- `docker/API_Frontend/_release/.env`
- `docker-compose.yml`
- `docker/WEB_AgentBackend/nginx.conf`

## 说明

该工程为平台的一个组件。仅收录文本源码与配置，图片、音频等二进制资源未收录；凭据、密钥、内网地址与域名已按实际内容替换。

