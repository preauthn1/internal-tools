########################################################################################
#!/bin/bash
# @Param $1 Namespace
# @Param $2 ReleaseName
# @Param $3 Chart
# @Param $4 BaseFolder
# @Param $5 BuildName
########################################################################################

BASEPATH=$(dirname $0)

. $BASEPATH/library.sh


########################################################################################
# ProtoFiles
# @Param $1 BaseFolder
# @Param $2 ReleaseName
# @Param $3 Namespace
########################################################################################
function ProtoFiles {
    local base="protodata"
    local folders=$(ls $1/$base)

    declare -a files

    for key in $folders
    do
        cmname=$2-proto-$key
        proto=$(printf "%s/%s" $1/$base $key)
        md5sum $proto/* > $proto/md5sum

        # kubectl -n $3 delete cm $cmname
        kubectl -n $3 create cm $cmname --from-file=$proto

        if [ "$?" = "1" ]; then
            kubectl -n $3 create cm $cmname --from-file=$proto -o yaml --dry-run=true | kubectl -n $3 replace -f -
        fi

        # protofiles=$(ls $1/$base/$key | grep .bytes)
        # for file in $protofiles
        # do  
        #     name=${file::${#file}-6}
        #     files[${#files[@]}]=$(printf "proto.%s.%s=%s/%s/%s.bytes\n" $key $name $1/$base $key $name)
        # done
    done

    echo ${files[@]}
}

function ValueFiles {
    local base="values"
    local names=$(ls $1/$base)

    declare -a files

    for file in $names
    do
        files[${#files[@]}]=$(printf "%s/%s" $1/$base ${file})
    done

    echo ${files[@]}
}

function EchoVersions {

    IFS=',' read -r -a array <<< "$2"
    if [ "${#array[@]}" != "0" ]; then
        for element in "${array[@]}"
        do
            IFS=':' read -r -a row <<< "$element"
            if [ "${#row[@]}" = "2" ]; then
                Fn_Echo_Version $1 ${row[0]} ${row[1]} $3
            fi
        done
    fi

}

function Versions {

    local groups=$(ls $1)

    declare -a versions

    local pattern="versions.%s.%s=%s"

    for grp in ${groups[@]}
    do
        local objects=$(ls $1/$grp)

        for obj in ${objects[@]}
        do
            versions[${#versions[@]}]=$(printf $pattern $grp $obj $(cat $1/$grp/$obj))
        done
    done

    echo ${versions[@]}
}

function GetImagesVersion {
    str=$(kubectl -n $1 get pod -l example.invalid/name=$2 -o go-template --template="{{range .items}}{{range .spec.containers}}{{.image}} {{end}}{{end}}")

    IFS=" " read -ra images <<< "$str"

    echo ${images[@]}
}

### Get Current Versions
Images=$(GetImagesVersion $1 $2)

for row in ${Images[@]}
do
    IFS=" " read -ra OBJ <<< "$(Fn_ParseImageTag "$row")"
    Fn_Echo_Version ${OBJ[2]} ${OBJ[0]} ${OBJ[1]} .version
done

# echo "===============  Current Versions  ==============="
# for row in $(Versions .version)
# do
#     echo $row
# done

EchoVersions $5 $CI_BUILD_TARGETS .version
ProtoFiles $4 $2 $1

values=$(ValueFiles $4)
sets=$(Versions .version)

echo "===============  New Versions  ==============="
for row in ${sets[@]}
do
    echo $row
done

sets[${#sets[@]}]=$(printf "build=%s" $5)

helm upgrade $2 $3 -i -n $1 --set $(Fn_Join , ${sets[@]}) -f $(Fn_Join , ${values[@]})
# helm template --debug $3 --set-file $(Fn_Join , ${setfiles[@]}) --set $(Fn_Join , ${sets[@]}) -f $(Fn_Join , ${values[@]})
# helm template --debug $3 --set $(Fn_Join , ${sets[@]}) -f $(Fn_Join , ${values[@]})