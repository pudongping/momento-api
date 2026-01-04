#!/usr/bin/env bash
set -euo pipefail

User_Path='/Users/pudongping'

Go_Version=go1.25.5
Go_Proxy=https://goproxy.cn,direct
Go_Root=$User_Path/go/sdk/$Go_Version
Go_Path=$User_Path/go/gomodule/momento-api

# goctl command path
GoCTL_CMD=$User_Path/go/gomodule/momento-api/bin/goctl

function _ensure_goctl_exec() {
  if [ ! -x "$GoCTL_CMD" ]; then
    echo "goctl not found or not executable: $GoCTL_CMD" >&2
    exit 1
  fi
}

function goctl_exec() {
  shift  # remove the "goctl" command name
  _ensure_goctl_exec
  exec "$GoCTL_CMD" "$@"
}

# 生成 go api 代码
function go_zero_gen_api() {
    GoCTL_CMD api go -api ./dsl/*.api -dir . --style=goZero --home ./goctlTemplates/1.9.2
}

# 生成 md 文档
function go_zero_gen_doc_md() {
    GoCTL_CMD api doc --dir ./
}

# 模版初始化
function go_zero_tpl_init() {
    GoCTL_CMD template init
}

# Dispatch based on first argument
cmd=${1:-}

case "$cmd" in
  genapi)
    go_zero_gen_api
    ;;
  mddoc)
    go_zero_gen_doc_md
    ;;
  goctl)
    goctl_exec "$@"
    ;;
  "")
    echo "Usage: ./local_run.sh {gengo|goctl ...|<forward to goctl>}" >&2
    exit 1
    ;;
  *)
    exit 1
    ;;
esac