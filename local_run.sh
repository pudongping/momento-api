#!/usr/bin/env bash
set -euo pipefail

# 切换到脚本所在目录（保证后续操作在脚本同级目录执行）
# 处理方式兼容被 symlink 的情况
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR" || {
  echo "failed to cd to script dir: $SCRIPT_DIR" >&2
  exit 1
}

# ---------- 用户和 Go 环境相关配置 ----------
# 注意：根据需要调整 User_Path 与 Go_Version
User_Path='/Users/pudongping'
Go_Version=go1.25.5
Go_Proxy=https://goproxy.cn,direct
Go_Root=$User_Path/go/sdk/$Go_Version
Go_Path=$User_Path/go/gomodule/momento-api

# goctl 二进制的绝对路径
GoCTL_CMD=$User_Path/go/gomodule/momento-api/bin/goctl

# ---------- 辅助函数 ----------
# 校验 goctl 是否存在且可执行；失败则退出
function _ensure_goctl_exec() {
  if [ ! -x "$GoCTL_CMD" ]; then
    echo "goctl not found or not executable: $GoCTL_CMD" >&2
    exit 1
  fi
}

# 执行 goctl 并把所有参数透传（第一个参数为 "goctl" 时使用）
# 使用 exec 替换当前 shell 进程，参数安全引用
function goctl_exec() {
  shift  # 移除命令名 "goctl"
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

function usage() {
    cat <<EOF
Usage: $0 {genapi|mddoc|tplinit|goctl ...}
Commands:
  genapi        生成 go api 代码
  mddoc        生成 markdown 文档
  tplinit      模版初始化
  goctl ...    直接透传到 goctl，可携带任意参数
EOF
}


# ---------- 命令分发 ----------
# 第一个参数决定要执行的子命令
cmd=${1:-}

case "$cmd" in
  genapi)
    # 生成 go api 代码
    go_zero_gen_api
    ;;
  mddoc)
    # 生成 markdown 文档
    go_zero_gen_doc_md
    ;;
  tplinit)
    # 模版初始化
    go_zero_tpl_init
    ;;
  goctl)
    # 直接透传到 goctl，可携带任意参数
    goctl_exec "$@"
    ;;
  "")
    # 未提供任何参数时，打印用法信息并退出非零
    usage
    exit 1
    ;;
  *)
    echo "Unknown command: $cmd" >&2
    exit 1
    ;;
esac