@echo off
setlocal enabledelayedexpansion

REM These two variables should be set in tandem to keep a consistent set of sources.
REM Last set Sat Jun 15 11:00:00 PDT 2024
set DEPOT_TOOLS_COMMIT=d431af39c998f4ba88bf376f500214d82ece766a
set SKIA_BRANCH=main

for %%a in (%*) do (
    set arg=%%a
    if "!arg!"=="--args" (set SHOW_ARGS=1)
    if "!arg!"=="--clean" (set CLEAN=restore)
    if "!arg!"=="--CLEAN" (set CLEAN=full)
    if "!arg!"=="--help" (
        echo %~n0 [options]
        echo   -a, --args  Display the available args list for the skia build (no build)
        echo   -c, --clean Remove the dist and skia/build directories (no build)
        echo   -C, --CLEAN Remove the dist and skia directories (no build)
        echo   -h, --help This help text
        exit /b 0
    )
    if "!arg!" neq "--args" if "!arg!" neq "--clean" if "!arg!" neq "--CLEAN" if "!arg!" neq "--help" (
        echo Invalid argument: !arg!
        exit /b 1
    )
)

if "!CLEAN!"=="full" (
    rmdir /s /q dist
    rmdir /s /q skia
    exit /b 0
)

if "!CLEAN!"=="restore" (
    rmdir /s /q dist
    rmdir /s /q skia\build
    if exist skia\skia (
        cd skia\skia
        git checkout -- .
        del include\sk_capi.h src\sk_capi.cpp
        cd ..\..
    )
    exit /b 0
)

if "!SHOW_ARGS!"=="1" (
    set PATH=%CD%\skia\depot_tools;%PATH%
    cd skia\skia
    bin\gn args ..\build --list --short
    cd ..\..
    exit /b 0
)

set BUILD_DIR=%CD%\skia\build
set DIST=%CD%\dist

REM As changes to Skia are made, these args may need to be adjusted.
REM Use 'bin/gn args %BUILD_DIR% --list' to see what args are available.
set COMMON_ARGS=is_debug=false is_official_build=true skia_enable_discrete_gpu=true skia_enable_fontmgr_android=false skia_enable_fontmgr_empty=false skia_enable_fontmgr_fuchsia=false skia_enable_fontmgr_win_gdi=false skia_enable_gpu=true skia_enable_pdf=true skia_enable_skottie=false skia_enable_skshaper=true skia_enable_skshaper_tests=false skia_enable_spirv_validation=false skia_enable_tools=false skia_enable_vulkan_debug_layers=false skia_use_angle=false skia_use_dawn=false skia_use_dng_sdk=false skia_use_egl=false skia_use_expat=false skia_use_ffmpeg=false skia_use_fixed_gamma_text=false skia_use_fontconfig=false skia_use_gl=true skia_use_harfbuzz=false skia_use_icu=false skia_use_libheif=false skia_use_libjxl_decode=false skia_use_lua=false skia_use_metal=false skia_use_piex=false skia_use_system_libjpeg_turbo=false skia_use_system_libpng=false skia_use_system_libwebp=false skia_use_system_zlib=false skia_use_vulkan=false skia_use_wuffs=true skia_use_xps=false skia_use_zlib=true

set OS_TYPE=windows
set LIB_NAME=skia.dll
set UNISON_LIB_NAME=skia_windows.dll
set PLATFORM_ARGS=is_component_build=true skia_enable_fontmgr_win=true skia_use_fonthost_mac=false skia_enable_fontmgr_fontconfig=false skia_use_fontconfig=false skia_use_freetype=false skia_use_x11=false clang_win="C:\Program Files\LLVM" extra_cflags=["-DSKIA_C_DLL","-UHAVE_NEWLOCALE","-UHAVE_XLOCALE_H","-UHAVE_UNISTD_H","-UHAVE_SYS_MMAN_H","-UHAVE_MMAP","-UHAVE_PTHREAD"] extra_ldflags=["/defaultlib:opengl32","/defaultlib:gdi32"]

REM Setup the Skia tree, pulling sources, if needed.
mkdir skia
cd skia

if not exist depot_tools (
    git clone https://chromium.googlesource.com/chromium/tools/depot_tools.git
    cd depot_tools
    git reset --hard %DEPOT_TOOLS_COMMIT%
    cd ..
)
set PATH=%CD%\depot_tools;%PATH%

if not exist skia (
    git clone https://github.com/google/skia.git
    cd skia
    git checkout %SKIA_BRANCH%
    python3 tools\git-sync-deps
    python3 bin\fetch-ninja
    cd ..
)

REM Apply our changes.
cd skia
rmdir /s /q src\c include\c
copy ..\..\capi\sk_capi.h include\
copy ..\..\capi\sk_capi.cpp src\
findstr /v "src\sk_capi.cpp" gn\core.gni | sed -e "s@skia_core_sources = \[@&\n  \"$_src/sk_capi.cpp\",@" > gn\core.gni.new
move gn\core.gni.new gn\core.gni
sed -e "s@^class SkData;$@#include \"include/core/SkData.h\"@" src\pdf\SkPDFSubsetFont.h > src\pdf\SkPDFSubsetFont.h.new
move src\pdf\SkPDFSubsetFont.h.new src\pdf\SkPDFSubsetFont.h

REM Perform the build
bin\gn gen %BUILD_DIR% --args="%COMMON_ARGS% %PLATFORM_ARGS%"
ninja -C %BUILD_DIR%

REM Copy the result into %DIST%
mkdir %DIST%\include
del %DIST%\include\*.h
copy include\sk_capi.h %DIST%\include\
mkdir %DIST%\lib\%OS_TYPE%
copy %BUILD_DIR%\%LIB_NAME% %DIST%\lib\%OS_TYPE%\

cd ..\..

REM If present, also copy the results into the unison build tree
if exist ..\unison (
    set RELATIVE_UNISON_DIR=..\unison\internal\skia
    mkdir %RELATIVE_UNISON_DIR%
    copy %DIST%\include\sk_capi.h %RELATIVE_UNISON_DIR%\
    copy %DIST%\lib\%OS_TYPE%\%LIB_NAME% %RELATIVE_UNISON_DIR%\%UNISON_LIB_NAME%
    echo Copied distribution to unison
)
