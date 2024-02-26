#!/bin/sh
set -e

CI=${CI:-""}
CHECKSUM_FILE=./checksum
CHECKSUM_NEW=$(find . -type f \( -name '*.go' -o -name '*.rego'  \) ! -path "./vendor/*" -exec sha256sum {} \; | sort -k 2 | sha256sum | cut -d  ' ' -f1)
CHECKSUM_OLD=$(cat $CHECKSUM_FILE || true)

if [ "$CHECKSUM_NEW" != "$CHECKSUM_OLD" ];
then
  echo "--> gatekeeper/rego"
  opa test gatekeeper/rego/
  
  echo "--> coverage"
  go test --count=1 -cover -coverprofile cover.out  \
    ./cache/... \
    ./cipher/... \
    ./clock/... \
    ./configuration/... \
    ./gatekeeper/... \
    ./gateway/... \
    ./healthcheck/background/... \
    ./idx/... \
    ./project/... \
    ./passport/... \
    ./persistence/... \
    ./safe/... \
    ./utils/... \
    ./validator/...
fi

if [ "$CI" = "" ];
then
  find . -type f \( -name '*.go' -o -name '*.rego'  \) ! -path "./vendor/*" -exec sha256sum {} \; | sort -k 2 | sha256sum | cut -d  ' ' -f1 > ./checksum
fi
