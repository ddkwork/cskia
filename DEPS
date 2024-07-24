use_relative_paths = True

vars = {
  # Three lines of non-changing comments so that
  # the commit queue can handle CLs rolling different
  # dependencies without interference from each other.
  'sk_tool_revision': 'git_revision:d77400c7282ac0939634f1a4c281dec0ec4506c0',

  # ninja CIPD package version.
  # https://chrome-infra-packages.appspot.com/p/infra/3pp/tools/ninja
  'ninja_version': 'version:2@1.12.1.chromium.4',

  # googlefonts_testdata CIPD package version
  # https://chrome-infra-packages.appspot.com/p/chromium/third_party/googlefonts_testdata/
  'googlefonts_testdata_version': 'version:20230913',
}

# If you modify this file, you will need to regenerate the Bazel version of this file (bazel/deps.bzl).
# To do so, run:
#     bazelisk run //bazel/deps_parser
#
# To apply the changes for the GN build, you will need to resync the git repositories using:
#     ./tools/git-sync-deps
deps = {
  "buildtools"                                   : "https://git.homegu.com/ddkwork/buildtools.git@b138e6ce86ae843c42a1a08f37903207bebcca75",
  "third_party/externals/angle2"                 : "https://git.homegu.com/ddkwork/angle.git@b98c4d810794f52a56b4afa39f5029e8508cf117",
  "third_party/externals/brotli"                 : "https://git.homegu.com/ddkwork/brotli.git@6d03dfbedda1615c4cba1211f8d81735575209c8",
  "third_party/externals/d3d12allocator"         : "https://git.homegu.com/ddkwork/D3D12MemoryAllocator.git@169895d529dfce00390a20e69c2f516066fe7a3b",
  # Dawn requires jinja2 and markupsafe for the code generator, tint for SPIRV compilation, and abseil for string formatting.
  # When the Dawn revision is updated these should be updated from the Dawn DEPS as well.
  "third_party/externals/dawn"                   : "https://git.homegu.com/ddkwork/dawn.git@90fdaa810322540281d565e57d3f473d69232e82",
  "third_party/externals/jinja2"                 : "https://git.homegu.com/ddkwork/jinja2@e2d024354e11cc6b041b0cff032d73f0c7e43a07",
  "third_party/externals/markupsafe"             : "https://git.homegu.com/ddkwork/markupsafe@0bad08bb207bbfc1d6f3bbc82b9242b0c50e5794",
  "third_party/externals/abseil-cpp"             : "https://git.homegu.com/ddkwork/abseil-cpp.git@65a55c2ba891f6d2492477707f4a2e327a0b40dc",
  "third_party/externals/dng_sdk"                : "https://git.homegu.com/ddkwork/dng_sdk.git@c8d0c9b1d16bfda56f15165d39e0ffa360a11123",
  "third_party/externals/egl-registry"           : "https://git.homegu.com/ddkwork/EGL-Registry@b055c9b483e70ecd57b3cf7204db21f5a06f9ffe",
  "third_party/externals/emsdk"                  : "https://git.homegu.com/ddkwork/emsdk.git@a896e3d066448b3530dbcaa48869fafefd738f57",
  "third_party/externals/expat"                  : "https://git.homegu.com/ddkwork/libexpat.git@441f98d02deafd9b090aea568282b28f66a50e36",
  "third_party/externals/freetype"               : "https://git.homegu.com/ddkwork/freetype2.git@73720c7c9958e87b3d134a7574d1720ad2d24442",
  "third_party/externals/harfbuzz"               : "https://git.homegu.com/ddkwork/harfbuzz.git@b74a7ecc93e283d059df51ee4f46961a782bcdb8",
  "third_party/externals/highway"                : "https://git.homegu.com/ddkwork/highway.git@424360251cdcfc314cfc528f53c872ecd63af0f0",
  "third_party/externals/icu"                    : "https://git.homegu.com/ddkwork/icu.git@364118a1d9da24bb5b770ac3d762ac144d6da5a4",
  "third_party/externals/icu4x"                  : "https://git.homegu.com/ddkwork/icu4x.git@bcf4f7198d4dc5f3127e84a6ca657c88e7d07a13",
  "third_party/externals/imgui"                  : "https://git.homegu.com/ddkwork/imgui.git@55d35d8387c15bf0cfd71861df67af8cfbda7456",
  "third_party/externals/libavif"                : "https://git.homegu.com/ddkwork/libavif.git@55aab4ac0607ab651055d354d64c4615cf3d8000",
  "third_party/externals/libgav1"                : "https://git.homegu.com/ddkwork/libgav1.git@5cf722e659014ebaf2f573a6dd935116d36eadf1",
  "third_party/externals/libgrapheme"            : "https://git.homegu.com/ddkwork/libgrapheme/@c0cab63c5300fa12284194fbef57aa2ed62a94c0",
  "third_party/externals/libjpeg-turbo"          : "https://git.homegu.com/ddkwork/libjpeg_turbo.git@ccfbe1c82a3b6dbe8647ceb36a3f9ee711fba3cf",
  "third_party/externals/libjxl"                 : "https://git.homegu.com/ddkwork/jpeg-xl.git@a205468bc5d3a353fb15dae2398a101dff52f2d3",
  "third_party/externals/libpng"                 : "https://git.homegu.com/ddkwork/libpng.git@ed217e3e601d8e462f7fd1e04bed43ac42212429",
  "third_party/externals/libwebp"                : "https://git.homegu.com/ddkwork/libwebp.git@845d5476a866141ba35ac133f856fa62f0b7445f",
  "third_party/externals/libyuv"                 : "https://git.homegu.com/ddkwork/libyuv.git@d248929c059ff7629a85333699717d7a677d8d96",
  "third_party/externals/microhttpd"             : "https://git.homegu.com/ddkwork/libmicrohttpd@748945ec6f1c67b7efc934ab0808e1d32f2fb98d",
  "third_party/externals/oboe"                   : "https://git.homegu.com/ddkwork/oboe.git@b02a12d1dd821118763debec6b83d00a8a0ee419",
  "third_party/externals/opengl-registry"        : "https://git.homegu.com/ddkwork/OpenGL-Registry@14b80ebeab022b2c78f84a573f01028c96075553",
  "third_party/externals/perfetto"               : "https://git.homegu.com/ddkwork/perfetto@93885509be1c9240bc55fa515ceb34811e54a394",
  "third_party/externals/piex"                   : "https://git.homegu.com/ddkwork/piex.git@bb217acdca1cc0c16b704669dd6f91a1b509c406",
  "third_party/externals/swiftshader"            : "https://git.homegu.com/ddkwork/SwiftShader@c4dfa69de7deecf52c6b53badbc8bb7be1a05e8c",
  "third_party/externals/vulkanmemoryallocator"  : "https://git.homegu.com/ddkwork/VulkanMemoryAllocator@a6bfc237255a6bac1513f7c1ebde6d8aed6b5191",
  # vulkan-deps is a meta-repo containing several interdependent Khronos Vulkan repositories.
  # When the vulkan-deps revision is updated, those repos (spirv-*, vulkan-*) should be updated as well.
  "third_party/externals/vulkan-deps"            : "https://git.homegu.com/ddkwork/vulkan-deps@8e90204125ac61dcb01f24bf8b221d5b388846a5",
  "third_party/externals/spirv-cross"            : "https://git.homegu.com/ddkwork/SPIRV-Cross@b8fcf307f1f347089e3c46eb4451d27f32ebc8d3",
  "third_party/externals/spirv-headers"          : "https://git.homegu.com/ddkwork/SPIRV-Headers.git@db5a00f8cebe81146cafabf89019674a3c4bf03d",
  "third_party/externals/spirv-tools"            : "https://git.homegu.com/ddkwork/SPIRV-Tools.git@a0817526b8e391732632e6a887134be256a20a18",
  "third_party/externals/vello"                  : "https://git.homegu.com/ddkwork/vello.git@3ee3bea02164c5a816fe6c16ef4e3a810edb7620",
  "third_party/externals/vulkan-headers"         : "https://git.homegu.com/ddkwork/Vulkan-Headers@fabe9e2672334fdb9a622d42a2e8f94578952082",
  "third_party/externals/vulkan-tools"           : "https://git.homegu.com/ddkwork/Vulkan-Tools@46df205dcad665b652f57ee580d78051925b296a",
  "third_party/externals/vulkan-utility-libraries": "https://git.homegu.com/ddkwork/Vulkan-Utility-Libraries@67522b34edde86dbb97e164280291f387ade55fc",
  "third_party/externals/unicodetools"           : "https://git.homegu.com/ddkwork/unicodetools@66a3fa9dbdca3b67053a483d130564eabc5fe095",
  #"third_party/externals/v8"                     : "https://git.homegu.com/ddkwork.git@5f1ae66d5634e43563b2d25ea652dfb94c31a3b4",
  "third_party/externals/wuffs"                  : "https://git.homegu.com/ddkwork.git@e3f919ccfe3ef542cfc983a82146070258fb57f8",
  "third_party/externals/zlib"                   : "https://git.homegu.com/ddkwork/zlib@646b7f569718921d7d4b5b8e22572ff6c76f2596",

  'bin': {
    'packages': [
      {
        'package': 'skia/tools/sk/${{platform}}',
        'version': Var('sk_tool_revision'),
      },
      {
        'package': 'infra/3pp/tools/ninja/${{platform}}',
        'version': Var('ninja_version'),
      }
    ],
    'dep_type': 'cipd',
  },
}
