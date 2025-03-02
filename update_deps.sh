#!/bin/bash

REPO_DIR=$(mktemp -d)
git clone https://github.com/dekinsq/whatsmeow $REPO_DIR
cd $REPO_DIR

LATEST_COMMIT=$(git log --pretty=format:'%H' -n 1)
COMMIT_DATE=$(git log -1 --format=%cd --date=format:"%Y%m%d%H%M%S")
PSEUDO_VERSION="v0.0.0-${COMMIT_DATE}-${LATEST_COMMIT:0:12}"

# 更新 go.mod
sed -i "s|replace go.mau.fi/whatsmeow => github.com/dekinsq/whatsmeow.*|replace go.mau.fi/whatsmeow => github.com/dekinsq/whatsmeow $PSEUDO_VERSION|" go.mod

# 清理并同步
go clean -modcache
GOPROXY=direct go get -u github.com/dekinsq/whatsmeow@$PSEUDO_VERSION
go mod tidy

rm -rf $REPO_DIR
