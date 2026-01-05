# 修改了官方默认的模版文件

> 要是不知道调整了哪些地方，可以查看 git commit 记录 3a26117c260e1ee32f94bcd05dddc5a709dc2eb6

## 修改一

调整了 `./1.9.2/model/model.tpl` 模版文件，扩展了一些额外的方法
对应的官方模版文件备份为 `./1.9.2/model/model.tpl.bak`

## 修改二

调整了 `./1.9.2/api/handler.tpl` 模版文件，调整了返回方法，需要请注意 import 中的变化
对应的官方模版文件备份为 `./1.9.2/api/handler.tpl.bak`

生成了 handler 文件之后，你需要手动将 `your-project-module-name` 替换为你自己的项目模块名（见项目根目录下 `go.mod` 中的 `module`）