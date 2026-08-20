.Phony: clear-swag swag-init

clear-swag:
	@echo "clearing docs dirrectory..."
	del /Q docs\* 2>nul || exit 0

swag-init:clear-swag
	@echo "Generate swagger docs..."
	swag init -g main.go
# 	netstat -ano | findstr :8082     # 找到占用 8082 的 PID                                          
#   taskkill /F /PID <PID>           # 强制终止                            