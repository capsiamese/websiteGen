---
title: "Makefile"
date: 2023-01-21T12:39:10+08:00
draft: true
tags: ['tools']
---


# 变量
变量声明时需要给一个初始值, 访问变量时要用`$变量名`也可以使用`${变量名}`, `$(变量名)`, 如果要使用$要通过`$$`转义.

make中使用`name = val`定义变量, 右侧可以是值也可以是变量, 并且右侧的变量是顺序无关的, 即可以使用还未定义的变量.

make中还有另一种定义变量的语法`name := val`这种不能使用未定义的变量

Q: 如何定义一个空格
A: make中的变量类似c语言中的宏, 所以不能用'', 必须先定义一个空字符串, 而make中#是注释开始, 也可以表示变量的终止.所以空字符和注释开始之间的就是空格了.
```make
nullstring := 
space := $(nullstring) #end of the line
```

`name ?= val` 如果name被定义过了什么都不做, 否则name就是val

## 变量值的替换
`${var:a=b}`意思是var中所有的以a结尾的替换成b

变量展开后开可以继续展开`$($(b))`也可以`$($(A)_$(B))`

追加变量值`objects += x.o`相当于`objects := $(objects) x.o`


# 条件判断
```make
ifeq (a,b)
    xxx
else
    xxx
endif
```

* ifeq
* ifneq
* ifdef
* ifndef

# 函数
`$(<function> <arguments>)`参数使用逗号分隔

* $(subst <from>,<to>,<text>)  字符串替换
* $(patsubst <pattern>,<replacement>,<text>) 模式字符串替换
* $(strip <string>) 去掉空格
* $(findstring <find>,<in>)
* $(filter <pattern...>,<text>)
* $(filter-out <pattern...>, <text>)
* $(sotr <list>)
* $(word n,text)
* $(wordlist ss,e,text)
* $(words text)
* $(firstword text)
* $(dir names...)
* $(notdir names...)
* $(suffix names...)
* $(basename names...)
* $(addprefix prefix, names...)
* $(join list1,list2)
* $(foreach var,list,text)
* $(if cond,then,else)
* $(call expression, params...) 其中params通过$(1),$(2)...表示
* $(origin var)
* $(shell)
* $(error text...)
* $(warning text...)



# GCC 命令



# References
[跟我一起写Makefile](https://seisman.github.io/how-to-write-makefile/index.html)
[gcc命令详解](https://blog.csdn.net/lu_embedded/article/details/117848775)
[GCC C compiler](https://www.rapidtables.com/code/linux/gcc.html#syntax)