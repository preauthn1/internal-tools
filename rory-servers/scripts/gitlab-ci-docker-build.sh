########################################################################################
#!/bin/bash
########################################################################################

BASEPATH=$(dirname $0)

. $BASEPATH/library.sh
. $BASEPATH/gitlab-ci-variables.sh

DOCKER_FOLDER="./build/docker"

# 先處理base
docker pull example.invalid/base
if [ "$?" = "1" ]; then
    docker build --target base -t example.invalid/base -f $DOCKER_FOLDER/Dockerfile-base $DOCKER_FOLDER
    docker push example.invalid/base
fi

IFS=',' read -r -a array <<< "$CI_BUILD_TARGETS"
if [ "${#array[@]}" != "0" ]; then
    for element in "${array[@]}"
    do
        IFS=':' read -r -a row <<< "$element"
        if [ "${#row[@]}" = "2" ]; then
            Fn_Build_Docker ${row[0]} ${row[1]} $DOCKER_REGISTRY_DOMAIN $RELEASE_FOLDER $DOCKER_FOLDER $BUILD_NAME $GIT_TAG
        fi
    done
fi
