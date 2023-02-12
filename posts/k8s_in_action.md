---
title: "K8s_in_action"
date: 2022-11-14T11:43:28+08:00
draft: true
tags: ["cicd"]
---

Dockerfile中ENTRYPOINT和CMD的区别
1. ENTRYPOINT 定义容器启动时被调用的可执行程序
2. CMD 指定传递给ENTRYPOINT的参数

shell与exec形式的区别, 主要在于命令是否在shell中被调用
1. shell形式 ENTRYPOINT node app.js
2. exec形式 ENTRYPOINT ["node", "app.js"]