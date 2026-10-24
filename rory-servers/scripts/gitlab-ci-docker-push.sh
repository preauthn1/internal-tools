########################################################################################
#!/bin/bash
########################################################################################

BASEPATH=$(dirname $0)

. $BASEPATH/library.sh
. $BASEPATH/gitlab-ci-variables.sh

IFS=',' read -r -a array <<< "$CI_BUILD_TARGETS"
if [ "${#array[@]}" != "0" ]; then
    for element in "${array[@]}"
    do
        IFS=':' read -r -a row <<< "$element"
        if [ "${#row[@]}" = "2" ]; then
            Fn_Push_Docker ${row[0]} ${row[1]} $DOCKER_REGISTRY_DOMAIN $DOCKER_REGISTRY_USERNAME $DOCKER_REGISTRY_PASSWORD $BUILD_NAME $GIT_TAG
        fi
    done
fi
