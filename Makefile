.PHONY: goctlenv
# 查看 goctl 配置信息
goctlenv:
	$(info ******************** goctlenv ********************)
	@echo "---------> process [goctlenv] \r\n"
	./local_run.sh goctl env
	@echo "\r\n---------> processed"

###############################################################

.PHONY: api
# 生成api文件
api:
	$(info ******************** api ********************)
	@echo "---------> process build [api] \r\n"
	./local_run.sh genapi
	@echo "\r\n---------> processed"