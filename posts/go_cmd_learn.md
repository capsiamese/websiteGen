---
title: "Go-command-line"
date: 2022-11-30T13:26:17+08:00
draft: false
tags: ["golang"]
---


# go build
`go build [-o output] [build flags] [packages]`

build, clearn, get, install, list, run, test 共享的参数

1. `-a` 强制重新构建已经是最新的包
2. `-n` 打印命令但不运行
3. `-p n` 并发程度, 默认是GOMAXPROCS的值
4. `-race` 启用数据竞争检测
5. `-msan` 启用内存错误检测(用的是google sanitizer)
6. `asan` 启用地址错误检测
7. `-v` 打印出要编译的所有包
8. `-work` 打印出编译时的临时目录并且编译完成后不删除
9. `-x` 打印出编译命令
10. `-asmflags '[pattern=]arg list'` 传递给 go tool asm调用的参数
11. `-buildmode mode` 见buildmode
12. `buildvsc` 是用版本控制信息标记二进制文件 [true | false | auto]
13. `-complier name` 指定编译器[gccgo | gc]
14. `gccgoflags '[pattern=]arg list'` gccgo编译器参数
15. `gcflags '[pattern=]arg list'` go tool compile编译参数
16. `-installsuffix suffix` 安装后缀
17. `-ldflags '[pattern=]arg list'` go tool link参数
18. `-linkshared` 链接共享库
19. `-mod mode` 下载mod模式 readonly, vendor, mod
20. `-modcacherw` 
21. `-modfile file` 
22. `-overlay file`
23. `-pkgdir dir`
24. `-tags tag,list`
25. `-trimpath`
26. `toolexec 'cmd args'`

# go clean

`go clean [clean flags] [build flags] [packages]`

# go env
` go env [-json] [-u] [-w] [var ...]`

# go fix
`go fix [-fix list] [packages]` 

Update packages to use new APIs

# go fmt
`go fmt [-n] [-x] [packages]`

# go generate
`go generate [-run regexp] [-n] [-v] [-x] [build flags] [file.go... | packages]`

# go get
`go get [-t] [-u] [-v] [build flags] [packages]`


# buildmode
`go help buildmode`

当使用`go build`和`go install`时传递-buildmode参数指出将要构建那种类型的obj文件.

1. -buildmode=archive 构建列出的没有main包的文件到.a文件中
2. -buildmode=c-archive 构建所有列出的main包加上所有导入的包到一个C archive文件中, 只有使用cgo注释的函数才能调用`//export`
3. -buildmode=c-shared 构建所有包到一个C shared库中, 也必须使用`//export`注释
4. -buildmode=default 列出的main包将会被构建入可执行文件中, 列出的非main包将构建入.a文件中
5. -buildmode=shared 构建所有non-main包到单个 shared library中, 构建时使用-linkshared选项
6. -buildmode=exe  构建所有列出的main包和所有导入的包到可执行文件中, 包名不是main的都将忽略
7. -buildmode=pie 构建到pie中(position independent executables)
8. -buildmode=plugin 构建go plugin

# Ref
[Cmd](https://pkg.go.dev/cmd/go)
[mod reference](https://go.dev/ref/mod)
[mod ref](https://go.dev/doc/modules/gomod-ref)