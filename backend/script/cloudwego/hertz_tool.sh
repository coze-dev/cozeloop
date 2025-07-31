#!/usr/bin/env bash

SCRIPT_PATH="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_PATH/../../.."
PROJECT_PATH=$(pwd)
THRIFT_PATH=${PROJECT_PATH}/idl/thrift


API_IDL=$THRIFT_PATH/coze/loop/apis/coze.loop.apis.thrift
API_PATH=${PROJECT_PATH}/backend/api

if [ ! -f "$API_IDL" ]; then
    echo "Error: API_IDL file not found at $API_IDL"
    exit 1
fi

cd "$API_PATH"

cloudwego_hz update -enable_extends -handler_dir=handler -model_dir=model -use=github.com/coze-dev/coze-loop/backend/kitex_gen --customize_package=tpl/package.yaml -idl $API_IDL

if ! command -v goimports &> /dev/null; then
    echo "Starting installation of goimports..."
    go install golang.org/x/tools/cmd/goimports@latest
fi

goimports -w $API_PATH/handler/coze/loop/apis