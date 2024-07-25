package main

import (
	"fmt"
	"os"
	"os/exec"
)

func InstallNinjaAndGn() {
	// 下载并解压gn
	err := downloadAndUnzip("https://chrome-infra-packages.appspot.com/dl/gn/gn/windows-amd64/+/latest", "gn.zip", "c:///gn")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 设置环境变量
	err = setEnvPath("c://gn")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 复制gn.exe到C:/Windows/System32/
	err = copyFile("c://gn.exe", "C:/Windows/System32/gn.exe")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 下载并解压ninja
	err = downloadAndUnzip("https://github.com/ninja-build/ninja/releases/download/v1.10.2/ninja-win.zip", "ninja.zip", "c://ninja")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 设置环境变量
	err = setEnvPath("c://ninja")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 复制ninja.exe到C:/Windows/System32/
	err = copyFile("c://ninja.exe", "C:/Windows/System32/ninja.exe")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func downloadAndUnzip(url, zipFile, destDir string) error {
	// 下载文件
	err := exec.Command("curl", "-L", url, "-o", zipFile).Run()
	if err != nil {
		return err
	}

	// 解压文件
	err = exec.Command("unzip", zipFile, "-d", destDir).Run()
	if err != nil {
		return err
	}

	return nil
}

func setEnvPath(dir string) error {
	// 获取当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	// 设置环境变量
	path := os.Getenv("PATH")
	newPath := fmt.Sprintf("%s;%s/%s", path, wd, dir)
	return exec.Command("setx", "/M", "PATH", newPath).Run()
}

func copyFile(src, dest string) error {
	// 复制文件
	return exec.Command("cp", src, dest).Run()
}

/*
name: Build and Publish Windows Skia DLL with GN and Ninja

on:
  push:
    branches:
      - main
  pull_request:
    branches:
      - main

jobs:
  build:
    runs-on: windows-latest

    steps:
      - name: Checkout Skia repository
        uses: actions/checkout@v2
        with:
          repository: google/skia  # Skia仓库地址
          path: skia
          ref: main  # 指定分支

      - name: Setup MSBuild
        uses: microsoft/setup-msbuild@v1

      - name: Setup Python
        uses: actions/setup-python@v2
        with:
          python-version: '3.x'

      - name: Install GN
        run: |
          cd skia
          python tools/git-sync-deps
          mkdir -p third_party/externals/gn
          curl -L https://chrome-infra-packages.appspot.com/dl/gn/gn/windows-amd64/+/latest -o gn.zip
          unzip gn.zip -d third_party/externals/gn
          setx /M PATH "%PATH%;$(pwd)/third_party/externals/gn"
          cp third_party/externals/gn/gn.exe C:/Windows/System32/

      - name: Install Ninja
        run: |
          cd skia
          mkdir -p third_party/externals/ninja
          curl -L https://github.com/ninja-build/ninja/releases/download/v1.10.2/ninja-win.zip -o ninja.zip
          unzip ninja.zip -d third_party/externals/ninja
          setx /M PATH "%PATH%;$(pwd)/third_party/externals/ninja"
          cp third_party/externals/ninja/ninja.exe C:/Windows/System32/
*/
