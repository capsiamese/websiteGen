---
title: "wireshark过滤器"
date: 2023-03-16T10:58:59+08:00
draft: true
tags: []
---

## 过滤器分类
wireshark中过滤器分为两种
1. 捕获过滤器, 在开始捕获网卡流量前进行设置
2. 显示过滤器, 在捕获网卡流量的过程中可以设置, 设置条件后会隐藏不符合条件的数据包

## 显示过滤器写法

### 过滤值比较符号和意义

| 英文 | 符号 | 描述 | 示例 |
|------|-----|------|------|
| eq | == | 等于 | ip.src==127.0.0.1 |
| ne | != | 不等于 | ip.src!=127.0.0.1 |
| gt | > | 大于 | frame.len > 10 |
| lt | < | 小于 | frame.len < 127 |
| ge | >= | 大于等于 | frame.len ge 100 |
| le | <= | 小于等于 | frame.len le 10 |
| contains | "string" | 包含 | sip.To contains "abc" |
| matches | "string" | 正则匹配 | host matches "aaa(org|com|io)" |
| bitwise_and | & | 位操作 | tcp.flags & 0x01 |


### 多表达式之间的组合

| 英文 | 符号 | 描述 | 示例 |
|------|-----|------|------|
| and | && | 逻辑与 | ip.src == 127.0.0.1 and tcp.flags.fin |
| or | \|\| | 逻辑或 | ip.src == 127.0.0.1 or ip.src==192.168.0.1 |
| xor | ^^ | 逻辑异或 | tr.dst[0:3]==0.6.29 xor tr.src[0:3] == 0.6.29 |
| not | ! | 逻辑非 | not llc |
| slice | [...] | 切片 | tr.dst[0:3] |
| set | in | 集合 | xxx

[官方文档过滤捕获](https://www.wireshark.org/docs/wsug_html_chunked/ChCapCaptureFilterSection.html)

[官方Wiki捕获过滤器](https://gitlab.com/wireshark/wireshark/-/wikis/CaptureFilters)

