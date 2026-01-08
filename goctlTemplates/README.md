# 修改了官方默认的模版文件

所有的修改变动都有对应的备份文件，文件名以 `.bak` 结尾，如果需要恢复到官方模版，那么可以：

1. 找出所有以 `.bak` 结尾的备份文件，eg： `./1.9.2/api/handler.tpl.bak`
2. 将备份文件覆盖对应的模版文件，eg： `mv ./1.9.2/api/handler.tpl.bak ./1.9.2/api/handler.tpl`
3. 删除备份文件（可选），eg： `rm ./1.9.2/api/handler.tpl.bak`
4. 重新运行 goctl 相关命令生成代码

## 变更记录

> 生成了 handler 文件之后，你需要手动将 `your-project-module-name` 替换为你自己的项目模块名（见项目根目录下 `go.mod` 中的 `module`）
> 
> 或者可以执行项目根目录下的 `local_run.sh` 脚本中的 `replace_module_api` 方法
> 默认当在项目根目录下，执行 `./local_run.sh genapi` 命令时，会自动调用 `replace_module_api` 方法，详见脚本内容

### 修改一

调整了 `./1.9.2/model/model.tpl` 模版文件，扩展了一些额外的方法
对应的官方模版文件备份为 `./1.9.2/model/model.tpl.bak`

### 修改二

调整了 `./1.9.2/model` 目录下的部分模版文件。详见以 `.bak` 结尾的备份文件