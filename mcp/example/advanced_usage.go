package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/tenz-io/notionapi/mcp"
)

func mainAdvanced() {
	// 加载环境变量
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// 从环境变量获取 Notion Token
	notionToken := os.Getenv("NOTION_TOKEN")
	if notionToken == "" {
		log.Fatal("NOTION_TOKEN environment variable is required")
	}

	// 创建自定义配置
	config := &mcp.MCPConfig{
		NotionToken:     notionToken,
		ServerName:      "my-notion-mcp-server",
		ServerVersion:   "1.0.0",
		DefaultPageSize: 20,
		MaxRetries:      5,
		Timeout:         60,
	}

	// 创建 SDK 实例
	sdk, err := mcp.NewNotionMCPSDK(config)
	if err != nil {
		log.Fatalf("Failed to create SDK: %v", err)
	}

	ctx := context.Background()

	// 示例 1: 高级搜索
	fmt.Println("=== 高级搜索示例 ===")

	// 搜索页面，按最后编辑时间降序排列
	searchParams := &mcp.NotionSearchParams{
		Query:     "项目",
		Filter:    "page",
		SortBy:    "last_edited_time",
		SortOrder: "descending",
		PageSize:  10,
	}

	searchResult, err := sdk.Search(ctx, searchParams)
	if err != nil {
		log.Printf("高级搜索失败: %v", err)
	} else {
		fmt.Printf("找到 %d 个页面结果\n", len(searchResult.Results))
		for i, item := range searchResult.Results {
			if i >= 5 { // 只显示前5个结果
				break
			}
			fmt.Printf("- %s (最后编辑: %s)\n", item.Title, item.LastEditedTime.Format("2006-01-02 15:04:05"))
		}
	}

	// 示例 2: 创建带属性的页面
	fmt.Println("\n=== 创建带属性的页面示例 ===")
	parentPageID := os.Getenv("NOTION_PARENT_PAGE_ID")
	if parentPageID != "" {
		// 创建带图标和属性的页面
		createParams := &mcp.NotionCreatePageParams{
			ParentID: parentPageID,
			Title:    "MCP SDK 高级测试页面",
			Content: `# MCP SDK 高级测试页面

这是一个通过 MCP SDK 创建的高级测试页面，支持 **Markdown** 格式。

## 功能特性

- ✅ 支持 Markdown 格式
- ✅ 支持多种属性类型
- ✅ 一次性创建页面和内容

### 代码示例

` + "```" + `go
// 创建页面示例
page, err := sdk.CreatePage(ctx, params)
` + "```" + `

> 这是一个引用块，展示 Markdown 功能。

- 列表项 1
- 列表项 2
- 列表项 3

1. 有序列表 1
2. 有序列表 2
3. 有序列表 3`,
			Properties: map[string]any{
				"status":   "进行中",
				"priority": "高",
				"progress": 75.5,
				"due_date": "2024-12-31",
				"created":  "2024-01-15",
			},
		}

		pageResult, err := sdk.CreatePage(ctx, createParams)
		if err != nil {
			log.Printf("创建带属性的页面失败: %v", err)
		} else {
			fmt.Printf("页面创建成功: %s\n", pageResult.URL)

			// 示例 3: 更新页面
			fmt.Println("\n=== 更新页面示例 ===")
			time.Sleep(2 * time.Second) // 等待一下

			updateParams := &mcp.NotionUpdatePageParams{
				PageID:  pageResult.ID,
				Title:   "MCP SDK 高级测试页面 - 已更新",
				Content: "页面内容已更新，添加了更多信息。",
				Properties: map[string]any{
					"status":     "已完成",
					"priority":   "中",
					"updated_at": time.Now().Format("2006-01-02 15:04:05"),
				},
			}

			updatedPage, err := sdk.UpdatePage(ctx, updateParams)
			if err != nil {
				log.Printf("更新页面失败: %v", err)
			} else {
				fmt.Printf("页面更新成功: %s\n", updatedPage.URL)
			}
		}
	} else {
		fmt.Println("跳过创建页面示例（需要设置 NOTION_PARENT_PAGE_ID 环境变量）")
	}

	// 示例 4: 批量添加不同类型的内容
	fmt.Println("\n=== 批量添加内容示例 ===")
	if parentPageID != "" {
		// 先创建一个页面用于演示
		pageResult, err := sdk.QuickCreatePage(ctx, parentPageID, "批量内容测试页面", "这个页面将用于演示批量添加不同类型的内容。")
		if err != nil {
			log.Printf("创建测试页面失败: %v", err)
		} else {
			fmt.Printf("测试页面创建成功: %s\n", pageResult.URL)

			// 定义要添加的内容
			contents := []struct {
				Type    string
				Content string
				Level   int
			}{
				{"heading", "项目概述", 1},
				{"text", "这是一个关于 MCP SDK 的项目概述。我们将演示如何使用 SDK 来操作 Notion 内容。", 0},
				{"heading", "功能特性", 2},
				{"bullet", "支持搜索 Notion 内容", 0},
				{"bullet", "支持创建和更新页面", 0},
				{"bullet", "支持添加各种类型的块内容", 0},
				{"heading", "代码示例", 2},
				{"code", "package main\n\nimport (\n    \"context\"\n    \"github.com/tenz-io/notionapi/mcp\"\n)\n\nfunc main() {\n    sdk, _ := mcp.NewNotionMCPSDKWithDefaults(\"your-token\")\n    result, _ := sdk.QuickSearch(context.Background(), \"test\")\n    fmt.Println(result)\n}", 0},
				{"heading", "待办事项", 2},
				{"todo", "完成 SDK 开发", 0},
				{"todo", "编写文档", 0},
				{"todo", "进行测试", 0},
				{"heading", "重要说明", 2},
				{"callout", "请确保在使用 SDK 之前正确配置 Notion Token 和权限。", 0},
				{"quote", "MCP SDK 让 Notion 集成变得简单而强大。", 0},
			}

			// 批量添加内容
			for i, item := range contents {
				var err error
				switch item.Type {
				case "heading":
					_, err = sdk.QuickAppendHeading(ctx, pageResult.ID, item.Content, item.Level)
				case "text":
					_, err = sdk.QuickAppendText(ctx, pageResult.ID, item.Content)
				case "bullet":
					_, err = sdk.QuickAppendBulletList(ctx, pageResult.ID, item.Content)
				case "code":
					_, err = sdk.QuickAppendCode(ctx, pageResult.ID, item.Content)
				case "todo":
					_, err = sdk.QuickAppendTodo(ctx, pageResult.ID, item.Content)
				case "callout":
					_, err = sdk.QuickAppendCallout(ctx, pageResult.ID, item.Content)
				case "quote":
					_, err = sdk.QuickAppendQuote(ctx, pageResult.ID, item.Content)
				}

				if err != nil {
					log.Printf("添加第 %d 项内容失败: %v", i+1, err)
				} else {
					fmt.Printf("✓ 添加 %s: %s\n", item.Type, item.Content[:min(30, len(item.Content))])
				}

				// 添加小延迟避免 API 限制
				time.Sleep(100 * time.Millisecond)
			}

			fmt.Println("批量内容添加完成！")
		}
	}

	// 示例 5: 配置管理
	fmt.Println("\n=== 配置管理示例 ===")
	fmt.Printf("当前配置: %+v\n", sdk.GetConfig())

	// 修改配置
	sdk.SetDefaultPageSize(50)
	sdk.SetTimeout(120)
	sdk.SetMaxRetries(10)

	fmt.Printf("修改后配置: %+v\n", sdk.GetConfig())

	fmt.Println("\n=== 高级示例完成 ===")
}

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
