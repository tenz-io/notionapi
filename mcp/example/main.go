package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/tenz-io/notionapi/mcp"
)

// 加载环境变量
func loadEnvVars() {
	// 从 .env 文件加载环境变量
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
		log.Println("继续使用系统环境变量...")
	}

	// 显示加载的环境变量（隐藏敏感信息）
	notionToken := os.Getenv("NOTION_TOKEN")
	parentPageID := os.Getenv("NOTION_PARENT_PAGE_ID")

	fmt.Println("环境变量加载状态:")
	if notionToken != "" {
		fmt.Println("NOTION_TOKEN: 已设置")
	} else {
		fmt.Println("NOTION_TOKEN: 未设置")
	}

	if parentPageID != "" {
		fmt.Println("NOTION_PARENT_PAGE_ID: 已设置")
	} else {
		fmt.Println("NOTION_PARENT_PAGE_ID: 未设置")
	}
	fmt.Println()
}

// 运行基础示例
func runBasicExample() {
	fmt.Println("=== 运行基础示例 ===")

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
	searchResult, err := sdk.QuickSearch(ctx, "mcp")
	if err != nil {
		log.Printf("搜索失败: %v", err)
	} else {
		fmt.Printf("找到 %d 个结果\n", len(searchResult.Results))
		for i, item := range searchResult.Results {
			if i >= 3 { // 只显示前3个结果
				break
			}
			fmt.Printf("- %s (%s) %s\n", item.Title, item.Object, item.ID)
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

	// 示例 3: 显示可用工具和资源
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

	fmt.Println("\n=== 基础示例完成 ===")
}

// 运行高级示例
func runAdvancedExample() {
	fmt.Println("=== 运行高级示例 ===")

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

	// 示例 2: 配置管理
	fmt.Println("\n=== 配置管理示例 ===")
	fmt.Printf("当前配置: %+v\n", sdk.GetConfig())

	// 修改配置
	sdk.SetDefaultPageSize(50)
	sdk.SetTimeout(120)
	sdk.SetMaxRetries(10)

	fmt.Printf("修改后配置: %+v\n", sdk.GetConfig())

	fmt.Println("\n=== 高级示例完成 ===")
}

// 运行服务器示例
func runServerExample() {
	fmt.Println("=== 启动 HTTP 服务器示例 ===")

	// 从环境变量获取 Notion Token
	notionToken := os.Getenv("NOTION_TOKEN")
	if notionToken == "" {
		log.Fatal("NOTION_TOKEN environment variable is required")
	}

	// 创建 SDK 实例
	sdk, err := mcp.NewNotionMCPSDKWithDefaults(notionToken)
	if err != nil {
		log.Fatalf("Failed to create SDK: %v", err)
	}

	// 设置 HTTP 路由
	http.HandleFunc("/mcp", sdk.HandleHTTPRequest)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"healthy","service":"notion-mcp-server"}`)
	})

	// 启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("MCP 服务器启动在端口 %s\n", port)
	fmt.Printf("健康检查: http://localhost:%s/health\n", port)
	fmt.Printf("MCP 端点: http://localhost:%s/mcp\n", port)
	fmt.Println("按 Ctrl+C 停止服务器")

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// 声明函数以避免未定义错误
func runIntegrationTestsMain() {
	fmt.Println("运行集成测试...")
	// 这里可以调用 integration_test.go 中的 runIntegrationTests 函数
	// 或者直接在这里实现测试逻辑
}

func main() {
	// 加载环境变量
	loadEnvVars()
	if len(os.Args) < 2 {
		fmt.Println("用法:")
		fmt.Println("  go run main.go basic     - 运行基础示例")
		fmt.Println("  go run main.go advanced  - 运行高级示例")
		fmt.Println("  go run main.go server    - 启动 HTTP 服务器")
		fmt.Println("  go run main.go test      - 运行集成测试")
		return
	}

	switch os.Args[1] {
	case "basic":
		runBasicExample()
	case "advanced":
		runAdvancedExample()
	case "server":
		runServerExample()
	case "test":
		runIntegrationTestsMain()
	default:
		fmt.Printf("未知命令: %s\n", os.Args[1])
		fmt.Println("可用命令: basic, advanced, server, test")
	}
}
