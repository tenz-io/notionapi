package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// NotionMCPSDK 表示 Notion MCP SDK
type NotionMCPSDK struct {
	config *MCPConfig
	server *MCPServer
}

// NewNotionMCPSDK 创建新的 Notion MCP SDK
func NewNotionMCPSDK(config *MCPConfig) (*NotionMCPSDK, error) {
	server, err := NewMCPServer(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create MCP server: %w", err)
	}

	return &NotionMCPSDK{
		config: config,
		server: server,
	}, nil
}

// NewNotionMCPSDKWithDefaults 使用默认配置创建 SDK
func NewNotionMCPSDKWithDefaults(notionToken string) (*NotionMCPSDK, error) {
	config := DefaultMCPConfig()
	config.NotionToken = notionToken
	return NewNotionMCPSDK(config)
}

// Search 搜索 Notion 内容
func (sdk *NotionMCPSDK) Search(ctx context.Context, params *NotionSearchParams) (*NotionSearchResult, error) {
	// 构建工具调用参数
	args := map[string]interface{}{
		"pageSize": params.PageSize,
	}

	if params.Query != "" {
		args["query"] = params.Query
	}
	if params.Filter != "" {
		args["filter"] = params.Filter
	}
	if params.SortBy != "" {
		args["sortBy"] = params.SortBy
	}
	if params.SortOrder != "" {
		args["sortOrder"] = params.SortOrder
	}
	if params.StartCursor != "" {
		args["startCursor"] = params.StartCursor
	}

	// 调用搜索工具
	toolCall := ToolCall{
		Name:      "notion_search",
		Arguments: args,
	}

	result, err := sdk.callTool(ctx, toolCall)
	if err != nil {
		return nil, err
	}

	// 解析结果
	var searchResult NotionSearchResult
	if len(result.Content) > 0 {
		err = json.Unmarshal([]byte(result.Content[0].Text), &searchResult)
		if err != nil {
			return nil, fmt.Errorf("failed to parse search result: %w", err)
		}
	}

	return &searchResult, nil
}

// CreatePage 创建 Notion 页面
func (sdk *NotionMCPSDK) CreatePage(ctx context.Context, params *NotionCreatePageParams) (*NotionPageResult, error) {
	// 构建工具调用参数
	args := map[string]interface{}{
		"parentId": params.ParentID,
		"title":    params.Title,
	}

	if params.Content != "" {
		args["content"] = params.Content
	}
	if params.Properties != nil {
		args["properties"] = params.Properties
	}
	if params.Icon != nil {
		args["icon"] = params.Icon
	}
	if params.Cover != nil {
		args["cover"] = params.Cover
	}

	// 调用创建页面工具
	toolCall := ToolCall{
		Name:      "notion_create_page",
		Arguments: args,
	}

	result, err := sdk.callTool(ctx, toolCall)
	if err != nil {
		return nil, err
	}

	// 解析结果
	var pageResult NotionPageResult
	if len(result.Content) > 0 {
		err = json.Unmarshal([]byte(result.Content[0].Text), &pageResult)
		if err != nil {
			return nil, fmt.Errorf("failed to parse page result: %w", err)
		}
	}

	return &pageResult, nil
}

// UpdatePage 更新 Notion 页面
func (sdk *NotionMCPSDK) UpdatePage(ctx context.Context, params *NotionUpdatePageParams) (*NotionPageResult, error) {
	// 构建工具调用参数
	args := map[string]interface{}{
		"pageId": params.PageID,
	}

	if params.Title != "" {
		args["title"] = params.Title
	}
	if params.Content != "" {
		args["content"] = params.Content
	}
	if params.Properties != nil {
		args["properties"] = params.Properties
	}
	if params.Archived != nil {
		args["archived"] = *params.Archived
	}

	// 调用更新页面工具
	toolCall := ToolCall{
		Name:      "notion_update_page",
		Arguments: args,
	}

	result, err := sdk.callTool(ctx, toolCall)
	if err != nil {
		return nil, err
	}

	// 解析结果
	var pageResult NotionPageResult
	if len(result.Content) > 0 {
		err = json.Unmarshal([]byte(result.Content[0].Text), &pageResult)
		if err != nil {
			return nil, fmt.Errorf("failed to parse page result: %w", err)
		}
	}

	return &pageResult, nil
}

// AppendBlock 向页面添加块内容
func (sdk *NotionMCPSDK) AppendBlock(ctx context.Context, params *NotionAppendBlockParams) (*NotionBlockResult, error) {
	// 构建工具调用参数
	args := map[string]interface{}{
		"pageId":  params.PageID,
		"content": params.Content,
	}

	if params.BlockType != "" {
		args["blockType"] = params.BlockType
	}

	// 调用添加块工具
	toolCall := ToolCall{
		Name:      "notion_append_block",
		Arguments: args,
	}

	result, err := sdk.callTool(ctx, toolCall)
	if err != nil {
		return nil, err
	}

	// 解析结果
	var blockResult NotionBlockResult
	if len(result.Content) > 0 {
		err = json.Unmarshal([]byte(result.Content[0].Text), &blockResult)
		if err != nil {
			return nil, fmt.Errorf("failed to parse block result: %w", err)
		}
	}

	return &blockResult, nil
}

// GetWorkspaceInfo 获取工作区信息
func (sdk *NotionMCPSDK) GetWorkspaceInfo(ctx context.Context) (map[string]interface{}, error) {
	// 构建资源读取请求
	req := &Request{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "resources/read",
		Params: map[string]interface{}{
			"uri": "notion://workspace",
		},
	}

	// 处理请求
	resp := sdk.server.HandleRequest(ctx, req)
	if resp.Error != nil {
		return nil, fmt.Errorf("failed to get workspace info: %s", resp.Error.Message)
	}

	// 解析结果
	result, ok := resp.Result.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	contents, ok := result["contents"].([]interface{})
	if !ok || len(contents) == 0 {
		return nil, fmt.Errorf("no workspace info found")
	}

	content, ok := contents[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid content format")
	}

	text, ok := content["text"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid text content")
	}

	// 解析 JSON
	var workspaceInfo map[string]interface{}
	err := json.Unmarshal([]byte(text), &workspaceInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to parse workspace info: %w", err)
	}

	return workspaceInfo, nil
}

// GetTools 获取可用工具列表
func (sdk *NotionMCPSDK) GetTools() []Tool {
	return sdk.server.tools
}

// GetResources 获取可用资源列表
func (sdk *NotionMCPSDK) GetResources() []Resource {
	return sdk.server.resources
}

// HandleHTTPRequest 处理 HTTP 请求（用于集成到 HTTP 服务器）
func (sdk *NotionMCPSDK) HandleHTTPRequest(w http.ResponseWriter, r *http.Request) {
	// 只允许 POST 请求
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 读取请求体
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// 解析 MCP 请求
	var req Request
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// 处理请求
	resp := sdk.server.HandleRequest(r.Context(), &req)

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")

	// 发送响应
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// StartHTTPServer 启动 HTTP 服务器
func (sdk *NotionMCPSDK) StartHTTPServer(addr string) error {
	http.HandleFunc("/mcp", sdk.HandleHTTPRequest)
	return http.ListenAndServe(addr, nil)
}

// 辅助方法

// callTool 调用工具
func (sdk *NotionMCPSDK) callTool(ctx context.Context, toolCall ToolCall) (*ToolResult, error) {
	// 构建工具调用请求
	req := &Request{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "tools/call",
		Params:  toolCall,
	}

	// 处理请求
	resp := sdk.server.HandleRequest(ctx, req)
	if resp.Error != nil {
		return nil, fmt.Errorf("tool call failed: %s", resp.Error.Message)
	}

	// 解析结果
	result, ok := resp.Result.(*ToolResult)
	if !ok {
		return nil, fmt.Errorf("invalid tool result format")
	}

	return result, nil
}

// 便捷方法

// QuickSearch 快速搜索（使用默认参数）
func (sdk *NotionMCPSDK) QuickSearch(ctx context.Context, query string) (*NotionSearchResult, error) {
	params := &NotionSearchParams{
		Query:    query,
		PageSize: sdk.config.DefaultPageSize,
	}
	return sdk.Search(ctx, params)
}

// QuickCreatePage 快速创建页面（使用默认参数）
func (sdk *NotionMCPSDK) QuickCreatePage(ctx context.Context, parentID, title, content string) (*NotionPageResult, error) {
	params := &NotionCreatePageParams{
		ParentID: parentID,
		Title:    title,
		Content:  content,
	}
	return sdk.CreatePage(ctx, params)
}

// QuickAppendText 快速添加文本内容
func (sdk *NotionMCPSDK) QuickAppendText(ctx context.Context, pageID, content string) (*NotionBlockResult, error) {
	params := &NotionAppendBlockParams{
		PageID:    pageID,
		Content:   content,
		BlockType: "paragraph",
	}
	return sdk.AppendBlock(ctx, params)
}

// QuickAppendHeading 快速添加标题
func (sdk *NotionMCPSDK) QuickAppendHeading(ctx context.Context, pageID, content string, level int) (*NotionBlockResult, error) {
	var blockType string
	switch level {
	case 1:
		blockType = "heading_1"
	case 2:
		blockType = "heading_2"
	case 3:
		blockType = "heading_3"
	default:
		blockType = "heading_1"
	}

	params := &NotionAppendBlockParams{
		PageID:    pageID,
		Content:   content,
		BlockType: blockType,
	}
	return sdk.AppendBlock(ctx, params)
}

// QuickAppendCode 快速添加代码块
func (sdk *NotionMCPSDK) QuickAppendCode(ctx context.Context, pageID, content string) (*NotionBlockResult, error) {
	params := &NotionAppendBlockParams{
		PageID:    pageID,
		Content:   content,
		BlockType: "code",
	}
	return sdk.AppendBlock(ctx, params)
}

// QuickAppendQuote 快速添加引用
func (sdk *NotionMCPSDK) QuickAppendQuote(ctx context.Context, pageID, content string) (*NotionBlockResult, error) {
	params := &NotionAppendBlockParams{
		PageID:    pageID,
		Content:   content,
		BlockType: "quote",
	}
	return sdk.AppendBlock(ctx, params)
}

// QuickAppendCallout 快速添加标注
func (sdk *NotionMCPSDK) QuickAppendCallout(ctx context.Context, pageID, content string) (*NotionBlockResult, error) {
	params := &NotionAppendBlockParams{
		PageID:    pageID,
		Content:   content,
		BlockType: "callout",
	}
	return sdk.AppendBlock(ctx, params)
}

// QuickAppendTodo 快速添加待办事项
func (sdk *NotionMCPSDK) QuickAppendTodo(ctx context.Context, pageID, content string) (*NotionBlockResult, error) {
	params := &NotionAppendBlockParams{
		PageID:    pageID,
		Content:   content,
		BlockType: "to_do",
	}
	return sdk.AppendBlock(ctx, params)
}

// QuickAppendBulletList 快速添加无序列表
func (sdk *NotionMCPSDK) QuickAppendBulletList(ctx context.Context, pageID, content string) (*NotionBlockResult, error) {
	params := &NotionAppendBlockParams{
		PageID:    pageID,
		Content:   content,
		BlockType: "bulleted_list_item",
	}
	return sdk.AppendBlock(ctx, params)
}

// QuickAppendNumberedList 快速添加有序列表
func (sdk *NotionMCPSDK) QuickAppendNumberedList(ctx context.Context, pageID, content string) (*NotionBlockResult, error) {
	params := &NotionAppendBlockParams{
		PageID:    pageID,
		Content:   content,
		BlockType: "numbered_list_item",
	}
	return sdk.AppendBlock(ctx, params)
}

// 配置方法

// SetDefaultPageSize 设置默认页面大小
func (sdk *NotionMCPSDK) SetDefaultPageSize(size int) {
	if size > 0 && size <= 100 {
		sdk.config.DefaultPageSize = size
	}
}

// SetTimeout 设置超时时间
func (sdk *NotionMCPSDK) SetTimeout(seconds int) {
	if seconds > 0 {
		sdk.config.Timeout = seconds
	}
}

// SetMaxRetries 设置最大重试次数
func (sdk *NotionMCPSDK) SetMaxRetries(retries int) {
	if retries > 0 {
		sdk.config.MaxRetries = retries
	}
}

// GetConfig 获取当前配置
func (sdk *NotionMCPSDK) GetConfig() *MCPConfig {
	return sdk.config
}
