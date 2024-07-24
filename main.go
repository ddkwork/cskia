package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream"
)

func main() {
	buffer := stream.NewBuffer("skia/skia/DEPS")
	for _, s := range stream.ToLines("skia/skia/DEPS") {
		s = strings.TrimSpace(s)
		if strings.HasPrefix(s, "#") {
			continue
		}
		if strings.Contains(s, "https://") && strings.Contains(s, "@") {
			before, after, found := strings.Cut(s, ":")
			if found {
				if strings.HasPrefix(strings.TrimSpace(before), "#") {
					continue
				}
				after = strings.TrimSpace(after)
				after = strings.TrimPrefix(after, `"`)
				url := strings.Split(after, "@")
				src := url[0]

				index := strings.LastIndex(src, "/")
				src = src[:index+1]
				println(src)
				buffer.Replace(src, "https://git.homegu.com/ddkwork/", 1)
			}
		}
	}
	stream.WriteTruncate("skia/skia/DEPS", buffer.Bytes())

	buffer = stream.NewBuffer("skia\\skia\\gn\\BUILDCONFIG.gn")
	if !strings.Contains(buffer.String(), "win_vc = \"C:\\Program Files\\Microsoft Visual Studio\\2022\\Enterprise\\VC\"") {
		buffer.Replace(`if (target_os == "win") {`, `
win_sdk = "C:/Program Files (x86)/Windows Kits/10"
win_sdk_version = "10.0.26100.0"

win_vc = "C:\Program Files\Microsoft Visual Studio\2022\Enterprise\VC"
win_toolchain_version = "14.40.33807"

clang_win = "C:\Program Files\LLVM"
clang_win_version = "18.1.7"

if (target_os == "win") {`, 1)
		stream.WriteTruncate("skia\\skia\\gn\\BUILDCONFIG.gn", buffer.Bytes())
	}

	buffer = stream.NewBuffer("skia\\skia\\gn\\toolchain\\BUILD.gn")
	buffer.Replace(`  dlsymutil_pool_depth = exec_script("num_cpus.py", [], "value")`, `  dlsymutil_pool_depth = 8`, 1)
	stream.WriteTruncate("skia\\skia\\gn\\toolchain\\BUILD.gn", buffer.Bytes())

	//python tools/git-sync-deps
	//gn gen out/Shared --args="is_debug=false is_official_build=true skia_use_system_libjpeg_turbo=false skia_use_system_libpng=false skia_use_system_zlib=false skia_enable_tools=false"
	//bin/gn gen out/config --ide=json --json-ide-script=../../gn/gn_to_cmake.py
	//ninja -C out/Shared
}

//go:generate go build -o build_skia
func main2() {
	workDir := "skia"
	mylog.CheckIgnore(os.MkdirAll(workDir, 0755))
	depotToolsDir := filepath.Join(workDir, "depot_tools")
	if !stream.IsDirEx(depotToolsDir) {
		mylog.Check(exec.Command("git", "clone", "--progress", "https://chromium.googlesource.com/chromium/tools/depot_tools.git", depotToolsDir).Run())
	}
	mylog.Check(os.Setenv("PATH", fmt.Sprintf("%s;%s", depotToolsDir, os.Getenv("PATH"))))

	// 同步 Skia 依赖
	skiaDir := filepath.Join(workDir, "skia")
	if !stream.IsDirEx(skiaDir) {
		mylog.Check(exec.Command("git", "clone", "--progress", "https://github.com/google/skia.git", skiaDir).Run())
	}
	mylog.Check(os.Chdir(skiaDir))
	mylog.Check(exec.Command("python3", "tools/git-sync-deps").Run())
	mylog.Check(exec.Command("python3", "bin/fetch-ninja").Run())
	mylog.Check(os.RemoveAll("src/c"))
	mylog.Check(os.RemoveAll("include/c"))
	mylog.Check(os.Chdir(skiaDir))
	stream.CopyFile("../../capi/sk_capi.h", "include/sk_capi.h")
	stream.CopyFile("../../capi/sk_capi.cpp", "src/sk_capi.cpp")

	modifyCoreGni()
	modifySkPDFSubsetFontH()

	// 构建 Skia DLL
	buildDir := filepath.Join(skiaDir, "out", "Release")
	mylog.Check(exec.Command("bin/gn", "gen", buildDir, "--args=is_debug=false is_official_build=true skia_enable_gpu=true is_component_build=true").Run())
	mylog.Check(exec.Command("ninja", "-C", buildDir, "skia").Run())

	// 检查构建目录
	files := mylog.Check2(os.ReadDir(buildDir))

	for _, file := range files {
		fmt.Println("Build file:", file.Name())
	}

	distDir := "dist"
	mylog.Check(os.MkdirAll(distDir, 0755))
	mylog.Check(exec.Command("cp", filepath.Join(buildDir, "skia.dll"), filepath.Join(distDir, "skia.dll")).Run())

	mylog.Success("", "DLL published successfully!")
}

func copyFile(src, dst string) {
	input := mylog.Check2(os.ReadFile(src))
	mylog.Check(os.WriteFile(dst, input, 0644))
}

func modifyCoreGni() {
	filePath := "gn/core.gni"
	input := mylog.Check2(os.ReadFile(filePath))
	lines := string(input)
	re := regexp.MustCompile(`src/sk_capi.cpp`)
	lines = re.ReplaceAllString(lines, "")
	re = regexp.MustCompile(`skia_core_sources = \[`)
	lines = re.ReplaceAllString(lines, `skia_core_sources = [\n  "$_src/sk_capi.cpp",`)
	mylog.Check(os.WriteFile("gn/core.gni.new", []byte(lines), 0644))
	mylog.Check(os.Rename("gn/core.gni.new", "gn/core.gni"))
}

func modifySkPDFSubsetFontH() {
	filePath := "src/pdf/SkPDFSubsetFont.h"
	input := mylog.Check2(os.ReadFile(filePath))
	lines := string(input)
	re := regexp.MustCompile(`^class SkData;$`)
	lines = re.ReplaceAllString(lines, `#include "include/core/SkData.h"`)
	mylog.Check(os.WriteFile("src/pdf/SkPDFSubsetFont.h.new", []byte(lines), 0644))
	mylog.Check(os.Rename("src/pdf/SkPDFSubsetFont.h.new", "src/pdf/SkPDFSubsetFont.h"))
}
