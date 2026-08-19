.Phony: clear-swag swag-init

clear-swag:
	@echo "clearing docs dirrectory..."
	del /Q docs\* 2>nul || exit 0

swag-init:clear-swag
	@echo "Generate swagger docs..."
	swag init -g main.go