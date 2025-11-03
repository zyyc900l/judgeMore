package utils

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"judgeMore/config"
	"net/http"
	"os"
	"strings"
	"time"
)

// OpenAIRequest 定义 OpenAI 请求体
type OpenAIRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens,omitempty"`
}

// CompetitionInfo 使用大写字段名 + 正确 JSON 标签（修正）
type CompetitionInfo struct {
	Success      string `json:"success"`
	EventName    string `json:"event_name"`
	EventSponsor string `json:"event_sponsor"`
	EventTime    string `json:"event_time"`
	AwardLevel   string `json:"award_level"`
}

type Message struct {
	Role    string    `json:"role"`
	Content []Content `json:"content"`
}

type Content struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
}

type ImageURL struct {
	URL string `json:"url"` // data:image/...;base64,...
}

// OpenAIResponse 定义响应体
type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// CompetitionInfo 保持不变

// callOpenAIWithImage 调用 OpenAI GPT-4V API
func CallGLM4VWithImage(ctx context.Context, imagePath string, apiKey string) (*CompetitionInfo, error) {
	imgData, err := os.ReadFile(imagePath)
	if err != nil {
		return nil, fmt.Errorf("读取图片失败: %w", err)
	}
	mimeType := http.DetectContentType(imgData)
	supportedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
		"image/gif":  true,
	}
	if !supportedTypes[mimeType] {
		return nil, fmt.Errorf("不支持的图片格式: %s，仅支持 JPEG、PNG、WebP、GIF", mimeType)
	}

	encoded := base64.StdEncoding.EncodeToString(imgData)
	dataURL := fmt.Sprintf("data:%s;base64,%s", mimeType, encoded)

	prompt := `请从图片中提取以下信息，并以严格的 JSON 格式返回，不要包含任何额外说明、Markdown 或解释：
{
  "success": "true/false"
  "event_name": "竞赛全称",
  "event_sponsor": "主办单位",
  "event_time": "竞赛时间（如：2024年5月）",
  "award_level": "获得奖项（如：一等奖）"
}
如果某项无法识别，请留空对应字符串。请注意你需要对图片进行判断，如果你判断图片显然不是荣誉证书或者奖状，无法识别出任何赛事或奖项信息，"success"为"false"，否则返回"true"`

	// 注意：模型必须是 glm-4v（支持视觉）
	reqBody := OpenAIRequest{
		Model: config.OpenAI.ApiModel, // ✅ 正确模型名
		Messages: []Message{{
			Role: "user",
			Content: []Content{
				{Type: "image_url", ImageURL: &ImageURL{URL: dataURL}},
				{Type: "text", Text: prompt},
			},
		}},
		MaxTokens: 1000,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求体失败: %w", err)
	}

	client := &http.Client{Timeout: 60 * time.Second}
	// ✅ 修正：URL 末尾不能有空格！
	apiURL := config.OpenAI.ApiUrl
	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GLM API 返回错误 %d: %s", resp.StatusCode, string(b))
	}

	var result OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if len(result.Choices) == 0 || result.Choices[0].Message.Content == "" {
		return nil, fmt.Errorf("GLM 返回空内容")
	}

	jsonText := result.Choices[0].Message.Content

	// 清理 Markdown
	cleaned := strings.TrimSpace(jsonText)
	if strings.HasPrefix(cleaned, "```json") {
		cleaned = strings.TrimPrefix(cleaned, "```json")
		cleaned = strings.TrimSuffix(cleaned, "```")
	} else if strings.HasPrefix(cleaned, "```") {
		cleaned = strings.TrimPrefix(cleaned, "```")
		cleaned = strings.TrimSuffix(cleaned, "```")
	}
	cleaned = strings.TrimSpace(cleaned)

	var info CompetitionInfo
	if err := json.Unmarshal([]byte(cleaned), &info); err != nil {
		return nil, fmt.Errorf("解析 JSON 失败（原始响应: %q）: %w", jsonText, err)
	}

	return &info, nil
}
