# rory-servers

原始路径：`rory-servers`　|　文件：**1951**　|　主语言：Go 1780, YAML 41, Shell 16

## 目录构成

```
internal/game/  — 658
internal/dbserver/  — 276
pkg/slot/  — 140
pkg/netproto/  — 113
pkg/pachislot/  — 77
cmd/game/  — 68
internal/apiserver/  — 52
release-notes/  — 37
pkg/gamelogic/  — 32
pkg/poker/  — 31
cmd/slot88/  — 31
internal/cli-api-svc/  — 28
internal/hallserver/  — 27
pkg/gameserver/  — 23
```

## 关键文件

- CI/构建：`Makefile`、`.gitlab-ci.yml`、`internal/game/infc/Makefile`

## 主要文件

- `pkg/netproto/HallMsg.pb.go`
- `internal/game/gen/game_test.go`
- `internal/game/by/bylogic_test.go`
- `internal/game/by/bylogic.go`
- `internal/game/jlbpk/jlbpkLogic.go`
- `internal/game/ddz/ddzgamelogic.go`
- `internal/game/ddz/ddzrobotlogic.go`
- `internal/game/bjl/baijialelogic.go`

## 说明

该工程为平台的一个组件。仅收录文本源码与配置，图片、音频等二进制资源未收录；凭据、密钥、内网地址与域名已按实际内容替换。

