########################################################################################
#!/bin/bash
# Librarys, define Functions, Arguments
########################################################################################

BASEPATH=$(dirname $0)

EXTNAME=""
if [ "$GOOS" = "windows" ]; then
    EXTNAME=".exe"
fi

COMMONS=( dbserver apiserver hallserver notifyserver tfconfig tfserver )
SLOT_GAMES=( dfdc tgw ald zcjb wlzb xgpd xmbx jjbx dhzw wlzbjp csdjp wlhc wljr jgyf cfll wlyc jxry hyxm jkjl lkb faf qow s888 s88f s5ll bb dgno hhf)
PACHI_SLOT_GAMES=( gen )
POKER_GAMES=( bjl hhdz lhd jlbpk jdesyd dczz hdc dchjk hbsl fgsb cjl hljbp inrm )
OLD_GAMES=( brnn bjl by ddz dzpk ermj fkeb hhdz hhnn lhd qzniuniu scmj sg srniuniu sss ttz wuzhang zjh tbniuniu )

########################################################################################
# Fn_Build_Golang_Binary
#   Build Binary to Output Folder
#
# @Param $1 Type                    [all/common/game]
# @Param $2 Name                    
# @Param $3 Output base folder      
# @Param $4 ldflags
# @Param $5 Build Count
# @Param $6 Total Count
# 
# Usage: 
#     Fn_Build_Golang_Binary ${TYPE} ${NAME} ${FOLDER} ${FLAG_STR} 
# Example:
#     Fn_Build_Golang_Binary all all ./bin/20200204-build-1234 "$LDFLAGS"
#     Fn_Build_Golang_Binary common dbserver ./bin/20200204-build-1234 "$LDFLAGS"
#     Fn_Build_Golang_Binary common all ./bin/20200204-build-1234 "$LDFLAGS"
#     Fn_Build_Golang_Binary game dfdc ./bin/20200204-build-1234 "$LDFLAGS"
#     Fn_Build_Golang_Binary game slot ./bin/20200204-build-1234 "$LDFLAGS"
#     Fn_Build_Golang_Binary game all ./bin/20200204-build-1234 "$LDFLAGS"
#
########################################################################################
function Fn_Build_Golang_Binary() {
    if [ "$1" = "all" ]; then
        Fn_Build_Golang_Binary common all $3 "$4"
        Fn_Build_Golang_Binary game all $3 "$4"
    elif [ "$1" = "common" ]; then
        if [ "$2" = "all" ]; then
            echo "===== Build Commons [${#COMMONS[@]}] ====="
            for (( i = 0; i < ${#COMMONS[@]}; i++ ));
            do
                Fn_Build_Golang_Binary $1 ${COMMONS[$i]} $3 "$4" $(expr $i + 1) ${#COMMONS[@]}
            done
        else
            if [ "$5" = "" ] && [ "$6" = "" ]; then
                printf "Building %12s ...\r" $2
            else
                printf "Building %12s [%02d/%02d] ...\r" $2 $5 $6
            fi

            local start=$SECONDS

            go build -ldflags "$4" -i -o $3/$2/$2$EXTNAME ./cmd/$2 || { echo "go build $2 fail" ; exit 1; }

            local duration=$(($SECONDS - $start))

            if [ "$5" = "" ] && [ "$6" = "" ]; then
                printf "Built %15s Successful ... (%dm%ds)\n" $2 $(($duration/60)) $(($duration%60))
            else
                printf "Built %15s [%02d/%02d] Successful ... (%dm%ds)\n" $2 $5 $6 $(($duration/60)) $(($duration%60))
            fi
        fi
    elif [ "$1" = "game" ]; then
        if [ "$2" = "all" ]; then
            Fn_Build_Golang_Binary $1 slot $3 "$4"
            Fn_Build_Golang_Binary $1 pachislot $3 "$4"
            Fn_Build_Golang_Binary $1 poker $3 "$4"
            Fn_Build_Golang_Binary $1 old $3 "$4"
        elif [ "$2" = "slot" ]; then
            echo "===== Build Slot Games [${#SLOT_GAMES[@]}] ====="
            for (( i = 0; i < ${#SLOT_GAMES[@]}; i++ ));
            do
                Fn_Build_Golang_Binary $1 ${SLOT_GAMES[$i]} $3 "$4" $(expr $i + 1) ${#SLOT_GAMES[@]}
            done
        elif [ "$2" = "pachislot" ]; then
            echo "===== Build Pachi Slot Games [${#PACHI_SLOT_GAMES[@]}] ====="
            for (( i = 0; i < ${#PACHI_SLOT_GAMES[@]}; i++ ));
            do
                Fn_Build_Golang_Binary $1 ${PACHI_SLOT_GAMES[$i]} $3 "$4" $(expr $i + 1) ${#PACHI_SLOT_GAMES[@]}
            done
        elif [ "$2" = "poker" ]; then
            echo "===== Build Poker Games [${#POKER_GAMES[@]}] ====="
            for (( i = 0; i < ${#POKER_GAMES[@]}; i++ ));
            do
                Fn_Build_Golang_Binary $1 ${POKER_GAMES[$i]} $3 "$4" $(expr $i + 1) ${#POKER_GAMES[@]}
            done
        elif [ "$2" = "old" ]; then
            echo "===== Build old Games [${#OLD_GAMES[@]}] ====="
            for (( i = 0; i < ${#OLD_GAMES[@]}; i++ ));
            do
                Fn_Build_Golang_Binary $1 ${OLD_GAMES[$i]} $3 "$4" $(expr $i + 1) ${#OLD_GAMES[@]}
            done
        else
            if [ "$5" = "" ] && [ "$6" = "" ]; then
                printf "Building %12s ...\r" $2
            else
                printf "Building %12s [%02d/%02d] ...\r" $2 $5 $6
            fi

            local start=$SECONDS

            go build -ldflags "$4" -i -o $3/gameserver/$2/file/$2$EXTNAME ./cmd/game/$2 || { echo "go build $2 fail" ; exit 1; }

            local duration=$(($SECONDS - $start))

            if [ "$5" = "" ] && [ "$6" = "" ]; then
                printf "Built %15s Successful ... (%dm%ds)\n" $2 $(($duration/60)) $(($duration%60))
            else
                printf "Built %15s [%02d/%02d] Successful ... (%dm%ds)\n" $2 $5 $6 $(($duration/60)) $(($duration%60))
            fi
        fi
    else
        echo "Unknown Type, Please give [all/common/game]"
    fi
}

########################################################################################
# Fn_Build_Docker
#   Build Docker Images
#   1. copy binary file from Release folder to Docker folder
#   2. build docker image
#   3. delete binary file in Docker folder after built success
#
# @Param $1 Type                    [all/common/game]
# @Param $2 Name
# @Param $3 Repository Domain                    
# @Param $4 Release folder      
# @Param $5 Docker folder
# @Param $6 Build Number
# @Param $7 Git Tag
# @Param $8 Build Count
# @Param $9 Total Count
# 
# Usage: 
#     Fn_Build_Docker $TYPE $NAME $REGISTERY $RELEASE_FOLDER $DOCKER_FOLDER $BUILD_NAME $GIT_TAG
# Example:
#     Fn_Build_Docker all all $REGISTERY $RELEASE_FOLDER $DOCKER_FOLDER $BUILD_NAME $GIT_TAG
#     Fn_Build_Docker common dbserver $REGISTERY $RELEASE_FOLDER $DOCKER_FOLDER $BUILD_NAME $GIT_TAG
#     Fn_Build_Docker common all $REGISTERY $RELEASE_FOLDER $DOCKER_FOLDER $BUILD_NAME $GIT_TAG
#     Fn_Build_Docker game dfdc $REGISTERY $RELEASE_FOLDER $DOCKER_FOLDER $BUILD_NAME $GIT_TAG
#     Fn_Build_Docker game slot $REGISTERY $RELEASE_FOLDER $DOCKER_FOLDER $BUILD_NAME $GIT_TAG
#     Fn_Build_Docker game all $REGISTERY $RELEASE_FOLDER $DOCKER_FOLDER $BUILD_NAME $GIT_TAG
#
########################################################################################
function Fn_Build_Docker() {
    TAG_NAME=$7
    if [ "$TAG_NAME" = "" ] || [ "$TAG_NAME" = "-" ]; then
        TAG_NAME=$6
    fi

    LATEST_STR="latest"
    GIT_TAG="-"
    if [ "$7" = "" ] || [ "$7" = "-" ]; then
        LATEST_STR="beta"
    else
        GIT_TAG=$7
    fi

    if [ "$1" = "all" ]; then
        Fn_Build_Docker common all $3 $4 $5 $6 $GIT_TAG
        Fn_Build_Docker game all $3 $4 $5 $6 $GIT_TAG
    elif [ "$1" = "common" ]; then
        if [ "$2" = "all" ]; then
            echo "===== Docker Build Commons [${#COMMONS[@]}] ====="
            for (( i = 0; i < ${#COMMONS[@]}; i++ ));
            do
                Fn_Build_Docker $1 ${COMMONS[$i]} $3 $4 $5 $6 $GIT_TAG $(expr $i + 1) ${#COMMONS[@]}
            done
        else
            REPO_NAME=$2
            FULL_NAME="$3/$REPO_NAME"

            if [ "$8" = "" ] && [ "$9" = "" ]; then
                printf "Docker Building %12s ...\r" $2
            else
                printf "Docker Building %12s [%02d/%02d] ...\r" $2 $8 $9
            fi

            local start=$SECONDS

            # Copy File to build resource
            cp ${4}/$2/$2 $5/

            docker pull example.invalid/base

            if [ "$?" = "0" ]; then
                docker build --target app --cache-from example.invalid/base \
                    -t ${FULL_NAME}:$TAG_NAME \
                    -t ${FULL_NAME}:${LATEST_STR} \
                    --build-arg APP_TYPE=$2 \
                    --build-arg APP_NAME=$2 \
                    $5 || { echo "Docker Build $2 Fail" ; exit 1; }
            else
                docker build \
                    -t ${FULL_NAME}:$TAG_NAME \
                    -t ${FULL_NAME}:${LATEST_STR} \
                    --build-arg APP_TYPE=$2 \
                    --build-arg APP_NAME=$2 \
                    $5 || { echo "Docker Build $2 Fail" ; exit 1; }
            fi

            local duration=$(($SECONDS - $start))

            if [ "$8" = "" ] && [ "$9" = "" ]; then
                printf "Docker Built %15s Successful ... (%dm%ds) (%s)\n" $2 $(($duration/60)) $(($duration%60)) ${FULL_NAME}:$TAG_NAME
            else
                printf "Docker Built %15s [%02d/%02d] Successful ... (%dm%ds) (%s)\n" $2 $8 $9 $(($duration/60)) $(($duration%60)) ${FULL_NAME}:$TAG_NAME
            fi

            # Delete File from build resource
            rm $5/$2
        fi

    elif [ "$1" = "game" ]; then
        if [ "$2" = "all" ]; then
            Fn_Build_Docker game slot $3 $4 $5 $6 $GIT_TAG
            Fn_Build_Docker game pachislot $3 $4 $5 $6 $GIT_TAG
            Fn_Build_Docker game poker $3 $4 $5 $6 $GIT_TAG
            Fn_Build_Docker game old $3 $4 $5 $6 $GIT_TAG
        elif [ "$2" = "slot" ]; then
            echo "===== Docker Build Slot Games [${#SLOT_GAMES[@]}] ====="
            for (( i = 0; i < ${#SLOT_GAMES[@]}; i++ ));
            do
                Fn_Build_Docker $1 ${SLOT_GAMES[$i]} $3 $4 $5 $6 $GIT_TAG $(expr $i + 1) ${#SLOT_GAMES[@]}
            done
        elif [ "$2" = "pachislot" ]; then
            echo "===== Docker Pachi Slot Games [${#PACHI_SLOT_GAMES[@]}] ====="
            for (( i = 0; i < ${#PACHI_SLOT_GAMES[@]}; i++ ));
            do
                Fn_Build_Docker $1 ${PACHI_SLOT_GAMES[$i]} $3 $4 $5 $6 $GIT_TAG $(expr $i + 1) ${#PACHI_SLOT_GAMES[@]}
            done
        elif [ "$2" = "poker" ]; then
            echo "===== Docker Build Poker Games [${#POKER_GAMES[@]}] ====="
            for (( i = 0; i < ${#POKER_GAMES[@]}; i++ ));
            do
                Fn_Build_Docker $1 ${POKER_GAMES[$i]} $3 $4 $5 $6 $GIT_TAG $(expr $i + 1) ${#POKER_GAMES[@]}
            done
        elif [ "$2" = "old" ]; then
            echo "===== Docker Build Old Games [${#OLD_GAMES[@]}] ====="
            for (( i = 0; i < ${#OLD_GAMES[@]}; i++ ));
            do
                Fn_Build_Docker $1 ${OLD_GAMES[$i]} $3 $4 $5 $6 $GIT_TAG $(expr $i + 1) ${#OLD_GAMES[@]}
            done
        else
            REPO_NAME="gameserver/$2"
            FULL_NAME="$3/$REPO_NAME"

            if [ "$8" = "" ] && [ "$9" = "" ]; then
                printf "Docker Building %12s ...\r" $2
            else
                printf "Docker Building %12s [%02d/%02d] ...\r" $2 $8 $9
            fi

            local start=$SECONDS

            # Copy File to build resource
            cp ${4}/gameserver/$2/file/$2 $5/

            docker pull example.invalid/base

            if [ "$?" = "0" ]; then
                docker build --target app --cache-from example.invalid/base \
                    -t ${FULL_NAME}:$TAG_NAME \
                    -t ${FULL_NAME}:${LATEST_STR} \
                    --build-arg APP_TYPE=gameserver \
                    --build-arg APP_NAME=$2 \
                    $5 || { echo "Docker Build $2 Fail" ; exit 1; }
            else
                docker build \
                    -t ${FULL_NAME}:$TAG_NAME \
                    -t ${FULL_NAME}:${LATEST_STR} \
                    --build-arg APP_TYPE=gameserver \
                    --build-arg APP_NAME=$2 \
                    $5 || { echo "Docker Build $2 Fail" ; exit 1; }
            fi

            local duration=$(($SECONDS - $start))

            if [ "$8" = "" ] && [ "$9" = "" ]; then
                printf "Docker Built %15s Successful ... (%dm%ds) (%s)\n" $2 $(($duration/60)) $(($duration%60)) ${FULL_NAME}:$TAG_NAME
            else
                printf "Docker Built %15s [%02d/%02d] Successful ... (%dm%ds) (%s)\n" $2 $8 $9 $(($duration/60)) $(($duration%60)) ${FULL_NAME}:$TAG_NAME
            fi

            # Delete File from build resource
            rm $5/$2
        fi
    else
        echo "Unknown Type, Please give [all/common/game]"
    fi
}

########################################################################################
# Fn_Push_Docker
#   Push Docker Images to Registry
#
# @Param $1 Type                    [all/common/game]
# @Param $2 Name
# @Param $3 Repository Domain                    
# @Param $4 Repository Username      
# @Param $5 Repository Password
# @Param $6 Build Number
# @Param $7 Git Tag
# @Param $8 Build Count
# @Param $9 Total Count
# 
# Usage: 
#     Fn_Build_Docker $TYPE $NAME $REGISTERY $RELEASE_FOLDER $DOCKER_FOLDER $BUILD_NAME $GIT_TAG
# Example:
#     Fn_Build_Docker all all $REGISTERY $RELEASE_FOLDER $DOCKER_FOLDER $BUILD_NAME $GIT_TAG
#     Fn_Build_Docker common dbserver $REGISTERY $RELEASE_FOLDER $DOCKER_FOLDER $BUILD_NAME $GIT_TAG
#     Fn_Build_Docker common all $REGISTERY $RELEASE_FOLDER $DOCKER_FOLDER $BUILD_NAME $GIT_TAG
#     Fn_Build_Docker game dfdc $REGISTERY $RELEASE_FOLDER $DOCKER_FOLDER $BUILD_NAME $GIT_TAG
#     Fn_Build_Docker game slot $REGISTERY $RELEASE_FOLDER $DOCKER_FOLDER $BUILD_NAME $GIT_TAG
#     Fn_Build_Docker game all $REGISTERY $RELEASE_FOLDER $DOCKER_FOLDER $BUILD_NAME $GIT_TAG
#
########################################################################################
function Fn_Push_Docker() {
    TAG_NAME=$7
    if [ "$TAG_NAME" = "" ] || [ "$TAG_NAME" = "-" ]; then
        TAG_NAME=$6
    fi

    LATEST_STR="latest"
    GIT_TAG="-"
    if [ "$7" = "" ] || [ "$7" = "-" ]; then
        LATEST_STR="beta"
    else
        GIT_TAG=$7
    fi

    # docker login -u $4 -p $5 $3

    if [ "$1" = "all" ]; then
        Fn_Push_Docker common all $3 $4 $5 $6 $GIT_TAG
        Fn_Push_Docker game all $3 $4 $5 $6 $GIT_TAG
        return
    elif [ "$1" = "common" ]; then
        if [ "$2" = "all" ]; then
            echo "===== Docker Push Commons [${#COMMONS[@]}] ====="
            for (( i = 0; i < ${#COMMONS[@]}; i++ ));
            do
                Fn_Push_Docker $1 ${COMMONS[$i]} $3 $4 $5 $6 $GIT_TAG $(expr $i + 1) ${#COMMONS[@]}
            done
            return
        else
            REPO_NAME=$2
            FULL_NAME="$3/$REPO_NAME"
        fi
    elif [ "$1" = "game" ]; then
        if [ "$2" = "all" ]; then
            Fn_Push_Docker game slot $3 $4 $5 $6 $GIT_TAG
            Fn_Push_Docker game pachislot $3 $4 $5 $6 $GIT_TAG
            Fn_Push_Docker game poker $3 $4 $5 $6 $GIT_TAG
            Fn_Push_Docker game old $3 $4 $5 $6 $GIT_TAG
            return
        elif [ "$2" = "slot" ]; then
            echo "===== Docker Push Slot Games [${#SLOT_GAMES[@]}] ====="
            for (( i = 0; i < ${#SLOT_GAMES[@]}; i++ ));
            do
                Fn_Push_Docker $1 ${SLOT_GAMES[$i]} $3 $4 $5 $6 $GIT_TAG $(expr $i + 1) ${#SLOT_GAMES[@]}
            done
            return
        elif [ "$2" = "pachislot" ]; then
            echo "===== Docker Pachi Slot Games [${#PACHI_SLOT_GAMES[@]}] ====="
            for (( i = 0; i < ${#PACHI_SLOT_GAMES[@]}; i++ ));
            do
                Fn_Push_Docker $1 ${PACHI_SLOT_GAMES[$i]} $3 $4 $5 $6 $GIT_TAG $(expr $i + 1) ${#PACHI_SLOT_GAMES[@]}
            done
        elif [ "$2" = "poker" ]; then
            echo "===== Docker Push Poker Games [${#POKER_GAMES[@]}] ====="
            for (( i = 0; i < ${#POKER_GAMES[@]}; i++ ));
            do
                Fn_Push_Docker $1 ${POKER_GAMES[$i]} $3 $4 $5 $6 $GIT_TAG $(expr $i + 1) ${#POKER_GAMES[@]}
            done
            return
        elif [ "$2" = "old" ]; then
            echo "===== Docker Push Old Games [${#OLD_GAMES[@]}] ====="
            for (( i = 0; i < ${#OLD_GAMES[@]}; i++ ));
            do
                Fn_Push_Docker $1 ${OLD_GAMES[$i]} $3 $4 $5 $6 $GIT_TAG $(expr $i + 1) ${#OLD_GAMES[@]}
            done
            return
        else
            REPO_NAME="gameserver/$2"
            FULL_NAME="$3/$REPO_NAME"
        fi
    else
        echo "Unknown Type, Please give [all/common/game]"
        return
    fi

    if [ "$8" = "" ] && [ "$9" = "" ]; then
        printf "Docker Pushing %14s ...\r" $2
    else
        printf "Docker Pushing %14s [%02d/%02d] ...\r" $2 $8 $9
    fi

    local start=$SECONDS

    docker push ${FULL_NAME}:${TAG_NAME} || { echo "Docker Push ${FULL_NAME}:${TAG_NAME} Fail" ; exit 1; }
    docker push ${FULL_NAME}:${LATEST_STR} || { echo "Docker Push ${FULL_NAME}:${LATEST_STR} Fail" ; exit 1; }

    local duration=$(($SECONDS - $start))

    if [ "$8" = "" ] && [ "$9" = "" ]; then
        printf "Docker Pushed %15s Successful ... (%dm%ds)\n" $2 $(($duration/60)) $(($duration%60))
    else
        printf "Docker Pushed %15s [%02d/%02d] Successful ... (%dm%ds)\n" $2 $8 $9 $(($duration/60)) $(($duration%60))
    fi
}

########################################################################################
# Fn_Echo_Version
#   echo version to output folder
#
# @Param $1 Version     
# @Param $2 Type                    [all/common/game]
# @Param $3 Name
# @Param $4 Output Folder
# 
# Usage: 
#     Fn_Echo_Version $BUILD_NAME $TYPE $NAME $DIR
# Example:
#     Fn_Echo_Version $BUILD_NAME all all .version
#     Fn_Echo_Version $BUILD_NAME common dbserver .version
#     Fn_Echo_Version $BUILD_NAME common all .version
#     Fn_Echo_Version $BUILD_NAME game dfdc .version
#     Fn_Echo_Version $BUILD_NAME game slot .version
#     Fn_Echo_Version $BUILD_NAME game all .version
#
########################################################################################
function Fn_Echo_Version {

    if [ "$2" = "all" ]; then
        Fn_Echo_Version $1 common all $4
        Fn_Echo_Version $1 game all $4
        return
    elif [ "$2" = "common" ]; then
        if [ "$3" = "all" ]; then
            for (( i = 0; i < ${#COMMONS[@]}; i++ ));
            do
                Fn_Echo_Version $1 common ${COMMONS[$i]} $4
            done
            return
        fi
    elif [ "$2" = "game" ]; then
        if [ "$3" = "all" ]; then
            Fn_Echo_Version $1 game slot $4
            Fn_Echo_Version $1 game pachislot $4
            Fn_Echo_Version $1 game poker $4
            Fn_Echo_Version $1 game old $4
            return
        elif [ "$3" = "slot" ]; then
            for (( i = 0; i < ${#SLOT_GAMES[@]}; i++ ));
            do
                Fn_Echo_Version $1 game ${SLOT_GAMES[$i]} $4
            done
            return
        elif [ "$3" = "pachislot" ]; then
            for (( i = 0; i < ${#PACHI_SLOT_GAMES[@]}; i++ ));
            do
                Fn_Echo_Version $1 game ${PACHI_SLOT_GAMES[$i]} $4
            done
            return
        elif [ "$3" = "poker" ]; then
            for (( i = 0; i < ${#POKER_GAMES[@]}; i++ ));
            do
                Fn_Echo_Version $1 game ${POKER_GAMES[$i]} $4
            done
            return
        elif [ "$3" = "old" ]; then
            for (( i = 0; i < ${#OLD_GAMES[@]}; i++ ));
            do
                Fn_Echo_Version $1 game ${OLD_GAMES[$i]} $4
            done
            return
        fi
    fi

    mkdir -p $4/$2

    echo $1 > $4/$2/$3
}

########################################################################################
# Fn_Join
#   Join Array to String with give Sep
#
# @Param $1 Sep
# @Param $2 Array
########################################################################################
function Fn_Join {
    local IFS="$1";
    shift;
    echo "$*"
}

########################################################################################
# Fn_Split
#   Split String to Array with give Sep
#
# @Param $1 Sep
# @Param $2 String
########################################################################################
function Fn_Split {
    IFS="$1" # space is set as delimiter
    read -ra ADDR <<< "$2" # str is read into an array as tokens separated by IFS
    echo ${ADDR[@]}
    # for i in "${ADDR[@]}"; do # access each element of array
    #     echo "$i"
    # done
}

########################################################################################
# Fn_ParseImageTag
#   Parse Image tag to [Type], [Name], [Tag]
#
# @Param $1 Image Tag String
########################################################################################
function Fn_ParseImageTag {
    local TYPE=""
    local NAME=""
    local TAG=""

    # echo $1

    IFS=":"
    read -ra STR1 <<< "$1"
    # echo ${#STR1[@]}
    # echo ${STR1[@]}

    TAG=${STR1[1]}

    IFS="/"
    read -ra STR2 <<< "${STR1[0]}"
    # echo ${#STR2[@]}
    # echo ${STR2[@]}

    if [ "${STR2[1]}" = "gameserver" ]; then
        TYPE="game"
        NAME=${STR2[2]}
    else
        TYPE="common"
        NAME=${STR2[1]}
    fi

    # echo TYPE = $TYPE
    # echo NAME = $NAME
    # echo TAG = $TAG

    echo $TYPE $NAME $TAG
}