---
title: "Redis_note"
date: 2022-11-09T15:27:21+08:00
draft: false
tags: ["db"]
---

# 数据结构
1. String
2. Hash
3. List
4. Set
5. SortedSet

底层实现有8种
1. SDS(simple dynamic string) 支持自动扩容的字节数组
2. list 双向链表
3. dict 使用双hash实现, 支持平滑扩容的字典
4. zskiplist 附加了向后指针的跳跃表
5. intset 用于存储整数的集合
6. ziplist 用于存储任意数据的有序序列, 类似TLV(type-length-value)
7. quicklist 以ziplist作为节点的上香链表
8. zipmap 用于小规模场景的轻量级字典

#### 散列冲突

## String
redis中最基本的数据结构, 一个Key对应一个Value, 是二进制安全的, 使用sds实现

***使用场景***
1. 缓存, 将常用信息, 字符串, 图片, 视频缓存到redis中, mysql做持久化, 降低mysql的读写压力
2. 计数器, redis是单线程的一个命令执行完才会执行另一个
3. session, redis实现session共享

## Hash
是一个HashMap 通过Key找到这个HashMap然后通过FieldKey找到Value

***使用场景***
1. 缓存, 更直观, 比String更省空间

## List
实现为双向链表
***使用场景***
1. 栈
2. 队列
3. 有限列表
4. 消息队列

## Set
***特点***
1. 不能重复
2. 无序
3. 支持集合操作

*** 使用场景***
1. 标签
2. 点赞

## SortedSet
有序集合, 成员不能重复, 但是分数可以重复
*** 使用场景***
1. 排行榜

# 持久化
redis提供两种持久化方式
1. AOF 每次记录服务器操作, 重启时重新执行这些操作
2. RDB 指定间隔创建一次快照

### RDB
`save`命令会阻塞服务, 直到持久化完成
`bgsave`不会阻塞服务

### AOF

# 缓存淘汰算法

# 主从复制