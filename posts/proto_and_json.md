---
title: "protobuf序列化和反序列化json的坑"
date: 2023-04-11T20:30:23+08:00
draft: true
tags: []
---


jsonpb.Unmarshal使用的是protobuf tag中的json name而不是 json tag中的name