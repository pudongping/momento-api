#!/usr/bin/env bash
# set -u 设置如果引用未定义的变量则报错退出，管道中任一命令失败则报错退出
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

# 将数据库配置信息修改成你自己的
DB_Host=127.0.0.1
DB_Port=3306
DB_Name=momento
DB_Username=root
DB_Password=123456
# 表生成的模型存放的目录
Model_Dir=./model

# 这是我在 goctl 1.9.2 版本的模版文件中修改的 module 占位名称
Tpl_GoCTL_Module_Name="your-project-module-name"

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

# 从 go.mod 中获取 module 名称
function get_module_from_go_mod() {
  # 输出 module 名称到 stdout，失败返回非 0
  if [[ "$(uname)" != "Darwin" ]]; then
    echo "error: only supported on macOS" >&2
    return 2
  fi
  if [[ ! -f "go.mod" ]]; then
    echo "error: go.mod not found in $PWD" >&2
    return 3
  fi
  local module
  module=$(awk '/^module[[:space:]]+/ {print $2; exit}' go.mod || true)
  if [[ -z "$module" ]]; then
    echo "error: failed to parse module from go.mod" >&2
    return 4
  fi
  printf '%s' "$module"
}

# 在指定目录下递归替换所有文件中的指定字符串
function replace_in_dir() {
  # 参数：目录 path, module, word
  local dir="${1:-}"
  local module="${2:-}"
  local word="${3:-}"
  if [[ -z "$dir" || -z "$module" || -z "$word" ]]; then
    echo "usage: replace_in_dir <dir> <module> <word>" >&2
    return 1
  fi
  if [[ "$(uname)" != "Darwin" ]]; then
    echo "error: only supported on macOS" >&2
    return 2
  fi
  if [[ ! -d "$dir" ]]; then
    echo "error: directory \`$dir\` not found" >&2
    return 3
  fi

  # 转义 replacement 中可能出现的 &，避免 sed 将其当作匹配字符串的引用
  local esc_module="${module//&/\\&}"
  local count=0

  # 遍历目录下的所有文件（\`find\` 输出 NUL 分隔，避免文件名问题）
  while IFS= read -r -d '' file; do
    # 只在包含目标字符串时执行替换（避免把所有文件都 touch）
    if grep -q -- $word "$file" 2>/dev/null; then
      echo "replacing in: $file"
      sed -i '' "s|$word|$esc_module|g" "$file"
      count=$((count + 1))
    fi
  done < <(find "$dir" -type f -print0)

  if [[ $count -eq 0 ]]; then
    echo "no occurrences of '$word' found in \`$dir\`"
  else
    echo "replacement completed. files modified: $count"
  fi
}

# 替换 internal/handler 目录下的 module 名称
function replace_module_api() {
  module="$(get_module_from_go_mod)"
  replace_in_dir ./internal/handler "$module" "$Tpl_GoCTL_Module_Name"
}

# 替换模型目录下的 module 名称
function replace_module_model() {
  module="$(get_module_from_go_mod)"
  replace_in_dir "$Model_Dir" "$module" "$Tpl_GoCTL_Module_Name"
}

# 生成 go api 代码
function go_zero_gen_api() {
  _ensure_goctl_exec
  $GoCTL_CMD api go -api ./dsl/*.api -dir . --style=goZero --home ./goctlTemplates/1.9.2
  sleep 1
  replace_module_api
}

# 生成 md 文档
function go_zero_gen_doc_md() {
  _ensure_goctl_exec
  $GoCTL_CMD api doc --dir ./
}

# 模版初始化
function go_zero_tpl_init() {
  _ensure_goctl_exec
  $GoCTL_CMD template init
}

# goctl 生成模型文件
function go_zero_gen_model() {
  _ensure_goctl_exec

  # 判断存放模型的目录是否存在，不存在则创建
  if [[ ! -d "$Model_Dir" ]]; then
    mkdir -p "$Model_Dir"
  fi

  db_table="${1:-}"
  if [[ -z "$db_table" ]]; then
    echo -e "\033[31m 请提供表名参数！ \033[0m"
    exit 1
  fi

  echo -e "\033[32m 正在对数据库【 $DB_Name 】中的 【 $db_table 】数据表创建模型 \033[0m"

  datasource_url="${DB_Username}:${DB_Password}@tcp(${DB_Host}:${DB_Port})/${DB_Name}"

  # --cache=true 代表生成带缓存的模型
  # --prefix='momento_api:cache:' 缓存前缀
  # --ignore-columns 需要忽略的字段，插入和更新时需要同时忽略的字段，如 create_time、update_time 等等
  $GoCTL_CMD model mysql datasource --url="$datasource_url" \
    --table="$db_table" \
    --dir="$Model_Dir" \
    --home=./goctlTemplates/1.9.2 \
    --style=goZero \
    --ignore-columns=''

  sleep 1
  replace_module_model
}

function usage() {
    cat <<EOF
Usage: ./local_run.sh {model|genapi|mddoc|tplinit|goctl ...}
Commands:
  model <table_name>    生成模型文件，第二个参数为表名 比如： ./local_run.sh model users
  genapi                生成 go api 代码
  mddoc                 生成 markdown 文档
  tplinit               模版初始化
  goctl ...             直接透传到 goctl，可携带任意参数，比如： ./local_run.sh goctl env
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
  model)
    # 生成模型文件，第二个参数为表名
    go_zero_gen_model $2
    ;;
  *)
    # 未提供任何参数时，打印用法信息并退出非零
    usage
    exit 1
    ;;
esac