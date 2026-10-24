#!/bin/bash
source .go/.variables
echo $LDFLAGS
EXTNAME=""
if [ "$GOOS" = "windows" ]; then
    EXTNAME=".exe"
fi

mkdir -p /mnt/release
mount -t cifs -o username="release",password=[REDACTED]',vers=2.0 //0.0.0.0/99.release /mnt/release

### 老遊戲
echo -ne '##### build brnn                     (0%)\r'
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/gameserver/brnn/file/brnn${EXTNAME} ./cmd/game/brnn || { echo 'go build fail' ; exit 1; }
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/gameserver/bjl/file/bjl${EXTNAME} ./cmd/game/bjl || { echo 'go build fail' ; exit 1; }
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/gameserver/by/file/by${EXTNAME} ./cmd/game/by || { echo 'go build fail' ; exit 1; }
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/gameserver/ddz/file/ddz${EXTNAME} ./cmd/game/ddz || { echo 'go build fail' ; exit 1; }
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/gameserver/dzpk/file/dzpk${EXTNAME} ./cmd/game/dzpk || { echo 'go build fail' ; exit 1; }
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/gameserver/ermj/file/ermj${EXTNAME} ./cmd/game/ermj || { echo 'go build fail' ; exit 1; }
echo -ne '########## build fkeb                (15%)\r'
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/gameserver/fkeb/file/fkeb${EXTNAME} ./cmd/game/fkeb || { echo 'go build fail' ; exit 1; }
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/gameserver/hhdz/file/hhdz${EXTNAME} ./cmd/game/hhdz || { echo 'go build fail' ; exit 1; }
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/gameserver/hhnn/file/hhnn${EXTNAME} ./cmd/game/hhnn || { echo 'go build fail' ; exit 1; }
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/gameserver/lhd/file/lhd${EXTNAME} ./cmd/game/lhd || { echo 'go build fail' ; exit 1; }
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/gameserver/qzniuniu/file/qzniuniu${EXTNAME} ./cmd/game/qzniuniu || { echo 'go build fail' ; exit 1; }
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/gameserver/scmj/file/scmj${EXTNAME} ./cmd/game/scmj || { echo 'go build fail' ; exit 1; }
echo -ne '############### build sg             (30%)\r'
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/gameserver/sg/file/sg${EXTNAME} ./cmd/game/sg || { echo 'go build fail' ; exit 1; }
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/gameserver/srniuniu/file/srniuniu${EXTNAME} ./cmd/game/srniuniu || { echo 'go build fail' ; exit 1; }
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/gameserver/sss/file/sss${EXTNAME} ./cmd/game/sss || { echo 'go build fail' ; exit 1; }
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/gameserver/ttz/file/ttz${EXTNAME} ./cmd/game/ttz || { echo 'go build fail' ; exit 1; }
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/gameserver/wuzhang/file/wuzhang${EXTNAME} ./cmd/game/wuzhang || { echo 'go build fail' ; exit 1; }
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/gameserver/zjh/file/zjh${EXTNAME} ./cmd/game/zjh || { echo 'go build fail' ; exit 1; }
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/gameserver/tbniuniu/file/tbniuniu${EXTNAME} ./cmd/game/tbniuniu || { echo 'go build fail' ; exit 1; }

### Slot
echo -ne '#################### build dfdc      (45%)\r'
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/gameserver/common-slot/file/common-slot${EXTNAME} ./cmd/game/common-slot || { echo 'go build fail' ; exit 1; }
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/gameserver/dfdc/file/dfdc${EXTNAME} ./cmd/game/dfdc || { echo 'go build fail' ; exit 1; }
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/gameserver/tgw/file/tgw${EXTNAME} ./cmd/game/tgw || { echo 'go build fail' ; exit 1; }
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/gameserver/ald/file/ald${EXTNAME} ./cmd/game/ald || { echo 'go build fail' ; exit 1; }
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/gameserver/zcjb/file/zcjb${EXTNAME} ./cmd/game/zcjb || { echo 'go build fail' ; exit 1; }
echo -ne '######################### build wlzb (60%)\r'
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/gameserver/wlzb/file/wlzb${EXTNAME} ./cmd/game/wlzb || { echo 'go build fail' ; exit 1; }
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/gameserver/xgpd/file/xgpd${EXTNAME} ./cmd/game/xgpd || { echo 'go build fail' ; exit 1; }
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/gameserver/xmbx/file/xmbx${EXTNAME} ./cmd/game/xmbx || { echo 'go build fail' ; exit 1; }
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/gameserver/jjbx/file/jjbx${EXTNAME} ./cmd/game/jjbx || { echo 'go build fail' ; exit 1; }
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/gameserver/dhzw/file/dhzw${EXTNAME} ./cmd/game/dhzw || { echo 'go build fail' ; exit 1; }
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/gameserver/wlzbjp/file/wlzbjp${EXTNAME} ./cmd/game/wlzbjp || { echo 'go build fail' ; exit 1; }
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/gameserver/csdjp/file/csdjp${EXTNAME} ./cmd/game/csdjp || { echo 'go build fail' ; exit 1; }
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/gameserver/wlhc/file/wlhc${EXTNAME} ./cmd/game/wlhc || { echo 'go build fail' ; exit 1; }
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/gameserver/wljr/file/wljr${EXTNAME} ./cmd/game/wljr || { echo 'go build fail' ; exit 1; }
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/gameserver/jgyf/file/jgyf${EXTNAME} ./cmd/game/jgyf || { echo 'go build fail' ; exit 1; }
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/gameserver/cfll/file/cfll${EXTNAME} ./cmd/game/cfll || { echo 'go build fail' ; exit 1; }
echo -ne '##############################       (75%)\r'

### 棋牌類
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/gameserver/jlbpk/file/jlbpk${EXTNAME} ./cmd/game/jlbpk || { echo 'go build fail' ; exit 1; }
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/gameserver/jdesyd/file/jdesyd${EXTNAME} ./cmd/game/jdesyd || { echo 'go build fail' ; exit 1; }
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/gameserver/dczz/file/dczz${EXTNAME} ./cmd/game/dczz || { echo 'go build fail' ; exit 1; }
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/gameserver/hdc/file/hdc${EXTNAME} ./cmd/game/hdc || { echo 'go build fail' ; exit 1; }
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/gameserver/dchjk/file/dchjk${EXTNAME} ./cmd/game/dchjk || { echo 'go build fail' ; exit 1; }
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/gameserver/hbsl/file/hbsl${EXTNAME} ./cmd/game/hbsl || { echo 'go build fail' ; exit 1; }
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/gameserver/fgsb/file/fgsb${EXTNAME} ./cmd/game/fgsb || { echo 'go build fail' ; exit 1; }
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/gameserver/cjl/file/cjl{EXTNAME} ./cmd/game/cjl || { echo 'go build fail' ; exit 1; }
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/gameserver/hljbp/file/hljbp${EXTNAME} ./cmd/game/hljbp || { echo 'go build fail' ; exit 1; }
echo -ne '###################################  (99%)\r'