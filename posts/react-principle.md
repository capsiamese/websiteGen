---
title: "React原理"
date: 2022-07-03T14:30:17+08:00
draft: true
tags: ["js"]
---

# [宏观结构](https://7kms.github.io/react-illustration-series/main/macro-structure)

## 基础包结构
源代码位于[react/packages](https://github.com/facebook/react/tree/17.0.2/packages)下

1. react基础包, 提供了组件的基础定义
2. react-dom 渲染器
3. react-reconciler 协调其他包的配合, 将输入信号转换成输出信号传递给渲染器
4. scheduler 调度机制的核心实现

## 宏观总览

react应用整体可分为两层, 一层api层, 一层core层.

接口层提供了平时使用的绝大多数api, react启动后正常可以改变渲染的操作有3个

1. 函数组件中的setState()
2. 函数组件中的hooks并发起dispatchAction改变hook对象
3. 改变context

内核层由3部分组成

1. scheduler内部有一个回调队列, 把react-reconciler提供的回调函数包装到任务对象中执行.
2. react-reconciler负责装在渲染器, 接受render和setState发起的更新请求, 将fiber的构造过程传入scheduler等待调用
3. react-dom负责生成dom节点

![内部关系](img/react-core-package.png)

# [两大工作循环](https://7kms.github.io/react-illustration-series/main/workloop)