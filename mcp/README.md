# Notion MCP SDK

基于 [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) 的 Notion API SDK，提供简洁易用的接口来搜索和写入 Notion 内容。

## 特性

- 🚀 **易于集成**: 提供简洁的 API 接口，支持默认配置
- 🔍 **强大的搜索**: 支持全文搜索、过滤和排序
- ✍️ **丰富的写入**: 支持创建页面、更新内容、添加各种块类型
- 🛠️ **MCP 兼容**: 完全符合 MCP 协议规范
- 📦 **开箱即用**: 提供默认配置和便捷方法
- 🔧 **高度可配置**: 支持自定义配置和高级选项

## 快速开始

### 安装

```bash
go get github.com/tenz-io/notionapi/mcp
```

### 基础使用

```go
package main

import (
    "context"
    "log"
    
    "github.com/tenz-io/notionapi/mcp"
)

func main() {
    // 创建 SDK 实例
    sdk, err := mcp.NewNotionMCPSDKWithDefaults("your-notion-token")
    if err != nil {
        log.Fatal(err)
    }
    
    ctx := context.Background()
    
    // 搜索内容
    result, err := sdk.QuickSearch(ctx, "项目")
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("找到 %d 个结果", len(result.Results))
}
```

## 核心功能

### 1. 搜索功能

```go
// 简单搜索
result, err := sdk.QuickSearch(ctx, "查询内容")

// 高级搜索
searchParams := &mcp.NotionSearchParams{
    Query:       "项目",
    Filter:      "page",           // "page" 或 "database"
    SortBy:      "last_edited_time",
    SortOrder:   "descending",     // "ascending" 或 "descending"
    PageSize:    20,
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
    Title:    "页面标题",
    Content:  "页面内容",
    Icon: &mcp.Icon{
        Type:  "emoji",
        Emoji: "🚀",
    },
    Properties: map[string]any{
        "status": "进行中",
    },
}
page, err := sdk.CreatePage(ctx, createParams)
```

### 3. 添加内容

```go
// 添加文本段落
sdk.QuickAppendText(ctx, pageID, "段落内容")

// 添加标题
sdk.QuickAppendHeading(ctx, pageID, "标题", 2) // 1-3 级标题

// 添加代码块
sdk.QuickAppendCode(ctx, pageID, "fmt.Println(\"Hello\")")

// 添加引用
sdk.QuickAppendQuote(ctx, pageID, "引用内容")

// 添加标注
sdk.QuickAppendCallout(ctx, pageID, "重要信息")

// 添加待办事项
sdk.QuickAppendTodo(ctx, pageID, "完成任务")

// 添加列表
sdk.QuickAppendBulletList(ctx, pageID, "无序列表项")
sdk.QuickAppendNumberedList(ctx, pageID, "有序列表项")
```

### 4. 更新页面

```go
updateParams := &mcp.NotionUpdatePageParams{
    PageID:  pageID,
    Title:   "新标题",
    Content: "新内容",
    Properties: map[string]any{
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

## 配置选项

### 默认配置

```go
sdk, err := mcp.NewNotionMCPSDKWithDefaults("your-token")
```

### 自定义配置

```go
config := &mcp.MCPConfig{
    NotionToken:    "your-notion-token",
    ServerName:     "my-notion-mcp-server",
    ServerVersion:  "1.0.0",
    DefaultPageSize: 20,
    MaxRetries:     5,
    Timeout:        60, // 秒
}
sdk, err := mcp.NewNotionMCPSDK(config)
```

### 运行时配置

```go
sdk.SetDefaultPageSize(50)
sdk.SetTimeout(120)
sdk.SetMaxRetries(10)
```

## HTTP 服务器集成

SDK 可以轻松集成到 HTTP 服务器中：

```go
package main

import (
    "net/http"
    "github.com/tenz-io/notionapi/mcp"
)

func main() {
    sdk, _ := mcp.NewNotionMCPSDKWithDefaults("your-token")
    
    // 设置 MCP 端点
    http.HandleFunc("/mcp", sdk.HandleHTTPRequest)
    
    // 启动服务器
    http.ListenAndServe(":8080", nil)
}
```

## MCP 协议支持

SDK 完全支持 MCP 协议，提供以下工具：

- `notion_search` - 搜索 Notion 内容
- `notion_create_page` - 创建页面
- `notion_update_page` - 更新页面
- `notion_append_block` - 添加块内容

### MCP 请求示例

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "notion_search",
    "arguments": {
      "query": "项目",
      "filter": "page",
      "pageSize": 10
    }
  }
}
```

## 错误处理

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

### 常见错误代码

- `-32602`: 无效参数
- `-32603`: 内部错误
- `-32000`: Notion API 错误
- `-32001`: 认证错误
- `-32002`: 速率限制错误

## 示例和测试

查看 `example/` 目录获取完整示例：

- `basic_usage.go` - 基础使用示例
- `advanced_usage.go` - 高级使用示例
- `http_server.go` - HTTP 服务器示例
- `integration_test.go` - 集成测试

### 运行示例

```bash
cd example
export NOTION_TOKEN="your-token"
export NOTION_PARENT_PAGE_ID="your-parent-page-id"
go run basic_usage.go
```

## 环境变量

- `NOTION_TOKEN` - Notion 集成令牌（必需）
- `NOTION_PARENT_PAGE_ID` - 父页面 ID（用于创建页面示例）

## 注意事项

1. **API 限制**: Notion API 有速率限制，建议在批量操作之间添加延迟
2. **权限**: 确保集成有足够的权限访问所需的页面和数据库
3. **页面 ID**: 页面 ID 是 32 位字符串，不包含连字符
4. **内容格式**: 支持 Markdown 格式的内容，但某些复杂格式可能需要特殊处理

## 故障排除

### 常见问题

1. **401 Unauthorized**: 检查 Notion Token 是否正确
2. **403 Forbidden**: 检查集成权限和页面访问权限
3. **429 Too Many Requests**: 减少请求频率，增加延迟
4. **400 Bad Request**: 检查参数格式和页面 ID

### 调试技巧

1. 启用详细日志
2. 检查网络连接
3. 验证 Notion 集成设置
4. 使用 Notion API 文档验证参数格式

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License

## 相关链接

- [Notion API 文档](https://developers.notion.com/)
- [MCP 协议规范](https://modelcontextprotocol.io/)
- [项目 GitHub 仓库](https://github.com/tenz-io/notionapi)
