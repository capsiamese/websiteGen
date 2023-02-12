---
title: "Go Mistakes"
description: "100 go mistakes and how to avoid them"
date: 2022-12-06T13:09:55+08:00
draft: true
tags: ["golang"]
---

# 1. Go: Sample to learn but hard to master

本章覆盖了以下主题
1. 是什么使golang成为如此的高效, 可拓展, 可用在生产生的语言.
2. 解释为什么golang是简单的但是难于精通
3. 介绍开发者常见的类型错误

犯错误是每个人生命中的一部分, 正如爱因斯坦曾经说过的, 一个人从来不犯错误那他将永远不会去尝试任何新的东西.

重要的不是我们最终犯错误的数量, 而是我们从犯的错误中学习的能量. 这个主张也适用于程序员. 我们从语言中获得的资历不是一个神奇的过程, 而是
我们犯了多少错误和从中学到的. 这本书的宗旨就是围绕着这个点. 这本书将会帮助你变成一个更专业的golang程序员通过学习这100个开发者在这门语言
的众多领域通常会犯的错误.

这章将会先向你呈现为什么golang在过去几年变成主流. 我们将会讨论为什么尽管golang被设计的尽可能简单易学, 但是他的细微差别将会成为挑战.
最终, 我们将会介绍这本书封面的核心.

## 1.1 Go outline

如果你正在阅读这本书, 那么你大概率已经掌握了golang, 因此, 这节将提供一个简短的提醒告诉你go为什么是一门如此强大的语言.

软件工程在过去的几十年中取得了很大的发展, 许多现代化的系统不在只被一个人实现, 而是由多个程序员有时甚至是成百上千个, 而现在
代码必须是容易被阅读的, 容易表达的, 可维护的来保证系统在几年内的持续更新, 与此同时在我们这个快节奏的时代, 最大限度的提高
敏捷性和最短的时间内上市对大多数组织来说是至关重要的. 程序员也应该跟随这个趋势. 并且公司的竞争确保软件工程师们写出的生产代码
尽可能可读可写.

功能方面go没有继承, 没有异常, 没有宏, 没有偏函数, 不支持惰性或者及早求值, 没有运算符重载, 没有模式匹配等等...
为什么go语言没有这么多特性[官方](https://go.dev/doc/faq)给我们了一些见解.

为什么go没有xxx特性? 没有你最喜欢的特性因为它不合适, 因为它影响编译速度或者设计的清晰度, 又或者因为它将使基本类型系统变得太复杂.

判断一个语言的质量通过特性的数量或许是不准确的衡量. 最后go也没有对象, 实际上go利用一些必不可少的特性来实现上规模的组织, 包括
1. 稳定性--尽管go接受频繁的更新(包括改进和安全补丁), 但它仍然保持语言的稳定. 这甚至是语言应该考虑的最佳特性
2. 表现力--我们可以通过编写和阅读代码的自然和直观程度来定义编程语言的表现力, 一个减少关键字和限制解决常见问题的方法使go成为大型代码库的实现语言
3. 编译--作为开发者还有什么比等待构建来测试我们的应用程序更令人恼火的? 目标的快速编译时间已经被有意识的作为了语言设计的目标, 这将分过来启用生产力
4. 安全行--go是健壮, 静态类型的语言, 因此它有严格的编译期规则来确保代码的类型安全在大多数情况下.

go语言是建立在优秀的并发原语和goroutine和channel的坚实基础上的. 这些不是强依赖外部库来建立高效的并发程序, 观察并发的重要性已经证明了为什么
go是如此的合适用来呈现可见的未来.

## 1.2 Simple doesn't mean easy

简单和容易之间有一些微妙的不同, 简单同样也使用于技术方面, 意味着学起来或者理解起来不复杂, 然而容易意味着我们可以构建任何东西不需要太多的努力, go
是容易学习的但不一定是容易精通的.

让我们以并发为例, 例如在2019年一个学术焦点在并发错误被出版了: "understanding Real-World Concurrency Bugs in Go" 这个学术讨论首次系统化的
分析了并发的bug, 它的焦点在多个受欢迎的go库中, 例如docker, grpc, k8s, 这项研究最重要的收获之一是大多数的阻塞错误是由不正确的使用channel传递
消息的范例引起的, 尽管相信消息传递比共享内存的错误更容易处理.

对这样的外部声音应该有什么适当的反应? 例如我们考虑语言的设计是错的关于消息传递? 我们是否应该重新考虑如何处理并发在我们的程序中? 当然不是.

问题不在于面对消息传递和共享内存之间这是赢家, 然而这激起了我们go开发者有责任彻底理解如何使用并发. 这意味着在现代处理器上什么时候更偏爱一种方法,
并且如何避免一些陷阱, 这个例子突出的抛出了一个观点, 例如channel和goroutine是容易学习的, 但并不是一个简单的主题在实践中.

简单并不意味着容易这个主旋律可以一般性的用在go的各个方面, 不仅仅是并发, 因此要成为一个更专业的go程序员我们必须了解语言的各个方面, 这需要大量的时间
和努力和犯大量的错误.

这本书的目标是加快帮助我们的读者通过研究这100个错误变得更专业.

## 1.3 100 Go mistakes

为什么我们读这本关于常常会犯的100个错误的书? 问什么不加深我们的知识通过一个更通常的方式.

在一个2011年的文章中, 神经科学家提到了一个大脑成长的最好时机是当我们面对错误时. 难道我们都没有经历过从错误中学习和几个月后或者几年后回忆之前的场合
当联想到一些上下文时?  现在另一个文章告诉我们犯错误是有促进作用的, 主观点是我们可以记起不仅仅是错误并且是围绕错误的上下问. 这就是另一个为什么我们
从错误中学习更高效的原因.

为了加强这个促进作用, 这本书中的每个例子都紧紧围绕着生活中的错误. 这本书不仅仅是理论, 还帮助我们更好的避免错误并且如何更好的改进. 
有意识的决定, 因为我们现在已经了解了它们背后的基本原理

Tell me and I forget, Teach me and I remember. Involve me and I learn.  -- 鲁迅

这本书呈现了七个主要类别的错误, 这些错误整体可以分为一下几种
1. Bugs
2. Needless complexity
3. Weaker readability
4. Suboptimal or unidiomatic organization  次优的或者单一的组织
5. Lack of API convenience
6. Under-optimized code       优化不足的代码
7. Lake of productivity

接下来我们将会介绍每个类别

### 1.3.1 Bugs

第一个类型的错误也是绝大多数的显然是软件的bug, 在2020年, 一个研究表明由软件bug引起的花费估计达到了$2 trillion

这本书涵盖了很多关于这些沉重的软件bug的各种案例, 包括数据竞争, 数据泄露, 逻辑错误, 和一些其他缺陷. 尽管准确的测试应该尽可能
早的发现这些bug, 我们某些时候错过了一些测试案例由于不同的工厂函数例如时间限制或者复杂性, 因此, 作为一个go开发者, 确保
避开这些错误是必不可少的.

### 1.3.2 Needless complexity

下一个关于错误的分类涉及到一些不必要的复杂性, 这是关于软件复杂性最显著的部分, 作为开发者我们努力思考关于虚构的未来, 而不是立刻解决实际问题.
构建可以解决未来出现的任何问题的进化软件可能很诱人, 这可能导致缺点大于效益在大多数情况下, 因为这可能使理解代码库变得更复杂.

回到go上, 我们可以思考丰富的用力在那些开发者可能被引诱设计抽象的特性, 例如接口或者泛型, 这本书讨论的主题当我们应该保持小心的不损坏代码库
的前提下增加不必要的复杂性.

### 1.3.3 Weaker readability

另一种类型的错误使很弱的可读性, 就行Clean Code中写道的一样, 在读代码和写代码的比率上读要超过写的10倍, 在大多数我们的个人项目中, 可读性可能
不是很重要, 然而今天的软件工程增加了一个时间维度, 确保我们能继续工作和保持软件维护几个月几年甚至十几年之久.

当我们使用go编程时, 可能做了许多有损可读性的事, 可能包括, 签到代码, 数据类型交涉, 后者未使用命名返回参数在一些例子中, 整本书我们将学习如何写出
可读的代码和照顾未来的读者(包括未来的自己)

### 1.3.4 Suboptimal or unidiomatic organization

无论是从事新项目时, 或者我们获得了不准确的反馈, 另一个错误是组织我们的代码和一个项目的次优或者单一, 这些问题可以使一个项目很难保持的原因, 
这本书涵盖了一些这样的问题, 例如我们将看到如何构建一个项目和如何处理工具包和init函数, 总而言之, 学习这些案例将会帮助我们组织我们的代码
项目更有效和通用.

### 1.3.5 Lack of API convenience

犯下削弱api对我们客户的便利性是另一种类型的错误, 如果一个api不是用户有好的, 将缺乏表现力, 因此会更加难理解且更容易出错.

我们能想象关于许多情况, 例如过度使用any类型, 使用错误的创建模式来处理可选值, 或者盲目的适用基本做法在oop上, 影响使用我们的api, 这本书涵盖了
阻止我们向用户提供方便的api的常见错误

### 1.3.6 Under-optimized code

优化不足的代码是另一种可能的错误, 这可能由各种原因导致, 例如没有理解语言的提醒或者缺乏基本的知识, 性能是一个重要的且明显的影响对于这个错误.

我们可以想想关于优化的另一个目标, 例如精度, 这本书提供了一些通用的技术来确保浮点运算是准确的, 与此同时, 还涵盖了许多消极影响性能的代码
因为匮乏的并行执行, 不知道如何教导内存分配或者影响数据对齐, 例如我们将解决优化通过不同方面.

### 1.3.7 Lack of productivity

在大多数例子中, 什么是最佳的语言我们可以选择当我们进行一个新的项目? 其中一个是最有生产力的, 变得舒适的和一个语言如何工作, 并且利用它
充分利用达到熟练的程度至关重要.

这本书我们将涵盖许多坚实的例子将会帮助我们更有生产力, 实际上, 我们将看到如何写出高效的测试确保我们的代码工作, 依赖标准库更有效,
并且获得最佳的剖析工具和格式化工具, 现在是时候深入到这100个错误中了.

## Summary
* go是现代的语言启用生产力, 这是至关重要的对于公司
* go是易学的但是难以精通, 这是为什么我们需要深入我们的知识确保有效的
* 学习通过错误和坚实的例子是一个强有力的途径去精通语言, 这本书将加快我们精通的路径


# 2.Code and project organization

本章包含以下内容
* 组织我们的代码变得通用
* 如何高效的处理抽象, 接口, 泛型
* 如何通过最佳实践构建项目

组织一个干净的, 通用的, 可维护的go代码库一点也不容易, 这需要经验和甚至是失败来理解最佳实践涉及到的代码和项目组织, 陷阱应该如何避免
例如 变量覆盖和嵌套代码滥用, 我们应该如何组织包结构, 什么时间什么地点我们应该使用接口或者泛型, init函数和工具包, 在这章, 我们
将了解常见的组织错误.

## 2.1 Unintended variable shadowing

变量的范围是指变量可以被引用的地方, 换句话说, 一个名字绑定一个有效的值在程序的部分, go中, 一个变量名声明在一个块中可以被重新声明在
内部块中, 原理是变量遮罩是容易犯的错误.

下面这个例子展示了一个无意识的副作用, 因为一个被遮盖的变量, 它创建了一个HTTP client在两个不同的地方, 依赖一个`tracing`的布尔变量.
```go
var client *http.Client
if tracing {
	client, err := createClientWithTracing()
	if err != nil {
		return err
    }
	log.Println(client)
} else {
	client, err := createDefaultClient()
	if err != nil {
		return err
    }
	log.Println(client)
}
// Use client
```
在这个例子中, 我们声明了一个client变量, 然后我们使用短变量声明在内部块中将结果赋值给内部的client变量而不是外部的, 这样做的结果是外部的client
始终是nil

我们如何确保变量被赋值到原始的client变量上? 有两种做法

第一种做法是使用一个临时变量在内部块中.
```go
var client *http.Client
if tracing {
	c, err := createClientWithTracing()
	if err != nil {
		return err
    }
	client = c
} else {
	// same logic
}
```
这次我们将结果赋值给了一个临时变量c, 他的作用域只在内部块中, 然后我们将它赋值给client变量, 同时else也一样.

第二种方式是使用赋值运算符在内部代码块中, 将结果直接赋值给client变量, 然而这需要再创建一个error变量, 因为赋值运算符只能用在变量已经被声明过的地方
```go
var client *http.Client
var err error
if tracing {
	client, err = createClientWithTracing()
	if err != nil {
		return err
    }
} else {
	// same logic
}
```
对比第一种方式我们可以直接赋值结果给client变量.

这两种方式都可以避免, 这两种可选的方式主要的区别在于我们我们只少做了一次赋值运算在第二种方式中, 这种可能更容易阅读, 同时在第二种方式中, 我们可以
将err处理放到条件判断外边.
```go
var client *http.Client
var err error
if tracing {
	client, err = createClientWithTracing()
} else {
	client, err = createDefaultClient()
}
if err != nil {
	return err
}
```
变量遮挡发生在一个变量名重复声明在内部代码块中, 但是通过这个实践可以很容易的避免. 应该明确一个规则来禁止依赖个人的味道. 例如有时可能很方便的
重用一个变量名像err对于error一样, 当然, 通常来讲, 我们应该谨慎的保留这些, 因为我们现在知道我们能面对谨慎的代码编译. 但是接受一个值不是预期的
对于一个变量, 本章后面, 我们将看到如何检测代码遮罩, 这将帮助我们发现可能的bugs.

接下来的小节为我们展示为什么避免嵌套代码滥用是重要的

## 2.2 Unnecessary nested code

一个思考模型应用在软件是一个内部表示一个系统的行为, 编程时, 我们需要保持这个思维模式(关于总体的代码交互和功能实现). 代码是具备可读基于多个角度的
例如名字, 一致性, 格式等等, 可读的代码需要跟少的认知感来保持思维模式, 因此这是容易保持和读的.

一个值得批判的方面对于可读性是嵌套的深度, 让我们做个练习, 想想我们正在做一个新的项目, 并且需要理解下面的代码.
```go
func join(s1, s2 string, max int) (string, error) {
	if s1 == "" {
		return "", errors.New("s1 is empty")
    } else {
        if s2 == "" {
            return "", errors.New("s2 is empty")
        } else {
			concat, err := concatenate(s1, s2)
			if err != nil {
				return "", err
            } else {
				if len(concat) > max {
					return concat[:max], nil
                } else {
					return concat, nil
                }
            }
        }
    }
}
func concatenate(s1, s2 string) (string, error) {/*...*/}
```

join函数将连个字符串连接然后返回一个子串如果长度大于max的话, 同时还检查s1和s2是否为空.

从一种实现上来说, 这个函数是正确的, 然而构建一个思维模式, 包含所有的不同因为大概率不是一个线性过程, 为什么, 因为包含了很多嵌套.

现在我们再次试着做个练习.
```go
func join(s1, s2 string, max int) (string ,error) {
	if s1 == "" {
		return "", errors.New("s1 is empty")
    }
    if s2 == "" {
        return "", errors.New("s2 is empty")
    }
	concat, err := concatenate(s1, s2)
	if err != nil {
		return "", err
    }
	if len(concat) > max {
		return concat[:max], nil
    }
	return concat, nil
}
```

你或许可能注意到这里构建的思维模式了, 在这个全新的版本中, 需要更少的认知尽管在相同的代码中, 这里我们仅仅有两层嵌套.

这里很难讨论预期的执行流程在第一个版本中, 因为嵌套逻辑. 相反的第二个版本只需要向下扫描一列, 期望的历程就可以看到边界条件如何处理.

通常来讲, 越多的层级嵌套会带来更多的复杂度去读和理解, 让我们看看一些不同的程序规则如何优化代码的可读性.

* 当if块返回时, 我们删除else块
```go
if foo() {
	// ...
	return true
} else {
    // ...	
}
```
然而我们可以删除else块
```go
if foo() {
	// ...
	return true
}
// ...
```
在这个新的版本中之前的else代码逻辑被提升到了顶层代码块中使之更容易阅读.

* 我们也可以跟随下面的逻辑通过一个不快乐的方式.
```go
if s != "" {
	// ...
} else {
	return errors.New("empty string")
}
```
这里一个空的表示不快乐的路线, 因此我们应该反转条件
```go
if s == "" {
	return errors.New("empty string")
}
```
写出可阅读的代码使重要的挑战对于每个开发者, 与消除嵌套代码块斗争, 对齐快乐的路径在左侧, 尽可能以简单的方式是具体的提升我们代码可读性.

下一节我们讨论滥用init函数

## 2.3 misusing init functions

有时我们滥用init函数, 潜在的后果是错误管理或者代码难以理解. 让我们刷新思维关于init函数, 然后我们将看到如何使用.

### 2.3.1 comcepts

一个init函数被用来初始化状态在应用中, 它没有参数也没有返回值, 当一个包初始化时, 所有的常量和变量声明在这个包中都将被求值.
然后init函数将会执行.
```go
package main
import "fmt"
var a = func() int {
	fmt.Println("var")
	return 0
}()
func init() {
	fmt.Println("init")
}
func main() {
	fmt.Println("main")
}
```
代码执行结果是
```text
var
init
main
```

一个init函数在包初始时执行, 在下面的例子中, 我们定义了两个包, main和redis, main依赖redis.
```go
package main
import "fmt"
import "redis"
func init() {
	//...
}
func main() {
	err := redis.Store("foo", "bar")
	// ...
}
```
```go
package redis
//imports
func init() {
	// ...
}
func Store(key, value string) error {
	// ...
}
```
因为main包依赖redis包, 所以redis包中的init函数先执行, 接着是main包中的init函数, 然后是main函数自身.

我们可以在每个包中定义init函数, 当我们这么做时init函数的执行顺序在包中是字母顺序的, 例如一个包中包含a.go和b.go都有init函数, 则a.go中的会
先执行.

我们不能依赖包中init函数的执行顺序, 事实上, 这可能是危险的当源文件可能被重命名, 潜在的影响执行顺序.

同时我们也可以定义多个init函数在一个源文件中.
```go
package main
import "fmt"
func init() {
	fmt.Println("init 1")
}
func init() {
	fmt.Println("init 2")
}
func main() {}
```
```text
init 1
init 2
```

同时我们可以让init函数作为副作用, 下一个例子中, 我们定义一个main包, 但是没有强依赖foo(例如这里没有直接使用公共函数), 然而这要求foo包是被初始化过的, 我们可以通过`_`运算符
实现.
```go
package main
import "fmt"
import _ "foo"
func main() {}
```
这样foo包将会被初始化在main之前, 所以foo的init函数将会被执行.

另一方面, 一个init函数不能被直接执行.
```go
package main

func init(){}

func main() {
	init()
}
```
这个代码将会出现以下编译错误
```shell
$ go build .
./main.go:6:2: undefined: init
```

现在我们有了新的认识关于init函数如何工作的. 下面让我们看看应该什么时候使用, 什么时候不要使用init函数.

### 2.3.1 When to use init functions
首先, 让我们看一个使用init函数被认为是不合时宜的例子. 保持一个数据库连接池, 在这个init函数中, 我们打开了一个数据库用`sql.Open` 然后我们将这个变量变成一个全局的.
```go
var db sql.DB
func init() {
	dsn := os.Getenv("MYSQL_DATA_SOURCE_NAME")
	d, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
    }
	err = d.Ping()
	if err != nil{
		panic(err)
	}
	db = d
}
```









# Reference
[100 go mistakes and how to avoid them](https://www.manning.com/books/100-go-mistakes-and-how-to-avoid-them)
[100 go mistakes and how to avoid them](https://livebook.manning.com/book/100-go-mistakes-and-how-to-avoid-them/chapter-1/)
































