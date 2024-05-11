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
  go test -timeout 30s --count=1 -cover -coverprofile cover.out $(go list ./... | grep github.com/kanthorlabs/common | grep -v 'github.com/kanthorlabs/common/\(commands\|logging\|testdata\|testify\)')
fi

if [ "$CI" = "" ];
then
  find . -type f \( -name '*.go' -o -name '*.rego'  \) ! -path "./vendor/*" -exec sha256sum {} \; | sort -k 2 | sha256sum | cut -d  ' ' -f1 > ./checksum
fi
