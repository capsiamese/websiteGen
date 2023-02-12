---
title: "Golang_note"
date: 2022-11-09T13:48:02+08:00
draft: false
tags: ["golang"]
---

# advice

1. 使用time.Duration代替int64
2. 为整数常量实现Stringer
3. 每个阻塞或者io操作都应该有可取消或可超时
4. 不要给枚举使用类型别名这会打破类型安全, 要重新定义类型
5. 两个空结构体的地址可能相同但是比较时可能不同 (UB)
6. slice, map, chan零值是nil, 内置函数是类型安全的
7. 文件操作使用os.FileMode而不是使用 0777
8. 不要依赖计算顺序, return res, json.Unmarshal(b, &res)不要这样做
9. 在结构体定义中加入`_ struct{}`可以防止纯值初始化`A{1, 2, 3}`这样会失败
10. 在结构体中加入`_ [0]func()`可以防止使用`==`比较两个结构体的值
11. 管道应该由创建者关闭, 向一个关闭的管道发送数据会`panic`, 但是读不会.
12. 总是要关闭http body defer r.Body.Close() 否则会导致泄露, 因为会服用goroutine所以泄露的少
13. 不要过度使用fmt.Sprintf("%s%s", a, b)
14. 如果不要body考虑io.Copy(ioutil.Discard, resp.Body) http客户端不会重用链接直到body被读完.
15. 不要在循环中使用defer
16. 不要忘记停止time.Ticker
17. 清空map的方式1. for k in m delete(m, k)  2. m = make(map)
18. 


## 内存逃逸

使用`go build -gcflags=-m`

### 引起逃逸的情况
1. 在方法内将局部变量通过指针返回
2. 发送指针或带有指针的值到channel中, 编译时没法知道那个goroutine会在channel上接收数据, 所以编译器没法知道变量什么时候被释放
3. 在一个切片上存储指针或者带指针的值, 会导致切片的内容逃逸到堆上
4. slice数组被重新分配时
5. 在interface类型上调用方法, 在接口类型上调用方法都是动态调度, 方法真正的实现只有在运行时知道

## 字符串转成byte数组是否会发生内存拷贝
字符串转成切片会发生内存拷贝, 只要发生强制类型转换都会发生内存拷贝

### 如何在转换时不发生内存拷贝?
```go
package main

import "reflect"
import "unsafe"

func main() {
	a := "aaa"
	stringHeader := (*reflect.StringHeader)(unsafe.Pointer(&a))
	b := *(*[]byte)(unsafe.Pointer(stringHeader))
	_ = b
}
```

## GMP
GMP代表了三个角色 
G: Goroutine 通过go关键字创建的执行体, 对应一个结构体g, 结构体保存了goroutine的堆栈信息
P: Process 代表处理器, 建立起了G和M的联系
M: Machine 代表操作系统的线程

goroutine由一个runtime.go(type g struct)的结构体表示, 主要存储执行的堆栈状态, 当前占用的线程, 调度等数据

goroutine的调度数据存储在sched.go, 在协程切换, 恢复上下文时用到

M对应操作系统的线程数, 由GOMAXPROCS控制(type m struct) 结构体中两个比较重要的东西 g0 表示深度参与运行时的调度过程, 比如创建, 内存非陪等. curg表示正在执行的goroutine, p表示正在运行代码的处理器, nextp, oldp 前一个和后一个处理器.

processor负责G和M的连接, 提供线程需要的上下文环境, 也能分配G到对应的线程上执行(type p struct)存储了性能追踪, 垃圾回收, 处理器的待运行队列, 计时器等.
```text
   GlobalQ   <-  g g g
   g         g           g
   g         g           g
   Q         Q           Q  // 当p的队列满时g加入全局的队列
   ↓         ↓           ↓
   p         p           p 
   ↓         ↓           ↓
   m         m           m
```
### 调度策略
1. 调度时机
2. 调度策略
3. 切换机制

#### CAS操作和ABA问题
cas步骤: 1.内存地址V 2.旧的预期值A 3.新值B
当且仅当地址V的值是预期值A时才更新成B值
ABA问题: 指在cas操作时其他线程将变量从A改成B但又改回成A, 此时其他线程在执行cas时又会改成B了.
解决方法: 增加版本号每次更新时版本号+1, 即A->B->C 变成 1A->2B->3C


# http request 中Form的解析
|动作           | 函数                | 是否读取url | 是否读取body | 支持文本 | 支持二进制 |
|Form          |ParseForm()          | Y         | Y          | Y       |          |
|PostFrom      |ParseForm()          |           | Y          | Y       |          |
|FormValue     |ParseForm()          | Y         | Y          | Y       |          |
|PostFormValue |ParseForm()          |           | Y          | Y       |          |
|MultipartForm |ParseMultipartForm() |           | Y          | Y       | Y        |
|FormFile      |ParseMultipartForm() |           | Y          |         | Y        |

`x-www-form-urlencoded` 每个键值对在body中的形式和url中的query一样

`form-data`: 每个键值对在body中都会用boundary隔开, body中的格式为`----xxxxx`开头`Content-Disposition: form-data; name=xxx\n\n`
结尾是`------xxxx--`, boundary在`Content-Type`字段中
