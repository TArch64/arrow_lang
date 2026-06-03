#!/usr/bin/env just --justfile

llvm_version := shell("llvm-config --version | awk -F'.' '{ print $1 }'")
export CC := "clang"
#export PATH := quote(shell("llvm-config --prefix") + x"/bin:$PATH")
#export CGO_CPPFLAGS := quote(trim(replace(shell("llvm-config --cppflags"), "\n", " ")))
#export CGO_LDFLAGS := quote(trim(replace(shell("llvm-config --ldflags --libs --system-libs all"), "\n", " ")))
export PATH := shell("llvm-config --prefix") + x"/bin:$PATH"
export CGO_CPPFLAGS := shell("llvm-config --cppflags")
export CGO_LDFLAGS := shell("llvm-config --ldflags --libs --system-libs all")

compile_test input output: 
    go run -tags=llvm{{llvm_version}} . -i {{input}} -o {{output}} --debug

export_envs:
    env | grep -E '^(CC|PATH|CGO_CPPFLAGS|CGO_LDFLAGS)=' > build.env

update:
  go get -u
  go mod tidy -v
