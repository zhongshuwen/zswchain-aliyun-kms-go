#!/bin/sh
goreleaser build --single-target --rm-dist --snapshot -o ./dist/zswchain-aliyun-kms-go

source ./.env.sh
./dist/zswchain-aliyun-kms-go
