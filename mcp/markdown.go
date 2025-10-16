package mcp

import (
	"regexp"
	"strings"

	"github.com/tenz-io/notionapi"
)

// MarkdownToBlocks 将 markdown 文本转换为 Notion Block 数组
func MarkdownToBlocks(markdown string) []notionapi.Block {
	if strings.TrimSpace(markdown) == "" {
		return nil
	}

	lines := strings.Split(markdown, "\n")
	var blocks []notionapi.Block
	var currentParagraph []string
	var inCodeBlock bool
	var codeBlockContent []string
	var codeLanguage string

	for _, line := range lines {
		line = strings.TrimRight(line, "\r")

		// 处理代码块
		if strings.HasPrefix(line, "```") {
			if inCodeBlock {
				// 结束代码块
				if len(codeBlockContent) > 0 {
					blocks = append(blocks, createCodeBlock(strings.Join(codeBlockContent, "\n"), codeLanguage))
				}
				inCodeBlock = false
				codeBlockContent = nil
				codeLanguage = ""
			} else {
				// 开始代码块
				if len(currentParagraph) > 0 {
					blocks = append(blocks, createParagraph(strings.Join(currentParagraph, "\n")))
					currentParagraph = nil
				}
				inCodeBlock = true
				codeLanguage = strings.TrimSpace(strings.TrimPrefix(line, "```"))
			}
			continue
		}

		if inCodeBlock {
			codeBlockContent = append(codeBlockContent, line)
			continue
		}

		// 处理标题
		if strings.HasPrefix(line, "# ") {
			if len(currentParagraph) > 0 {
				blocks = append(blocks, createParagraph(strings.Join(currentParagraph, "\n")))
				currentParagraph = nil
			}
			blocks = append(blocks, createHeading1(strings.TrimSpace(strings.TrimPrefix(line, "# "))))
			continue
		}

		if strings.HasPrefix(line, "## ") {
			if len(currentParagraph) > 0 {
				blocks = append(blocks, createParagraph(strings.Join(currentParagraph, "\n")))
				currentParagraph = nil
			}
			blocks = append(blocks, createHeading2(strings.TrimSpace(strings.TrimPrefix(line, "## "))))
			continue
		}

		if strings.HasPrefix(line, "### ") {
			if len(currentParagraph) > 0 {
				blocks = append(blocks, createParagraph(strings.Join(currentParagraph, "\n")))
				currentParagraph = nil
			}
			blocks = append(blocks, createHeading3(strings.TrimSpace(strings.TrimPrefix(line, "### "))))
			continue
		}

		// 处理列表项
		if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") {
			if len(currentParagraph) > 0 {
				blocks = append(blocks, createParagraph(strings.Join(currentParagraph, "\n")))
				currentParagraph = nil
			}
			blocks = append(blocks, createBulletedListItem(strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(line, "- "), "* "))))
			continue
		}

		if matched, _ := regexp.MatchString(`^\d+\. `, line); matched {
			if len(currentParagraph) > 0 {
				blocks = append(blocks, createParagraph(strings.Join(currentParagraph, "\n")))
				currentParagraph = nil
			}
			content := regexp.MustCompile(`^\d+\. `).ReplaceAllString(line, "")
			blocks = append(blocks, createNumberedListItem(strings.TrimSpace(content)))
			continue
		}

		// 处理引用
		if strings.HasPrefix(line, "> ") {
			if len(currentParagraph) > 0 {
				blocks = append(blocks, createParagraph(strings.Join(currentParagraph, "\n")))
				currentParagraph = nil
			}
			blocks = append(blocks, createQuote(strings.TrimSpace(strings.TrimPrefix(line, "> "))))
			continue
		}

		// 处理分隔线
		if strings.TrimSpace(line) == "---" || strings.TrimSpace(line) == "***" {
			if len(currentParagraph) > 0 {
				blocks = append(blocks, createParagraph(strings.Join(currentParagraph, "\n")))
				currentParagraph = nil
			}
			blocks = append(blocks, createDivider())
			continue
		}

		// 处理任务列表
		if strings.HasPrefix(line, "- [ ] ") {
			if len(currentParagraph) > 0 {
				blocks = append(blocks, createParagraph(strings.Join(currentParagraph, "\n")))
				currentParagraph = nil
			}
			blocks = append(blocks, createToDo(strings.TrimSpace(strings.TrimPrefix(line, "- [ ] ")), false))
			continue
		}

		if strings.HasPrefix(line, "- [x] ") {
			if len(currentParagraph) > 0 {
				blocks = append(blocks, createParagraph(strings.Join(currentParagraph, "\n")))
				currentParagraph = nil
			}
			blocks = append(blocks, createToDo(strings.TrimSpace(strings.TrimPrefix(line, "- [x] ")), true))
			continue
		}

		// 处理空行
		if strings.TrimSpace(line) == "" {
			if len(currentParagraph) > 0 {
				blocks = append(blocks, createParagraph(strings.Join(currentParagraph, "\n")))
				currentParagraph = nil
			}
			continue
		}

		// 普通文本行
		currentParagraph = append(currentParagraph, line)
	}

	// 处理最后的段落
	if len(currentParagraph) > 0 {
		blocks = append(blocks, createParagraph(strings.Join(currentParagraph, "\n")))
	}

	return blocks
}

// 创建段落块
func createParagraph(text string) notionapi.Block {
	return &notionapi.ParagraphBlock{
		BasicBlock: notionapi.BasicBlock{
			Object: notionapi.ObjectTypeBlock,
			Type:   notionapi.BlockTypeParagraph,
		},
		Paragraph: notionapi.Paragraph{
			RichText: []notionapi.RichText{
				{
					Type:      notionapi.RichTextTypeText,
					Text:      &notionapi.Text{Content: text},
					PlainText: text,
				},
			},
		},
	}
}

// 创建标题1块
func createHeading1(text string) notionapi.Block {
	return &notionapi.Heading1Block{
		BasicBlock: notionapi.BasicBlock{
			Object: notionapi.ObjectTypeBlock,
			Type:   notionapi.BlockTypeHeading1,
		},
		Heading1: notionapi.Heading{
			RichText: []notionapi.RichText{
				{
					Type:      notionapi.RichTextTypeText,
					Text:      &notionapi.Text{Content: text},
					PlainText: text,
				},
			},
		},
	}
}

// 创建标题2块
func createHeading2(text string) notionapi.Block {
	return &notionapi.Heading2Block{
		BasicBlock: notionapi.BasicBlock{
			Object: notionapi.ObjectTypeBlock,
			Type:   notionapi.BlockTypeHeading2,
		},
		Heading2: notionapi.Heading{
			RichText: []notionapi.RichText{
				{
					Type:      notionapi.RichTextTypeText,
					Text:      &notionapi.Text{Content: text},
					PlainText: text,
				},
			},
		},
	}
}

// 创建标题3块
func createHeading3(text string) notionapi.Block {
	return &notionapi.Heading3Block{
		BasicBlock: notionapi.BasicBlock{
			Object: notionapi.ObjectTypeBlock,
			Type:   notionapi.BlockTypeHeading3,
		},
		Heading3: notionapi.Heading{
			RichText: []notionapi.RichText{
				{
					Type:      notionapi.RichTextTypeText,
					Text:      &notionapi.Text{Content: text},
					PlainText: text,
				},
			},
		},
	}
}

// 创建无序列表项块
func createBulletedListItem(text string) notionapi.Block {
	return &notionapi.BulletedListItemBlock{
		BasicBlock: notionapi.BasicBlock{
			Object: notionapi.ObjectTypeBlock,
			Type:   notionapi.BlockTypeBulletedListItem,
		},
		BulletedListItem: notionapi.ListItem{
			RichText: []notionapi.RichText{
				{
					Type:      notionapi.RichTextTypeText,
					Text:      &notionapi.Text{Content: text},
					PlainText: text,
				},
			},
		},
	}
}

// 创建有序列表项块
func createNumberedListItem(text string) notionapi.Block {
	return &notionapi.NumberedListItemBlock{
		BasicBlock: notionapi.BasicBlock{
			Object: notionapi.ObjectTypeBlock,
			Type:   notionapi.BlockTypeNumberedListItem,
		},
		NumberedListItem: notionapi.ListItem{
			RichText: []notionapi.RichText{
				{
					Type:      notionapi.RichTextTypeText,
					Text:      &notionapi.Text{Content: text},
					PlainText: text,
				},
			},
		},
	}
}

// 创建引用块
func createQuote(text string) notionapi.Block {
	return &notionapi.QuoteBlock{
		BasicBlock: notionapi.BasicBlock{
			Object: notionapi.ObjectTypeBlock,
			Type:   notionapi.BlockTypeQuote,
		},
		Quote: notionapi.Quote{
			RichText: []notionapi.RichText{
				{
					Type:      notionapi.RichTextTypeText,
					Text:      &notionapi.Text{Content: text},
					PlainText: text,
				},
			},
		},
	}
}

// 创建代码块
func createCodeBlock(code, language string) notionapi.Block {
	return &notionapi.CodeBlock{
		BasicBlock: notionapi.BasicBlock{
			Object: notionapi.ObjectTypeBlock,
			Type:   notionapi.BlockTypeCode,
		},
		Code: notionapi.Code{
			RichText: []notionapi.RichText{
				{
					Type:      notionapi.RichTextTypeText,
					Text:      &notionapi.Text{Content: code},
					PlainText: code,
				},
			},
			Language: language,
		},
	}
}

// 创建任务块
func createToDo(text string, checked bool) notionapi.Block {
	return &notionapi.ToDoBlock{
		BasicBlock: notionapi.BasicBlock{
			Object: notionapi.ObjectTypeBlock,
			Type:   notionapi.BlockTypeToDo,
		},
		ToDo: notionapi.ToDo{
			RichText: []notionapi.RichText{
				{
					Type:      notionapi.RichTextTypeText,
					Text:      &notionapi.Text{Content: text},
					PlainText: text,
				},
			},
			Checked: checked,
		},
	}
}

// 创建分隔线块
func createDivider() notionapi.Block {
	return &notionapi.DividerBlock{
		BasicBlock: notionapi.BasicBlock{
			Object: notionapi.ObjectTypeBlock,
			Type:   notionapi.BlockTypeDivider,
		},
		Divider: notionapi.Divider{},
	}
}
