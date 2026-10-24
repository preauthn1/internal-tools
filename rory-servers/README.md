# rory-servers

> 本目录原始 README 见文末，本节为快照补充信息。

原始路径：`rory-servers`　|　文件 **1951**　|　主语言：Go 1780, YAML 41, Shell 16

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
```

## 关键文件

- CI/构建：`Makefile`、`.gitlab-ci.yml`、`internal/game/infc/Makefile`

---

## 原始说明

# Servers

[![pipeline status](http://example.invalid:1688/server-team/servers/badges/master/pipeline.svg)](http://example.invalid:1688/server-team/servers/commits/master)
[![coverage report](http://example.invalid:1688/server-team/servers/badges/master/coverage.svg)](http://example.invalid:1688/server-team/servers/commits/master)

## Packages

* Poker(v1)([v2](pkg/poker/v2/README.md))
* Slot([v1](pkg/slot/README.md))

## Games

#### Poker
* [印度拉米(INRM)](internal/game/inrm/README.md)

#### Slot 
* [多福多財(DFDC), 招財進寶(ZCJB)](internal/game/dfdc/README.md)
* [糖果屋(TGW)](internal/game/candyhouse/README.md)
* [星光派對(XGPD), 阿拉丁(ALD)](internal/game/aladdin/README.md)
* [五龍爭霸(WLZB)](internal/game/wlzb/README.md)
* [五龍爭霸_彩金版(WLZBJP), 財神到_採金版(CSDJP)](internal/game/wlzbjp/README.md)
* [五龍皇朝(WLHC)](internal/game/wlhc/README.md)
* [百萬金輪五龍爭霸(WLJR)](internal/game/wljr/README.md)
* [熊貓報喜(XMBX)](internal/game/xmbx/README.md)
* [擊鼓迎福(JGYF)](internal/game/jgyf/README.md)
* [吉祥如意(JXRY)](internal/game/jxry/README.md)
* [金吉報喜(JJBX), 財富連連(CFLL)](internal/game/jjbx/README.md)
* [東海之王(DHZW)](internal/game/dhzw/README.md)

#### PachiSlot
* [將軍(GEN)](internal/game/gen/README.md)

#### Arena (對戰)
* [印度切水果(INFC)](internal/game/infc/README.md)

#### Slot88 合作
* [88 Fortune(S88F)](internal/game/s88f/README.md)
* [Joker's Jewels(JKJL)](internal/game/jkjl/README.md)
* [Lucky Bats(LKB)](internal/game/lkb/README.md)
* [Frogs And Flies(FAF)](internal/game/faf/README.md)
* [Queen of Wands(QOW)](internal/game/qow/README.md)
* [Buffalo Bonus(BB)](internal/game/bb/README.md)
* [888(S888)](internal/game/s888/README.md)
* [5 Lucky Lions(s5ll)](internal/game/s5ll/README.md)

> 若使用`slot88`架構的Client需要透過[Forward](cmd/slot88/forward/README.md)來轉換

