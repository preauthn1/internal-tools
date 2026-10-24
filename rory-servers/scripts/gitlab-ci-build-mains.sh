#!/bin/bash
source .go/.variables
echo $LDFLAGS
EXTNAME=""
if [ "$GOOS" = "windows" ]; then
    EXTNAME=".exe"
fi

go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/apiserver/apiserver${EXTNAME} ./cmd/apiserver || { echo 'go build apiserver fail' ; exit 1; }
echo 'build apiserver successful'
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/dbserver/dbserver${EXTNAME} ./cmd/dbserver || { echo 'go build dbserver fail' ; exit 1; }
echo 'build dbserver successful'
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/dbserver/dbserver-phnum${EXTNAME} ./cmd/dbserver-phnum || { echo 'go build dbserver-phnum fail' ; exit 1; }
echo 'build dbserver-phnum successful'
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/notifyserver/notifyserver${EXTNAME} ./cmd/notifyserver || { echo 'go build notifyserver fail' ; exit 1; }
echo 'build notifyserver successful'
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/hallserver/hallserver${EXTNAME} ./cmd/hallserver || { echo 'go build hallserver fail' ; exit 1; }
echo 'build hallserver successful'
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/tfconfig/tfconfig${EXTNAME} ./cmd/tfconfig/main.go || { echo 'go build tfconfig fail' ; exit 1; }
echo 'build tfconfig successful'
go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/tfserver/tfserver${EXTNAME} ./cmd/tfserver/main.go || { echo 'go build tfserver fail' ; exit 1; }
echo 'build tfserver successful'

# go build -ldflags "$LDFLAGS" -i -o ${RELEASE_FOLDER}/newbydbserve/newbydbserver${EXTNAME} tools/newbydbserver/newbydbserver.go || { echo 'go build newbydbserver fail' ; exit 1; }

