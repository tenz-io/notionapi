package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/tenz-io/notionapi/mcp"
)

// TestNotionMCPSDK 集成测试
func TestNotionMCPSDK(t *testing.T) {
	// 加载环境变量
	err := godotenv.Load()
	if err != nil {
		t.Logf("Warning: Error loading .env file: %v", err)
	}

	// 检查环境变量
	notionToken := os.Getenv("NOTION_TOKEN")
	if notionToken == "" {
		t.Skip("NOTION_TOKEN environment variable is required for integration tests")
	}

	parentPageID := os.Getenv("NOTION_PARENT_PAGE_ID")
	if parentPageID == "" {
		t.Skip("NOTION_PARENT_PAGE_ID environment variable is required for integration tests")
	}

	// 创建 SDK 实例
	sdk, err := mcp.NewNotionMCPSDKWithDefaults(notionToken)
	if err != nil {
		t.Fatalf("Failed to create SDK: %v", err)
	}

	ctx := context.Background()

	// 测试 1: 搜索功能
	t.Run("Search", func(t *testing.T) {
		searchResult, err := sdk.QuickSearch(ctx, "测试")
		if err != nil {
			t.Errorf("Search failed: %v", err)
			return
		}

		if searchResult.Object != "list" {
			t.Errorf("Expected object type 'list', got '%s'", searchResult.Object)
		}

		t.Logf("Found %d search results", len(searchResult.Results))
	})

	// 测试 2: 获取工作区信息
	t.Run("GetWorkspaceInfo", func(t *testing.T) {
		workspaceInfo, err := sdk.GetWorkspaceInfo(ctx)
		if err != nil {
			t.Errorf("GetWorkspaceInfo failed: %v", err)
			return
		}

		if workspaceInfo == nil {
			t.Error("Workspace info should not be nil")
			return
		}

		t.Logf("Workspace info: %+v", workspaceInfo)
	})

	// 测试 3: 创建页面
	t.Run("CreatePage", func(t *testing.T) {
		pageResult, err := sdk.QuickCreatePage(ctx, parentPageID, "MCP SDK 集成测试页面", "这是一个集成测试页面。")
		if err != nil {
			t.Errorf("CreatePage failed: %v", err)
			return
		}

		if pageResult.ID == "" {
			t.Error("Page ID should not be empty")
		}

		if pageResult.URL == "" {
			t.Error("Page URL should not be empty")
		}

		t.Logf("Created page: %s", pageResult.URL)

		// 测试 4: 更新页面
		t.Run("UpdatePage", func(t *testing.T) {
			updateParams := &mcp.NotionUpdatePageParams{
				PageID:  pageResult.ID,
				Title:   "MCP SDK 集成测试页面 - 已更新",
				Content: "页面内容已更新。",
			}

			updatedPage, err := sdk.UpdatePage(ctx, updateParams)
			if err != nil {
				t.Errorf("UpdatePage failed: %v", err)
				return
			}

			if updatedPage.ID != pageResult.ID {
				t.Error("Updated page ID should match original page ID")
			}

			t.Logf("Updated page: %s", updatedPage.URL)
		})

		// 测试 5: 添加各种类型的块内容
		t.Run("AppendBlocks", func(t *testing.T) {
			blockTests := []struct {
				name      string
				content   string
				blockType string
				method    func(context.Context, string, string) (*mcp.NotionBlockResult, error)
			}{
				{"Text", "这是一个段落内容。", "paragraph", sdk.QuickAppendText},
				{"Heading1", "这是一级标题", "heading_1", func(ctx context.Context, pageID, content string) (*mcp.NotionBlockResult, error) {
					return sdk.QuickAppendHeading(ctx, pageID, content, 1)
				}},
				{"Heading2", "这是二级标题", "heading_2", func(ctx context.Context, pageID, content string) (*mcp.NotionBlockResult, error) {
					return sdk.QuickAppendHeading(ctx, pageID, content, 2)
				}},
				{"Code", "func main() {\n    fmt.Println(\"Hello, World!\")\n}", "code", sdk.QuickAppendCode},
				{"Quote", "这是一个引用内容。", "quote", sdk.QuickAppendQuote},
				{"Callout", "这是一个重要的标注。", "callout", sdk.QuickAppendCallout},
				{"Todo", "完成集成测试", "to_do", sdk.QuickAppendTodo},
				{"BulletList", "列表项 1", "bulleted_list_item", sdk.QuickAppendBulletList},
				{"NumberedList", "有序列表项 1", "numbered_list_item", sdk.QuickAppendNumberedList},
			}

			for _, test := range blockTests {
				t.Run(test.name, func(t *testing.T) {
					result, err := test.method(ctx, pageResult.ID, test.content)
					if err != nil {
						t.Errorf("Append %s failed: %v", test.name, err)
						return
					}

					if result.Object != "list" {
						t.Errorf("Expected object type 'list', got '%s'", result.Object)
					}

					if len(result.Results) == 0 {
						t.Error("Expected at least one result")
					}

					t.Logf("Successfully appended %s block", test.name)
				})

				// 添加小延迟避免 API 限制
				time.Sleep(100 * time.Millisecond)
			}
		})

		// 测试 6: 高级搜索
		t.Run("AdvancedSearch", func(t *testing.T) {
			searchParams := &mcp.NotionSearchParams{
				Query:     "MCP",
				Filter:    "page",
				SortBy:    "last_edited_time",
				SortOrder: "descending",
				PageSize:  5,
			}

			searchResult, err := sdk.Search(ctx, searchParams)
			if err != nil {
				t.Errorf("Advanced search failed: %v", err)
				return
			}

			t.Logf("Advanced search found %d results", len(searchResult.Results))
		})
	})

	// 测试 7: 配置管理
	t.Run("ConfigManagement", func(t *testing.T) {
		originalConfig := sdk.GetConfig()

		// 测试配置修改
		sdk.SetDefaultPageSize(25)
		sdk.SetTimeout(60)
		sdk.SetMaxRetries(5)

		newConfig := sdk.GetConfig()
		if newConfig.DefaultPageSize != 25 {
			t.Error("DefaultPageSize should be 25")
		}
		if newConfig.Timeout != 60 {
			t.Error("Timeout should be 60")
		}
		if newConfig.MaxRetries != 5 {
			t.Error("MaxRetries should be 5")
		}

		// 恢复原始配置
		sdk.SetDefaultPageSize(originalConfig.DefaultPageSize)
		sdk.SetTimeout(originalConfig.Timeout)
		sdk.SetMaxRetries(originalConfig.MaxRetries)

		t.Log("Config management test passed")
	})

	// 测试 8: 工具和资源列表
	t.Run("ToolsAndResources", func(t *testing.T) {
		tools := sdk.GetTools()
		if len(tools) == 0 {
			t.Error("Should have at least one tool")
		}

		resources := sdk.GetResources()
		if len(resources) == 0 {
			t.Error("Should have at least one resource")
		}

		t.Logf("Found %d tools and %d resources", len(tools), len(resources))
	})
}

// BenchmarkNotionMCPSDK 性能测试
func BenchmarkNotionMCPSDK(b *testing.B) {
	// 检查环境变量
	notionToken := os.Getenv("NOTION_TOKEN")
	if notionToken == "" {
		b.Skip("NOTION_TOKEN environment variable is required for benchmarks")
	}

	// 创建 SDK 实例
	sdk, err := mcp.NewNotionMCPSDKWithDefaults(notionToken)
	if err != nil {
		b.Fatalf("Failed to create SDK: %v", err)
	}

	ctx := context.Background()

	// 基准测试搜索功能
	b.Run("Search", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := sdk.QuickSearch(ctx, "测试")
			if err != nil {
				b.Errorf("Search failed: %v", err)
			}
		}
	})

	// 基准测试获取工作区信息
	b.Run("GetWorkspaceInfo", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := sdk.GetWorkspaceInfo(ctx)
			if err != nil {
				b.Errorf("GetWorkspaceInfo failed: %v", err)
			}
		}
	})
}

// ExampleNotionMCPSDK 示例函数
func ExampleNotionMCPSDK() {
	// 创建 SDK 实例
	sdk, err := mcp.NewNotionMCPSDKWithDefaults("your-notion-token")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// 搜索内容
	searchResult, err := sdk.QuickSearch(ctx, "项目")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("找到 %d 个结果\n", len(searchResult.Results))

	// 获取工作区信息
	workspaceInfo, err := sdk.GetWorkspaceInfo(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("工作区信息: %+v\n", workspaceInfo)

	// 输出:
	// 找到 X 个结果
	// 工作区信息: map[recentItems:... totalDatabases:X totalPages:Y]
}

// 辅助函数：运行集成测试
func runIntegrationTests() {
	fmt.Println("运行 Notion MCP SDK 集成测试...")

	// 加载环境变量
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Warning: Error loading .env file: %v\n", err)
	}

	// 检查环境变量
	notionToken := os.Getenv("NOTION_TOKEN")
	if notionToken == "" {
		fmt.Println("❌ NOTION_TOKEN 环境变量未设置")
		return
	}

	parentPageID := os.Getenv("NOTION_PARENT_PAGE_ID")
	if parentPageID == "" {
		fmt.Println("❌ NOTION_PARENT_PAGE_ID 环境变量未设置")
		return
	}

	// 创建 SDK 实例
	sdk, err := mcp.NewNotionMCPSDKWithDefaults(notionToken)
	if err != nil {
		fmt.Printf("❌ 创建 SDK 失败: %v\n", err)
		return
	}

	ctx := context.Background()

	// 测试搜索功能
	fmt.Println("🔍 测试搜索功能...")
	searchResult, err := sdk.QuickSearch(ctx, "测试")
	if err != nil {
		fmt.Printf("❌ 搜索失败: %v\n", err)
	} else {
		fmt.Printf("✅ 搜索成功，找到 %d 个结果\n", len(searchResult.Results))
	}

	// 测试获取工作区信息
	fmt.Println("📊 测试获取工作区信息...")
	workspaceInfo, err := sdk.GetWorkspaceInfo(ctx)
	if err != nil {
		fmt.Printf("❌ 获取工作区信息失败: %v\n", err)
	} else {
		fmt.Printf("✅ 获取工作区信息成功\n")
		fmt.Printf("   总页面数: %v\n", workspaceInfo["totalPages"])
		fmt.Printf("   总数据库数: %v\n", workspaceInfo["totalDatabases"])
	}

	// 测试创建页面
	fmt.Println("📝 测试创建页面...")
	pageResult, err := sdk.QuickCreatePage(ctx, parentPageID, "MCP SDK 集成测试", "这是一个集成测试页面。")
	if err != nil {
		fmt.Printf("❌ 创建页面失败: %v\n", err)
	} else {
		fmt.Printf("✅ 创建页面成功: %s\n", pageResult.URL)
	}

	// 测试添加内容
	if pageResult != nil {
		fmt.Println("📄 测试添加内容...")

		// 添加标题
		_, err = sdk.QuickAppendHeading(ctx, pageResult.ID, "测试标题", 2)
		if err != nil {
			fmt.Printf("❌ 添加标题失败: %v\n", err)
		} else {
			fmt.Println("✅ 添加标题成功")
		}

		// 添加段落
		_, err = sdk.QuickAppendText(ctx, pageResult.ID, "这是一个测试段落。")
		if err != nil {
			fmt.Printf("❌ 添加段落失败: %v\n", err)
		} else {
			fmt.Println("✅ 添加段落成功")
		}

		// 添加代码块
		_, err = sdk.QuickAppendCode(ctx, pageResult.ID, "fmt.Println(\"Hello, MCP SDK!\")")
		if err != nil {
			fmt.Printf("❌ 添加代码块失败: %v\n", err)
		} else {
			fmt.Println("✅ 添加代码块成功")
		}
	}

	fmt.Println("🎉 集成测试完成！")
}

func mainTest() {
	// 如果直接运行此文件，执行集成测试
	if len(os.Args) > 1 && os.Args[1] == "test" {
		runIntegrationTests()
		return
	}

	// 否则运行 Go 测试
	fmt.Println("使用 'go test' 运行测试，或使用 'go run integration_test.go test' 运行集成测试")
}
