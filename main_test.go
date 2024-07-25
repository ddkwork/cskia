package main

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestName(t *testing.T) {
	//githubWorkspace := os.Getenv("GITHUB_WORKSPACE")
	//GitHub Actions workspace directory: D:\a\cskia\cskia
	//D:\a\cskia\cskia\capi\sk_capi.cpp
	//D:\a\cskia\cskia\capi\sk_capi.h
	//D:\a\cskia\cskia\capi\sk_capi.h
	githubWorkspace := "D:\\a\\cskia\\cskia"
	assert.Equal(t, "D:\\a\\cskia\\cskia\\capi\\sk_capi.h", filepath.Join(githubWorkspace, "capi/sk_capi.h"))
}
