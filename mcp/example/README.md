# Notion MCP SDK 示例

这个目录包含了 Notion MCP SDK 的完整示例和集成测试。

## 文件说明

- `basic_usage.go` - 基础使用示例
- `advanced_usage.go` - 高级使用示例
- `http_server.go` - HTTP 服务器示例
- `integration_test.go` - 集成测试
- `go.mod` - Go 模块文件

## 环境设置

### 方法 1: 使用 .env 文件 (推荐)

1. 复制环境变量示例文件：
```bash
cp env.example .env
```

2. 编辑 `.env` 文件，填入实际的配置值：
```bash
# Notion 集成令牌 (必需)
NOTION_TOKEN=your_notion_integration_token_here

# 父页面 ID (可选，用于创建页面示例)
NOTION_PARENT_PAGE_ID=your_parent_page_id_here

# 数据库 ID (可选，用于数据库操作示例)
NOTION_DATABASE_ID=your_database_id_here

# HTTP 服务器端口 (可选，默认 8080)
PORT=8080
```

### 方法 2: 使用系统环境变量

```bash
export NOTION_TOKEN="your_notion_integration_token"
export NOTION_PARENT_PAGE_ID="your_parent_page_id"  # 可选，用于创建页面示例
```

### 获取 Notion Token

1. 访问 [Notion Developers](https://www.notion.so/my-integrations)
2. 创建新的集成
3. 复制内部集成令牌

### 获取父页面 ID

1. 在 Notion 中打开要作为父页面的页面
2. 从 URL 中复制页面 ID（32位字符串，用连字符分隔）

## 运行示例

### 安装依赖

```bash
cd example
go mod tidy
```

### 基础使用示例

```bash
# 使用 .env 文件 (推荐)
go run main.go basic

# 或者直接运行
go run basic_usage.go
```

### 高级使用示例

```bash
# 使用 .env 文件 (推荐)
go run main.go advanced

# 或者直接运行
go run advanced_usage.go
```

### HTTP 服务器示例

```bash
# 使用 .env 文件 (推荐)
go run main.go server

# 或者直接运行
go run http_server.go
```

服务器将在 `http://localhost:8080` 启动，提供以下端点：

- `GET /health` - 健康检查
- `POST /mcp` - MCP 协议端点
- `GET /tools` - 工具列表
- `GET /resources` - 资源列表
- `GET /search?q=query` - 搜索示例
- `GET /workspace` - 工作区信息

### 集成测试

```bash
# 使用 .env 文件运行测试
go test -v

# 或者使用统一入口运行集成测试
go run main.go test

# 或者直接运行集成测试脚本
go run integration_test.go test
```

## 功能演示

### 1. 搜索功能

```go
// 简单搜索
result, err := sdk.QuickSearch(ctx, "项目")

// 高级搜索
searchParams := &mcp.NotionSearchParams{
    Query:       "项目",
    Filter:      "page",
    SortBy:      "last_edited_time",
    SortOrder:   "descending",
    PageSize:    10,
}
result, err := sdk.Search(ctx, searchParams)
```

### 2. 创建页面

```go
// 简单创建
page, err := sdk.QuickCreatePage(ctx, parentID, "标题", "内容")

// 高级创建
createParams := &mcp.NotionCreatePageParams{
    ParentID: parentID,
    Title:    "标题",
    Content:  "内容",
    Icon: &mcp.Icon{
        Type:  "emoji",
        Emoji: "🚀",
    },
    Properties: map[string]interface{}{
        "status": "进行中",
    },
}
page, err := sdk.CreatePage(ctx, createParams)
```

### 3. 添加内容

```go
// 添加文本
sdk.QuickAppendText(ctx, pageID, "段落内容")

// 添加标题
sdk.QuickAppendHeading(ctx, pageID, "标题", 2)

// 添加代码块
sdk.QuickAppendCode(ctx, pageID, "fmt.Println(\"Hello\")")

// 添加引用
sdk.QuickAppendQuote(ctx, pageID, "引用内容")

// 添加标注
sdk.QuickAppendCallout(ctx, pageID, "重要信息")

// 添加待办事项
sdk.QuickAppendTodo(ctx, pageID, "完成任务")

// 添加列表
sdk.QuickAppendBulletList(ctx, pageID, "列表项")
sdk.QuickAppendNumberedList(ctx, pageID, "有序列表项")
```

### 4. 更新页面

```go
updateParams := &mcp.NotionUpdatePageParams{
    PageID:  pageID,
    Title:   "新标题",
    Content: "新内容",
    Properties: map[string]interface{}{
        "status": "已完成",
    },
}
page, err := sdk.UpdatePage(ctx, updateParams)
```

### 5. 获取工作区信息

```go
workspaceInfo, err := sdk.GetWorkspaceInfo(ctx)
fmt.Printf("总页面数: %v\n", workspaceInfo["totalPages"])
fmt.Printf("总数据库数: %v\n", workspaceInfo["totalDatabases"])
```

## 配置管理

```go
// 创建自定义配置
config := &mcp.MCPConfig{
    NotionToken:    "your-token",
    ServerName:     "my-server",
    ServerVersion:  "1.0.0",
    DefaultPageSize: 20,
    MaxRetries:     5,
    Timeout:        60,
}
sdk, err := mcp.NewNotionMCPSDK(config)

// 运行时修改配置
sdk.SetDefaultPageSize(50)
sdk.SetTimeout(120)
sdk.SetMaxRetries(10)
```

## 错误处理

SDK 提供了详细的错误信息：

```go
result, err := sdk.QuickSearch(ctx, "查询")
if err != nil {
    // 检查是否是 MCP 错误
    if mcpErr, ok := err.(*mcp.MCPError); ok {
        fmt.Printf("MCP 错误: %d - %s\n", mcpErr.Code, mcpErr.Message)
    } else {
        fmt.Printf("其他错误: %v\n", err)
    }
}
```

## 注意事项

1. **API 限制**: Notion API 有速率限制，建议在批量操作之间添加延迟
2. **权限**: 确保集成有足够的权限访问所需的页面和数据库
3. **页面 ID**: 页面 ID 是 32 位字符串，不包含连字符
4. **内容格式**: 支持 Markdown 格式的内容，但某些复杂格式可能需要特殊处理

## 故障排除

### 常见错误

1. **401 Unauthorized**: 检查 Notion Token 是否正确
2. **403 Forbidden**: 检查集成权限和页面访问权限
3. **429 Too Many Requests**: 减少请求频率，增加延迟
4. **400 Bad Request**: 检查参数格式和页面 ID

### 调试技巧

1. 启用详细日志
2. 检查网络连接
3. 验证 Notion 集成设置
4. 使用 Notion API 文档验证参数格式

## 更多信息

- [Notion API 文档](https://developers.notion.com/)
- [MCP 协议规范](https://modelcontextprotocol.io/)
- [项目 GitHub 仓库](https://github.com/tenz-io/notionapi)
