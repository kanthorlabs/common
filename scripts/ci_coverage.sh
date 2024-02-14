#!/bin/sh
set -e

COVERAGE_EXPECTED=${COVERAGE_EXPECTED:-"80.0"}
COVERAGE_FILE=${COVERAGE_FILE:-"cover.out"}

if test -f $COVERAGE_FILE; then
  go tool cover -func $COVERAGE_FILE | grep total | awk '{print substr($3, 1, length($3)-1)}' > coverage.out

  COVERAGE_ACUTAL=$(cat coverage.out)

  if [ $(echo "${COVERAGE_ACUTAL} < ${COVERAGE_EXPECTED}" | bc) -eq 1 ]; 
  then
    echo "actual:$COVERAGE_ACUTAL < expected:$COVERAGE_EXPECTED"
    exit 1
  fi
else
  echo "$COVERAGE_FILE is not found"
fi