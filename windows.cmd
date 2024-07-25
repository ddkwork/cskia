				   python tools/git-sync-deps

			set COMMON_ARGS=is_debug=false is_official_build=true skia_enable_discrete_gpu=true skia_enable_fontmgr_android=false skia_enable_fontmgr_empty=false skia_enable_fontmgr_fuchsia=false skia_enable_fontmgr_win_gdi=false skia_enable_gpu=true skia_enable_pdf=true skia_enable_skottie=false skia_enable_skshaper=true skia_enable_skshaper_tests=false skia_enable_spirv_validation=false skia_enable_tools=false skia_enable_vulkan_debug_layers=false skia_use_angle=false skia_use_dawn=false skia_use_dng_sdk=false skia_use_egl=false skia_use_expat=false skia_use_ffmpeg=false skia_use_fixed_gamma_text=false skia_use_fontconfig=false skia_use_gl=true skia_use_harfbuzz=false skia_use_icu=false skia_use_libheif=false skia_use_libjxl_decode=false skia_use_lua=false skia_use_metal=false skia_use_piex=false skia_use_system_libjpeg_turbo=false skia_use_system_libpng=false skia_use_system_libwebp=false skia_use_system_zlib=false skia_use_vulkan=false skia_use_wuffs=true skia_use_xps=false skia_use_zlib=true

			set PLATFORM_ARGS=is_component_build=true skia_enable_fontmgr_win=true skia_use_fonthost_mac=false skia_enable_fontmgr_fontconfig=false skia_use_fontconfig=false skia_use_freetype=false skia_use_x11=false clang_win=\"C:\\Program Files\\LLVM\" extra_cflags=[\"-DSKIA_C_DLL\", \"-UHAVE_NEWLOCALE\", \"-UHAVE_XLOCALE_H\", \"-UHAVE_UNISTD_H\", \"-UHAVE_SYS_MMAN_H\", \"-UHAVE_MMAP\", \"-UHAVE_PTHREAD\"] extra_ldflags=[\"/defaultlib:opengl32\", \"/defaultlib:gdi32\"]

		 gn gen out/Shared --args="%COMMON_ARGS% %PLATFORM_ARGS%"

		 D:\fork\skia\cskia\skia\skia\bin\gn.exe gen out/config --ide=json --json-ide-script=../../gn/gn_to_cmake.py
				   ninja -C out/Shared