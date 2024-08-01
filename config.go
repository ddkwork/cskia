package main

import (
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"
)

type config struct {
	win_vc                string
	win_toolchain_version string
	win_sdk_version       string
	clang_win_version     string
	clang_win             string
}

func (c config) Win_sdk_version() string {
	return "win_sdk_version=" + c.win_sdk_version
}

func (c config) Win_vc() string {
	return "win_vc=" + c.win_vc
}

func (c config) Win_toolchain_version() string {
	return "win_toolchain_version=" + c.win_toolchain_version
}

func (c config) Clang_win_version() string {
	return "clang_win_version=" + c.clang_win_version
}

func (c config) Clang_win() string {
	return "clang_win=" + c.clang_win
}

func (c config) Valid() {
	if c.win_vc == "" {
		panic("win_vc is empty")
	}
	if c.win_toolchain_version == "" {
		panic("win_toolchain_version is empty")
	}
	if c.win_sdk_version == "" {
		panic("win_sdk_version is empty")
	}
	if c.clang_win_version == "" {
		panic("clang_win_version is empty")
	}
	if c.clang_win == "" {
		panic("clang_win is empty")
	}
	join := strings.Join([]string{
		c.Win_vc(),
		c.Win_toolchain_version(),
		c.Win_sdk_version(),
		c.Clang_win_version(),
		c.Clang_win(),
	}, "\n")
	mylog.Json("config", join)
}

func newConfig() *config {
	path := "C:\\Program Files\\Microsoft Visual Studio\\2022\\Enterprise\\VC\\Tools\\MSVC"
	path = "C:\\Program Files\\Microsoft Visual Studio"
	clPath := ""
	filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if filepath.Base(path) == "ml64.exe" {
			clPath = path
			return err
		}
		return err
	})
	before, _, found := strings.Cut(clPath, "bin")
	if !found {
		panic("can not find win_toolchain_version")
	}
	before = filepath.Dir(before)
	before = filepath.Base(before)
	win_toolchain_version := strconv.Quote(before)

	before, _, found = strings.Cut(clPath, "Tools")
	if !found {
		panic("can not find win_vc")
	}
	win_vc := strconv.Quote(before)

	//C:\Program Files (x86)Windows Kits\10\bin\10.0.26100.0\x64\mc.exe
	mcPath := ""
	filepath.Walk("C:\\Program Files (x86)\\Windows Kits\\10\\bin", func(path string, info fs.FileInfo, err error) error {
		if filepath.Base(path) == "mc.exe" {
			mcPath = path
			return err
		}
		return err
	})

	before = filepath.Dir(mcPath)
	before = filepath.Dir(before)
	before = filepath.Base(before)
	win_sdk_version := strconv.Quote(before)

	out := stream.RunCommandArgs("clang", "-v").Output.String()
	//clang version 18.1.7
	//Target: x86_64-pc-windows-msvc
	//Thread model: posix
	//InstalledDir: C:\Program Files\LLVM\bin
	split := strings.Split(out, "\n")[1:]
	s := split[3]
	s = strings.TrimPrefix(s, "InstalledDir: ")
	clang_win := strconv.Quote(filepath.Dir(s))

	s = split[0]
	s = strings.TrimPrefix(s, "clang version ")
	clang_win_version := strings.Split(s, ".")[0]

	//win_vc="C:\Program Files\Microsoft Visual Studio\2022\Enterprise\VC"
	// win_toolchain_version="14.40.33807"
	//win_sdk_version="10.0.26100.0"

	//clang_win_version=18
	//clang_win="C:\Program Files\LLVM"
	return &config{
		win_vc:                win_vc,
		win_toolchain_version: win_toolchain_version,
		win_sdk_version:       win_sdk_version,
		clang_win_version:     clang_win_version,
		clang_win:             clang_win,
	}
}
