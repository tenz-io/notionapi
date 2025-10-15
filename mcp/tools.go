package mcp

// getDefaultTools 返回默认工具列表
func getDefaultTools() []Tool {
	return []Tool{
		{
			Name:        "notion_search",
			Description: "在 Notion 工作区中搜索页面和数据库",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query": map[string]interface{}{
						"type":        "string",
						"description": "搜索查询字符串",
					},
					"filter": map[string]interface{}{
						"type":        "string",
						"description": "过滤结果类型：'page' 或 'database'",
						"enum":        []string{"page", "database"},
					},
					"sortBy": map[string]interface{}{
						"type":        "string",
						"description": "排序字段：'last_edited_time'",
						"enum":        []string{"last_edited_time"},
					},
					"sortOrder": map[string]interface{}{
						"type":        "string",
						"description": "排序顺序：'ascending' 或 'descending'",
						"enum":        []string{"ascending", "descending"},
					},
					"startCursor": map[string]interface{}{
						"type":        "string",
						"description": "分页起始游标",
					},
					"pageSize": map[string]interface{}{
						"type":        "integer",
						"description": "每页结果数量（最大100）",
						"minimum":     1,
						"maximum":     100,
					},
				},
			},
		},
		{
			Name:        "notion_create_page",
			Description: "在 Notion 中创建新页面",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"parentId": map[string]interface{}{
						"type":        "string",
						"description": "父页面或数据库的 ID",
					},
					"title": map[string]interface{}{
						"type":        "string",
						"description": "页面标题",
					},
					"content": map[string]interface{}{
						"type":        "string",
						"description": "页面内容（可选）",
					},
					"properties": map[string]interface{}{
						"type":        "object",
						"description": "页面属性（可选）",
					},
					"icon": map[string]interface{}{
						"type":        "object",
						"description": "页面图标（可选）",
						"properties": map[string]interface{}{
							"type": map[string]interface{}{
								"type":        "string",
								"description": "图标类型：'emoji' 或 'external'",
								"enum":        []string{"emoji", "external"},
							},
							"emoji": map[string]interface{}{
								"type":        "string",
								"description": "表情符号（当 type 为 'emoji' 时）",
							},
							"url": map[string]interface{}{
								"type":        "string",
								"description": "图标 URL（当 type 为 'external' 时）",
							},
						},
					},
					"cover": map[string]interface{}{
						"type":        "object",
						"description": "页面封面（可选）",
						"properties": map[string]interface{}{
							"type": map[string]interface{}{
								"type":        "string",
								"description": "封面类型：'external'",
								"enum":        []string{"external"},
							},
							"url": map[string]interface{}{
								"type":        "string",
								"description": "封面图片 URL",
							},
						},
					},
				},
				"required": []string{"parentId", "title"},
			},
		},
		{
			Name:        "notion_update_page",
			Description: "更新 Notion 页面",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"pageId": map[string]interface{}{
						"type":        "string",
						"description": "要更新的页面 ID",
					},
					"title": map[string]interface{}{
						"type":        "string",
						"description": "新的页面标题（可选）",
					},
					"content": map[string]interface{}{
						"type":        "string",
						"description": "要添加的页面内容（可选）",
					},
					"properties": map[string]interface{}{
						"type":        "object",
						"description": "要更新的页面属性（可选）",
					},
					"archived": map[string]interface{}{
						"type":        "boolean",
						"description": "是否归档页面（可选）",
					},
				},
				"required": []string{"pageId"},
			},
		},
		{
			Name:        "notion_append_block",
			Description: "向 Notion 页面添加块内容",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"pageId": map[string]interface{}{
						"type":        "string",
						"description": "目标页面 ID",
					},
					"content": map[string]interface{}{
						"type":        "string",
						"description": "要添加的内容",
					},
					"blockType": map[string]interface{}{
						"type":        "string",
						"description": "块类型",
						"enum": []string{
							"paragraph",
							"heading_1",
							"heading_2",
							"heading_3",
							"bulleted_list_item",
							"numbered_list_item",
							"to_do",
							"code",
							"quote",
							"callout",
						},
					},
				},
				"required": []string{"pageId", "content"},
			},
		},
	}
}

// getDefaultResources 返回默认资源列表
func getDefaultResources() []Resource {
	return []Resource{
		{
			URI:         "notion://workspace",
			Name:        "Notion Workspace",
			Description: "Notion 工作区信息，包括页面和数据库统计",
			MimeType:    "application/json",
		},
	}
}
