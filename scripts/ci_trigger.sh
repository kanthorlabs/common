#!/bin/sh
set -e

SCRIPT_DIR="$( cd "$( dirname "$0" )" && pwd )";

find . -type f -name 'checksum' -exec git add {} \;

NOW=$(date +%Y.%-m%d.%-H%M)
HAHSTAGS=${HAHSTAGS:-""}

git commit --allow-empty -m "ci($NOW): ✨🐛🚨 $HAHSTAGS"

TARGET=${1:-origin}
echo "---------------------------"
printf "Pushing... $NOW --> %s\n" "$TARGET"
echo "---------------------------"
git push "$TARGET"  