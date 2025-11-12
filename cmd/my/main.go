package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/pelletier/go-toml/v2"
)

const (
	mappingURL = "https://raw.githubusercontent.com/eastLaugh/mygo/main/my.toml"
)

type MY map[string]string

func main() {
	if len(os.Args) < 3 {
		fmt.Println("用法: my go [-n] <package-name>...")
		os.Exit(1)
	}

	if os.Args[1] != "go" {
		fmt.Println("用法: my go [-n] <package-name>...")
		os.Exit(1)
	}

	// 解析参数，检查是否有 -n 或 --no 标志
	args := os.Args[2:]
	noExecute := false
	packageNames := []string{}

	for _, arg := range args {
		if arg == "-n" || arg == "--no" {
			noExecute = true
		} else {
			packageNames = append(packageNames, arg)
		}
	}

	if len(packageNames) == 0 {
		fmt.Println("需要包名")
		os.Exit(1)
	}

	my, err := fetch()
	if err != nil {
		fmt.Printf("获取映射失败: %v\n", err)
		os.Exit(1)
	}

	// 处理每个包名
	var failedPackages []string
	var successCount int

	for _, packageName := range packageNames {
		if packageName == "" {
			continue
		}

		importPath, exists := my[packageName]
		if !exists {
			fmt.Printf("%s: 未找到\n", packageName)
			failedPackages = append(failedPackages, packageName)
			continue
		}

		if err := goGet(importPath, noExecute); err != nil {
			fmt.Printf("%s: 失败\n", packageName)
			failedPackages = append(failedPackages, packageName)
			continue
		}

		successCount++
	}

	if len(failedPackages) > 0 {
		if successCount == 0 {
			for name := range my {
				fmt.Println(name)
			}
		}
		os.Exit(1)
	}
}

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

func goGet(importPath string, noExecute bool) error {
	fmt.Printf("go get %s\n", importPath)
	if noExecute {
		return errors.New("已忽略")
	}
	cmd := exec.Command("go", "get", importPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
