#!/bin/bash

PKG_LIST=$(go list ./... | 
    grep -v /vendor/ | 
    grep -v /cmd | 
    grep -v /game/bjl | 
    grep -v /game/brnn | 
    grep -v /game/by | 
    grep -v /game/cd21 |
    grep -v /game/dchjk| 
    grep -v /game/dczz| 
    grep -v /game/ddz | 
    grep -v /game/ermj | 
    grep -v /game/hbsl | 
    grep -v /game/hhnn | 
    grep -v /game/hhdz | 
    grep -v /game/jdesyd | 
    grep -v /game/jlbpk | 
    grep -v /game/lhd | 
    grep -v /game/sg | 
    grep -v /pkg/mylog
)

go test -count=1 -short -v -coverprofile=.cover ${PKG_LIST} > output
cat output | grep -iE '^(ok|FAIL|\?)|FAIL'

go tool cover -func=.cover > cover
cat cover | grep total

if grep -q FAIL output ; then
  echo "Fail"
  exit 1
fi

echo "Success"
exit 0