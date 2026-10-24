#!/bin/bash
export VAR_BUILD_NAME=".go/build-name"
source $VAR_BUILD_NAME

### environment variables ###
export GOARCH="amd64"
export GIT_TAG=$CI_COMMIT_TAG
export BUILD_TIME=$(TZ=Asia/Taipei date +%FT%T%z)
export GIT_HASH=$CI_COMMIT_SHORT_SHA
export BUILD_NUM=$CI_PIPELINE_ID
export PKG_NAME="example.invalid/ServerTeam/Servers/pkg/conf"
LDFLAGS_1="-X ${PKG_NAME}.BuildTime=${BUILD_TIME} -X ${PKG_NAME}.GitHash=${GIT_HASH} -X ${PKG_NAME}.BuildNum=${BUILD_NAME}"
LDFLAGS_2="-X ${PKG_NAME}.GitTag=${GIT_TAG} -X ${PKG_NAME}.BuildTime=${BUILD_TIME} -X ${PKG_NAME}.GitHash=${GIT_HASH} -X ${PKG_NAME}.BuildNum=${BUILD_NAME}"
if [ "$GIT_TAG" == "" ]; then export LDFLAGS=${LDFLAGS_1}; else export LDFLAGS=${LDFLAGS_2}; fi
export RELEASE_FOLDER=${CI_PROJECT_DIR}/bin/${GOOS}

### output to ${GO_VARIABLES} for after script
export GO_VARIABLES=".go/.variables"
echo 'export BUILD_NAME="'${BUILD_NAME}'"' >> ${GO_VARIABLES}
echo 'export CGO_ENABLED=0'
echo 'export GOOS="'${GOOS}'"' >> ${GO_VARIABLES}
echo 'export GOARCH="'${GOARCH}'"' >> ${GO_VARIABLES}
echo 'export GIT_TAG="'${GIT_TAG}'"' >> ${GO_VARIABLES}
echo 'export BUILD_TIME="'${BUILD_TIME}'"' >> ${GO_VARIABLES}
echo 'export GIT_HASH="'${GIT_HASH}'"' >> ${GO_VARIABLES}
echo 'export BUILD_NUM="'${BUILD_NUM}'"' >> ${GO_VARIABLES}
echo 'export PKG_NAME="'${PKG_NAME}'"' >> ${GO_VARIABLES}
echo 'export LDFLAGS="'${LDFLAGS}'"' >> ${GO_VARIABLES}
echo 'export RELEASE_FOLDER="'${RELEASE_FOLDER}'"' >> ${GO_VARIABLES}
cat ${GO_VARIABLES}

### release folder
mkdir -p ${RELEASE_FOLDER}

### version info
echo '{"GitTag":"'${GIT_TAG}'","BuildNum":"'${BUILD_NAME}'","BuildTime":"'${BUILD_TIME}'","GitHash":"'${GIT_HASH}'"}' > ${RELEASE_FOLDER}/version.txt
