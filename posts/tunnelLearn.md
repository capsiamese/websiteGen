---
title: "Tunnel原理学习"
date: 2023-02-22T13:48:45+08:00
draft: true
tags: ['golang', 'network']
---

# 原理
正常情况下本地计算机可以直接通过网络与其他服务器建立链接, 而由于某种原因这个过程要受到防火墙的检查,
防火墙会过滤掉不想让你发送的内容, 所以代理服务就有了, 大致原理就是本地和远程各有一个服务器充当代理
的角色, 他们之间的通讯是加密的所以可以绕过防火墙的检查, 本地的其他主机通过访问本地的代理服务器, 代理
服务器将流量加密转发到远程代理上, 然后远程代理解密流量再将流量发送到目标服务器上.

# socks5介绍
一般代理服务器都建立在socks5协议上, [socks5](https://zh.wikipedia.org/zh-cn/SOCKS)是tcp/ip栈上网络层的代理协议.
[rfc1928](https://www.rfc-editor.org/rfc/rfc1928)

## socks5实现
### 第一步
创建与socks5服务器的tcp连接后客户端需要先发送请求确认协议版本以及认证方式

|ver | nmethods | methods |
|----|----------|---------|
|1B  |  1B      | 1-255B  |

* ver是socks的版本, socks5这里应该是0x05
* nmethods是后面methods的长度
* methods是客户端支持的认证方式列表, 每个方法一个字节
  * 0x00  不用认证
  * 0x01  GSSAPI
  * 0x02  用户名, 密码认证
  * 0x03-0x7f由IANA分配  
    * 0x03  握手挑战认证协议
    * 0x04  未分派
    * 0x05  响应挑战认证协议
    * 0x06  传输层安全
    * 0x07  NDS认证
    * 0x08  多认证框架
    * 0x09  JSON参数块
    * ....
  * 0x8f - 0xfe为私人方法保留
  * 0xff  不可接受的方法

### 第二步
服务器从客户端提供的方法中选择一个并通过一下消息通知客户端
|ver | method|
|---| --|
|1B | 1B|

* ver是返回的socks版本
* method是服务器选中的方法, 如果是0xff则客户端要断开连接

之后客户端和服务器根据选定的认证发发执行相应的认证

### 第三步
如果使用用户名和密码认证后, 客户端发出用户名和密码
| 鉴定协议版本 | 用户名长度 | 用户名 | 密码长度 | 密码 |
|---          |-----------|--------|---------|-----|
|1B           | 1B        | dyn    |  1B     | dyn |

目前协议鉴定版本是0x01
服务器鉴定完后返回如下
| 鉴定协议版本 | 鉴定状态 |
|-------------|---------|
| 1B          |   1B    |
其中鉴定状态0x00表示成功  0x01表示失败


### 第三步
认证结束后客户端就可以发送请求消息了, 如果认证方法有特殊封装请求, 请求必须按照方法定义的方式进行封装.

socks5请求格式
|ver | cmd | rsv | atyp | dst.addr | dst.port |
| ---|---- | ----| ----| ----| ----|
|1B  | 1B  | 0x00| 1B  |  dyn | 2B |

* ver表示socks版本
* cmd表示命令码
  * 0x01 表示connect请求
  * 0x02 表示bind请求
  * 0x03 表示udp转发
* rsv 保留
* atyp 标识dst.addr的长度
  * 0x01 ipv4地址 dst.addr 4B
  * 0x03 域名, dst.addr的第一个字节表示域名长度, 后面是域名
  * 0x04 ipv6 dst.addr 16B

服务器按以下格式回复
| ver | rep | rsv | atyp | bnd.addr | bnd.port |
|-----|-----|-----|------|----------|----------|
|1B   | 1B  | 0x00| 1B   | dyn      | 2B       |

* rep应答字段
  * 0x00 成功
  * 0x01 普通socks服务器连接失败
  * 0x02 现有规则不允许连接
  * 0x03 网络不可达
  * 0x04 主机不可达
  * 0x05 连接被拒
  * 0x06 TTL超时
  * 0x07 不支持的命令
  * 0x08 不支持的地址类型
  * 0x09 - 0xFF 未定义




# References
[Lightsocks](https://github.com/gwuhaolin/lightsocks)
[你也能写一个Shadowsocks](https://wuhaolin.cn/2017/11/03/%E4%BD%A0%E4%B9%9F%E8%83%BD%E5%86%99%E4%B8%AA%20Shadowsocks/)
[gost](https://github.com/ginuerzh/gost)