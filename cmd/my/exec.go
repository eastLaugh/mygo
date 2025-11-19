package main

import (
	"fmt"
	"os"
	"os/exec"
)

func goGet(pkg string) error {
	fmt.Printf("go get %s\n", pkg)
	if noExecute {
		return IgnoreError{pkg: pkg}
	}
	cmd := exec.Command("go", "get", pkg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

type IgnoreError struct {
	pkg string
}

func (e IgnoreError) Error() string {
	return fmt.Sprintf("已忽略安装 %s", e.pkg)
}
