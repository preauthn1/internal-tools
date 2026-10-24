########################################################################################
#! /bin/bash
#
# Determine ENV from gitlab CI/CD, use "-" as default argument
#
# Usage: 
#     sh gitlab-ci-variables.sh ${CI_PIPELINE_ID} ${CI_PROJECT_DIR} ${CI_COMMIT_TAG} ${CI_COMMIT_SHORT_SHA}
# Example:
#     sh gitlab-ci-variables.sh - ${DIR} 1.2.0 -        # use default CI_PIPELINE_ID, CI_COMMIT_SHORT_SHA
#        
# @Param $1 CI_PIPELINE_ID          default: current time (HH:mm:ss)
# @Param $2 CI_PROJECT_DIR          default: current path
# @Param $3 CI_COMMIT_TAG           default: ""
# @Param $4 CI_COMMIT_SHORT_SHA     default: git rev-parse --short HEAD
########################################################################################

if [[ "$CI_PIPELINE_ID" == "" ]]; then 
    if [[ "$1" == "" || "$1" == "-" ]]; then
        CI_PIPELINE_ID=$(TZ=Asia/Taipei date +%H%M%S)
    else
        CI_PIPELINE_ID=$1
    fi
fi

if [[ "$CI_PROJECT_DIR" == "" ]]; then 
    if [[ "$2" == "" || "$2" == "-" ]]; then
        CI_PROJECT_DIR="."
    else
        CI_PROJECT_DIR=$2
    fi
fi

if [[ "$CI_COMMIT_TAG" == "" ]]; then 
    if [[ "$3" == "" || "$3" == "-" ]]; then
        CI_COMMIT_TAG=""
    else
        CI_COMMIT_TAG=$3
    fi
fi

if [[ "$CI_COMMIT_SHORT_SHA" == "" ]]; then 
    if [[ "$4" == "" || "$4" == "-" ]]; then
        CI_COMMIT_SHORT_SHA=$(git rev-parse --short HEAD)
    else
        CI_COMMIT_SHORT_SHA=$4
    fi
fi

if [[ "${CI_COMMIT_REF_NAME}" =~ "phnum" ]]; then
    BUILD_NAME="$(TZ=Asia/Taipei date +%Y%m%d)-phnum-${CI_PIPELINE_ID}"
else
    BUILD_NAME="$(TZ=Asia/Taipei date +%Y%m%d)-build-${CI_PIPELINE_ID}"
fi

GOARCH="amd64"
GIT_TAG=$CI_COMMIT_TAG
BUILD_TIME=$(TZ=Asia/Taipei date +%FT%T%z)
GIT_HASH=$CI_COMMIT_SHORT_SHA
BUILD_NUM=$CI_PIPELINE_ID
PKG_NAME="example.invalid/ServerTeam/Servers/pkg/conf"
LDFLAGS_1="-X ${PKG_NAME}.GitTag=${CI_PIPELINE_ID} -X ${PKG_NAME}.BuildTime=${BUILD_TIME} -X ${PKG_NAME}.GitHash=${GIT_HASH} -X ${PKG_NAME}.BuildNum=${BUILD_NAME}"
LDFLAGS_2="-X ${PKG_NAME}.GitTag=${GIT_TAG} -X ${PKG_NAME}.BuildTime=${BUILD_TIME} -X ${PKG_NAME}.GitHash=${GIT_HASH} -X ${PKG_NAME}.BuildNum=${BUILD_NAME}"
if [ "$GIT_TAG" == "" ]; then LDFLAGS=${LDFLAGS_1}; else LDFLAGS=${LDFLAGS_2}; fi
RELEASE_FOLDER=bin/${BUILD_NAME}/${GOOS}

function Var_Fn_Print() {
    echo "========  Gitlab CI/CD  ========"
    echo "CI_PIPELINE_ID = $CI_PIPELINE_ID"
    echo "CI_PROJECT_DIR = $CI_PROJECT_DIR"
    echo "CI_COMMIT_TAG = $CI_COMMIT_TAG"
    echo "CI_COMMIT_SHORT_SHA = $CI_COMMIT_SHORT_SHA"
    echo ""
    echo "========  Build Init  ========"
    echo "BUILD_NAME = $BUILD_NAME"
    echo "GOARCH = $GOARCH"
    echo "GIT_TAG = $GIT_TAG"
    echo "BUILD_TIME = $BUILD_TIME"
    echo "BUILD_NUM = $BUILD_NUM"
    echo "PKG_NAME = $PKG_NAME"
    echo "LDFLAGS_1 = $LDFLAGS_1"
    echo "LDFLAGS_2 = $LDFLAGS_2"
    echo "RELEASE_FOLDER = $RELEASE_FOLDER"
}

# Var_Fn_Print;