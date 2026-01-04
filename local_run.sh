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

function goctl() {
  shift  # remove the "goctl" command name
  _ensure_goctl_exec
  exec "$GoCTL_CMD" "$@"
}

function gengo() {
    GoCTL_CMD api go -api ./dsl/*.api -dir . --style=goZero
}

# Dispatch based on first argument
cmd=${1:-}

case "$cmd" in
  gengo)
    gengo "$@"
    ;;
  goctl)
    goctl "$@"
    ;;
  "")
    echo "Usage: ./local_run.sh {gengo|goctl ...|<forward to goctl>}" >&2
    exit 1
    ;;
  *)
    exit 1
    ;;
esac