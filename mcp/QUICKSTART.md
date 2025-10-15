# Notion MCP SDK 快速开始指南

## 快速设置

### 1. 设置环境文件

```bash
cd mcp/example
make setup-env
```

这将创建 `.env` 文件，然后编辑它：

```bash
# 编辑 .env 文件
nano .env
```

填入您的 Notion 配置：

```bash
NOTION_TOKEN=your_notion_integration_token_here
NOTION_PARENT_PAGE_ID=your_parent_page_id_here
```

### 2. 安装依赖

```bash
go mod tidy
```

### 3. 运行示例

```bash
# 基础示例
make example

# 高级示例
make example-advanced

# HTTP 服务器
make example-server

# 集成测试
make example-test
```

## 获取 Notion Token

1. 访问 [Notion Developers](https://www.notion.so/my-integrations)
2. 创建新的集成
3. 复制内部集成令牌到 `.env` 文件

## 获取页面 ID

1. 在 Notion 中打开页面
2. 从 URL 复制页面 ID（32位字符串）

## 验证设置

```bash
# 检查环境配置
make env

# 运行基础测试
make example
```

## 使用统一入口

```bash
# 使用 main.go 统一入口
go run main.go basic      # 基础示例
go run main.go advanced   # 高级示例
go run main.go server     # HTTP 服务器
go run main.go test       # 集成测试
```

## 故障排除

### 常见问题

1. **NOTION_TOKEN 未设置**
   - 确保 `.env` 文件存在且包含正确的 token
   - 检查 token 是否有效

2. **权限错误**
   - 确保集成有访问所需页面的权限
   - 在 Notion 中邀请集成到页面

3. **依赖问题**
   - 运行 `go mod tidy` 安装依赖

### 获取帮助

```bash
make help
```

## 下一步

- 查看 [完整文档](README.md)
- 运行 [集成测试](example/integration_test.go)
- 探索 [高级功能](example/advanced_usage.go)
