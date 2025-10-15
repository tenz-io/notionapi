package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/tenz-io/notionapi"
)

// MCPServer 表示 MCP 服务器
type MCPServer struct {
	config       *MCPConfig
	notionClient *notionapi.Client
	tools        []Tool
	resources    []Resource
}

// NewMCPServer 创建新的 MCP 服务器
func NewMCPServer(config *MCPConfig) (*MCPServer, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	// 创建 Notion 客户端
	notionClient := notionapi.NewClient(notionapi.Token(config.NotionToken))

	server := &MCPServer{
		config:       config,
		notionClient: notionClient,
		tools:        getDefaultTools(),
		resources:    getDefaultResources(),
	}

	return server, nil
}

// HandleRequest 处理 MCP 请求
func (s *MCPServer) HandleRequest(ctx context.Context, req *Request) *Response {
	// 设置超时
	if s.config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(s.config.Timeout)*time.Second)
		defer cancel()
	}

	switch req.Method {
	case "initialize":
		return s.handleInitialize(req)
	case "tools/list":
		return s.handleToolsList(req)
	case "tools/call":
		return s.handleToolsCall(ctx, req)
	case "resources/list":
		return s.handleResourcesList(req)
	case "resources/read":
		return s.handleResourcesRead(ctx, req)
	default:
		return &Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &Error{
				Code:    ErrCodeMethodNotFound,
				Message: ErrMsgMethodNotFound,
			},
		}
	}
}

// handleInitialize 处理初始化请求
func (s *MCPServer) handleInitialize(req *Request) *Response {
	result := map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities": map[string]interface{}{
			"tools": map[string]interface{}{
				"listChanged": false,
			},
			"resources": map[string]interface{}{
				"subscribe":   false,
				"listChanged": false,
			},
		},
		"serverInfo": map[string]interface{}{
			"name":    s.config.ServerName,
			"version": s.config.ServerVersion,
		},
	}

	return &Response{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	}
}

// handleToolsList 处理工具列表请求
func (s *MCPServer) handleToolsList(req *Request) *Response {
	return &Response{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"tools": s.tools,
		},
	}
}

// handleToolsCall 处理工具调用请求
func (s *MCPServer) handleToolsCall(ctx context.Context, req *Request) *Response {
	// 解析工具调用参数
	var toolCall ToolCall

	// 检查 Params 的类型并正确解析
	if req.Params == nil {
		return &Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &Error{
				Code:    ErrCodeInvalidParams,
				Message: "Missing params",
			},
		}
	}

	// 如果 Params 已经是 ToolCall 类型，直接使用
	if tc, ok := req.Params.(ToolCall); ok {
		toolCall = tc
	} else {
		// 否则尝试从 JSON 解析
		paramsBytes, err := json.Marshal(req.Params)
		if err != nil {
			return &Response{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error: &Error{
					Code:    ErrCodeParseError,
					Message: ErrMsgParseError,
				},
			}
		}

		if err := json.Unmarshal(paramsBytes, &toolCall); err != nil {
			return &Response{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error: &Error{
					Code:    ErrCodeParseError,
					Message: ErrMsgParseError,
				},
			}
		}
	}

	// 根据工具名称调用相应的方法
	var result *ToolResult
	var err error

	switch toolCall.Name {
	case "notion_search":
		result, err = s.handleNotionSearch(ctx, toolCall.Arguments)
	case "notion_create_page":
		result, err = s.handleNotionCreatePage(ctx, toolCall.Arguments)
	case "notion_update_page":
		result, err = s.handleNotionUpdatePage(ctx, toolCall.Arguments)
	case "notion_append_block":
		result, err = s.handleNotionAppendBlock(ctx, toolCall.Arguments)
	default:
		return &Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &Error{
				Code:    ErrCodeMethodNotFound,
				Message: fmt.Sprintf("Tool '%s' not found", toolCall.Name),
			},
		}
	}

	if err != nil {
		return &Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &Error{
				Code:    ErrCodeInternalError,
				Message: err.Error(),
			},
		}
	}

	return &Response{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	}
}

// handleResourcesList 处理资源列表请求
func (s *MCPServer) handleResourcesList(req *Request) *Response {
	return &Response{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"resources": s.resources,
		},
	}
}

// handleResourcesRead 处理资源读取请求
func (s *MCPServer) handleResourcesRead(ctx context.Context, req *Request) *Response {
	// 解析资源读取参数
	var params map[string]interface{}

	// 检查 Params 的类型并正确解析
	if req.Params == nil {
		return &Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &Error{
				Code:    ErrCodeInvalidParams,
				Message: "Missing params",
			},
		}
	}

	// 如果 Params 已经是 map[string]interface{} 类型，直接使用
	if p, ok := req.Params.(map[string]interface{}); ok {
		params = p
	} else {
		// 否则尝试从 JSON 解析
		paramsBytes, err := json.Marshal(req.Params)
		if err != nil {
			return &Response{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error: &Error{
					Code:    ErrCodeParseError,
					Message: ErrMsgParseError,
				},
			}
		}

		if err := json.Unmarshal(paramsBytes, &params); err != nil {
			return &Response{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error: &Error{
					Code:    ErrCodeParseError,
					Message: ErrMsgParseError,
				},
			}
		}
	}

	uri, ok := params["uri"].(string)
	if !ok {
		return &Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &Error{
				Code:    ErrCodeInvalidParams,
				Message: "Invalid params: uri is required",
			},
		}
	}

	// 根据 URI 处理不同的资源
	var content []Content
	var err error

	switch uri {
	case "notion://workspace":
		content, err = s.handleWorkspaceResource(ctx)
	default:
		return &Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &Error{
				Code:    ErrCodeMethodNotFound,
				Message: fmt.Sprintf("Resource '%s' not found", uri),
			},
		}
	}

	if err != nil {
		return &Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &Error{
				Code:    ErrCodeInternalError,
				Message: err.Error(),
			},
		}
	}

	return &Response{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"contents": content,
		},
	}
}

// handleNotionSearch 处理 Notion 搜索
func (s *MCPServer) handleNotionSearch(ctx context.Context, args map[string]interface{}) (*ToolResult, error) {
	// 解析搜索参数
	params := &NotionSearchParams{
		PageSize: s.config.DefaultPageSize,
	}

	if query, ok := args["query"].(string); ok {
		params.Query = query
	}
	if filter, ok := args["filter"].(string); ok {
		params.Filter = filter
	}
	if sortBy, ok := args["sortBy"].(string); ok {
		params.SortBy = sortBy
	}
	if sortOrder, ok := args["sortOrder"].(string); ok {
		params.SortOrder = sortOrder
	}
	if startCursor, ok := args["startCursor"].(string); ok {
		params.StartCursor = startCursor
	}
	if pageSize, ok := args["pageSize"].(float64); ok {
		params.PageSize = int(pageSize)
	}

	// 构建搜索请求
	// 由于 SearchFilter 是结构体而不是指针，我们需要使用一个技巧
	// 创建一个自定义的搜索请求结构体，其中 Filter 是指针类型
	type CustomSearchRequest struct {
		Query       string                  `json:"query,omitempty"`
		Sort        *notionapi.SortObject   `json:"sort,omitempty"`
		Filter      *notionapi.SearchFilter `json:"filter,omitempty"`
		StartCursor notionapi.Cursor        `json:"start_cursor,omitempty"`
		PageSize    int                     `json:"page_size,omitempty"`
	}

	customReq := &CustomSearchRequest{
		Query:       params.Query,
		StartCursor: notionapi.Cursor(params.StartCursor),
		PageSize:    params.PageSize,
	}

	// 设置过滤器 - 只有当确实需要过滤器时才设置
	if params.Filter == "page" {
		customReq.Filter = &notionapi.SearchFilter{
			Property: "object",
			Value:    "page",
		}
	} else if params.Filter == "database" {
		customReq.Filter = &notionapi.SearchFilter{
			Property: "object",
			Value:    "database",
		}
	}
	// 注意：如果 params.Filter 为空，customReq.Filter 保持为 nil，不会被序列化
	// 注意：如果 params.Filter 为空，customReq.Filter 保持为 nil，不会被序列化

	// 设置排序
	if params.SortBy != "" && params.SortOrder != "" {
		sortOrder := notionapi.SortOrderASC
		if params.SortOrder == "descending" {
			sortOrder = notionapi.SortOrderDESC
		}
		customReq.Sort = &notionapi.SortObject{
			Direction: sortOrder,
			Timestamp: notionapi.TimestampLastEdited,
		}
	}

	// 创建一个自定义的搜索请求结构体，其中 Filter 是指针类型
	type CustomSearchRequestWithPointer struct {
		Query       string                  `json:"query,omitempty"`
		Sort        *notionapi.SortObject   `json:"sort,omitempty"`
		Filter      *notionapi.SearchFilter `json:"filter,omitempty"`
		StartCursor notionapi.Cursor        `json:"start_cursor,omitempty"`
		PageSize    int                     `json:"page_size,omitempty"`
	}

	// 将 customReq 转换为指针版本
	pointerReq := &CustomSearchRequestWithPointer{
		Query:       customReq.Query,
		Sort:        customReq.Sort,
		Filter:      customReq.Filter,
		StartCursor: customReq.StartCursor,
		PageSize:    customReq.PageSize,
	}

	// 序列化为 JSON
	jsonData, err := json.Marshal(pointerReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal search request: %w", err)
	}

	// 反序列化为标准的 SearchRequest
	var searchReq notionapi.SearchRequest
	if err := json.Unmarshal(jsonData, &searchReq); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search request: %w", err)
	}

	// 执行搜索
	searchResp, err := s.notionClient.Search.Do(ctx, &searchReq)
	if err != nil {
		return nil, fmt.Errorf("failed to search Notion: %w", err)
	}

	// 转换结果
	result := &NotionSearchResult{
		Object:     string(searchResp.Object),
		HasMore:    searchResp.HasMore,
		NextCursor: string(searchResp.NextCursor),
		Results:    make([]NotionItem, len(searchResp.Results)),
	}

	for i, item := range searchResp.Results {
		result.Results[i] = s.convertToNotionItem(item)
	}

	// 返回结果
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal search result: %w", err)
	}

	return &ToolResult{
		Content: []Content{
			{
				Type: "text",
				Text: string(resultJSON),
			},
		},
	}, nil
}

// handleNotionCreatePage 处理创建 Notion 页面
func (s *MCPServer) handleNotionCreatePage(ctx context.Context, args map[string]interface{}) (*ToolResult, error) {
	// 解析创建页面参数
	params := &NotionCreatePageParams{}

	if parentID, ok := args["parentId"].(string); ok {
		params.ParentID = parentID
	} else {
		return nil, fmt.Errorf("parentId is required")
	}

	if title, ok := args["title"].(string); ok {
		params.Title = title
	} else {
		return nil, fmt.Errorf("title is required")
	}

	if content, ok := args["content"].(string); ok {
		params.Content = content
	}

	if properties, ok := args["properties"].(map[string]interface{}); ok {
		params.Properties = properties
	}

	if icon, ok := args["icon"].(map[string]interface{}); ok {
		params.Icon = &Icon{
			Type:  icon["type"].(string),
			Emoji: icon["emoji"].(string),
			URL:   icon["url"].(string),
		}
	}

	if cover, ok := args["cover"].(map[string]interface{}); ok {
		params.Cover = &Cover{
			Type: cover["type"].(string),
			URL:  cover["url"].(string),
		}
	}

	// 构建创建页面请求
	createReq := &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			Type:   notionapi.ParentTypePageID,
			PageID: notionapi.PageID(params.ParentID),
		},
		Properties: notionapi.Properties{
			"title": notionapi.TitleProperty{
				Type: notionapi.PropertyTypeTitle,
				Title: []notionapi.RichText{
					{
						Type: notionapi.RichTextTypeText,
						Text: &notionapi.Text{
							Content: params.Title,
						},
					},
				},
			},
		},
	}

	// 添加其他属性
	for key, value := range params.Properties {
		// 这里需要根据具体的属性类型进行转换
		// 为了简化，这里只处理基本的文本属性
		if strValue, ok := value.(string); ok {
			createReq.Properties[key] = notionapi.RichTextProperty{
				Type: notionapi.PropertyTypeRichText,
				RichText: []notionapi.RichText{
					{
						Type: notionapi.RichTextTypeText,
						Text: &notionapi.Text{
							Content: strValue,
						},
					},
				},
			}
		}
	}

	// 添加图标
	if params.Icon != nil {
		if params.Icon.Type == "emoji" {
			emoji := notionapi.Emoji(params.Icon.Emoji)
			createReq.Icon = &notionapi.Icon{
				Type:  "emoji",
				Emoji: &emoji,
			}
		}
	}

	// 添加封面
	if params.Cover != nil {
		createReq.Cover = &notionapi.Image{
			Type: notionapi.FileTypeExternal,
			External: &notionapi.FileObject{
				URL: params.Cover.URL,
			},
		}
	}

	// 执行创建页面
	page, err := s.notionClient.Page.Create(ctx, createReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create Notion page: %w", err)
	}

	// 如果有内容，添加块
	if params.Content != "" {
		blocks := s.parseContentToBlocks(params.Content)
		if len(blocks) > 0 {
			appendReq := &notionapi.AppendBlockChildrenRequest{
				Children: blocks,
			}
			_, err = s.notionClient.Block.AppendChildren(ctx, notionapi.BlockID(page.ID), appendReq)
			if err != nil {
				log.Printf("Warning: failed to append content blocks: %v", err)
			}
		}
	}

	// 转换结果
	result := s.convertToNotionPageResult(page)

	// 返回结果
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal page result: %w", err)
	}

	return &ToolResult{
		Content: []Content{
			{
				Type: "text",
				Text: string(resultJSON),
			},
		},
	}, nil
}

// handleNotionUpdatePage 处理更新 Notion 页面
func (s *MCPServer) handleNotionUpdatePage(ctx context.Context, args map[string]interface{}) (*ToolResult, error) {
	// 解析更新页面参数
	params := &NotionUpdatePageParams{}

	if pageID, ok := args["pageId"].(string); ok {
		params.PageID = pageID
	} else {
		return nil, fmt.Errorf("pageId is required")
	}

	if title, ok := args["title"].(string); ok {
		params.Title = title
	}

	if content, ok := args["content"].(string); ok {
		params.Content = content
	}

	if properties, ok := args["properties"].(map[string]interface{}); ok {
		params.Properties = properties
	}

	if archived, ok := args["archived"].(bool); ok {
		params.Archived = &archived
	}

	// 构建更新页面请求
	updateReq := &notionapi.PageUpdateRequest{}

	// 更新标题
	if params.Title != "" {
		updateReq.Properties = notionapi.Properties{
			"title": notionapi.TitleProperty{
				Type: notionapi.PropertyTypeTitle,
				Title: []notionapi.RichText{
					{
						Type: notionapi.RichTextTypeText,
						Text: &notionapi.Text{
							Content: params.Title,
						},
					},
				},
			},
		}
	}

	// 更新其他属性
	for key, value := range params.Properties {
		if strValue, ok := value.(string); ok {
			if updateReq.Properties == nil {
				updateReq.Properties = make(notionapi.Properties)
			}
			updateReq.Properties[key] = notionapi.RichTextProperty{
				Type: notionapi.PropertyTypeRichText,
				RichText: []notionapi.RichText{
					{
						Type: notionapi.RichTextTypeText,
						Text: &notionapi.Text{
							Content: strValue,
						},
					},
				},
			}
		}
	}

	// 设置归档状态
	if params.Archived != nil {
		updateReq.Archived = *params.Archived
	}

	// 执行更新页面
	page, err := s.notionClient.Page.Update(ctx, notionapi.PageID(params.PageID), updateReq)
	if err != nil {
		return nil, fmt.Errorf("failed to update Notion page: %w", err)
	}

	// 如果有内容，添加块
	if params.Content != "" {
		blocks := s.parseContentToBlocks(params.Content)
		if len(blocks) > 0 {
			appendReq := &notionapi.AppendBlockChildrenRequest{
				Children: blocks,
			}
			_, err = s.notionClient.Block.AppendChildren(ctx, notionapi.BlockID(page.ID), appendReq)
			if err != nil {
				log.Printf("Warning: failed to append content blocks: %v", err)
			}
		}
	}

	// 转换结果
	result := s.convertToNotionPageResult(page)

	// 返回结果
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal page result: %w", err)
	}

	return &ToolResult{
		Content: []Content{
			{
				Type: "text",
				Text: string(resultJSON),
			},
		},
	}, nil
}

// handleNotionAppendBlock 处理添加块内容
func (s *MCPServer) handleNotionAppendBlock(ctx context.Context, args map[string]interface{}) (*ToolResult, error) {
	// 解析添加块参数
	params := &NotionAppendBlockParams{}

	if pageID, ok := args["pageId"].(string); ok {
		params.PageID = pageID
	} else {
		return nil, fmt.Errorf("pageId is required")
	}

	if content, ok := args["content"].(string); ok {
		params.Content = content
	} else {
		return nil, fmt.Errorf("content is required")
	}

	if blockType, ok := args["blockType"].(string); ok {
		params.BlockType = blockType
	} else {
		params.BlockType = "paragraph" // 默认段落
	}

	// 解析内容为块
	blocks := s.parseContentToBlocksWithType(params.Content, params.BlockType)

	// 构建添加块请求
	appendReq := &notionapi.AppendBlockChildrenRequest{
		Children: blocks,
	}

	// 执行添加块
	result, err := s.notionClient.Block.AppendChildren(ctx, notionapi.BlockID(params.PageID), appendReq)
	if err != nil {
		return nil, fmt.Errorf("failed to append blocks: %w", err)
	}

	// 转换结果
	blockResult := &NotionBlockResult{
		Object:  string(result.Object),
		Results: make([]BlockItem, len(result.Results)),
	}

	for i, block := range result.Results {
		blockResult.Results[i] = s.convertToBlockItem(block)
	}

	// 返回结果
	resultJSON, err := json.Marshal(blockResult)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal block result: %w", err)
	}

	return &ToolResult{
		Content: []Content{
			{
				Type: "text",
				Text: string(resultJSON),
			},
		},
	}, nil
}

// handleWorkspaceResource 处理工作区资源
func (s *MCPServer) handleWorkspaceResource(ctx context.Context) ([]Content, error) {
	// 获取工作区信息（通过搜索所有页面和数据库）
	// 使用与 handleNotionSearch 相同的自定义搜索请求结构体
	type CustomSearchRequest struct {
		Query       string                  `json:"query,omitempty"`
		Sort        *notionapi.SortObject   `json:"sort,omitempty"`
		Filter      *notionapi.SearchFilter `json:"filter,omitempty"`
		StartCursor notionapi.Cursor        `json:"start_cursor,omitempty"`
		PageSize    int                     `json:"page_size,omitempty"`
	}

	customReq := &CustomSearchRequest{
		PageSize: 100,
	}

	// 将自定义请求转换为标准的 SearchRequest
	jsonData, err := json.Marshal(customReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal search request: %w", err)
	}

	var searchReq notionapi.SearchRequest
	if err := json.Unmarshal(jsonData, &searchReq); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search request: %w", err)
	}

	searchResp, err := s.notionClient.Search.Do(ctx, &searchReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get workspace info: %w", err)
	}

	// 构建工作区信息
	workspaceInfo := map[string]interface{}{
		"totalPages":     0,
		"totalDatabases": 0,
		"recentItems":    make([]NotionItem, 0),
	}

	for _, item := range searchResp.Results {
		notionItem := s.convertToNotionItem(item)
		workspaceInfo["recentItems"] = append(workspaceInfo["recentItems"].([]NotionItem), notionItem)

		if item.GetObject() == notionapi.ObjectTypePage {
			workspaceInfo["totalPages"] = workspaceInfo["totalPages"].(int) + 1
		} else if item.GetObject() == notionapi.ObjectTypeDatabase {
			workspaceInfo["totalDatabases"] = workspaceInfo["totalDatabases"].(int) + 1
		}
	}

	workspaceJSON, err := json.Marshal(workspaceInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal workspace info: %w", err)
	}

	return []Content{
		{
			Type: "text",
			Text: string(workspaceJSON),
		},
	}, nil
}

// 辅助方法

// convertToNotionItem 转换 Notion 对象为 NotionItem
func (s *MCPServer) convertToNotionItem(item notionapi.Object) NotionItem {
	notionItem := NotionItem{
		Object:    string(item.GetObject()),
		ID:        "",    // 需要根据具体类型获取
		Archived:  false, // 需要根据具体类型获取
		URL:       "",
		PublicURL: "",
	}

	// 根据类型设置特定字段
	switch obj := item.(type) {
	case *notionapi.Page:
		notionItem.ID = string(obj.ID)
		notionItem.Archived = obj.Archived
		notionItem.CreatedTime = obj.CreatedTime
		notionItem.LastEditedTime = obj.LastEditedTime
		notionItem.CreatedBy = s.convertToUser(&obj.CreatedBy)
		notionItem.LastEditedBy = s.convertToUser(&obj.LastEditedBy)
		notionItem.Properties = s.convertProperties(obj.Properties)
		notionItem.Parent = s.convertParent(&obj.Parent)
		notionItem.URL = obj.URL
		notionItem.PublicURL = obj.PublicURL
		notionItem.Icon = s.convertIcon(obj.Icon)
		notionItem.Cover = s.convertCover(obj.Cover)
		notionItem.Title = s.extractTitle(obj.Properties)

	case *notionapi.Database:
		notionItem.ID = string(obj.ID)
		notionItem.Archived = obj.Archived
		notionItem.CreatedTime = obj.CreatedTime
		notionItem.LastEditedTime = obj.LastEditedTime
		notionItem.CreatedBy = s.convertToUser(&obj.CreatedBy)
		notionItem.LastEditedBy = s.convertToUser(&obj.LastEditedBy)
		notionItem.Properties = s.convertPropertyConfigs(obj.Properties)
		notionItem.Parent = s.convertParent(&obj.Parent)
		notionItem.Icon = s.convertIcon(obj.Icon)
		notionItem.Cover = s.convertCover(obj.Cover)
		notionItem.Title = s.extractTitleFromConfigs(obj.Properties)
	}

	return notionItem
}

// convertToNotionPageResult 转换页面为结果
func (s *MCPServer) convertToNotionPageResult(page *notionapi.Page) *NotionPageResult {
	return &NotionPageResult{
		Object:         string(page.Object),
		ID:             string(page.ID),
		CreatedTime:    page.CreatedTime,
		LastEditedTime: page.LastEditedTime,
		CreatedBy:      s.convertToUser(&page.CreatedBy),
		LastEditedBy:   s.convertToUser(&page.LastEditedBy),
		Archived:       page.Archived,
		Properties:     s.convertProperties(page.Properties),
		Parent:         s.convertParent(&page.Parent),
		URL:            page.URL,
		PublicURL:      page.PublicURL,
		Icon:           s.convertIcon(page.Icon),
		Cover:          s.convertCover(page.Cover),
	}
}

// convertToBlockItem 转换块为 BlockItem
func (s *MCPServer) convertToBlockItem(block notionapi.Block) BlockItem {
	return BlockItem{
		Object:         string(block.GetObject()),
		ID:             string(block.GetID()),
		Type:           string(block.GetType()),
		CreatedTime:    *block.GetCreatedTime(),
		LastEditedTime: *block.GetLastEditedTime(),
		CreatedBy:      s.convertToUser(block.GetCreatedBy()),
		LastEditedBy:   s.convertToUser(block.GetLastEditedBy()),
		HasChildren:    block.GetHasChildren(),
		Archived:       block.GetArchived(),
		Parent:         s.convertParent(block.GetParent()),
		Content:        map[string]interface{}{"text": block.GetRichTextString()},
	}
}

// convertToUser 转换用户
func (s *MCPServer) convertToUser(user *notionapi.User) *User {
	if user == nil || user.ID == "" {
		return nil
	}
	return &User{
		Object:    string(user.Object),
		ID:        string(user.ID),
		Name:      user.Name,
		AvatarURL: user.AvatarURL,
		Type:      string(user.Type),
	}
}

// convertParent 转换父级对象
func (s *MCPServer) convertParent(parent *notionapi.Parent) *Parent {
	if parent == nil {
		return nil
	}
	return &Parent{
		Type:       string(parent.Type),
		PageID:     string(parent.PageID),
		DatabaseID: string(parent.DatabaseID),
		BlockID:    string(parent.BlockID),
		Workspace:  parent.Workspace,
	}
}

// convertIcon 转换图标
func (s *MCPServer) convertIcon(icon *notionapi.Icon) *Icon {
	if icon == nil {
		return nil
	}
	result := &Icon{
		Type: string(icon.Type),
	}
	if icon.Emoji != nil {
		result.Emoji = string(*icon.Emoji)
	}
	if icon.External != nil {
		result.URL = icon.External.URL
	}
	return result
}

// convertCover 转换封面
func (s *MCPServer) convertCover(cover *notionapi.Image) *Cover {
	if cover == nil {
		return nil
	}
	return &Cover{
		Type: string(cover.Type),
		URL:  cover.GetURL(),
	}
}

// convertProperties 转换属性
func (s *MCPServer) convertProperties(properties notionapi.Properties) map[string]interface{} {
	result := make(map[string]interface{})
	for key, prop := range properties {
		result[key] = map[string]interface{}{
			"type": string(prop.GetType()),
			"id":   prop.GetID(),
		}
	}
	return result
}

// convertPropertyConfigs 转换属性配置
func (s *MCPServer) convertPropertyConfigs(configs notionapi.PropertyConfigs) map[string]interface{} {
	result := make(map[string]interface{})
	for key, config := range configs {
		result[key] = map[string]interface{}{
			"type": string(config.GetType()),
			"id":   config.GetID(),
		}
	}
	return result
}

// extractTitle 提取标题
func (s *MCPServer) extractTitle(properties notionapi.Properties) string {
	for _, prop := range properties {
		if prop.GetType() == notionapi.PropertyTypeTitle {
			// 尝试转换为 TitleProperty
			if titleProp, ok := prop.(notionapi.TitleProperty); ok && len(titleProp.Title) > 0 {
				return titleProp.Title[0].PlainText
			}
		}
	}
	return ""
}

// extractTitleFromConfigs 从属性配置中提取标题
func (s *MCPServer) extractTitleFromConfigs(configs notionapi.PropertyConfigs) string {
	for _, config := range configs {
		if config.GetType() == notionapi.PropertyConfigTypeTitle {
			// 对于数据库，标题通常存储在 Title 字段中，而不是属性中
			// 这里返回空字符串，实际标题应该从数据库的 Title 字段获取
			return ""
		}
	}
	return ""
}

// parseContentToBlocks 解析内容为块
func (s *MCPServer) parseContentToBlocks(content string) []notionapi.Block {
	// 简单的解析逻辑，将内容按行分割并创建段落块
	lines := []string{content} // 简化处理，将整个内容作为一个段落
	blocks := make([]notionapi.Block, len(lines))

	for i, line := range lines {
		blocks[i] = &notionapi.ParagraphBlock{
			BasicBlock: notionapi.BasicBlock{
				Type: notionapi.BlockTypeParagraph,
			},
			Paragraph: notionapi.Paragraph{
				RichText: []notionapi.RichText{
					{
						Type: notionapi.RichTextTypeText,
						Text: &notionapi.Text{
							Content: line,
						},
						PlainText: line,
					},
				},
			},
		}
	}

	return blocks
}

// parseContentToBlocksWithType 根据类型解析内容为块
func (s *MCPServer) parseContentToBlocksWithType(content string, blockType string) []notionapi.Block {
	blocks := make([]notionapi.Block, 1)

	switch blockType {
	case "heading_1":
		blocks[0] = &notionapi.Heading1Block{
			BasicBlock: notionapi.BasicBlock{
				Type: notionapi.BlockTypeHeading1,
			},
			Heading1: notionapi.Heading{
				RichText: []notionapi.RichText{
					{
						Type: notionapi.RichTextTypeText,
						Text: &notionapi.Text{
							Content: content,
						},
						PlainText: content,
					},
				},
			},
		}
	case "heading_2":
		blocks[0] = &notionapi.Heading2Block{
			BasicBlock: notionapi.BasicBlock{
				Type: notionapi.BlockTypeHeading2,
			},
			Heading2: notionapi.Heading{
				RichText: []notionapi.RichText{
					{
						Type: notionapi.RichTextTypeText,
						Text: &notionapi.Text{
							Content: content,
						},
						PlainText: content,
					},
				},
			},
		}
	case "heading_3":
		blocks[0] = &notionapi.Heading3Block{
			BasicBlock: notionapi.BasicBlock{
				Type: notionapi.BlockTypeHeading3,
			},
			Heading3: notionapi.Heading{
				RichText: []notionapi.RichText{
					{
						Type: notionapi.RichTextTypeText,
						Text: &notionapi.Text{
							Content: content,
						},
						PlainText: content,
					},
				},
			},
		}
	case "bulleted_list_item":
		blocks[0] = &notionapi.BulletedListItemBlock{
			BasicBlock: notionapi.BasicBlock{
				Type: notionapi.BlockTypeBulletedListItem,
			},
			BulletedListItem: notionapi.ListItem{
				RichText: []notionapi.RichText{
					{
						Type: notionapi.RichTextTypeText,
						Text: &notionapi.Text{
							Content: content,
						},
						PlainText: content,
					},
				},
			},
		}
	case "numbered_list_item":
		blocks[0] = &notionapi.NumberedListItemBlock{
			BasicBlock: notionapi.BasicBlock{
				Type: notionapi.BlockTypeNumberedListItem,
			},
			NumberedListItem: notionapi.ListItem{
				RichText: []notionapi.RichText{
					{
						Type: notionapi.RichTextTypeText,
						Text: &notionapi.Text{
							Content: content,
						},
						PlainText: content,
					},
				},
			},
		}
	case "to_do":
		blocks[0] = &notionapi.ToDoBlock{
			BasicBlock: notionapi.BasicBlock{
				Type: notionapi.BlockTypeToDo,
			},
			ToDo: notionapi.ToDo{
				RichText: []notionapi.RichText{
					{
						Type: notionapi.RichTextTypeText,
						Text: &notionapi.Text{
							Content: content,
						},
						PlainText: content,
					},
				},
				Checked: false,
			},
		}
	case "code":
		blocks[0] = &notionapi.CodeBlock{
			BasicBlock: notionapi.BasicBlock{
				Type: notionapi.BlockTypeCode,
			},
			Code: notionapi.Code{
				RichText: []notionapi.RichText{
					{
						Type: notionapi.RichTextTypeText,
						Text: &notionapi.Text{
							Content: content,
						},
						PlainText: content,
					},
				},
				Language: "plain text",
			},
		}
	case "quote":
		blocks[0] = &notionapi.QuoteBlock{
			BasicBlock: notionapi.BasicBlock{
				Type: notionapi.BlockTypeQuote,
			},
			Quote: notionapi.Quote{
				RichText: []notionapi.RichText{
					{
						Type: notionapi.RichTextTypeText,
						Text: &notionapi.Text{
							Content: content,
						},
						PlainText: content,
					},
				},
			},
		}
	case "callout":
		blocks[0] = &notionapi.CalloutBlock{
			BasicBlock: notionapi.BasicBlock{
				Type: notionapi.BlockTypeCallout,
			},
			Callout: notionapi.Callout{
				RichText: []notionapi.RichText{
					{
						Type: notionapi.RichTextTypeText,
						Text: &notionapi.Text{
							Content: content,
						},
						PlainText: content,
					},
				},
				Icon: &notionapi.Icon{
					Type:  "emoji",
					Emoji: func() *notionapi.Emoji { e := notionapi.Emoji("💡"); return &e }(),
				},
			},
		}
	default: // paragraph
		blocks[0] = &notionapi.ParagraphBlock{
			BasicBlock: notionapi.BasicBlock{
				Type: notionapi.BlockTypeParagraph,
			},
			Paragraph: notionapi.Paragraph{
				RichText: []notionapi.RichText{
					{
						Type: notionapi.RichTextTypeText,
						Text: &notionapi.Text{
							Content: content,
						},
						PlainText: content,
					},
				},
			},
		}
	}

	return blocks
}
