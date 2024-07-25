package main

import (
	"archive/zip"
	"fmt"
	"github.com/ddkwork/golibrary/stream"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
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
	err := downloadFile(url, zipFile)
	if err != nil {
		return err
	}

	// 解压文件
	err = unzip(zipFile, destDir)
	if err != nil {
		return err
	}

	return nil
}

func downloadFile(url, filepath string) error {
	// 创建目标文件
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// 获取数据
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 将响应体写入文件
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func unzip(zipFile, destDir string) error {
	// 打开ZIP文件
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer r.Close()

	// 创建目标目录
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	// 遍历ZIP文件中的每个文件/目录
	for _, f := range r.File {
		err := unzipFile(f, destDir)
		if err != nil {
			return err
		}
	}

	return nil
}

func unzipFile(f *zip.File, destDir string) error {
	// 打开ZIP文件中的文件
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	// 构建解压后的文件路径
	path := filepath.Join(destDir, f.Name)

	// 如果是目录，创建目录
	if f.FileInfo().IsDir() {
		os.MkdirAll(path, f.Mode())
		return nil
	}

	// 创建解压后的文件
	outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer outFile.Close()

	// 将ZIP文件中的内容复制到解压后的文件
	_, err = io.Copy(outFile, rc)
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
	stream.CopyFile(src, dest)
	return nil
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
