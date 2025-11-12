package main

import (
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
		fmt.Println("用法: my go [-n|--no] <package-name> [package-name2] ...")
		fmt.Println("示例: my go gorm")
		fmt.Println("示例: my go gorm gin fiber")
		fmt.Println("示例: my go -n gorm  (只显示命令，不执行)")
		os.Exit(1)
	}

	if os.Args[1] != "go" {
		fmt.Println("用法: my go [-n|--no] <package-name> [package-name2] ...")
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
		fmt.Println("错误: 至少需要一个包名")
		os.Exit(1)
	}

	// 获取映射文件
	my, err := fetch()
	if err != nil {
		fmt.Printf("错误: 无法获取包映射信息: %v\n", err)
		os.Exit(1)
	}

	// 处理每个包名
	var failedPackages []string
	var successCount int

	for _, packageName := range packageNames {
		if packageName == "" {
			fmt.Println("警告: 跳过空的包名")
			continue
		}

		// 查找包路径
		importPath, exists := my[packageName]
		if !exists {
			fmt.Printf("错误: 未找到包 '%s' 的映射\n", packageName)
			failedPackages = append(failedPackages, packageName)
			continue
		}

		// 执行 go get 命令
		fmt.Printf("正在安装: %s -> %s\n", packageName, importPath)
		if err := executeGoGet(importPath, noExecute); err != nil {
			fmt.Printf("错误: 执行 go get 失败 (%s): %v\n", packageName, err)
			failedPackages = append(failedPackages, packageName)
			continue
		}

		if !noExecute {
			fmt.Printf("✓ 成功安装 %s\n", importPath)
		}
		successCount++
	}

	// 显示总结
	if len(failedPackages) > 0 {
		fmt.Printf("\n失败: %d 个包安装失败\n", len(failedPackages))
		if successCount == 0 {
			fmt.Println("可用的包:")
			for name := range my {
				fmt.Printf("  - %s\n", name)
			}
		}
		os.Exit(1)
	}

	if successCount > 0 {
		fmt.Printf("\n✓ 成功安装 %d 个包\n", successCount)
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

func executeGoGet(importPath string, noExecute bool) error {
	fmt.Printf(">>> go get %s\n", importPath)
	if noExecute {
		return nil // 只输出命令，不执行
	}
	cmd := exec.Command("go", "get", importPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
