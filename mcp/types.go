package mcp

import (
	"time"
)

// MCP 协议相关类型定义

// Request 表示 MCP 请求
type Request struct {
	JSONRPC string `json:"jsonrpc"`
	ID      any    `json:"id"`
	Method  string `json:"method"`
	Params  any    `json:"params,omitempty"`
}

// Response 表示 MCP 响应
type Response struct {
	JSONRPC string `json:"jsonrpc"`
	ID      any    `json:"id"`
	Result  any    `json:"result,omitempty"`
	Error   *Error `json:"error,omitempty"`
}

// Error 表示 MCP 错误
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Tool 表示 MCP 工具定义
type Tool struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	InputSchema map[string]any `json:"inputSchema"`
}

// Resource 表示 MCP 资源定义
type Resource struct {
	URI         string         `json:"uri"`
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	MimeType    string         `json:"mimeType,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}

// ToolCall 表示工具调用
type ToolCall struct {
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments"`
}

// ToolResult 表示工具调用结果
type ToolResult struct {
	Content []Content `json:"content"`
	IsError bool      `json:"isError,omitempty"`
}

// Content 表示内容项
type Content struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
	Data any    `json:"data,omitempty"`
}

// NotionSearchParams 表示 Notion 搜索参数
type NotionSearchParams struct {
	Query       string `json:"query,omitempty"`
	Filter      string `json:"filter,omitempty"`    // "page" 或 "database"
	SortBy      string `json:"sortBy,omitempty"`    // "last_edited_time"
	SortOrder   string `json:"sortOrder,omitempty"` // "ascending" 或 "descending"
	StartCursor string `json:"startCursor,omitempty"`
	PageSize    int    `json:"pageSize,omitempty"`
}

// NotionCreatePageParams 表示创建 Notion 页面参数
type NotionCreatePageParams struct {
	ParentID   string         `json:"parentId"`
	Title      string         `json:"title"`
	Content    string         `json:"content,omitempty"`    // 支持 markdown 格式
	Properties map[string]any `json:"properties,omitempty"` // 简化的属性定义
}

// NotionUpdatePageParams 表示更新 Notion 页面参数
type NotionUpdatePageParams struct {
	PageID     string         `json:"pageId"`
	Title      string         `json:"title,omitempty"`
	Content    string         `json:"content,omitempty"`
	Properties map[string]any `json:"properties,omitempty"`
	Archived   *bool          `json:"archived,omitempty"`
}

// NotionAppendBlockParams 表示添加块内容参数
type NotionAppendBlockParams struct {
	PageID    string `json:"pageId"`
	Content   string `json:"content"`
	BlockType string `json:"blockType,omitempty"` // "paragraph", "heading_1", "heading_2", "heading_3", "bulleted_list_item", "numbered_list_item", "to_do", "code", "quote", "callout"
}

// Icon 表示页面图标
type Icon struct {
	Type  string `json:"type"` // "emoji" 或 "external"
	Emoji string `json:"emoji,omitempty"`
	URL   string `json:"url,omitempty"`
}

// Cover 表示页面封面
type Cover struct {
	Type string `json:"type"` // "external"
	URL  string `json:"url"`
}

// NotionSearchResult 表示 Notion 搜索结果
type NotionSearchResult struct {
	Object     string       `json:"object"`
	Results    []NotionItem `json:"results"`
	HasMore    bool         `json:"has_more"`
	NextCursor string       `json:"next_cursor"`
}

// NotionItem 表示 Notion 项目（页面或数据库）
type NotionItem struct {
	Object         string         `json:"object"`
	ID             string         `json:"id"`
	CreatedTime    time.Time      `json:"created_time"`
	LastEditedTime time.Time      `json:"last_edited_time"`
	CreatedBy      *User          `json:"created_by,omitempty"`
	LastEditedBy   *User          `json:"last_edited_by,omitempty"`
	Archived       bool           `json:"archived"`
	Properties     map[string]any `json:"properties,omitempty"`
	Parent         *Parent        `json:"parent,omitempty"`
	URL            string         `json:"url"`
	PublicURL      string         `json:"public_url,omitempty"`
	Icon           *Icon          `json:"icon,omitempty"`
	Cover          *Cover         `json:"cover,omitempty"`
	Title          string         `json:"title,omitempty"`
}

// User 表示用户信息
type User struct {
	Object    string `json:"object"`
	ID        string `json:"id"`
	Name      string `json:"name,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
	Type      string `json:"type,omitempty"`
}

// Parent 表示父级对象
type Parent struct {
	Type       string `json:"type,omitempty"`
	PageID     string `json:"page_id,omitempty"`
	DatabaseID string `json:"database_id,omitempty"`
	BlockID    string `json:"block_id,omitempty"`
	Workspace  bool   `json:"workspace,omitempty"`
}

// NotionPageResult 表示 Notion 页面操作结果
type NotionPageResult struct {
	Object         string         `json:"object"`
	ID             string         `json:"id"`
	CreatedTime    time.Time      `json:"created_time"`
	LastEditedTime time.Time      `json:"last_edited_time"`
	CreatedBy      *User          `json:"created_by,omitempty"`
	LastEditedBy   *User          `json:"last_edited_by,omitempty"`
	Archived       bool           `json:"archived"`
	Properties     map[string]any `json:"properties"`
	Parent         *Parent        `json:"parent"`
	URL            string         `json:"url"`
	PublicURL      string         `json:"public_url,omitempty"`
	Icon           *Icon          `json:"icon,omitempty"`
	Cover          *Cover         `json:"cover,omitempty"`
}

// NotionBlockResult 表示 Notion 块操作结果
type NotionBlockResult struct {
	Object  string      `json:"object"`
	Results []BlockItem `json:"results"`
}

// BlockItem 表示块项目
type BlockItem struct {
	Object         string         `json:"object"`
	ID             string         `json:"id"`
	Type           string         `json:"type"`
	CreatedTime    time.Time      `json:"created_time"`
	LastEditedTime time.Time      `json:"last_edited_time"`
	CreatedBy      *User          `json:"created_by,omitempty"`
	LastEditedBy   *User          `json:"last_edited_by,omitempty"`
	HasChildren    bool           `json:"has_children"`
	Archived       bool           `json:"archived"`
	Parent         *Parent        `json:"parent,omitempty"`
	Content        map[string]any `json:"content,omitempty"`
}

// MCPConfig 表示 MCP 服务器配置
type MCPConfig struct {
	NotionToken     string `json:"notionToken"`
	ServerName      string `json:"serverName,omitempty"`
	ServerVersion   string `json:"serverVersion,omitempty"`
	DefaultPageSize int    `json:"defaultPageSize,omitempty"`
	MaxRetries      int    `json:"maxRetries,omitempty"`
	Timeout         int    `json:"timeout,omitempty"` // 秒
}

// DefaultMCPConfig 返回默认配置
func DefaultMCPConfig() *MCPConfig {
	return &MCPConfig{
		ServerName:      "notion-mcp-server",
		ServerVersion:   "1.0.0",
		DefaultPageSize: 10,
		MaxRetries:      3,
		Timeout:         30,
	}
}

// Validate 验证配置
func (c *MCPConfig) Validate() error {
	if c.NotionToken == "" {
		return &MCPError{
			Code:    -32602,
			Message: "Invalid params: notionToken is required",
		}
	}
	if c.DefaultPageSize <= 0 {
		c.DefaultPageSize = 10
	}
	if c.MaxRetries <= 0 {
		c.MaxRetries = 3
	}
	if c.Timeout <= 0 {
		c.Timeout = 30
	}
	return nil
}

// MCPError 表示 MCP 错误
type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (e *MCPError) Error() string {
	return e.Message
}

// 常用错误代码
const (
	ErrCodeParseError     = -32700
	ErrCodeInvalidRequest = -32600
	ErrCodeMethodNotFound = -32601
	ErrCodeInvalidParams  = -32602
	ErrCodeInternalError  = -32603
	ErrCodeNotionAPIError = -32000
	ErrCodeAuthError      = -32001
	ErrCodeRateLimitError = -32002
)

// 常用错误消息
const (
	ErrMsgParseError     = "Parse error"
	ErrMsgInvalidRequest = "Invalid Request"
	ErrMsgMethodNotFound = "Method not found"
	ErrMsgInvalidParams  = "Invalid params"
	ErrMsgInternalError  = "Internal error"
	ErrMsgNotionAPIError = "Notion API error"
	ErrMsgAuthError      = "Authentication error"
	ErrMsgRateLimitError = "Rate limit exceeded"
)
