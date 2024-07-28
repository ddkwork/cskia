package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream"
)

// cmd := exec.Command("vswhere", "-latest", "-property", "installationVersion")

func main() {
	mylog.ChdirToGithubWorkspace()
	//C:\Users\Admin\Desktop>where python
	//C:\Users\Admin\AppData\Local\Microsoft\WindowsApps\python.exe
	//C:\Users\Admin\AppData\Local\Programs\Python\Python312\python.exe
	//stream.RunCommand("where python")//todo bug

	// 查找 python 的路径
	out, err := exec.Command("where", "python").Output()
	if err != nil {
		fmt.Println("Error finding python:", err)
		return
	}

	index := 0
	pythonPaths := strings.Split(string(out), "\n")
	for i, path := range pythonPaths {
		if path != "" && !strings.Contains(path, "WindowsApps") {
			index = i
		}
		pythonPaths[i] = strings.TrimSpace(path)
	}

	pythonPath := pythonPaths[index]
	pythonPath = filepath.Dir(pythonPath)
	println(pythonPath)
	AppendPathEnvWindows(pythonPath)
	if !isAdmin() {
		runMeElevated()
		return
	}
	// 使用 setx 命令将路径添加到系统的 PATH 环境变量中
	//cmd := exec.Command("setx", "PATH", fmt.Sprintf("%s;%s", "%PATH%", pythonPath), "/M")
	//
	//err = cmd.Run()
	//if err != nil {
	//	fmt.Println("Error setting environment variable:", err)
	//	return
	//}

	out, err = exec.Command("cmd", "/C", "echo %PATH%").Output()
	if err != nil {
		fmt.Println("Error getting PATH:", err)
		return
	}
	currentPath := strings.TrimSpace(string(out))

	// 假设pythonPath是你想要添加的Python路径

	// 将Python路径添加到PATH的最前面
	newPath := fmt.Sprintf("%s;%s", pythonPath, currentPath)

	// 使用setx命令设置新的PATH
	cmd := exec.Command("setx", "PATH", newPath, "/M")
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error setting PATH:", err)
		return
	}

	fmt.Println("Path added successfully to PATH:", pythonPath)

	//mylog.Check(addPathToSystemPath(filepath.Dir(pythonPath)))

	// 使用 setx 命令将路径设置为系统的环境变量
	//cmd := exec.Command("setx", "PYTHON_PATH", pythonPath, "/M")
	//
	//err = cmd.Run()
	//if err != nil {
	//	fmt.Println("Error setting environment variable:", err)
	//	return
	//}
	//
	//fmt.Println("PYTHON_PATH set successfully to:", pythonPath)
	//

	//setx PYTHONPATH "C:\new_path_demo;$env:PYTHONPATH"
	//return

	stream.CopyFile("gn.exe", "c:/windows/gn.exe")
	stream.CopyFile("ninja.exe", "c:/windows/ninja.exe")
	stream.RunCommand("ninja --version")
	stream.RunCommand("python --version")
	stream.RunCommand("gn.exe --version")

	if mylog.IsAction {
		stream.RunCommand("git clone --progress https://chromium.googlesource.com/chromium/tools/depot_tools.git")
		stream.RunCommand("git clone --progress -b chrome/m110 https://github.com/google/skia.git")
	}
	//fixGn() //C:\ProgramData\Chocolatey\bin\vswhere.exe -products * -requires Microsoft.Component.MSBuild -property installationPath -latest
	AppendPathEnvWindows("depot_tools")

	mylog.Info("num cpu", runtime.NumCPU())
	//stream.CopyFile("DEPS_github", "skia/DEPS")

	buffer := stream.NewBuffer("skia\\gn\\toolchain\\BUILD.gn")
	buffer.Replace(`  dlsymutil_pool_depth = exec_script("num_cpus.py", [], "value")`, `  dlsymutil_pool_depth = `+fmt.Sprint(runtime.NumCPU()), 1)
	stream.WriteTruncate("skia\\gn\\toolchain\\BUILD.gn", buffer.Bytes())

	stream.CopyFile("capi/sk_capi.h", "skia/include/sk_capi.h")
	stream.CopyFile("capi/sk_capi.cpp", "skia/src/sk_capi.cpp")

	gni := stream.NewBuffer("skia/gn/core.gni")
	if !gni.Contains("sk_capi.cpp") {
		gni.Replace(`skia_core_sources = [`, `skia_core_sources = [
  "$_src/sk_capi.cpp",`, 1)
		stream.WriteTruncate("skia/gn/core.gni", gni.Bytes())
	}

	font := stream.NewBuffer("skia/src/pdf/SkPDFSubsetFont.h")
	if !font.Contains("include/core/SkData.h") {
		font.Replace(`#include "include/docs/SkPDFDocument.h"`, `#include "include/core/SkData.h"
#include "include/docs/SkPDFDocument.h"`, 1)
		stream.WriteTruncate("skia/src/pdf/SkPDFSubsetFont.h", font.Bytes())
	}
	log.Println("add c api files")

	mylog.Check(os.Chdir("skia"))
	mylog.Info("Chdir to", stream.RunDir())

	if mylog.IsAction {
		stream.RunCommand("python tools/git-sync-deps")
	}
	//stream.RunCommand("python fetch-ninja")

	buildDir := "out/Static"
	mylog.Check(stream.CreatDirectory(buildDir))
	args := fmt.Sprintf("--args=%s %s", COMMON_ARGS, PLATFORM_ARGS)
	cmd = exec.Command("gn", "gen", "out/Static", "--args=", args)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error running gn:", err)
		fmt.Println("gn output:", string(output))
		return
	}
	log.Println("gn gen args success")

	//gn.exe gen out/config --ide=json --json-ide-script=../../gn/gn_to_cmake.py
	//gn.exe gen out/config --ide=vs --json-ide-script=../../gn/gn_meta_sln.py
	//  ide="vs2022"
	//  sln="skia"

	//  #dlsymutil_pool_depth exec_script("num_cpus.py", [], "value")
	//dlsymutil_pool_depth =8
	// todo numberof cpu  test action

	cmd = exec.Command("ninja.exe", "-C", "out/Static", "-v")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Println("ninja run")
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error building Skia:", err)
		fmt.Println("ninja output:", string(output))
		return
	}
	log.Println("ninja build success")

	filepath.Walk("out/Static", func(path string, info fs.FileInfo, err error) error {
		if strings.Contains(path, "obj") {
			return err
		}
		println(path)
		return err
	})

	stream.CopyFile("out/Static/skia.dll", "../skia.dll")
	mylog.Check(os.Chdir(".."))
}

func AppendPathEnvWindows(newPath string) {
	mylog.Check(os.Setenv("PATH", os.Getenv("PATH")+";"+newPath))
}

const (
	COMMON_ARGS = ` 
  win_vc="C:\Program Files\Microsoft Visual Studio\2022\Enterprise\VC"
  clang_win_version=18
  dlsymutil_pool_depth=8
  win_toolchain_version="14.40.33807"
  win_sdk_version="10.0.26100.0"
  is_debug=false 
  is_official_build=true 
  skia_enable_discrete_gpu=true 
  skia_enable_fontmgr_android=false 
  skia_enable_fontmgr_empty=false 
  skia_enable_fontmgr_fuchsia=false 
  skia_enable_fontmgr_win_gdi=false 
  skia_enable_gpu=true 
  skia_enable_pdf=true 
  skia_enable_skottie=false 
  skia_enable_skshaper=true 
  skia_enable_skshaper_tests=false 
  skia_enable_spirv_validation=false 
  skia_enable_tools=false 
  skia_enable_vulkan_debug_layers=false 
  skia_use_angle=false 
  skia_use_dawn=false 
  skia_use_dng_sdk=false 
  skia_use_egl=false 
  skia_use_expat=false 
  skia_use_ffmpeg=false 
  skia_use_fixed_gamma_text=false 
  skia_use_fontconfig=false 
  skia_use_gl=true 
  skia_use_harfbuzz=false 
  skia_use_icu=false 
  skia_use_libheif=false 
  skia_use_libjxl_decode=false 
  skia_use_lua=false 
  skia_use_metal=false 
  skia_use_piex=false 
  skia_use_system_libjpeg_turbo=false 
  skia_use_system_libpng=false 
  skia_use_system_libwebp=false 
  skia_use_system_zlib=false 
  skia_use_vulkan=false 
  skia_use_wuffs=true 
  skia_use_xps=false 
  skia_use_zlib=false 
`
	PLATFORM_ARGS = ` 
is_component_build=true
skia_enable_fontmgr_win=true 
skia_use_fonthost_mac=false 
skia_enable_fontmgr_fontconfig=false 
skia_use_fontconfig=false 
skia_use_freetype=false 
skia_use_x11=false 
clang_win="C:\Program Files\LLVM" 
extra_cflags=[ 
"-DSKIA_C_DLL", 
"-DCRC32_SIMD_SSE42_PCLMUL", 
"-DDEFLATE_FILL_WINDOW_SSE2", 
"-UHAVE_NEWLOCALE", 
"-UHAVE_XLOCALE_H", 
"-UHAVE_UNISTD_H", 
"-UHAVE_SYS_MMAN_H", 
"-UHAVE_MMAP", 
"-UHAVE_PTHREAD",
] 
extra_ldflags=[ 
"/defaultlib:opengl32", 
"/defaultlib:gdi32" 
] 
`
)

func fixGn() {
	log.Println("fix gn")
	//stream.CopyFile("DEPS_github", "skia/skia/DEPS")
	//skia_viewer
	buffer := stream.NewBuffer("skia\\gn\\BUILDCONFIG.gn")
	if !strings.Contains(buffer.String(), "win_vc = \"C:\\Program Files\\Microsoft Visual Studio\\2022\\Enterprise\\VC\"") {
		buffer.Replace(`assert(!(is_debug && is_official_build))`, `#assert(!(is_debug && is_official_build))`, 1) //todo bug

		buffer.Replace(`if (target_os == "win") {`, `
win_sdk = "C:/Program Files (x86)/Windows Kits/10"
win_sdk_version = "10.0.26100.0"

win_vc = "C:\Program Files\Microsoft Visual Studio\2022\Enterprise\VC"
win_toolchain_version = "14.40.33807"

clang_win = "C:\Program Files\LLVM"
clang_win_version = "18.1.7"

if (target_os == "win") {`, 1)
		stream.WriteTruncate("skia\\gn\\BUILDCONFIG.gn", buffer.Bytes())
	}

	buffer = stream.NewBuffer("skia\\gn\\toolchain\\BUILD.gn")
	buffer.Replace(`  dlsymutil_pool_depth = exec_script("num_cpus.py", [], "value")`, `  dlsymutil_pool_depth = 8`, 1)
	stream.WriteTruncate("skia\\gn\\toolchain\\BUILD.gn", buffer.Bytes())
}

//buffer := stream.NewBuffer("skia/skia/DEPS")
//for _, s := range stream.ToLines("skia/skia/DEPS") {
//	break //todo copy for commit id,因为github导入的仓库不知道怎么同步，有点仓库有8000多个分支，烦
//	s = strings.TrimSpace(s)
//	if strings.HasPrefix(s, "#") {
//		continue
//	}
//	if strings.Contains(s, "https://") && strings.Contains(s, "@") {
//		before, after, found := strings.Cut(s, ":")
//		if found {
//			if strings.HasPrefix(strings.TrimSpace(before), "#") {
//				continue
//			}
//			after = strings.TrimSpace(after)
//			after = strings.TrimPrefix(after, `"`)
//			url := strings.Split(after, "@")
//			src := url[0]
//
//			index := strings.LastIndex(src, "/")
//			src = src[:index+1]
//			println(src)
//			buffer.Replace(src, "https://git.homegu.com/ddkwork/", 1)
//		}
//	}
//}
//stream.WriteTruncate("skia/skia/DEPS", buffer.Bytes())

/* working
、编译skia（编译器：LLVM）
（1）编译64位的Lib（x64版本）
首先进入VS 2022的64位命令行编译环境: （x64 Native Tools Command Prompt for VS 2022）
进入skia源码目录：
> cd /d D:\develop\skia
编译skia静态库（Release版，官方版，最小依赖）
.\bin\gn.exe gen out/LLVM.x64.Release --ide="vs2022" --sln="skia" --args="target_cpu=\"x64\" cc=\"clang\" cxx=\"clang++\" clang_win=\"C:/LLVM\" clang_win_version=\"18\" is_trivial_abi=false is_official_build=true skia_use_system_libpng=false skia_use_system_libjpeg_turbo=false skia_use_system_zlib=false skia_use_icu=false skia_use_expat=false skia_use_libwebp_decode=false skia_use_libwebp_encode=false skia_use_xps=false skia_enable_pdf=false is_debug=false"
ninja -C out/LLVM.x64.Release
编译skia静态库（Debug版，非官方版，完整）
.\bin\gn.exe gen out/LLVM.x64.Debug --ide="vs2022" --sln="skia" --args="target_cpu=\"x64\" cc=\"clang\" cxx=\"clang++\" clang_win=\"C:/LLVM\" clang_win_version=\"18\" is_trivial_abi=false is_debug=true extra_cflags=[\"/MTd\"]"
ninja -C out/LLVM.x64.Debug 如果需要重新编译，可用以下命令清理编译的中间文件等文件：
ninja -C out/LLVM.x64.Debug -t clean
*/
