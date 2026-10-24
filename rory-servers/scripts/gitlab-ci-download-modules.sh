#!/bin/bash

export VAR_BUILD_NAME=".go/build-name"

## global variables ###

if [[ "${CI_COMMIT_REF_NAME}" =~ "phnum" ]]; then
    export BUILD_NAME="$(TZ=Asia/Taipei date +%Y%m%d)-phnum-${CI_PIPELINE_ID}"
else
    export BUILD_NAME="$(TZ=Asia/Taipei date +%Y%m%d)-build-${CI_PIPELINE_ID}"
fi
echo 'export BUILD_NAME="'${BUILD_NAME}'"' >> ${VAR_BUILD_NAME}
cat ${VAR_BUILD_NAME}
echo ${CI_COMMIT_REF_NAME}

export GO111MODULE="on"
export GOPROXY="http://0.0.0.0:3000"

# download go modules
go mod download