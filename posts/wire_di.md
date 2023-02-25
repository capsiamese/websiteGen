---
title: "wire依赖注入学习"
date: 2023-02-25T11:33:41+08:00
draft: false
tags: []
---

# wire用户指南

# 基础
wire只有两个核心概念

1. providers
2. injectors

## 定义providers

**wire中最主要的机制就是provider, 主要作用是用来提供一个对象, provider可以是一个函数并返回一个值.**

```go
package foobarbaz

type Foo struct {
    X int
}

// ProvideFoo a Foo provider
func ProvideFoo() Foo {
    return Foo{X: 1}
}
```

provider函数必须是导出的(public)以供其他包使用.

**同时provider函数可以通过参数传入依赖.**

```go
package foobarbaz

type Bar struct {
    Y int
}

// ProvideBar provide a Bar depends on Foo
func ProvideBar(foo Foo) Bar {
    return Bar{Y: foo.X + 1}
}
```

**并且provider函数可以返回一个错误.**

```go
package foobarbaz

type Baz struct {
    Z int
}

func ProvideBaz(ctx context.Context, bar Bar) (Baz, error) {
    if bar.Y < 0 {
        return Baz{}, errors.New("can not provide baz, bar must be positive")
    }
    return Baz{Z: bar.Y + 2}, nil
}
```

**将provider组织成一组依赖.**provider可以被放入一个依赖集合中, 在某些场景中非常有用,
要创建一组依赖集合需要使用`wire.NewSet`函数. provider集合中的handler是顺序无关的, 所以
可以不用关注集合中handler的依赖关系. 并且集合是可以嵌套的.

```go
package foobarbaz

import (
    // ...
    "github.com/google/wire"
)

var SuperSet = wire.NewSet(ProvideBaz, ProvideBar, ProvideFoo)
var SuperSuperSet = wire.NewSet(SuperSet, ProvideOther)
```

## injectors

一个应用程序使用injector将这些provider连接(wires up)起来. injector是一个将这些provider按依赖顺序
组织, 我们只需要提供injector的签名, 然后wire会通过代码生成的方式来实现具体的过程. 这种实现方式主要是
通过go语言提供的build tags特性(如果文件以`//go:build wireinject`开头, build时没有指定tag的话该文件就不会参与构建)
然后在injector中调用wire.Build将需要依赖的providerSet或者provider们填进去, injector 的返回值也无关紧要*只要是正确的类型就行了*

```go
//go:build wireinject
package main

import (
    "context"
    "github.com/google/wire"
    "foobarbaz"
)

func newBaz(ctx context.Context) (foobarbaz.Baz, error) {
    wire.Build(foobarbaz.SuperSuperSet/*, xxx, xxx, xxx*/)
    return foobazbar.Baz{}, nil
}
```

# 高级特性

## 绑定接口

自然的, 依赖注入可以用来绑定一个接口的具体实现, wire通过[类型标识](https://go.dev/ref/spec#Type_identity)选择一个符合接口类型的provider.
然而, 这并不是惯用法(idiomatic), go的最佳实践是返回一个具体的类型, 事实上我们可以定义一个接口如何绑定一个具体的类型.

```go
type Fooer interface {
    Foo() string
}
type MyFooer string

func (m *MyFooer) Foo() string {
    return string(*m)
}

func provideMyFooer() *MyFooer {
    m := "Hello"
    return (*MyFooer)(&m)
}

type Bar string

func provideBar(f Fooer) string {
    return f.Foo()
}

var Set = wire.NewSet(
    provideMyFooer,
    wire.Bind(new(Fooer), new(*MyFooer)),
    provideBar,
)
```

`wire.Bind`的第一个参数是需要用到的接口, 而第二个参数是实现这个接口的具体类型的指针.
任何一个包含接口绑定的集合必须至有一个提供实现了这个接口的provider.

## 结构体providers

结构体可以通过provider提供的值来构造, 使用`wire.Struct`函数来构造一个结构体类型, 然后告诉injector那些字段应该被注入.
injector将填充每一个声明了的字段.

```go
type Foo int
type Bar int

func ProvideFoo() Foo {return 1}
func ProvideBar() Bar {return 2}

type FooBar struct {
    MyFoo Foo
    MyBar Bar
}

var Set = wire.Set(
    ProvideFoo, ProvideBar,
    wire.Struct(new(FooBar), "MyFoo", "MyBar"),
)
```

这里第一个参数是要将值注入的类型的指针, 后面是需要注入的字段的名字, 如果需要将所有字段注入可以直接使用`"*"`=>`wire.Struct(new(FooBar), "*")`.

使用`*`时如果要避免某个字段不被注入可以使用结构体标签`wire: "-"`

## 绑定值

如果要给一个字段绑定一个具体值的话可以使用`wire.Value(Foo{X: 1})`, 对于接口的话使用`wire.InterfaceValue(new(io.Reader), os.Stdin)`

## 将结构体的某个字段作为provider

有时provider想要使用结构体的某个字段, 我们可以定义一个额外的provider, 但是这通常是不必要的, 可以使用`wire.FieldsOf`.

`wire.Build(provideFoo, wire.FieldsOf(new(Foo), "S"))`

## clean up 函数

如果一个provide创建的值需要在不用时销毁, 我们可以给provider的返回值添加一个clean up函数.

```go
func provideFile(log Logger, path Path) (*os.File, func(), error) {
    f, err := os.Open(string(path))
    if err != nil {
        return nil, nil, err
    }
    cleanup = func() {
        if err := f.Close(); err != nil {
            log.Log(err)
        }
    }
    return f, cleanup, nil
}
```

如果我们使用injector时不想处理错误可以使用

```go
func injectFoo() Foo {
    panic(wire.Build(ProvideFoo, ProvideBar/*...*/))
}
```

# 参考

[wire guide](https://github.com/google/wire/blob/main/docs/guide.md)
