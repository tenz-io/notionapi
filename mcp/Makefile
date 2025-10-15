# Notion MCP SDK Makefile

.PHONY: help build test example clean lint format

# 默认目标
help:
	@echo "可用的命令:"
	@echo "  setup-env - 设置环境文件 (.env)"
	@echo "  build     - 构建项目"
	@echo "  test      - 运行测试"
	@echo "  example   - 运行基础示例"
	@echo "  example-advanced - 运行高级示例"
	@echo "  example-server - 启动 HTTP 服务器"
	@echo "  example-test - 运行集成测试"
	@echo "  lint      - 运行代码检查"
	@echo "  format    - 格式化代码"
	@echo "  clean     - 清理构建文件"
	@echo "  install   - 安装依赖"
	@echo "  env       - 显示环境信息"

# 构建项目
build:
	@echo "构建 Notion MCP SDK..."
	go build ./...

# 运行测试
test:
	@echo "运行测试..."
	go test -v ./...

# 运行集成测试
test-integration:
	@echo "运行集成测试..."
	@if [ -z "$$NOTION_TOKEN" ]; then \
		echo "错误: 请设置 NOTION_TOKEN 环境变量"; \
		exit 1; \
	fi
	cd example && go test -v -run TestNotionMCPSDK

# 运行示例
example:
	@echo "运行基础示例..."
	@if [ ! -f "example/.env" ]; then \
		echo "警告: 未找到 .env 文件，请确保设置了 NOTION_TOKEN 环境变量"; \
		if [ -z "$$NOTION_TOKEN" ]; then \
			echo "错误: 请设置 NOTION_TOKEN 环境变量或创建 .env 文件"; \
			exit 1; \
		fi; \
	fi
	cd example && go run main.go basic

# 运行高级示例
example-advanced:
	@echo "运行高级示例..."
	@if [ ! -f "example/.env" ]; then \
		echo "警告: 未找到 .env 文件，请确保设置了 NOTION_TOKEN 环境变量"; \
		if [ -z "$$NOTION_TOKEN" ]; then \
			echo "错误: 请设置 NOTION_TOKEN 环境变量或创建 .env 文件"; \
			exit 1; \
		fi; \
	fi
	cd example && go run main.go advanced

# 启动 HTTP 服务器示例
example-server:
	@echo "启动 HTTP 服务器示例..."
	@if [ ! -f "example/.env" ]; then \
		echo "警告: 未找到 .env 文件，请确保设置了 NOTION_TOKEN 环境变量"; \
		if [ -z "$$NOTION_TOKEN" ]; then \
			echo "错误: 请设置 NOTION_TOKEN 环境变量或创建 .env 文件"; \
			exit 1; \
		fi; \
	fi
	cd example && go run main.go server

# 运行集成测试脚本
example-test:
	@echo "运行集成测试脚本..."
	@if [ ! -f "example/.env" ]; then \
		echo "警告: 未找到 .env 文件，请确保设置了 NOTION_TOKEN 环境变量"; \
		if [ -z "$$NOTION_TOKEN" ]; then \
			echo "错误: 请设置 NOTION_TOKEN 环境变量或创建 .env 文件"; \
			exit 1; \
		fi; \
	fi
	cd example && go run main.go test

# 代码检查
lint:
	@echo "运行代码检查..."
	golangci-lint run

# 格式化代码
format:
	@echo "格式化代码..."
	go fmt ./...
	goimports -w .

# 清理构建文件
clean:
	@echo "清理构建文件..."
	go clean ./...

# 安装依赖
install:
	@echo "安装依赖..."
	go mod download
	go mod tidy

# 生成文档
docs:
	@echo "生成文档..."
	godoc -http=:6060

# 检查代码覆盖率
coverage:
	@echo "检查代码覆盖率..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

# 基准测试
benchmark:
	@echo "运行基准测试..."
	go test -bench=. ./...

# 检查依赖更新
deps-check:
	@echo "检查依赖更新..."
	go list -u -m all

# 更新依赖
deps-update:
	@echo "更新依赖..."
	go get -u ./...
	go mod tidy

# 验证模块
verify:
	@echo "验证模块..."
	go mod verify

# 运行所有检查
check: lint test format
	@echo "所有检查完成"

# 完整构建流程
all: clean install format lint test build
	@echo "完整构建流程完成"

# 开发环境设置
dev-setup:
	@echo "设置开发环境..."
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "安装 golangci-lint..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.54.2; \
	fi
	@if ! command -v goimports >/dev/null 2>&1; then \
		echo "安装 goimports..."; \
		go install golang.org/x/tools/cmd/goimports@latest; \
	fi
	@echo "开发环境设置完成"

# 发布准备
release-prep:
	@echo "准备发布..."
	@if [ -z "$$VERSION" ]; then \
		echo "错误: 请设置 VERSION 环境变量"; \
		exit 1; \
	fi
	git tag -a v$$VERSION -m "Release v$$VERSION"
	git push origin v$$VERSION

# 设置环境文件
setup-env:
	@echo "设置环境文件..."
	@if [ ! -f "example/.env" ]; then \
		cp example/env.example example/.env; \
		echo "已创建 example/.env 文件，请编辑并填入实际值"; \
	else \
		echo "example/.env 文件已存在"; \
	fi

# 显示环境信息
env:
	@echo "Go 版本: $$(go version)"
	@echo "Go 路径: $$(go env GOPATH)"
	@echo "Go 模块: $$(go env GOMOD)"
	@echo "当前目录: $$(pwd)"
	@echo "环境变量:"
	@echo "  NOTION_TOKEN: $$(if [ -n "$$NOTION_TOKEN" ]; then echo "已设置"; else echo "未设置"; fi)"
	@echo "  NOTION_PARENT_PAGE_ID: $$(if [ -n "$$NOTION_PARENT_PAGE_ID" ]; then echo "已设置"; else echo "未设置"; fi)"
	@echo "  .env 文件: $$(if [ -f "example/.env" ]; then echo "已存在"; else echo "不存在"; fi)"
