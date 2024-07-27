package main

import (
	"fmt"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
	"os"
	"strings"
	"syscall"
)

// 检查是否以管理员权限运行
func isAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	return err == nil
}

// 请求管理员权限
func runMeElevated() {
	verb := "runas"
	exe, _ := os.Executable()
	cwd, _ := os.Getwd()
	args := strings.Join(os.Args[1:], " ")

	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString(exe)
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	argPtr, _ := syscall.UTF16PtrFromString(args)

	var showCmd int32 = 1 // 通常窗口

	err := windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd)
	if err != nil {
		fmt.Println(err)
	}
}

// 添加路径到系统的 PATH 环境变量
func addPathToSystemPath(newPath string) error {
	// 打开系统的 PATH 环境变量
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer key.Close()

	// 读取当前的 PATH 值
	currentPath, _, err := key.GetStringValue("PATH")
	if err != nil {
		return err
	}

	// 检查新路径是否已经存在
	if !strings.Contains(currentPath, newPath) {
		// 更新 PATH 值
		newPathValue := currentPath + ";" + newPath
		err = key.SetStringValue("PATH", newPathValue)
		if err != nil {
			return err
		}
	}

	return nil
}
