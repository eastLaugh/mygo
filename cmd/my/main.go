package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

var noExecute bool

func usage() {
	fmt.Println("用法: my go [-n] <package-name>...")
}

func main() {
	if len(os.Args) < 3 {
		usage()
		os.Exit(1)
	}

	if os.Args[1] != "go" {
		usage()
		os.Exit(1)
	}

	// 创建新的 FlagSet，跳过 "go" 子命令
	fs := flag.NewFlagSet("go", flag.ExitOnError)
	fs.BoolVar(&noExecute, "n", false, "只显示命令，不执行")
	fs.BoolVar(&noExecute, "no", false, "只显示命令，不执行")

	// 解析从 os.Args[2:] 开始的参数（跳过 "my" 和 "go"）
	if err := fs.Parse(os.Args[2:]); err != nil {
		usage()
		os.Exit(1)
	}

	// 获取位置参数（包名）
	packageNames := fs.Args()
	if len(packageNames) == 0 {
		fmt.Println("需要包名")
		os.Exit(1)
	}

	my, err := fetch()
	if err != nil {
		fmt.Printf("获取映射失败: %v\n", err)
		os.Exit(1)
	}

	processPackages(my, packageNames)
}

func processPackages(my MY, packageNames []string) {
	var failure []string

	for _, packageName := range packageNames {
		pkg, ok := my[packageName]

		if !ok {
			fmt.Printf("%s: 未找到\n", packageName)
			failure = append(failure, packageName)
			continue
		}

		if err := goGet(pkg); errors.As(err, &IgnoreError{}) {
			println(err.Error())
		} else if err != nil {
			println(err.Error())
			failure = append(failure, packageName)
		}

	}

	print("failure: ")
	for _, name := range failure {
		print(name, " ")
	}
	println()
}
