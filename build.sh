#! /usr/bin/env bash

SKIA_BRANCH=main

if [ "$SHOW_ARGS"x == "1x" ]; then
	export PATH="${PWD}/skia/depot_tools:${PATH}"
	cd skia/skia
	bin/gn args ../build --list --short
	exit 0
fi

BUILD_DIR=${PWD}/skia/build
DIST=${PWD}/dist

# Setup the Skia tree, pulling sources, if needed.
mkdir -p skia
cd skia

if [ ! -e depot_tools ]; then
	git clone https://chromium.googlesource.com/chromium/tools/depot_tools.git
	cd depot_tools
	cd ..
fi
export PATH="${PWD}/depot_tools:${PATH}"

if [ ! -e skia ]; then
	git clone https://github.com/google/skia.git
	cd skia
	git checkout "${SKIA_BRANCH}"
	python3 tools/git-sync-deps
	python3 bin/fetch-ninja
	cd ..
fi

# Apply our changes.
cd skia
/bin/rm -rf src/c include/c
cp ../../capi/sk_capi.h include/
cp ../../capi/sk_capi.cpp src/
grep -v src/sk_capi.cpp gn/core.gni | sed -e 's@skia_core_sources = \[@&\
  "$_src/sk_capi.cpp",@' >gn/core.gni.new
/bin/mv gn/core.gni.new gn/core.gni
sed -e 's@^class SkData;$@#include "include/core/SkData.h"@' src/pdf/SkPDFSubsetFont.h >src/pdf/SkPDFSubsetFont.h.new
/bin/mv src/pdf/SkPDFSubsetFont.h.new src/pdf/SkPDFSubsetFont.h

# Perform the build
bin/gn gen "${BUILD_DIR}" --args="${COMMON_ARGS} ${PLATFORM_ARGS}"
ninja -C "${BUILD_DIR}"

# Copy the result into ${DIST}
mkdir -p "${DIST}/include"
/bin/rm -f ${DIST}/include/*.h
cp include/sk_capi.h "${DIST}/include/"
mkdir -p "${DIST}/lib/${OS_TYPE}"
cp "${BUILD_DIR}/${LIB_NAME}" "${DIST}/lib/${OS_TYPE}/"

cd ../..

