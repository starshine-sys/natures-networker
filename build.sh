#!/bin/sh
CGO_ENABLED=0 go build -o natures-networker -v -ldflags="-buildid= -X github.com/starshine-sys/natures-networker/common.Version=`git rev-parse --short HEAD`"
strip natures-networker
