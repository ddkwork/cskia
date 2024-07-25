package main

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream"
)

func AddToPath(newPath string) {
	mylog.Check(os.Setenv("PATH", os.Getenv("PATH")+";skia/depot_tools"))
}

func main() {
	mylog.Check(stream.CreatDirectory("skia"))
	stream.RunCommand("git clone --progress https://chromium.googlesource.com/chromium/tools/depot_tools.git")

	AddToPath("skia/depot_tools")

	stream.RunCommand("git clone --progress -b chrome/m128 https://github.com/google/skia.git")
	mylog.Check(os.Chdir("skia"))
	stream.RunCommand("python3 tools/git-sync-deps")
	stream.RunCommand("python3 bin/fetch-ninja")
	mylog.Check(os.Chdir(".."))

	stream.CopyFile("capi/sk_capi.h", "skia/include/sk_capi.h")
	stream.CopyFile("capi/sk_capi.cpp", "skia/src/sk_capi.cpp")

	gni := stream.NewBuffer("skia/gn/core.gni")
	gni.Replace(`skia_core_sources = [`, `skia_core_sources = [
  "$_src/sk_capi.cpp",`, 1)
	stream.WriteTruncate("skia/gn/core.gni", gni.Bytes())

	font := stream.NewBuffer("skia/src/pdf/SkPDFSubsetFont.h")
	font.Replace(`#include "include/docs/SkPDFDocument.h"`, `#include "include/core/SkData.h"
#include "include/docs/SkPDFDocument.h"`, 1)
	stream.WriteTruncate("skia/src/pdf/SkPDFSubsetFont.h", font.Bytes())

	mylog.Check(os.Chdir("skia"))
	genCmd := "${" + COMMON_ARGS + "} ${" + PLATFORM_ARGS + "}"

	buildDir := filepath.Join("skia", "build")
	mylog.Check(stream.CreatDirectory(buildDir))
	stream.RunCommand("python3 tools/git-sync-deps")
	stream.RunCommand("python3 bin/fetch-gn")
	stream.RunCommandArgs("bin/gen", strconv.Quote(buildDir), "--args=", strconv.Quote(genCmd))
	mylog.Check(stream.CreatDirectory(buildDir))
	stream.RunCommandArgs("ninja -C -v", buildDir)

	stream.CopyFile("skia/build/skia.dll", "./skia.dll")
	return

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

	stream.CopyFile("DEPS_github", "skia/skia/DEPS")

	buffer := stream.NewBuffer("skia\\skia\\gn\\BUILDCONFIG.gn")
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

	/*
				   python tools/git-sync-deps

			set COMMON_ARGS=is_debug=false is_official_build=true skia_enable_discrete_gpu=true skia_enable_fontmgr_android=false skia_enable_fontmgr_empty=false skia_enable_fontmgr_fuchsia=false skia_enable_fontmgr_win_gdi=false skia_enable_gpu=true skia_enable_pdf=true skia_enable_skottie=false skia_enable_skshaper=true skia_enable_skshaper_tests=false skia_enable_spirv_validation=false skia_enable_tools=false skia_enable_vulkan_debug_layers=false skia_use_angle=false skia_use_dawn=false skia_use_dng_sdk=false skia_use_egl=false skia_use_expat=false skia_use_ffmpeg=false skia_use_fixed_gamma_text=false skia_use_fontconfig=false skia_use_gl=true skia_use_harfbuzz=false skia_use_icu=false skia_use_libheif=false skia_use_libjxl_decode=false skia_use_lua=false skia_use_metal=false skia_use_piex=false skia_use_system_libjpeg_turbo=false skia_use_system_libpng=false skia_use_system_libwebp=false skia_use_system_zlib=false skia_use_vulkan=false skia_use_wuffs=true skia_use_xps=false skia_use_zlib=true

			set PLATFORM_ARGS=is_component_build=true skia_enable_fontmgr_win=true skia_use_fonthost_mac=false skia_enable_fontmgr_fontconfig=false skia_use_fontconfig=false skia_use_freetype=false skia_use_x11=false clang_win=\"C:\\Program Files\\LLVM\" extra_cflags=[\"-DSKIA_C_DLL\", \"-UHAVE_NEWLOCALE\", \"-UHAVE_XLOCALE_H\", \"-UHAVE_UNISTD_H\", \"-UHAVE_SYS_MMAN_H\", \"-UHAVE_MMAP\", \"-UHAVE_PTHREAD\"] extra_ldflags=[\"/defaultlib:opengl32\", \"/defaultlib:gdi32\"]

		 gn gen out/Shared --args="%COMMON_ARGS% %PLATFORM_ARGS%"

		 D:\fork\skia\cskia\skia\skia\bin\gn.exe gen out/config --ide=json --json-ide-script=../../gn/gn_to_cmake.py
				   ninja -C out/Shared

	*/
}

const (
	COMMON_ARGS = ` \
is_debug=false \
is_official_build=true \
skia_enable_discrete_gpu=true \
skia_enable_fontmgr_android=false \
skia_enable_fontmgr_empty=false \
skia_enable_fontmgr_fuchsia=false \
skia_enable_fontmgr_win_gdi=false \
skia_enable_gpu=true \
skia_enable_pdf=true \
skia_enable_skottie=false \
skia_enable_skshaper=true \
skia_enable_skshaper_tests=false \
skia_enable_spirv_validation=false \
skia_enable_tools=false \
skia_enable_vulkan_debug_layers=false \
skia_use_angle=false \
skia_use_dawn=false \
skia_use_dng_sdk=false \
skia_use_egl=false \
skia_use_expat=false \
skia_use_ffmpeg=false \
skia_use_fixed_gamma_text=false \
skia_use_fontconfig=false \
skia_use_gl=true \
skia_use_harfbuzz=false \
skia_use_icu=false \
skia_use_libheif=false \
skia_use_libjxl_decode=false \
skia_use_lua=false \
skia_use_metal=false \
skia_use_piex=false \
skia_use_system_libjpeg_turbo=false \
skia_use_system_libpng=false \
skia_use_system_libwebp=false \
skia_use_system_zlib=false \
skia_use_vulkan=false \
skia_use_wuffs=true \
skia_use_xps=false \
skia_use_zlib=true \
`
	PLATFORM_ARGS = ` \
is_component_build=true \
skia_enable_fontmgr_win=true \
skia_use_fonthost_mac=false \
skia_enable_fontmgr_fontconfig=false \
skia_use_fontconfig=false \
skia_use_freetype=false \
skia_use_x11=false \
clang_win=\"C:\\Program Files\\LLVM\" \
extra_cflags=[ \
\"-DSKIA_C_DLL\", \
\"-UHAVE_NEWLOCALE\", \
\"-UHAVE_XLOCALE_H\", \
\"-UHAVE_UNISTD_H\", \
\"-UHAVE_SYS_MMAN_H\", \
\"-UHAVE_MMAP\", \
\"-UHAVE_PTHREAD\" \
] \
extra_ldflags=[ \
\"/defaultlib:opengl32\", \
\"/defaultlib:gdi32\" \
] \
`
)
