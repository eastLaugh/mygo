package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/pelletier/go-toml/v2"
)

const (
	mappingURL = "https://raw.githubusercontent.com/eastLaugh/mygo/main/my.toml"
)

type MY map[string]string

// fetch 从 GitHub Raw Content API 获取包映射
func fetch() (MY, error) {
	resp, err := http.Get(mappingURL)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP 状态码: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var mapping MY
	if err := toml.Unmarshal(body, &mapping); err != nil {
		return nil, fmt.Errorf("解析 TOML 失败: %w", err)
	}
	return mapping, nil
}
