#!/usr/bin/env just --justfile
# https://cheatography.com/linux-china/cheat-sheets/justfile/

llvm_version := shell("llvm-config --version | awk -F'.' '{ print $1 }'")
export CC := "clang"
export PATH := shell("llvm-config --prefix") + x"/bin:$PATH"
export CGO_CPPFLAGS := shell("llvm-config --cppflags")
export CGO_LDFLAGS := shell("llvm-config --ldflags --libs --system-libs all")

go_tags := f'"llvm{{llvm_version}}"'
cc_quoted := quote(CC)
path_quoted := quote(PATH)
cgo_cppflags_quoted := quote(trim(replace(CGO_CPPFLAGS, "\n", " ")))
cgo_ldflags_quoted := quote(trim(replace(CGO_LDFLAGS, "\n", " ")))

compile_test input output: 
    go run -tags={{go_tags}} . -i {{input}} -o {{output}} --debug

export_envs:
    @rm -f build.env
    @touch build.env
    @echo "CC={{cc_quoted}}" >> build.env
    @echo "PATH={{path_quoted}}" >> build.env
    @echo "CGO_CPPFLAGS={{cgo_cppflags_quoted}}" >> build.env
    @echo "CGO_LDFLAGS={{cgo_ldflags_quoted}}" >> build.env

test:
    go test -tags={{go_tags}} ./...

update:
  go get -u
  go mod tidy -v
