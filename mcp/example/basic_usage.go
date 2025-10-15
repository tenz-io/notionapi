package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/tenz-io/notionapi/mcp"
)

func mainBasic() {
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

	// 创建 SDK 实例（使用默认配置）
	sdk, err := mcp.NewNotionMCPSDKWithDefaults(notionToken)
	if err != nil {
		log.Fatalf("Failed to create SDK: %v", err)
	}

	ctx := context.Background()

	// 示例 1: 搜索 Notion 内容
	fmt.Println("=== 搜索示例 ===")
	searchResult, err := sdk.QuickSearch(ctx, "测试")
	if err != nil {
		log.Printf("搜索失败: %v", err)
	} else {
		fmt.Printf("找到 %d 个结果\n", len(searchResult.Results))
		for i, item := range searchResult.Results {
			if i >= 3 { // 只显示前3个结果
				break
			}
			fmt.Printf("- %s (%s)\n", item.Title, item.Object)
		}
	}

	// 示例 2: 获取工作区信息
	fmt.Println("\n=== 工作区信息 ===")
	workspaceInfo, err := sdk.GetWorkspaceInfo(ctx)
	if err != nil {
		log.Printf("获取工作区信息失败: %v", err)
	} else {
		fmt.Printf("工作区统计: %+v\n", workspaceInfo)
	}

	// 示例 3: 创建新页面（需要有效的父页面ID）
	fmt.Println("\n=== 创建页面示例 ===")
	parentPageID := os.Getenv("NOTION_PARENT_PAGE_ID")
	if parentPageID != "" {
		pageResult, err := sdk.QuickCreatePage(ctx, parentPageID, "MCP SDK 测试页面", "这是一个通过 MCP SDK 创建的测试页面。")
		if err != nil {
			log.Printf("创建页面失败: %v", err)
		} else {
			fmt.Printf("页面创建成功: %s\n", pageResult.URL)

			// 示例 4: 向页面添加内容
			fmt.Println("\n=== 添加内容示例 ===")

			// 添加标题
			_, err = sdk.QuickAppendHeading(ctx, pageResult.ID, "这是二级标题", 2)
			if err != nil {
				log.Printf("添加标题失败: %v", err)
			}

			// 添加段落
			_, err = sdk.QuickAppendText(ctx, pageResult.ID, "这是一个段落内容，通过 MCP SDK 添加。")
			if err != nil {
				log.Printf("添加段落失败: %v", err)
			}

			// 添加代码块
			_, err = sdk.QuickAppendCode(ctx, pageResult.ID, "func main() {\n    fmt.Println(\"Hello, MCP SDK!\")\n}")
			if err != nil {
				log.Printf("添加代码块失败: %v", err)
			}

			// 添加引用
			_, err = sdk.QuickAppendQuote(ctx, pageResult.ID, "这是一个引用内容。")
			if err != nil {
				log.Printf("添加引用失败: %v", err)
			}

			// 添加待办事项
			_, err = sdk.QuickAppendTodo(ctx, pageResult.ID, "完成 MCP SDK 集成测试")
			if err != nil {
				log.Printf("添加待办事项失败: %v", err)
			}

			// 添加无序列表
			_, err = sdk.QuickAppendBulletList(ctx, pageResult.ID, "列表项 1")
			if err != nil {
				log.Printf("添加无序列表失败: %v", err)
			}

			// 添加有序列表
			_, err = sdk.QuickAppendNumberedList(ctx, pageResult.ID, "有序列表项 1")
			if err != nil {
				log.Printf("添加有序列表失败: %v", err)
			}

			// 添加标注
			_, err = sdk.QuickAppendCallout(ctx, pageResult.ID, "这是一个重要的标注内容。")
			if err != nil {
				log.Printf("添加标注失败: %v", err)
			}

			fmt.Println("所有内容添加完成！")
		}
	} else {
		fmt.Println("跳过创建页面示例（需要设置 NOTION_PARENT_PAGE_ID 环境变量）")
	}

	// 示例 5: 显示可用工具和资源
	fmt.Println("\n=== 可用工具 ===")
	tools := sdk.GetTools()
	for _, tool := range tools {
		fmt.Printf("- %s: %s\n", tool.Name, tool.Description)
	}

	fmt.Println("\n=== 可用资源 ===")
	resources := sdk.GetResources()
	for _, resource := range resources {
		fmt.Printf("- %s: %s\n", resource.URI, resource.Description)
	}

	fmt.Println("\n=== 示例完成 ===")
}
