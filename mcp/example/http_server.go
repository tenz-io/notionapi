package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/tenz-io/notionapi/mcp"
)

func mainServer() {
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

	// 创建 SDK 实例
	sdk, err := mcp.NewNotionMCPSDKWithDefaults(notionToken)
	if err != nil {
		log.Fatalf("Failed to create SDK: %v", err)
	}

	// 设置 HTTP 路由
	http.HandleFunc("/mcp", sdk.HandleHTTPRequest)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/tools", toolsHandler(sdk))
	http.HandleFunc("/resources", resourcesHandler(sdk))
	http.HandleFunc("/search", searchHandler(sdk))
	http.HandleFunc("/workspace", workspaceHandler(sdk))

	// 启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("MCP 服务器启动在端口 %s\n", port)
	fmt.Printf("健康检查: http://localhost:%s/health\n", port)
	fmt.Printf("MCP 端点: http://localhost:%s/mcp\n", port)
	fmt.Printf("工具列表: http://localhost:%s/tools\n", port)
	fmt.Printf("资源列表: http://localhost:%s/resources\n", port)
	fmt.Printf("搜索示例: http://localhost:%s/search?q=test\n", port)
	fmt.Printf("工作区信息: http://localhost:%s/workspace\n", port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// healthHandler 健康检查处理器
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"service":   "notion-mcp-server",
	})
}

// toolsHandler 工具列表处理器
func toolsHandler(sdk *mcp.NotionMCPSDK) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		tools := sdk.GetTools()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"tools": tools,
		})
	}
}

// resourcesHandler 资源列表处理器
func resourcesHandler(sdk *mcp.NotionMCPSDK) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		resources := sdk.GetResources()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"resources": resources,
		})
	}
}

// searchHandler 搜索处理器
func searchHandler(sdk *mcp.NotionMCPSDK) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		query := r.URL.Query().Get("q")
		if query == "" {
			http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
			return
		}

		filter := r.URL.Query().Get("filter")
		pageSize := 10
		if ps := r.URL.Query().Get("pageSize"); ps != "" {
			if n, err := fmt.Sscanf(ps, "%d", &pageSize); err != nil || n != 1 {
				pageSize = 10
			}
		}

		// 构建搜索参数
		searchParams := &mcp.NotionSearchParams{
			Query:    query,
			Filter:   filter,
			PageSize: pageSize,
		}

		// 执行搜索
		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		result, err := sdk.Search(ctx, searchParams)
		if err != nil {
			http.Error(w, fmt.Sprintf("Search failed: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}

// workspaceHandler 工作区信息处理器
func workspaceHandler(sdk *mcp.NotionMCPSDK) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// 获取工作区信息
		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		workspaceInfo, err := sdk.GetWorkspaceInfo(ctx)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get workspace info: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(workspaceInfo)
	}
}
